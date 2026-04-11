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
}

// AIWorker polls the transactions table for rows with ai_status='pending',
// calls the configured ReceiptExtractor, and writes back the result.
type AIWorker struct {
	db        *gorm.DB
	extractor ReceiptExtractor
	storage   ImageDownloader
	cfg       AIWorkerConfig
}

// NewAIWorker constructs a worker. Pass nil cfg fields to get defaults.
func NewAIWorker(db *gorm.DB, extractor ReceiptExtractor, storage ImageDownloader, cfg AIWorkerConfig) *AIWorker {
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = 10 * time.Second
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 5
	}
	return &AIWorker{
		db:        db,
		extractor: extractor,
		storage:   storage,
		cfg:       cfg,
	}
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

	for _, txn := range pending {
		if ctx.Err() != nil {
			return
		}
		w.processOne(ctx, txn.ID)
	}
}

// processOne drives a single transaction through claim → call LLM → write-back.
// Returns without panicking on any error; errors are recorded on the row.
func (w *AIWorker) processOne(ctx context.Context, txnID string) {
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

	// Stage 2: load the first image and call the LLM with no DB connection held.
	imgBytes, mimeType, err := w.loadFirstImage(ctx, txnID)
	if err != nil {
		log.Warn("ai worker load image failed", "error", err)
		w.writeFailure(ctx, txnID, friendlyExtractError(err))
		return
	}

	result, err := w.extractor.Extract(ctx, imgBytes, mimeType)
	if err != nil {
		// Don't write back if the context was cancelled mid-flight: let the
		// next startup's recovery pass pick it up.
		if errors.Is(err, context.Canceled) {
			log.Warn("ai worker cancelled mid-call", "elapsed", time.Since(start))
			return
		}
		log.Warn("ai worker extract failed", "error", err)
		w.writeFailure(ctx, txnID, friendlyExtractError(err))
		return
	}

	// Stage 3: write success back inside a short db.Transaction.
	if err := w.writeSuccess(ctx, txnID, result); err != nil {
		log.Error("ai worker write-back failed", "error", err)
		w.writeFailure(ctx, txnID, "failed to save result")
		return
	}

	log.Info("ai worker completed", "items", len(result.Items), "elapsed", time.Since(start))
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
func (w *AIWorker) writeSuccess(ctx context.Context, txnID string, data *ReceiptData) error {
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"ai_status": aiStatusCompleted,
			"ai_error":  "",
		}
		if data.Date != nil {
			updates["date"] = *data.Date
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
