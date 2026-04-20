package services

import (
	"context"
	"errors"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"lovelion/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const (
	aiStatusPending    = "pending"
	aiStatusProcessing = "processing"
	aiStatusCompleted  = "completed"
	aiStatusFailed     = "failed"
)

// ImageDownloader is the minimum the worker needs from the storage layer.
// This allows tests to inject a fake without pulling in the S3 client.
type ImageDownloader interface {
	DownloadByURL(ctx context.Context, fullURL string) ([]byte, string, error)
}

// AIWorkerConfig controls polling cadence. Zero values get sensible defaults.
type AIWorkerConfig struct {
	PollInterval time.Duration // default 10s
	BatchSize    int           // default 5
	MaxRetries   int           // default 3
}

// AIWorker polls the transactions table for rows with ai_status='pending',
// calls the configured ReceiptExtractor, and writes back the result.
type AIWorker struct {
	db            *gorm.DB
	extractor     ReceiptExtractor
	textExtractor TextExtractor // optional; required only for no-image rows
	storage       ImageDownloader
	cfg           AIWorkerConfig
	// In-memory retry state per transaction ID. Resets on restart, which is
	// fine — startup recovery re-queues processing→pending anyway.
	retries map[string]int
	// retryAt tracks when a transaction is allowed to be retried. The worker
	// skips rows whose retryAt is still in the future, giving a simple
	// exponential backoff: 30s, 60s, 120s for attempts 1/2/3.
	retryAt map[string]time.Time
}

// NewAIWorker constructs a worker. Pass nil cfg fields to get defaults.
func NewAIWorker(db *gorm.DB, extractor ReceiptExtractor, storage ImageDownloader, cfg AIWorkerConfig) *AIWorker {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 10 * time.Second
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 5
	}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 3
	}
	return &AIWorker{
		db:        db,
		extractor: extractor,
		storage:   storage,
		cfg:       cfg,
		retries:   make(map[string]int),
		retryAt:   make(map[string]time.Time),
	}
}

// WithTextExtractor enables the text-based quick-entry path. If unset, rows
// without images will be marked failed instead of falling through to text.
func (w *AIWorker) WithTextExtractor(te TextExtractor) *AIWorker {
	w.textExtractor = te
	return w
}

// Run executes the worker loop until ctx is cancelled. It performs a startup
// recovery pass (processing→pending) so rows orphaned by a previous shutdown
// get re-queued.
func (w *AIWorker) Run(ctx context.Context) {
	if err := w.recoverStuck(ctx); err != nil {
		slog.Error("ai worker recovery failed", "error", err)
	}

	ticker := time.NewTicker(w.cfg.PollInterval)
	defer ticker.Stop()

	slog.Info("ai worker started", "poll_interval", w.cfg.PollInterval, "batch_size", w.cfg.BatchSize)

	// Run one pass immediately so a freshly-queued row doesn't wait a whole tick.
	w.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			slog.Info("ai worker shutting down")
			return
		case <-ticker.C:
			w.tick(ctx)
		}
	}
}

// recoverStuck re-queues rows left in processing by a prior shutdown. Since we
// run a single worker instance, any processing row at startup must be orphaned.
func (w *AIWorker) recoverStuck(ctx context.Context) error {
	result := w.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Where("ai_status = ?", aiStatusProcessing).
		Update("ai_status", aiStatusPending)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		slog.Info("ai worker recovered stuck rows", "count", result.RowsAffected)
	}
	return nil
}

// tick fetches pending rows and processes them one at a time.
func (w *AIWorker) tick(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}

	var pending []models.Transaction
	err := w.db.WithContext(ctx).
		Where("ai_status = ?", aiStatusPending).
		Order("created_at ASC").
		Limit(w.cfg.BatchSize).
		Find(&pending).Error
	if err != nil {
		slog.Error("ai worker pick pending failed", "error", err)
		return
	}

	now := time.Now()
	for _, txn := range pending {
		if ctx.Err() != nil {
			return
		}
		// Skip rows that are in backoff from a previous retryable failure.
		if t, ok := w.retryAt[txn.ID]; ok && now.Before(t) {
			continue
		}
		w.processOne(ctx, txn.ID, txn.Title)
	}
}

