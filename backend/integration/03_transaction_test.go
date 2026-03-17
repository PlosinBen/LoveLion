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

	t.Run("personal transactions", func(t *testing.T) {
		ae := authExpect(t, token)

		// 星巴克
		txn := ae.POST("/api/spaces/{id}/transactions", personalID).
			WithJSON(map[string]interface{}{
				"title":          "星巴克",
				"payer":          "Antigravity",
				"date":           now.Add(-2 * time.Hour).Format(time.RFC3339),
				"currency":       "TWD",
				"exchange_rate":  1,
				"category":       "餐飲",
				"payment_method": "信用卡",
				"note":           "跟同事下午茶",
				"items": []map[string]interface{}{
					{"name": "特大杯拿鐵", "unit_price": 155, "quantity": 2},
				},
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		txn.Value("title").IsEqual("星巴克")
		txn.Value("total_amount").IsEqual("310") // 155*2

		// 捷運定期票
		ae.POST("/api/spaces/{id}/transactions", personalID).
			WithJSON(map[string]interface{}{
				"title":          "捷運定期票",
				"payer":          "Antigravity",
				"date":           now.Add(-24 * time.Hour).Format(time.RFC3339),
				"currency":       "TWD",
				"exchange_rate":  1,
				"category":       "交通",
				"payment_method": "現金",
				"items": []map[string]interface{}{
					{"name": "捷運定期票", "unit_price": 1200, "quantity": 1},
				},
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object().Value("total_amount").IsEqual("1200")

		// list
		ae.GET("/api/spaces/{id}/transactions", personalID).
			Expect().
			Status(http.StatusOK).
			JSON().Array().Length().IsEqual(2)
	})

	t.Run("trip transactions", func(t *testing.T) {
		ae := authExpect(t, token)

		// 利木津巴士
		txn := ae.POST("/api/spaces/{id}/transactions", tripID).
			WithJSON(map[string]interface{}{
				"title":          "利木津巴士",
				"payer":          "Antigravity",
				"date":           now.AddDate(0, 1, 0).Format(time.RFC3339),
				"currency":       "JPY",
				"exchange_rate":  0.216,
				"billing_amount": 1944,
				"handling_fee":   29.16,
				"category":       "交通",
				"payment_method": "信用卡",
				"items": []map[string]interface{}{
					{"name": "成人票", "unit_price": 3000, "quantity": 3},
				},
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		txn.Value("total_amount").IsEqual("9000") // 3000*3
		txnID := txn.Value("id").String().Raw()

		// 一蘭拉麵
		ae.POST("/api/spaces/{id}/transactions", tripID).
			WithJSON(map[string]interface{}{
				"title":          "一蘭拉麵",
				"payer":          "Antigravity",
				"date":           now.AddDate(0, 1, 1).Format(time.RFC3339),
				"currency":       "JPY",
				"exchange_rate":  0.216,
				"billing_amount": 1253,
				"handling_fee":   0,
				"category":       "飲食",
				"payment_method": "現金",
				"items": []map[string]interface{}{
					{"name": "天然豚骨拉麵", "unit_price": 980, "quantity": 3},
					{"name": "加麵", "unit_price": 210, "quantity": 2},
					{"name": "生啤酒", "unit_price": 580, "quantity": 3},
					{"name": "半熟鹽味蛋", "unit_price": 140, "quantity": 5},
				},
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object().Value("items").Array().Length().IsEqual(4)

		// get single
		ae.GET("/api/spaces/{id}/transactions/{txn_id}", tripID, txnID).
			Expect().
			Status(http.StatusOK).
			JSON().Object().Value("title").IsEqual("利木津巴士")
	})

	t.Run("update and delete", func(t *testing.T) {
		ae := authExpect(t, token)

		// create temp transaction
		tmp := ae.POST("/api/spaces/{id}/transactions", personalID).
			WithJSON(map[string]interface{}{
				"title": "待刪除",
				"items": []map[string]interface{}{
					{"name": "test", "unit_price": 100, "quantity": 1},
				},
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()
		tmpID := tmp.Value("id").String().Raw()

		// update
		ae.PUT("/api/spaces/{id}/transactions/{txn_id}", personalID, tmpID).
			WithJSON(map[string]interface{}{
				"title": "已更新",
				"items": []map[string]interface{}{
					{"name": "updated", "unit_price": 50, "quantity": 3},
				},
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().Value("total_amount").IsEqual("150")

		// delete
		ae.DELETE("/api/spaces/{id}/transactions/{txn_id}", personalID, tmpID).
			Expect().
			Status(http.StatusOK)

		ae.GET("/api/spaces/{id}/transactions/{txn_id}", personalID, tmpID).
			Expect().
			Status(http.StatusNotFound)
	})
}
