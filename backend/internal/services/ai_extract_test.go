package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeGemini wraps an httptest.Server so tests can drive a canned response.
type fakeGemini struct {
	server      *httptest.Server
	lastRequest geminiRequest
	lastPath    string
}

func newFakeGemini(t *testing.T, status int, respBody string) *fakeGemini {
	t.Helper()
	f := &fakeGemini{}
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &f.lastRequest)
		f.lastPath = r.URL.Path + "?" + r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(respBody))
	}))
	t.Cleanup(f.server.Close)
	return f
}

// wrapCandidate builds a Gemini response whose candidate text field contains
// the given JSON string (as a structured output would).
func wrapCandidate(t *testing.T, inner string) string {
	t.Helper()
	resp := geminiResponse{
		Candidates: []struct {
			Content      geminiContent `json:"content"`
			FinishReason string        `json:"finishReason"`
		}{{
			Content: geminiContent{
				Parts: []geminiPart{{Text: inner}},
			},
			FinishReason: "STOP",
		}},
	}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	return string(b)
}

func TestGeminiExtract_Success(t *testing.T) {
	canned := wrapCandidate(t, `{
        "date": "2025-07-14",
        "items": [
            {"name": "牛肉麵", "unit_price": 180, "quantity": 1},
            {"name": "飲料",   "unit_price": 40,  "quantity": 2}
        ]
    }`)
	fake := newFakeGemini(t, 200, canned)

	ext := NewGeminiReceiptExtractor("test-key", "gemini-2.5-flash", fake.server.URL)
	data, err := ext.Extract(context.Background(), []byte("fake-image-bytes"), "image/jpeg")

	require.NoError(t, err)
	require.NotNil(t, data)
	require.NotNil(t, data.Date)
	assert.Equal(t, "2025-07-14", data.Date.Format("2006-01-02"))

	require.Len(t, data.Items, 2)
	assert.Equal(t, "牛肉麵", data.Items[0].Name)
	assert.True(t, decimal.NewFromInt(180).Equal(data.Items[0].UnitPrice))
	assert.True(t, decimal.NewFromInt(1).Equal(data.Items[0].Quantity))
	assert.Equal(t, "飲料", data.Items[1].Name)
	assert.True(t, decimal.NewFromInt(2).Equal(data.Items[1].Quantity))

	// Verify request shape: model in path, api key as query, inline base64 image.
	assert.Contains(t, fake.lastPath, "gemini-2.5-flash:generateContent")
	assert.Contains(t, fake.lastPath, "key=test-key")
	require.NotNil(t, fake.lastRequest.SystemInstruction)
	require.Len(t, fake.lastRequest.Contents, 1)
	parts := fake.lastRequest.Contents[0].Parts
	require.GreaterOrEqual(t, len(parts), 2)
	var foundInline bool
	for _, p := range parts {
		if p.InlineData != nil {
			foundInline = true
			assert.Equal(t, "image/jpeg", p.InlineData.MimeType)
			assert.NotEmpty(t, p.InlineData.Data)
		}
	}
	assert.True(t, foundInline, "expected inlineData part")
	assert.Equal(t, "application/json", fake.lastRequest.GenerationConfig.ResponseMimeType)
}

func TestGeminiExtract_NullDate(t *testing.T) {
	canned := wrapCandidate(t, `{"date": null, "items": [{"name":"A","unit_price":10,"quantity":1}]}`)
	fake := newFakeGemini(t, 200, canned)

	ext := NewGeminiReceiptExtractor("k", "", fake.server.URL)
	data, err := ext.Extract(context.Background(), []byte("img"), "image/png")
	require.NoError(t, err)
	assert.Nil(t, data.Date)
	require.Len(t, data.Items, 1)
}

func TestGeminiExtract_DiscountItem(t *testing.T) {
	canned := wrapCandidate(t, `{
        "items": [
            {"name":"商品","unit_price":100,"quantity":1},
            {"name":"折扣","unit_price":-20,"quantity":1}
        ]
    }`)
	fake := newFakeGemini(t, 200, canned)
	ext := NewGeminiReceiptExtractor("k", "", fake.server.URL)
	data, err := ext.Extract(context.Background(), []byte("img"), "image/jpeg")
	require.NoError(t, err)
	require.Len(t, data.Items, 2)
	assert.Equal(t, "折扣", data.Items[1].Name)
	assert.True(t, data.Items[1].UnitPrice.IsNegative())
}

func TestGeminiExtract_HTTPError(t *testing.T) {
	fake := newFakeGemini(t, 500, `{"error":{"code":500,"message":"boom","status":"INTERNAL"}}`)
	ext := NewGeminiReceiptExtractor("k", "", fake.server.URL)
	_, err := ext.Extract(context.Background(), []byte("img"), "image/jpeg")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gemini http 500")
}

func TestGeminiExtract_EmptyCandidates(t *testing.T) {
	fake := newFakeGemini(t, 200, `{"candidates":[]}`)
	ext := NewGeminiReceiptExtractor("k", "", fake.server.URL)
	_, err := ext.Extract(context.Background(), []byte("img"), "image/jpeg")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no candidates")
}

func TestGeminiExtract_MalformedJSON(t *testing.T) {
	canned := wrapCandidate(t, `not-json`)
	fake := newFakeGemini(t, 200, canned)
	ext := NewGeminiReceiptExtractor("k", "", fake.server.URL)
	_, err := ext.Extract(context.Background(), []byte("img"), "image/jpeg")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse receipt json")
}

func TestGeminiExtract_MissingAPIKey(t *testing.T) {
	ext := NewGeminiReceiptExtractor("", "", "")
	_, err := ext.Extract(context.Background(), []byte("img"), "image/jpeg")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "api key")
}

func TestGeminiExtract_EmptyImage(t *testing.T) {
	ext := NewGeminiReceiptExtractor("k", "", "")
	_, err := ext.Extract(context.Background(), []byte{}, "image/jpeg")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty image")
}

func TestGeminiExtract_ContextCancelled(t *testing.T) {
	// Server that sleeps longer than our ctx timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	ext := NewGeminiReceiptExtractor("k", "", server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, err := ext.Extract(ctx, []byte("img"), "image/jpeg")
	require.Error(t, err)
	// Error should mention deadline or context
	assert.True(t, strings.Contains(err.Error(), "context") || strings.Contains(err.Error(), "deadline"))
}