// processOne drives a single transaction through claim → call LLM → write-back.
// Returns without panicking on any error; errors are recorded on the row.
func (w *AIWorker) processOne(ctx context.Context, txnID, title string) {
	start := time.Now()
	log := slog.With("txn_id", txnID)

	// Stage 1: claim via conditional UPDATE.
	claim := w.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Where("id = ? AND ai_status = ?", txnID, aiStatusPending).
		Update("ai_status", aiStatusProcessing)
	if claim.Error != nil {
		log.Error("ai worker claim failed", "error", claim.Error)
		return
	}
	if claim.RowsAffected == 0 {
		log.Debug("ai worker row no longer pending, skipping")
		return
	}

	// Stage 2: decide image vs text based on whether any image is attached.
	hasImage, err := w.hasTransactionImage(ctx, txnID)
	if err != nil {
		log.Error("ai worker image lookup failed", "error", err)
		w.writeFailure(ctx, txnID, "辨識失敗")
		return
	}

	var result *ReceiptData
	if hasImage {
		imgBytes, mimeType, loadErr := w.loadFirstImage(ctx, txnID)
		if loadErr != nil {
			log.Warn("ai worker load image failed", "error", loadErr)
			w.writeFailure(ctx, txnID, friendlyExtractError(loadErr))
			return
		}
		result, err = w.extractor.Extract(ctx, imgBytes, mimeType)
	} else {
		if w.textExtractor == nil {
			log.Warn("ai worker no text extractor configured")
			w.writeFailure(ctx, txnID, "辨識服務未啟用")
			return
		}
		if strings.TrimSpace(title) == "" {
			w.writeFailure(ctx, txnID, "無文字可辨識")
			return
		}
		result, err = w.textExtractor.ExtractText(ctx, title)
	}

	if err != nil {
		// Don't write back if the context was cancelled mid-flight: let the
		// next startup's recovery pass pick it up.
		if errors.Is(err, context.Canceled) {
			log.Warn("ai worker cancelled mid-call", "elapsed", time.Since(start))
			return
		}

		if isRetryableError(err) {
			w.retries[txnID]++
			attempt := w.retries[txnID]
			if attempt <= w.cfg.MaxRetries {
				// Exponential backoff: 30s, 60s, 120s.
				backoff := time.Duration(30<<(attempt-1)) * time.Second
				w.retryAt[txnID] = time.Now().Add(backoff)
				log.Warn("ai worker retryable error, re-queuing",
					"error", err, "attempt", attempt, "max", w.cfg.MaxRetries,
					"backoff", backoff)
				w.writeRetry(ctx, txnID)
				return
			}
			log.Warn("ai worker max retries reached", "error", err, "attempts", attempt)
		}

		log.Warn("ai worker extract failed", "error", err)
		w.writeFailure(ctx, txnID, friendlyExtractError(err))
		delete(w.retries, txnID)
		delete(w.retryAt, txnID)
		return
	}

	// Stage 3: write success back inside a short db.Transaction.
	// For text-extraction rows we also overwrite the original raw input title
	// with the cleaned item name so the ledger reads naturally.
	if err := w.writeSuccess(ctx, txnID, result, !hasImage); err != nil {
		log.Error("ai worker write-back failed", "error", err)
		w.writeFailure(ctx, txnID, "failed to save result")
		delete(w.retries, txnID)
		delete(w.retryAt, txnID)
		return
	}

	delete(w.retries, txnID)
	delete(w.retryAt, txnID)
	log.Info("ai worker completed", "items", len(result.Items), "elapsed", time.Since(start))
}

// hasTransactionImage reports whether any image is attached to the given
// transaction. Used to pick the image- vs text-extraction path.
func (w *AIWorker) hasTransactionImage(ctx context.Context, txnID string) (bool, error) {
	var count int64
	err := w.db.WithContext(ctx).
		Model(&models.Image{}).
		Where("entity_type = ? AND entity_id = ?", "transaction", txnID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// loadFirstImage fetches the earliest (sort_order ASC) image for a transaction
// and downloads its bytes from storage.
func (w *AIWorker) loadFirstImage(ctx context.Context, txnID string) ([]byte, string, error) {
	var img models.Image
	err := w.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", "transaction", txnID).
		Order("sort_order ASC, created_at ASC").
		First(&img).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("no image attached")
		}
		return nil, "", err
	}

	data, ct, err := w.storage.DownloadByURL(ctx, img.FilePath)
	if err != nil {
		return nil, "", err
	}
	if ct == "" {
		ct = guessMimeFromURL(img.FilePath)
	}
	return data, ct, nil
}

