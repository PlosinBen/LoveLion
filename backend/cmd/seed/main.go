package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"lovelion/internal/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var apiBase string

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	apiBase = os.Getenv("API_BASE")
	if apiBase == "" {
		apiBase = "http://localhost:8080/api"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Check if already seeded
	var existing models.User
	if err := db.Where("username = ?", "dev").First(&existing).Error; err == nil {
		fmt.Println("✓ Database already seeded (user 'dev' exists)")
		return
	}

	// === Step 1: Create users via DB (only exception) ===
	users := map[string]*models.User{
		"dev":  {ID: uuid.New(), Username: "dev", DisplayName: "Antigravity"},
		"ming": {ID: uuid.New(), Username: "ming", DisplayName: "小明"},
		"mei":  {ID: uuid.New(), Username: "mei", DisplayName: "小美"},
	}
	for _, u := range users {
		u.SetPassword(u.Username + "123")
		db.Create(u)
	}
	fmt.Println("✓ Created users: dev, ming, mei")

	// === Step 2: Login to get tokens ===
	devToken := login("dev", "dev123")
	mingToken := login("ming", "ming123")
	meiToken := login("mei", "mei123")
	fmt.Println("✓ Obtained auth tokens")

	// === Step 3: Create personal space via API ===
	personalSpace := apiPost(devToken, "/spaces", map[string]any{
		"name":            "日常開銷",
		"type":            "personal",
		"base_currency":   "TWD",
		"currencies":      []string{"TWD", "JPY", "USD"},
		"categories":      []string{"餐飲", "交通", "購物", "娛樂", "生活"},
		"payment_methods": []string{"現金", "信用卡", "Line Pay"},
		"is_pinned":       true,
	})
	personalSpaceID := personalSpace["id"].(string)
	fmt.Printf("✓ Created space: 日常開銷 (%s)\n", personalSpaceID)

	// === Step 4: Create personal transactions ===
	now := time.Now()
	apiPost(devToken, fmt.Sprintf("/spaces/%s/transactions", personalSpaceID), map[string]any{
		"title":          "星巴克",
		"payer":          "Antigravity",
		"date":           now.Add(-2 * time.Hour).Format(time.RFC3339),
		"currency":       "TWD",
		"exchange_rate":  1,
		"category":       "餐飲",
		"payment_method": "信用卡",
		"note":           "跟同事下午茶",
		"items": []map[string]any{
			{"name": "特大杯拿鐵", "unit_price": 155, "quantity": 2},
		},
	})
	fmt.Println("✓ Created transaction: 星巴克")

	apiPost(devToken, fmt.Sprintf("/spaces/%s/transactions", personalSpaceID), map[string]any{
		"title":          "捷運定期票",
		"payer":          "Antigravity",
		"date":           now.Add(-24 * time.Hour).Format(time.RFC3339),
		"currency":       "TWD",
		"exchange_rate":  1,
		"category":       "交通",
		"payment_method": "現金",
		"items": []map[string]any{
			{"name": "捷運定期票", "unit_price": 1200, "quantity": 1},
		},
	})
	fmt.Println("✓ Created transaction: 捷運定期票")

	// === Step 5: Create trip space via API ===
	tripSpace := apiPost(devToken, "/spaces", map[string]any{
		"name":            "2024 東京春櫻季",
		"description":     "5 天 4 夜 東京賞櫻團",
		"type":            "trip",
		"base_currency":   "TWD",
		"currencies":      []string{"TWD", "JPY"},
		"categories":      []string{"住宿", "交通", "飲食", "購物", "娛樂"},
		"payment_methods": []string{"現金", "信用卡"},
		"start_date":      now.AddDate(0, 1, 0).Format(time.RFC3339),
		"end_date":        now.AddDate(0, 1, 5).Format(time.RFC3339),
		"is_pinned":       true,
	})
	tripSpaceID := tripSpace["id"].(string)
	fmt.Printf("✓ Created space: 2024 東京春櫻季 (%s)\n", tripSpaceID)

	// === Step 6: Invite members via API ===
	invite := apiPost(devToken, fmt.Sprintf("/spaces/%s/invites", tripSpaceID), map[string]any{
		"is_one_time": false,
		"max_uses":    10,
	})
	inviteToken := invite["token"].(string)

	apiPost(mingToken, fmt.Sprintf("/invites/%s/join", inviteToken), nil)
	apiPost(meiToken, fmt.Sprintf("/invites/%s/join", inviteToken), nil)
	fmt.Println("✓ Added members: ming, mei")

	// Set member aliases
	apiPatch(devToken, fmt.Sprintf("/spaces/%s/members/%s", tripSpaceID, users["ming"].ID), map[string]any{"alias": "小明"})
	apiPatch(devToken, fmt.Sprintf("/spaces/%s/members/%s", tripSpaceID, users["mei"].ID), map[string]any{"alias": "小美"})
	fmt.Println("✓ Set member aliases")

	// === Step 7: Create comparison stores ===
	store1 := apiPost(devToken, fmt.Sprintf("/spaces/%s/stores", tripSpaceID), map[string]any{
		"name":           "唐吉軻德 澀谷店",
		"location":       "澀谷",
		"google_map_url": "https://maps.app.goo.gl/ShibuyaDonki",
	})
	store1ID := store1["id"].(string)

	store2 := apiPost(devToken, fmt.Sprintf("/spaces/%s/stores", tripSpaceID), map[string]any{
		"name":     "Bic Camera 新宿",
		"location": "新宿",
	})
	store2ID := store2["id"].(string)
	fmt.Println("✓ Created comparison stores")

	// === Step 8: Create comparison products ===
	apiPost(devToken, fmt.Sprintf("/spaces/%s/stores/%s/products", tripSpaceID, store1ID), map[string]any{
		"name": "一蘭拉麵泡麵", "price": 1850, "currency": "JPY",
	})
	apiPost(devToken, fmt.Sprintf("/spaces/%s/stores/%s/products", tripSpaceID, store2ID), map[string]any{
		"name": "一蘭拉麵泡麵", "price": 1980, "currency": "JPY",
	})
	apiPost(devToken, fmt.Sprintf("/spaces/%s/stores/%s/products", tripSpaceID, store1ID), map[string]any{
		"name": "Dyson 吹風機", "price": 45000, "currency": "JPY",
	})
	fmt.Println("✓ Created comparison products")

	// === Step 9: Create trip transactions ===
	apiPost(devToken, fmt.Sprintf("/spaces/%s/transactions", tripSpaceID), map[string]any{
		"title":          "利木津巴士",
		"payer":          "Antigravity",
		"date":           now.AddDate(0, 1, 0).Format(time.RFC3339),
		"currency":       "JPY",
		"exchange_rate":  0.216,
		"billing_amount": 1944,
		"handling_fee":   29.16,
		"category":       "交通",
		"payment_method": "信用卡",
		"items": []map[string]any{
			{"name": "成人票", "unit_price": 3000, "quantity": 3},
		},
	})
	fmt.Println("✓ Created transaction: 利木津巴士")

	apiPost(devToken, fmt.Sprintf("/spaces/%s/transactions", tripSpaceID), map[string]any{
		"title":          "一蘭拉麵",
		"payer":          "Antigravity",
		"date":           now.AddDate(0, 1, 1).Format(time.RFC3339),
		"currency":       "JPY",
		"exchange_rate":  0.216,
		"billing_amount": 1253,
		"handling_fee":   0,
		"category":       "飲食",
		"payment_method": "現金",
		"items": []map[string]any{
			{"name": "天然豚骨拉麵", "unit_price": 980, "quantity": 3},
			{"name": "加麵", "unit_price": 210, "quantity": 2},
			{"name": "生啤酒", "unit_price": 580, "quantity": 3},
			{"name": "半熟鹽味蛋", "unit_price": 140, "quantity": 5},
		},
	})
	fmt.Println("✓ Created transaction: 一蘭拉麵")

	fmt.Println("\n🎉 Seed completed successfully!")
	fmt.Println("   Login User: dev")
	fmt.Println("   Password:   dev123")
}

// --- API helpers ---

func login(username, password string) string {
	resp := apiPost("", "/users/login", map[string]any{
		"username": username,
		"password": password,
	})
	token, ok := resp["token"].(string)
	if !ok || token == "" {
		log.Fatalf("Failed to login as %s", username)
	}
	return token
}

func apiPost(token, path string, body map[string]any) map[string]any {
	return apiRequest("POST", token, path, body)
}

func apiPatch(token, path string, body map[string]any) map[string]any {
	return apiRequest("PATCH", token, path, body)
}

func apiRequest(method, token, path string, body map[string]any) map[string]any {
	var reqBody io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, apiBase+path, reqBody)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("API request failed: %s %s: %v", method, path, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		log.Fatalf("API error: %s %s → %d: %s", method, path, resp.StatusCode, string(respBody))
	}

	var result map[string]any
	json.Unmarshal(respBody, &result)
	return result
}
