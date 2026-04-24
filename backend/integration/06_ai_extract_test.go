//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"testing"
	"time"
)

// tinyPNG builds a 1x1 red PNG for multipart uploads. Generating it at runtime
// keeps the test self-contained — no fixture files on disk.
func tinyPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode tiny png: %v", err)
	}
	return buf.Bytes()
}

// expenseBodyJSON returns the `data` field payload used with multipart create,
// mirroring the shape of a normal CreateExpenseRequest. When aiExtract is true
// the backend expects at least one image to be attached alongside.
func expenseBodyJSON(t *testing.T, title string, aiExtract bool) string {
	t.Helper()
	body := map[string]interface{}{
		"title":    title,
		"date":     time.Now().Format(time.RFC3339),
		"currency": "TWD",
		"expense": map[string]interface{}{
			"category":       "餐飲",
			"exchange_rate":  1,
			"payment_method": "現金",
			"items": []map[string]interface{}{
				{"name": "手動品項", "unit_price": 100, "quantity": 1},
			},
		},
	}
	if aiExtract {
		body["ai_extract"] = true
	}
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	return string(raw)
}

func TestAIReceiptExtraction(t *testing.T) {
	token := loginUser(t, "dev", "dev123")
	personalID := findSpaceByName(t, token, "日常開銷")

	// IDs captured in one sub-test and reused later — the test uses a single
	// pending transaction to exercise the PUT-blocks-on-pending flow and the
	// cancel-then-edit flow in sequence.
	var pendingTxnID string

	t.Run("multipart create with image attached", func(t *testing.T) {
		ae := authExpect(t, token)
		obj := ae.POST("/api/spaces/{id}/expenses", personalID).
			WithMultipart().
			WithFormField("data", expenseBodyJSON(t, "Multipart 無 AI", false)).
			WithFileBytes("images", "receipt.png", tinyPNG(t)).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		obj.Value("title").IsEqual("Multipart 無 AI")
		// ai_status has json:",omitempty" on a *string so it disappears entirely
		// when the column is NULL — the row should look identical to a normal
		// manually-entered expense.
		obj.NotContainsKey("ai_status")
		obj.Value("images").Array().Length().IsEqual(1)
	})

	t.Run("multipart create with ai_extract=true seeds ai_status=pending", func(t *testing.T) {
		ae := authExpect(t, token)
		obj := ae.POST("/api/spaces/{id}/expenses", personalID).
			WithMultipart().
			WithFormField("data", expenseBodyJSON(t, "AI 處理中", true)).
			WithFileBytes("images", "receipt.png", tinyPNG(t)).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		obj.Value("ai_status").IsEqual("pending")
		obj.Value("images").Array().Length().IsEqual(1)
		pendingTxnID = obj.Value("id").String().Raw()
	})

	t.Run("ai_extract=true without image or title is rejected", func(t *testing.T) {
		// Text-based quick entry (added after the original image-only AI flow)
		// lets the worker run on a non-empty title when no image is attached.
		// So AI extract is only rejected when BOTH title and images are empty.
		ae := authExpect(t, token)
		ae.POST("/api/spaces/{id}/expenses", personalID).
			WithMultipart().
			WithFormField("data", expenseBodyJSON(t, "", true)).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("PUT on a pending transaction is blocked with 409", func(t *testing.T) {
		if pendingTxnID == "" {
			t.Fatal("prior sub-test did not produce a pending transaction id")
		}
		ae := authExpect(t, token)
		ae.PUT("/api/spaces/{id}/expenses/{txn_id}", personalID, pendingTxnID).
			WithJSON(map[string]interface{}{
				"title":    "嘗試編輯",
				"currency": "TWD",
				"expense": map[string]interface{}{
					"category":      "餐飲",
					"exchange_rate": 1,
					"items": []map[string]interface{}{
						{"name": "無效", "unit_price": 99, "quantity": 1},
					},
				},
			}).
			Expect().
			Status(http.StatusConflict)
	})

	t.Run("ai-cancel on pending clears the ai_status and reopens edits", func(t *testing.T) {
		if pendingTxnID == "" {
			t.Fatal("prior sub-test did not produce a pending transaction id")
		}
		ae := authExpect(t, token)
		ae.POST("/api/spaces/{id}/transactions/{txn_id}/ai-cancel", personalID, pendingTxnID).
			Expect().
			Status(http.StatusOK)

		// The subsequent GET should no longer carry ai_status (NULL + omitempty).
		detail := ae.GET("/api/spaces/{id}/transactions/{txn_id}", personalID, pendingTxnID).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		detail.NotContainsKey("ai_status")

		// And a normal PUT now succeeds — the row is back to being a plain
		// manually-entered transaction.
		ae.PUT("/api/spaces/{id}/expenses/{txn_id}", personalID, pendingTxnID).
			WithJSON(map[string]interface{}{
				"title":    "取消後已編輯",
				"currency": "TWD",
				"expense": map[string]interface{}{
					"category":      "餐飲",
					"exchange_rate": 1,
					"items": []map[string]interface{}{
						{"name": "手動更新", "unit_price": 200, "quantity": 1},
					},
				},
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().Value("title").IsEqual("取消後已編輯")
	})

	t.Run("ai-cancel on a non-pending transaction is rejected", func(t *testing.T) {
		if pendingTxnID == "" {
			t.Fatal("prior sub-test did not produce a transaction id")
		}
		// pendingTxnID was cancelled above so it is now NULL — calling cancel
		// again must 409, matching the service-level RowsAffected==0 guard.
		ae := authExpect(t, token)
		ae.POST("/api/spaces/{id}/transactions/{txn_id}/ai-cancel", personalID, pendingTxnID).
			Expect().
			Status(http.StatusConflict)
	})
}