// writeSuccess updates the transaction + replaces expense items in one tx.
// The WHERE ai_status='processing' guard lets a concurrent cancel cause the
// whole write to be a no-op.
func (w *AIWorker) writeSuccess(ctx context.Context, txnID string, data *ReceiptData, overwriteTitle bool) error {
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"ai_status": aiStatusCompleted,
			"ai_error":  "",
		}
		if data.Date != nil {
			updates["date"] = *data.Date
		}
		if overwriteTitle && len(data.Items) > 0 {
			if name := strings.TrimSpace(data.Items[0].Name); name != "" {
				updates["title"] = name
			}
		}

		result := tx.Model(&models.Transaction{}).
			Where("id = ? AND ai_status = ?", txnID, aiStatusProcessing).
			Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			// Cancelled between claim and now — discard everything.
			return nil
		}

		// Fetch the expense row to find the expense_id for items.
		var expense models.TransactionExpense
		if err := tx.Where("transaction_id = ?", txnID).First(&expense).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Non-expense transaction shouldn't get here, but be defensive.
				return nil
			}
			return err
		}

		// Replace items entirely: delete existing, then insert extracted.
		if err := tx.Where("expense_id = ?", expense.ID).Delete(&models.TransactionExpenseItem{}).Error; err != nil {
			return err
		}

		items, totalAmount := buildExpenseItems(expense.ID, convertReceiptItems(data.Items))
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		// Refresh total_amount to match the new items.
		return tx.Model(&models.Transaction{}).
			Where("id = ?", txnID).
			Update("total_amount", totalAmount).Error
	})
}

// writeRetry sets the row back to pending so the next poll picks it up.
// Uses the same conditional WHERE so a racing cancel wins.
func (w *AIWorker) writeRetry(ctx context.Context, txnID string) {
	err := w.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Where("id = ? AND ai_status = ?", txnID, aiStatusProcessing).
		Update("ai_status", aiStatusPending).Error
	if err != nil {
		slog.Error("ai worker retry re-queue failed", "txn_id", txnID, "error", err)
	}
}

// writeFailure flips the row to failed with an error message. Uses the same
// conditional WHERE so a racing cancel wins.
func (w *AIWorker) writeFailure(ctx context.Context, txnID, message string) {
	err := w.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Where("id = ? AND ai_status = ?", txnID, aiStatusProcessing).
		Updates(map[string]interface{}{
			"ai_status": aiStatusFailed,
			"ai_error":  truncateError(message, 500),
		}).Error
	if err != nil {
		slog.Error("ai worker mark failed", "txn_id", txnID, "error", err)
	}
}

// convertReceiptItems maps extractor output to the shape buildExpenseItems expects.
func convertReceiptItems(in []ReceiptItem) []ExpenseItemInput {
	out := make([]ExpenseItemInput, 0, len(in))
	for _, it := range in {
		qty := it.Quantity
		if qty.IsZero() {
			qty = decimal.NewFromInt(1)
		}
		out = append(out, ExpenseItemInput{
			Name:      it.Name,
			UnitPrice: it.UnitPrice,
			Quantity:  qty,
			Discount:  decimal.Zero,
		})
	}
	return out
}

// friendlyExtractError converts raw extractor errors into short user-facing strings.
// The full error is still in server logs via slog.Warn above.
func friendlyExtractError(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "context deadline exceeded"):
		return "辨識逾時，請稍後再試"
	case strings.Contains(msg, "gemini http 429"):
		return "辨識服務忙碌中"
	case strings.Contains(msg, "gemini http"):
		return "辨識服務錯誤"
	case strings.Contains(msg, "no image"):
		return "找不到可辨識的圖片"
	case strings.Contains(msg, "parse receipt json"):
		return "辨識結果格式錯誤"
	default:
		return "辨識失敗"
	}
}

// isRetryableError returns true for transient errors that are likely to resolve
// on their own (Gemini overload, rate limit, network timeout).
func isRetryableError(err error) bool {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "gemini http 503"):
		return true
	case strings.Contains(msg, "gemini http 429"):
		return true
	case strings.Contains(msg, "context deadline exceeded"):
		return true
	default:
		return false
	}
}

func truncateError(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// guessMimeFromURL is a cheap fallback when R2 doesn't return a content-type.
func guessMimeFromURL(url string) string {
	ext := strings.ToLower(filepath.Ext(url))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}
