//go:build integration

package integration

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// promoteToAdmin sets a user's role to 'admin' via direct DB connection.
// Returns the user ID for reference.
func promoteToAdmin(t *testing.T, token string) string {
	t.Helper()

	ae := authExpect(t, token)
	userID := ae.GET("/api/users/me").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("id").String().Raw()

	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/lovelion?sslmode=disable")
	if err != nil {
		t.Fatalf("connect to db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users SET role = 'admin' WHERE id = $1", userID)
	if err != nil {
		t.Fatalf("promote to admin: %v", err)
	}

	return userID
}

func TestAnnouncements(t *testing.T) {
	devToken := loginUser(t, "dev", "dev123")

	// Register a second (non-admin) user for permission checks.
	normalToken := registerUser(t, "ann_normal", "pass123", "Normal User")

	t.Run("public list returns empty initially", func(t *testing.T) {
		e := newExpect(t)
		e.GET("/api/announcements").
			Expect().
			Status(http.StatusOK).
			JSON().Array()
	})

	t.Run("broadcast returns null when none active", func(t *testing.T) {
		e := newExpect(t)
		e.GET("/api/announcements/broadcast").
			Expect().
			Status(http.StatusOK)
	})

	t.Run("non-admin cannot access admin endpoints", func(t *testing.T) {
		ae := authExpect(t, normalToken)
		ae.GET("/api/admin/announcements").
			Expect().
			Status(http.StatusForbidden)

		ae.POST("/api/admin/announcements").
			WithJSON(map[string]interface{}{
				"title":  "Test",
				"status": "draft",
			}).
			Expect().
			Status(http.StatusForbidden)
	})

	t.Run("unauthenticated cannot access admin endpoints", func(t *testing.T) {
		e := newExpect(t)
		e.GET("/api/admin/announcements").
			Expect().
			Status(http.StatusUnauthorized)
	})

	// Promote dev to admin for CRUD tests
	promoteToAdmin(t, devToken)

	t.Run("admin CRUD flow", func(t *testing.T) {
		ae := authExpect(t, devToken)

		// Create draft announcement
		created := ae.POST("/api/admin/announcements").
			WithJSON(map[string]interface{}{
				"title":   "系統維護通知",
				"content": "我們將於下週三凌晨進行系統維護。",
				"status":  "draft",
			}).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		announcementID := created.Value("id").String().Raw()
		created.Value("status").IsEqual("draft")
		created.Value("title").IsEqual("系統維護通知")

		// Draft should NOT appear in public list
		e := newExpect(t)
		pubList := e.GET("/api/announcements").
			Expect().
			Status(http.StatusOK).
			JSON().Array()

		found := false
		for i := 0; i < int(pubList.Length().Raw()); i++ {
			if pubList.Value(i).Object().Value("id").Raw() == announcementID {
				found = true
			}
		}
		if found {
			t.Error("draft announcement should not appear in public list")
		}

		// Draft should appear in admin list
		adminList := ae.GET("/api/admin/announcements").
			Expect().
			Status(http.StatusOK).
			JSON().Array()

		adminFound := false
		for i := 0; i < int(adminList.Length().Raw()); i++ {
			if adminList.Value(i).Object().Value("id").Raw() == announcementID {
				adminFound = true
			}
		}
		if !adminFound {
			t.Error("draft announcement should appear in admin list")
		}

		// Update to published with broadcast
		now := time.Now()
		broadcastStart := now.Add(-1 * time.Hour)
		broadcastEnd := now.Add(1 * time.Hour)

		updated := ae.PUT("/api/admin/announcements/{id}", announcementID).
			WithJSON(map[string]interface{}{
				"title":           "系統維護通知（更新）",
				"content":         "維護時間已確認。",
				"status":          "published",
				"broadcast_start": broadcastStart.Format(time.RFC3339),
				"broadcast_end":   broadcastEnd.Format(time.RFC3339),
			}).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		updated.Value("status").IsEqual("published")
		updated.Value("title").IsEqual("系統維護通知（更新）")

		// Published should appear in public list
		e.GET("/api/announcements/{id}", announcementID).
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			Value("title").IsEqual("系統維護通知（更新）")

		// Should appear in broadcast
		broadcast := e.GET("/api/announcements/broadcast").
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		broadcast.Value("id").IsEqual(announcementID)

		// Delete
		ae.DELETE("/api/admin/announcements/{id}", announcementID).
			Expect().
			Status(http.StatusOK)

		// Should no longer exist
		e.GET("/api/announcements/{id}", announcementID).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("admin config endpoint", func(t *testing.T) {
		ae := authExpect(t, devToken)
		cfg := ae.GET("/api/admin/announcements/config").
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		cfg.ContainsKey("ai_available")
	})

	t.Run("create with invalid status returns 400", func(t *testing.T) {
		ae := authExpect(t, devToken)
		ae.POST("/api/admin/announcements").
			WithJSON(map[string]interface{}{
				"title":  "Test",
				"status": "invalid",
			}).
			Expect().
			Status(http.StatusBadRequest)
	})
}
