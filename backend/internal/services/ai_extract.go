package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// ReceiptData is the parsed output of a receipt extraction call.
type ReceiptData struct {
	Date          *time.Time
	Title         string // store name (image) or cleaned description (text)
	Category      string // best-match from ExtractHints.Categories
	PaymentMethod string // best-match from ExtractHints.PaymentMethods
	Items         []ReceiptItem
}

// ReceiptItem represents a single line item on a receipt.
type ReceiptItem struct {
	Name      string
	UnitPrice decimal.Decimal
	Quantity  decimal.Decimal
}

// ExtractHints carries space-specific option lists so the LLM can pick from
// existing values rather than inventing its own.
type ExtractHints struct {
	Categories     []string
	PaymentMethods []string
}

// promptFragment builds the user-message addendum that lists available options.
func (h ExtractHints) promptFragment() string {
	var parts []string
	if len(h.Categories) > 0 {
		parts = append(parts, fmt.Sprintf("可用的分類清單：%s", strings.Join(h.Categories, "、")))
	}
	if len(h.PaymentMethods) > 0 {
		parts = append(parts, fmt.Sprintf("可用的付款方式清單：%s", strings.Join(h.PaymentMethods, "、")))
	}
	return strings.Join(parts, "\n")
}

// ReceiptExtractor abstracts the underlying LLM provider so we can swap
// implementations (or inject fakes in tests).
type ReceiptExtractor interface {
	Extract(ctx context.Context, image []byte, mimeType string, hints ExtractHints) (*ReceiptData, error)
}

// TextExtractor parses a single free-form text line (e.g. "停車費 100") into the
// same ReceiptData shape used by image extraction.
type TextExtractor interface {
	ExtractText(ctx context.Context, text string, hints ExtractHints) (*ReceiptData, error)
}

// --- Gemini implementation ---

const (
	geminiDefaultBaseURL = "https://generativelanguage.googleapis.com"
	geminiCallTimeout    = 30 * time.Second
)

//nolint:lll // prompt is intentionally a single literal block
const geminiSystemPrompt = `你是發票辨識助手。請從圖片擷取消費資訊，並以指定的 JSON Schema 回傳。

規則：
- title 為店家名稱。若發票上無法判讀店家名稱，填空字串。
- date 取發票上的消費日期與時間，格式 YYYY-MM-DD HH:mm。若只有日期沒有時間則填 YYYY-MM-DD。若無法判讀填 null。
- items 依發票上的順序列出。
- unit_price 是單價（不是小計）。quantity 是數量。
- 若有整單折扣，加一筆 name="折扣" 的品項，unit_price 為負數，quantity=1。
- category：根據消費內容判斷最適合的分類。若使用者有提供可用分類清單，優先從清單中選擇；若都不符合，可自行填寫。若無法判斷填空字串。
- payment_method：若發票上有標示付款方式（如信用卡、現金、Line Pay 等），請填寫。若使用者有提供可用付款方式清單，優先從清單中選擇；若都不符合，可自行填寫。若無法判讀填空字串。
- 不要輸出 total，呼叫端會自行計算。`

//nolint:lll // prompt is intentionally a single literal block
const geminiTextSystemPrompt = `你是記帳輸入解析助手。使用者會給你一小段中文／英文混合的簡短文字（例如「停車費 100」、「昨天 小七 買咖啡 55」、「午餐 250 刷卡」），請從中擷取消費資訊，並以指定的 JSON Schema 回傳。

規則：
- title 為店家名稱，若文字中有提及店家（如「小七」、「全聯」），填入店家名稱；若無填空字串。
- 將整筆消費視為「一個」品項。items 陣列只會有 1 筆。
- item.name 為消費名稱（例如「停車費」、「午餐」、「咖啡」），去除金額、日期與店家字樣。若僅有金額沒有名稱，填「未命名」。
- item.unit_price 為金額數字；quantity 固定為 1。
- 若文字中指出相對日期（例如「昨天」、「前天」），請以今天的日期推算並輸出 YYYY-MM-DD HH:mm（無時間則 YYYY-MM-DD）；若無法判讀日期則填 null。
- category：根據消費內容判斷最適合的分類。若使用者有提供可用分類清單，優先從清單中選擇；若都不符合，可自行填寫。若無法判斷填空字串。
- payment_method：若文字中有提及付款方式（如「刷卡」、「現金」、「Line Pay」），請填寫。若使用者有提供可用付款方式清單，優先從清單中選擇；若都不符合，可自行填寫。若無法判讀填空字串。
- 若完全無法解析出金額，仍請回傳 items=[] — 呼叫端會視為失敗。
- 不要輸出 total。`

// Response schema used to force structured JSON output from Gemini.
var geminiResponseSchema = json.RawMessage(`{
  "type": "object",
  "properties": {
    "title":          { "type": "string" },
    "date":           { "type": "string", "nullable": true },
    "category":       { "type": "string" },
    "payment_method": { "type": "string" },
    "items": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "name":       { "type": "string" },
          "unit_price": { "type": "number" },
          "quantity":   { "type": "number" }
        },
        "required": ["name", "unit_price", "quantity"]
      }
    }
  },
  "required": ["items"]
}`)

// GeminiReceiptExtractor calls Google Gemini generateContent via plain HTTP.
type GeminiReceiptExtractor struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

// NewGeminiReceiptExtractor constructs a Gemini-backed extractor.
// baseURL may be overridden to point at a test server.
func NewGeminiReceiptExtractor(apiKey, model, baseURL string) *GeminiReceiptExtractor {
	if baseURL == "" {
		baseURL = geminiDefaultBaseURL
	}
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &GeminiReceiptExtractor{
		apiKey:  apiKey,
		model:   model,
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: geminiCallTimeout + 5*time.Second},
	}
}

