package handlers

import (
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"lovelion/internal/middleware"
	"lovelion/internal/models"
	"lovelion/internal/services"
	"lovelion/internal/utils/errorx"

	"github.com/bbrks/go-blurhash"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	maxImageSize  = 5 * 1024 * 1024 // 5MB
	maxImageCount = 10
)

type ExpenseHandler struct {
	svc         *services.TransactionService
	aiRateLimit *middleware.AIRateLimiter
}

func NewExpenseHandler(svc *services.TransactionService, aiRateLimit *middleware.AIRateLimiter) *ExpenseHandler {
	return &ExpenseHandler{svc: svc, aiRateLimit: aiRateLimit}
}

type ExpenseItemRequest struct {
	Name      string          `json:"name" binding:"required"`
	UnitPrice decimal.Decimal `json:"unit_price"`
	Quantity  decimal.Decimal `json:"quantity"`
	Discount  decimal.Decimal `json:"discount"`
}

type DebtRequest struct {
	PayerName  string          `json:"payer_name" binding:"required"`
	PayeeName  string          `json:"payee_name" binding:"required"`
	Amount     decimal.Decimal `json:"amount"`
	IsSpotPaid bool            `json:"is_spot_paid"`
}

type ExpenseDetailRequest struct {
	Category      string               `json:"category"`
	ExchangeRate  decimal.Decimal      `json:"exchange_rate"`
	BillingAmount decimal.Decimal      `json:"billing_amount"`
	HandlingFee   decimal.Decimal      `json:"handling_fee"`
	PaymentMethod string               `json:"payment_method"`
	Items         []ExpenseItemRequest `json:"items"`
}

type CreateExpenseRequest struct {
	Date      *time.Time           `json:"date"`
	Currency  string               `json:"currency"`
	Title     string               `json:"title"`
	Note      string               `json:"note"`
	Expense   ExpenseDetailRequest `json:"expense"`
	Debts     []DebtRequest        `json:"debts"`
	AIExtract bool                 `json:"ai_extract"`
}

type UpdateExpenseRequest struct {
	Date      *time.Time           `json:"date"`
	Currency  string               `json:"currency"`
	Title     string               `json:"title"`
	Note      string               `json:"note"`
	Expense   ExpenseDetailRequest `json:"expense"`
	Debts     []DebtRequest        `json:"debts"`
	AIExtract bool                 `json:"ai_extract"`
}

func toExpenseItemInputs(reqs []ExpenseItemRequest) []services.ExpenseItemInput {
	if reqs == nil {
		return nil
	}
	inputs := make([]services.ExpenseItemInput, len(reqs))
	for i, r := range reqs {
		inputs[i] = services.ExpenseItemInput{
			Name:      r.Name,
			UnitPrice: r.UnitPrice,
			Quantity:  r.Quantity,
			Discount:  r.Discount,
		}
	}
	return inputs
}

func toDebtInputs(reqs []DebtRequest) []services.DebtInput {
	if reqs == nil {
		return nil
	}
	inputs := make([]services.DebtInput, len(reqs))
	for i, r := range reqs {
		inputs[i] = services.DebtInput{
			PayerName:  r.PayerName,
			PayeeName:  r.PayeeName,
			Amount:     r.Amount,
			IsSpotPaid: r.IsSpotPaid,
		}
	}
	return inputs
}

// Create handles both application/json (legacy) and multipart/form-data
// (with optional images and ai_extract flag) content types. The multipart
// form has two fields:
//   - data:   JSON string matching CreateExpenseRequest
//   - images: 0..N files (jpg/jpeg/png)
func (h *ExpenseHandler) Create(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)

	req, images, err := parseCreateExpenseRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.AIExtract {
		if len(images) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "AI extraction requires at least one image"})
			return
		}
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if !h.aiRateLimit.Allow(userID.(uuid.UUID)) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Daily AI extraction limit reached"})
			return
		}
	}

	txn, err := h.svc.CreateExpense(c.Request.Context(), space.ID, services.CreateExpenseInput{
		Date:     req.Date,
		Currency: req.Currency,
		Title:    req.Title,
		Note:     req.Note,
		Expense: services.ExpenseInput{
			Category:      req.Expense.Category,
			ExchangeRate:  req.Expense.ExchangeRate,
			BillingAmount: req.Expense.BillingAmount,
			HandlingFee:   req.Expense.HandlingFee,
			PaymentMethod: req.Expense.PaymentMethod,
			Items:         toExpenseItemInputs(req.Expense.Items),
		},
		Debts:     toDebtInputs(req.Debts),
		Images:    images,
		AIExtract: req.AIExtract,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, txn)
}

