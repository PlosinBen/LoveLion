//go:build integration

package integration

import (
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	e := newExpect(t)

	t.Run("register users", func(t *testing.T) {
		// dev — main test account
		registerUser(t, "dev", "dev123", "Antigravity")

		// ming, mei — trip members
		registerUser(t, "ming", "ming123", "小明")
		registerUser(t, "mei", "mei123", "小美")

		// duplicate
		e.POST("/api/users/register").
			WithJSON(map[string]string{
				"username":     "dev",
				"password":     "dev123",
				"display_name": "Dup",
			}).
			Expect().
			Status(http.StatusConflict)

		// wrong password
		e.POST("/api/users/login").
			WithJSON(map[string]string{
				"username": "dev",
				"password": "wrong",
			}).
			Expect().
			Status(http.StatusUnauthorized)
	})

	t.Run("me and update", func(t *testing.T) {
		token := loginUser(t, "dev", "dev123")
		ae := authExpect(t, token)

		user := ae.GET("/api/users/me").
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		user.Value("username").IsEqual("dev")
		user.NotContainsKey("password_hash")

		// unauthorized
		e.GET("/api/users/me").Expect().Status(http.StatusUnauthorized)

		// change password and verify
		ae.PUT("/api/users/me").
			WithJSON(map[string]string{
				"current_password": "dev123",
				"new_password":     "newpass123",
			}).
			Expect().
			Status(http.StatusOK)

		loginUser(t, "dev", "newpass123")

		// revert password
		ae2 := authExpect(t, loginUser(t, "dev", "newpass123"))
		ae2.PUT("/api/users/me").
			WithJSON(map[string]string{
				"current_password": "newpass123",
				"new_password":     "dev123",
			}).
			Expect().
			Status(http.StatusOK)
	})
}