// --- Gemini wire types ---

type geminiRequest struct {
	SystemInstruction *geminiContent        `json:"systemInstruction,omitempty"`
	Contents          []geminiContent       `json:"contents"`
	GenerationConfig  geminiGenerationCfg   `json:"generationConfig"`
	SafetySettings    []geminiSafetySetting `json:"safetySettings,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text       string            `json:"text,omitempty"`
	InlineData *geminiInlineData `json:"inlineData,omitempty"`
}

type geminiInlineData struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"` // base64
}

type geminiGenerationCfg struct {
	ResponseMimeType string          `json:"responseMimeType"`
	ResponseSchema   json.RawMessage `json:"responseSchema"`
	Temperature      float64         `json:"temperature,omitempty"`
}

type geminiSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type geminiResponse struct {
	Candidates []struct {
		Content      geminiContent `json:"content"`
		FinishReason string        `json:"finishReason"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

type receiptJSONPayload struct {
	Title         string  `json:"title"`
	Date          *string `json:"date"`
	Category      string  `json:"category"`
	PaymentMethod string  `json:"payment_method"`
	Items         []struct {
		Name      string          `json:"name"`
		UnitPrice decimal.Decimal `json:"unit_price"`
		Quantity  decimal.Decimal `json:"quantity"`
	} `json:"items"`
}

// Extract sends the image to Gemini and parses the structured response.
func (g *GeminiReceiptExtractor) Extract(ctx context.Context, image []byte, mimeType string, hints ExtractHints) (*ReceiptData, error) {
	if len(image) == 0 {
		return nil, errors.New("empty image")
	}

	userText := "請辨識這張發票。"
	if extra := hints.promptFragment(); extra != "" {
		userText += "\n" + extra
	}

	reqBody := geminiRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{{Text: geminiSystemPrompt}},
		},
		Contents: []geminiContent{{
			Role: "user",
			Parts: []geminiPart{
				{Text: userText},
				{InlineData: &geminiInlineData{
					MimeType: mimeType,
					Data:     base64.StdEncoding.EncodeToString(image),
				}},
			},
		}},
		GenerationConfig: geminiGenerationCfg{
			ResponseMimeType: "application/json",
			ResponseSchema:   geminiResponseSchema,
			Temperature:      0.1,
		},
	}

	return g.callAndParse(ctx, reqBody)
}

// ExtractText sends a short text line to Gemini and parses the same structured
// shape as Extract. Used by the quick-text-entry flow.
func (g *GeminiReceiptExtractor) ExtractText(ctx context.Context, text string, hints ExtractHints) (*ReceiptData, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.New("empty text")
	}

	today := time.Now().Format("2006-01-02")
	userMsg := fmt.Sprintf("今天是 %s。請解析以下這筆記帳輸入：\n%s", today, text)
	if extra := hints.promptFragment(); extra != "" {
		userMsg += "\n" + extra
	}

	reqBody := geminiRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{{Text: geminiTextSystemPrompt}},
		},
		Contents: []geminiContent{{
			Role: "user",
			Parts: []geminiPart{
				{Text: userMsg},
			},
		}},
		GenerationConfig: geminiGenerationCfg{
			ResponseMimeType: "application/json",
			ResponseSchema:   geminiResponseSchema,
			Temperature:      0.1,
		},
	}

	result, err := g.callAndParse(ctx, reqBody)
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, errors.New("parse text: no items extracted")
	}
	return result, nil
}

// callAndParse is the common Gemini HTTP call + JSON parse used by both
// image and text extraction paths.
func (g *GeminiReceiptExtractor) callAndParse(ctx context.Context, reqBody geminiRequest) (*ReceiptData, error) {
	if g.apiKey == "" {
		return nil, errors.New("gemini api key not configured")
	}

	callCtx, cancel := context.WithTimeout(ctx, geminiCallTimeout)
	defer cancel()

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", g.baseURL, g.model, g.apiKey)
	req, err := http.NewRequestWithContext(callCtx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini call: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini http %d: %s", resp.StatusCode, truncate(string(body), 300))
	}

	var parsed geminiResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("gemini api error: %s", parsed.Error.Message)
	}
	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("gemini returned no candidates")
	}

	text := parsed.Candidates[0].Content.Parts[0].Text
	if text == "" {
		return nil, errors.New("gemini returned empty text")
	}

	var receipt receiptJSONPayload
	if err := json.Unmarshal([]byte(text), &receipt); err != nil {
		return nil, fmt.Errorf("parse receipt json: %w (raw=%s)", err, truncate(text, 200))
	}

	out := &ReceiptData{
		Title:         strings.TrimSpace(receipt.Title),
		Category:      strings.TrimSpace(receipt.Category),
		PaymentMethod: strings.TrimSpace(receipt.PaymentMethod),
		Items:         make([]ReceiptItem, 0, len(receipt.Items)),
	}

	if receipt.Date != nil && *receipt.Date != "" {
		for _, layout := range []string{"2006-01-02 15:04", "2006-01-02"} {
			if t, err := time.Parse(layout, *receipt.Date); err == nil {
				out.Date = &t
				break
			}
		}
	}

	for _, it := range receipt.Items {
		name := strings.TrimSpace(it.Name)
		if name == "" {
			continue
		}
		qty := it.Quantity
		if qty.IsZero() {
			qty = decimal.NewFromInt(1)
		}
		out.Items = append(out.Items, ReceiptItem{
			Name:      name,
			UnitPrice: it.UnitPrice,
			Quantity:  qty,
		})
	}

	return out, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