// parseCreateExpenseRequest extracts the request body and any attached images
// from either a JSON payload or a multipart form.
func parseCreateExpenseRequest(c *gin.Context) (CreateExpenseRequest, []services.ImageUpload, error) {
	var req CreateExpenseRequest

	contentType := c.ContentType()
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		if err := c.ShouldBindJSON(&req); err != nil {
			return req, nil, err
		}
		return req, nil, nil
	}

	// multipart: `data` holds the JSON payload, `images` holds files.
	dataField := c.PostForm("data")
	if dataField == "" {
		return req, nil, errorx.New("BAD_REQUEST", "data field is required")
	}
	if err := json.Unmarshal([]byte(dataField), &req); err != nil {
		return req, nil, err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return req, nil, err
	}
	files := form.File["images"]
	if len(files) > maxImageCount {
		return req, nil, errorx.New("BAD_REQUEST", "too many images")
	}

	uploads, err := readImageUploads(files)
	if err != nil {
		return req, nil, err
	}
	return req, uploads, nil
}

// readImageUploads reads each file into memory, validates size + extension,
// and computes a blurhash. Returns a slice in the original form order.
func readImageUploads(files []*multipart.FileHeader) ([]services.ImageUpload, error) {
	uploads := make([]services.ImageUpload, 0, len(files))
	for _, fh := range files {
		if fh.Size > maxImageSize {
			return nil, errorx.New("BAD_REQUEST", "image exceeds 5MB: "+fh.Filename)
		}
		ext := strings.ToLower(filepath.Ext(fh.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return nil, errorx.New("BAD_REQUEST", "only jpg/jpeg/png allowed: "+fh.Filename)
		}

		f, err := fh.Open()
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			return nil, err
		}

		// Best-effort blurhash; non-fatal if the image can't be decoded.
		var blurHash string
		if decoded, _, decodeErr := image.Decode(newBytesReader(body)); decodeErr == nil {
			if h, hashErr := blurhash.Encode(4, 3, decoded); hashErr == nil {
				blurHash = h
			}
		} else {
			slog.Warn("blurhash decode failed", "file", fh.Filename, "error", decodeErr)
		}

		uploads = append(uploads, services.ImageUpload{
			FileName:    fh.Filename,
			Body:        body,
			ContentType: fh.Header.Get("Content-Type"),
			BlurHash:    blurHash,
		})
	}
	return uploads, nil
}

// newBytesReader wraps a byte slice so image.Decode can consume it without
// pulling in an extra bytes import via strings.NewReader.
func newBytesReader(b []byte) io.Reader {
	return &bytesReader{data: b}
}

type bytesReader struct {
	data []byte
	pos  int
}

func (r *bytesReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func (h *ExpenseHandler) Update(c *gin.Context) {
	spaceVal, _ := c.Get("space")
	space := spaceVal.(*models.Space)
	txnID := c.Param("txn_id")

	var req UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Re-running AI on a failed row needs to go through the same rate limit
	// as the create flow.
	if req.AIExtract {
		userID, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if !h.aiRateLimit.Allow(userID.(uuid.UUID)) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Daily AI extraction limit reached"})
			return
		}
	}

	txn, err := h.svc.UpdateExpense(c.Request.Context(), txnID, space.ID, services.UpdateExpenseInput{
		Date:     req.Date,
		Currency: req.Currency,
		Title:    req.Title,
		Note:     req.Note,
		Expense: services.ExpenseInput{
			Category:      req.Expense.Category,
			ExchangeRate:  req.Expense.ExchangeRate,
			BillingAmount: req.Expense.BillingAmount,
			HandlingFee:   req.Expense.HandlingFee,
			PaymentMethod: req.Expense.PaymentMethod,
			Items:         toExpenseItemInputs(req.Expense.Items),
		},
		Debts:     toDebtInputs(req.Debts),
		AIExtract: req.AIExtract,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, txn)
}
