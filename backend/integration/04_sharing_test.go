//go:build integration

package integration

import (
	"net/http"
	"testing"
)

func TestSharing(t *testing.T) {
	devToken := loginUser(t, "dev", "dev123")
	mingToken := loginUser(t, "ming", "ming123")
	meiToken := loginUser(t, "mei", "mei123")

	// find 東京春櫻季
	tripID := findSpaceByName(t, devToken, "2024 東京春櫻季")

	t.Run("invite members", func(t *testing.T) {
		ae := authExpect(t, devToken)
		mingE := authExpect(t, mingToken)
		meiE := authExpect(t, meiToken)

		// create invite
		invite := ae.POST("/api/spaces/{id}/invites", tripID).
			WithJSON(map[string]interface{}{"max_uses": 10}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		inviteToken := invite.Value("token").String().Raw()

		// public info
		newExpect(t).GET("/api/invites/{token}", inviteToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object().Value("space_name").IsEqual("2024 東京春櫻季")

		// ming joins
		mingE.POST("/api/invites/{token}/join", inviteToken).
			Expect().
			Status(http.StatusOK)

		// mei joins
		meiE.POST("/api/invites/{token}/join", inviteToken).
			Expect().
			Status(http.StatusOK)

		// list members = 3
		ae.GET("/api/spaces/{id}/members", tripID).
			Expect().
			Status(http.StatusOK).
			JSON().Array().Length().IsEqual(3)

		// set aliases
		mingID := authExpect(t, mingToken).GET("/api/users/me").Expect().
			Status(http.StatusOK).JSON().Object().Value("id").String().Raw()
		meiID := authExpect(t, meiToken).GET("/api/users/me").Expect().
			Status(http.StatusOK).JSON().Object().Value("id").String().Raw()

		ae.PATCH("/api/spaces/{id}/members/{user_id}", tripID, mingID).
			WithJSON(map[string]interface{}{"alias": "小明"}).
			Expect().
			Status(http.StatusOK)

		ae.PATCH("/api/spaces/{id}/members/{user_id}", tripID, meiID).
			WithJSON(map[string]interface{}{"alias": "小美"}).
			Expect().
			Status(http.StatusOK)

		// ming can access trip space
		mingE.GET("/api/spaces/{id}", tripID).
			Expect().
			Status(http.StatusOK).
			JSON().Object().Value("name").IsEqual("2024 東京春櫻季")

		// non-owner can't remove owner
		devID := ae.GET("/api/users/me").Expect().
			Status(http.StatusOK).JSON().Object().Value("id").String().Raw()

		mingE.DELETE("/api/spaces/{id}/members/{user_id}", tripID, devID).
			Expect().
			Status(http.StatusForbidden)
	})

	t.Run("remove and leave", func(t *testing.T) {
		ae := authExpect(t, devToken)
		mingE := authExpect(t, mingToken)
		meiE := authExpect(t, meiToken)

		// create a separate group space for removal tests
		space := ae.POST("/api/spaces").
			WithJSON(map[string]interface{}{"name": "移除測試", "type": "group"}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()
		sid := space.Value("id").String().Raw()

		// invite ming and mei
		inv := ae.POST("/api/spaces/{id}/invites", sid).
			WithJSON(map[string]interface{}{"max_uses": 10}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()
		tok := inv.Value("token").String().Raw()

		mingE.POST("/api/invites/{token}/join", tok).Expect().Status(http.StatusOK)
		meiE.POST("/api/invites/{token}/join", tok).Expect().Status(http.StatusOK)

		// owner removes mei
		meiID := authExpect(t, meiToken).GET("/api/users/me").Expect().
			Status(http.StatusOK).JSON().Object().Value("id").String().Raw()

		ae.DELETE("/api/spaces/{id}/members/{user_id}", sid, meiID).
			Expect().
			Status(http.StatusOK)

		// mei can no longer access
		meiE.GET("/api/spaces/{id}", sid).
			Expect().
			Status(http.StatusForbidden)

		// ming leaves voluntarily
		mingE.POST("/api/spaces/{id}/leave", sid).
			Expect().
			Status(http.StatusOK)

		// ming can no longer access
		mingE.GET("/api/spaces/{id}", sid).
			Expect().
			Status(http.StatusForbidden)

		// members = 1 (only owner left)
		ae.GET("/api/spaces/{id}/members", sid).
			Expect().
			Status(http.StatusOK).
			JSON().Array().Length().IsEqual(1)

		// cleanup
		ae.DELETE("/api/spaces/{id}", sid).
			Expect().
			Status(http.StatusOK)
	})
}
