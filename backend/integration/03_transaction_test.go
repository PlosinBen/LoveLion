//go:build integration

package integration

import (
	"net/http"
	"testing"
	"time"
)

func TestTransactions(t *testing.T) {
	token := loginUser(t, "dev", "dev123")

	// find 日常開銷
	personalID := findSpaceByName(t, token, "日常開銷")
	// find 東京春櫻季
	tripID := findSpaceByName(t, token, "2024 東京春櫻季")

	now := time.Now()

	t.Run("personal expenses", func(t *testing.T) {
		ae := authExpect(t, token)

		// 星巴克
		txn := ae.POST("/api/spaces/{id}/expenses", personalID).
			WithJSON(map[string]interface{}{
				"title":    "星巴克",
				"date":     now.Add(-2 * time.Hour).Format(time.RFC3339),
				"currency": "TWD",
				"note":     "跟同事下午茶",
				"expense": map[string]interface{}{
					"category":       "餐飲",
					"exchange_rate":  1,
					"payment_method": "信用卡",
					"items": []map[string]interface{}{
						{"name": "特大杯拿鐵", "unit_price": 155, "quantity": 2},
					},
				},
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		txn.Value("title").IsEqual("星巴克")
		txn.Value("total_amount").IsEqual("310") // 155*2
		txn.Value("type").IsEqual("expense")
		txn.Value("expense").Object().Value("category").IsEqual("餐飲")

		// 捷運定期票
		ae.POST("/api/spaces/{id}/expenses", personalID).
			WithJSON(map[string]interface{}{
				"title":    "捷運定期票",
				"date":     now.Add(-24 * time.Hour).Format(time.RFC3339),
				"currency": "TWD",
				"expense": map[string]interface{}{
					"category":       "交通",
					"exchange_rate":  1,
					"payment_method": "現金",
					"items": []map[string]interface{}{
						{"name": "捷運定期票", "unit_price": 1200, "quantity": 1},
					},
				},
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object().Value("total_amount").IsEqual("1200")

		// list — we just created 2
		list := ae.GET("/api/spaces/{id}/transactions", personalID).
			Expect().
			Status(http.StatusOK).
			JSON().Array()
		list.Length().IsEqual(2)
	})

	t.Run("trip expenses with debts", func(t *testing.T) {
		ae := authExpect(t, token)

		// 利木津巴士（三人均分，含 debts）
		txn := ae.POST("/api/spaces/{id}/expenses", tripID).
			WithJSON(map[string]interface{}{
				"title":    "利木津巴士",
				"date":     now.AddDate(0, 1, 0).Format(time.RFC3339),
				"currency": "JPY",
				"expense": map[string]interface{}{
					"category":       "交通",
					"exchange_rate":  0.216,
					"billing_amount": 1944,
					"handling_fee":   29.16,
					"payment_method": "信用卡",
					"items": []map[string]interface{}{
						{"name": "成人票", "unit_price": 3000, "quantity": 3},
					},
				},
				"debts": []map[string]interface{}{
					{"payer_name": "小明", "payee_name": "Antigravity", "amount": 3000},
					{"payer_name": "小美", "payee_name": "Antigravity", "amount": 3000},
				},
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		txn.Value("total_amount").IsEqual("9000") // 3000*3
		txn.Value("debts").Array().Length().IsEqual(2)
		txnID := txn.Value("id").String().Raw()

		// 一蘭拉麵（含 is_spot_paid debt）
		ramen := ae.POST("/api/spaces/{id}/expenses", tripID).
			WithJSON(map[string]interface{}{
				"title":    "一蘭拉麵",
				"date":     now.AddDate(0, 1, 1).Format(time.RFC3339),
				"currency": "JPY",
				"expense": map[string]interface{}{
					"category":       "飲食",
					"exchange_rate":  0.216,
					"billing_amount": 1253,
					"handling_fee":   0,
					"payment_method": "現金",
					"items": []map[string]interface{}{
						{"name": "天然豚骨拉麵", "unit_price": 980, "quantity": 3},
						{"name": "加麵", "unit_price": 210, "quantity": 2},
						{"name": "生啤酒", "unit_price": 580, "quantity": 3},
						{"name": "半熟鹽味蛋", "unit_price": 140, "quantity": 5},
					},
				},
				"debts": []map[string]interface{}{
					{"payer_name": "小明", "payee_name": "Antigravity", "amount": 1580, "is_spot_paid": true},
					{"payer_name": "小美", "payee_name": "Antigravity", "amount": 1580},
				},
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		ramen.Value("expense").Object().Value("items").Array().Length().IsEqual(4)
		ramen.Value("debts").Array().Length().IsEqual(2)

		// get single
		detail := ae.GET("/api/spaces/{id}/transactions/{txn_id}", tripID, txnID).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		detail.Value("title").IsEqual("利木津巴士")
		detail.Value("expense").Object().Value("category").IsEqual("交通")
	})

	t.Run("payment and mixed list", func(t *testing.T) {
		ae := authExpect(t, token)

		// Create a payment (小明付款)
		payment := ae.POST("/api/spaces/{id}/payments", tripID).
			WithJSON(map[string]interface{}{
				"title":        "小明付款給 Antigravity",
				"date":         now.AddDate(0, 1, 2).Format(time.RFC3339),
				"total_amount": 500,
				"payer_name":   "小明",
				"payee_name":   "Antigravity",
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		payment.Value("type").IsEqual("payment")
		payment.Value("total_amount").IsEqual("500")
		payment.Value("debts").Array().Length().IsEqual(1)
		// Verify debt details
		debt := payment.Value("debts").Array().Value(0).Object()
		debt.Value("payer_name").IsEqual("小明")
		debt.Value("payee_name").IsEqual("Antigravity")
		debt.Value("settled_amount").IsEqual("500")

		// Same payer/payee should fail
		ae.POST("/api/spaces/{id}/payments", tripID).
			WithJSON(map[string]interface{}{
				"title":        "invalid",
				"total_amount": 100,
				"payer_name":   "小明",
				"payee_name":   "小明",
			}).
			Expect().
			Status(http.StatusBadRequest)

		// List trip transactions — should contain both expenses and payment
		// (2 expenses from "trip expenses with debts" + 1 payment = 3)
		tripList := ae.GET("/api/spaces/{id}/transactions", tripID).
			Expect().
			Status(http.StatusOK).
			JSON().Array()
		tripList.Length().IsEqual(3)

		// Verify both types exist in the list
		hasExpense := false
		hasPayment := false
		for i := 0; i < 3; i++ {
			txnType := tripList.Value(i).Object().Value("type").String().Raw()
			if txnType == "expense" {
				hasExpense = true
			}
			if txnType == "payment" {
				hasPayment = true
			}
		}
		if !hasExpense {
			t.Error("Expected at least one expense in trip transactions list")
		}
		if !hasPayment {
			t.Error("Expected at least one payment in trip transactions list")
		}
	})

	t.Run("update and delete expense", func(t *testing.T) {
		ae := authExpect(t, token)

		// create temp expense
		tmp := ae.POST("/api/spaces/{id}/expenses", personalID).
			WithJSON(map[string]interface{}{
				"title":    "待刪除",
				"currency": "TWD",
				"expense": map[string]interface{}{
					"category":      "其他",
					"exchange_rate": 1,
					"items": []map[string]interface{}{
						{"name": "test", "unit_price": 100, "quantity": 1},
					},
				},
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()
		tmpID := tmp.Value("id").String().Raw()

		// update expense
		ae.PUT("/api/spaces/{id}/expenses/{txn_id}", personalID, tmpID).
			WithJSON(map[string]interface{}{
				"title":    "已更新",
				"currency": "TWD",
				"expense": map[string]interface{}{
					"category":      "其他",
					"exchange_rate": 1,
					"items": []map[string]interface{}{
						{"name": "updated", "unit_price": 50, "quantity": 3},
					},
				},
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().Value("total_amount").IsEqual("150")

		// delete via /transactions
		ae.DELETE("/api/spaces/{id}/transactions/{txn_id}", personalID, tmpID).
			Expect().
			Status(http.StatusOK)

		ae.GET("/api/spaces/{id}/transactions/{txn_id}", personalID, tmpID).
			Expect().
			Status(http.StatusNotFound)
	})
}
