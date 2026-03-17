//go:build integration

package integration

import (
	"net/http"
	"testing"
)

func TestComparison(t *testing.T) {
	token := loginUser(t, "dev", "dev123")

	// find 東京春櫻季
	tripID := findSpaceByName(t, token, "2024 東京春櫻季")

	t.Run("stores and products", func(t *testing.T) {
		ae := authExpect(t, token)

		// 唐吉軻德 澀谷店
		store1 := ae.POST("/api/spaces/{id}/stores", tripID).
			WithJSON(map[string]interface{}{
				"name":           "唐吉軻德 澀谷店",
				"location":       "澀谷",
				"google_map_url": "https://maps.app.goo.gl/ShibuyaDonki",
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()
		store1.Value("name").IsEqual("唐吉軻德 澀谷店")
		store1ID := store1.Value("id").String().Raw()

		// Bic Camera 新宿
		store2 := ae.POST("/api/spaces/{id}/stores", tripID).
			WithJSON(map[string]interface{}{
				"name":     "Bic Camera 新宿",
				"location": "新宿",
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()
		store2ID := store2.Value("id").String().Raw()

		// list stores
		ae.GET("/api/spaces/{id}/stores", tripID).
			Expect().
			Status(http.StatusOK).
			JSON().Array().Length().IsEqual(2)

		// products — 一蘭拉麵泡麵 at both stores
		ae.POST("/api/spaces/{id}/stores/{store_id}/products", tripID, store1ID).
			WithJSON(map[string]interface{}{
				"name": "一蘭拉麵泡麵", "price": 1850, "currency": "JPY",
			}).
			Expect().
			Status(http.StatusCreated)

		ae.POST("/api/spaces/{id}/stores/{store_id}/products", tripID, store2ID).
			WithJSON(map[string]interface{}{
				"name": "一蘭拉麵泡麵", "price": 1980, "currency": "JPY",
			}).
			Expect().
			Status(http.StatusCreated)

		// Dyson 吹風機
		product := ae.POST("/api/spaces/{id}/stores/{store_id}/products", tripID, store1ID).
			WithJSON(map[string]interface{}{
				"name": "Dyson 吹風機", "price": 45000, "currency": "JPY",
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()
		productID := product.Value("id").String().Raw()

		// list all products across stores
		ae.GET("/api/spaces/{id}/products", tripID).
			Expect().
			Status(http.StatusOK).
			JSON().Array().Length().IsEqual(3)

		// update product
		ae.PUT("/api/spaces/{id}/stores/{store_id}/products/{product_id}", tripID, store1ID, productID).
			WithJSON(map[string]interface{}{"price": 43000}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().Value("price").IsEqual("43000")

		// update store
		ae.PUT("/api/spaces/{id}/stores/{store_id}", tripID, store1ID).
			WithJSON(map[string]interface{}{"name": "唐吉軻德 澀谷本店"}).
			Expect().
			Status(http.StatusOK).
			JSON().Object().Value("name").IsEqual("唐吉軻德 澀谷本店")
	})
}
