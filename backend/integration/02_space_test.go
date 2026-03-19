//go:build integration

package integration

import (
	"net/http"
	"testing"
	"time"
)

func TestSpaces(t *testing.T) {
	token := loginUser(t, "dev", "dev123")

	t.Run("create spaces", func(t *testing.T) {
		ae := authExpect(t, token)

		// personal space — 日常開銷
		space := ae.POST("/api/spaces").
			WithJSON(map[string]interface{}{
				"name":            "日常開銷",
				"type":            "personal",
				"base_currency":   "TWD",
				"currencies":      []string{"TWD", "JPY", "USD"},
				"categories":      []string{"餐飲", "交通", "購物", "娛樂", "生活"},
				"payment_methods": []string{"現金", "信用卡", "Line Pay"},
				"is_pinned":       true,
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		space.Value("name").IsEqual("日常開銷")
		space.Value("type").IsEqual("personal")
		space.Value("base_currency").IsEqual("TWD")
		space.Value("is_pinned").IsEqual(true)

		// trip space — 東京春櫻季
		now := time.Now()
		trip := ae.POST("/api/spaces").
			WithJSON(map[string]interface{}{
				"name":            "2024 東京春櫻季",
				"description":     "5 天 4 夜 東京賞櫻團",
				"type":            "trip",
				"base_currency":   "TWD",
				"currencies":      []string{"TWD", "JPY"},
				"categories":      []string{"住宿", "交通", "飲食", "購物", "娛樂"},
				"payment_methods": []string{"現金", "信用卡"},
				"split_members":   []string{"我", "老婆", "老媽"},
				"start_date":      now.AddDate(0, 1, 0).Format(time.RFC3339),
				"end_date":        now.AddDate(0, 1, 5).Format(time.RFC3339),
				"is_pinned":       true,
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		trip.Value("name").IsEqual("2024 東京春櫻季")
		trip.Value("type").IsEqual("trip")
		trip.Value("start_date").NotNull()
		trip.Value("end_date").NotNull()
	})

	t.Run("list and filter", func(t *testing.T) {
		ae := authExpect(t, token)

		// list all (default + 日常開銷 + 東京春櫻季)
		ae.GET("/api/spaces").
			Expect().
			Status(http.StatusOK).
			JSON().Array().Length().Ge(3)

		// filter by trip
		trips := ae.GET("/api/spaces").WithQuery("type", "trip").
			Expect().
			Status(http.StatusOK).
			JSON().Array()
		trips.Length().IsEqual(1)
		trips.Value(0).Object().Value("type").IsEqual("trip")
	})

	t.Run("update and delete", func(t *testing.T) {
		ae := authExpect(t, token)

		// create a temp space to test update/delete
		tmp := ae.POST("/api/spaces").
			WithJSON(map[string]interface{}{"name": "Temp Space", "type": "personal"}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()
		tmpID := tmp.Value("id").String().Raw()

		ae.PUT("/api/spaces/{id}", tmpID).
			WithJSON(map[string]interface{}{"name": "Updated Temp"}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().Value("name").IsEqual("Updated Temp")

		ae.DELETE("/api/spaces/{id}", tmpID).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("access control", func(t *testing.T) {
		ae := authExpect(t, token)

		spaces := ae.GET("/api/spaces").
			Expect().
			Status(http.StatusOK).
			JSON().Array()
		spaceID := spaces.Value(0).Object().Value("id").String().Raw()

		// ming can't access dev's personal space
		mingToken := loginUser(t, "ming", "ming123")
		authExpect(t, mingToken).GET("/api/spaces/{id}", spaceID).
			Expect().
			Status(http.StatusForbidden)

		// unauthenticated
		newExpect(t).GET("/api/spaces").
			Expect().
			Status(http.StatusUnauthorized)
	})
}
