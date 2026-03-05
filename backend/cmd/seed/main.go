package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"lovelion/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
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

	// 1. Create test user
	user := &models.User{
		ID:          uuid.New(),
		Username:    "dev",
		DisplayName: "Developer",
	}
	user.SetPassword("dev123")
	db.Create(user)
	fmt.Println("✓ Created user: dev / dev123")

	// 2. Create personal space (ledger)
	personalSpace := &models.Ledger{
		ID:             uuid.New(),
		UserID:         user.ID,
		Name:           "我的帳本",
		Type:           "personal",
		BaseCurrency:   "TWD",
		Currencies:     datatypes.JSON([]byte(`["TWD", "JPY", "USD"]`)),
		Categories:     datatypes.JSON([]byte(`["食物", "交通", "購物", "娛樂"]`)),
		PaymentMethods: datatypes.JSON([]byte(`["現金", "信用卡", "Line Pay"]`)),
		MemberNames:    datatypes.JSON([]byte(`[]`)),
	}
	db.Create(personalSpace)
	fmt.Println("✓ Created space: 我的帳本 (personal)")

	// 3. Create sample transactions for personal space
	now := time.Now()
	transactions := []models.Transaction{
		{
			ID:            "txn01",
			LedgerID:      personalSpace.ID,
			Title:         "午餐",
			Payer:         "Me",
			Date:          now.AddDate(0, 0, -2),
			Currency:      "TWD",
			ExchangeRate:  decimal.NewFromInt(1),
			TotalAmount:   decimal.NewFromInt(150),
			Category:      "食物",
			PaymentMethod: "現金",
			Note:          "午餐",
		},
		{
			ID:            "txn02",
			LedgerID:      personalSpace.ID,
			Title:         "捷運儲值",
			Payer:         "Me",
			Date:          now.AddDate(0, 0, -1),
			Currency:      "TWD",
			ExchangeRate:  decimal.NewFromInt(1),
			TotalAmount:   decimal.NewFromInt(35),
			Category:      "交通",
			PaymentMethod: "悠遊卡",
			Note:          "捷運",
		},
	}
	for _, txn := range transactions {
		db.Create(&txn)
	}
	fmt.Println("✓ Created 2 sample transactions for personal space")

	// 4. Create a trip space (unified Ledger + Trip)
	tripSpace := &models.Ledger{
		ID:             uuid.New(),
		UserID:         user.ID,
		Name:           "2024 日本旅行",
		Description:    "東京 + 大阪 5天4夜",
		Type:           "trip",
		BaseCurrency:   "TWD",
		Currencies:     datatypes.JSON([]byte(`["TWD", "JPY"]`)),
		Categories:     datatypes.JSON([]byte(`["住宿", "交通", "飲食", "購物", "門票"]`)),
		PaymentMethods: datatypes.JSON([]byte(`["現金", "信用卡"]`)),
		MemberNames:    datatypes.JSON([]byte(`["我", "小明", "小美"]`)),
		StartDate:      timePtr(now.AddDate(0, 1, 0)),
		EndDate:        timePtr(now.AddDate(0, 1, 5)),
		IsPinned:       true,
	}
	db.Create(tripSpace)
	fmt.Println("✓ Created space: 2024 日本旅行 (trip)")

	// 5. Add ledger members for trip space
	members := []models.LedgerMember{
		{ID: uuid.New(), LedgerID: tripSpace.ID, UserID: user.ID, Role: "owner"},
		{ID: uuid.New(), LedgerID: tripSpace.ID, Role: "member", Alias: "小明"},
		{ID: uuid.New(), LedgerID: tripSpace.ID, Role: "member", Alias: "小美"},
	}
	// Note: Currently LedgerMember requires UserID, but for seeding unnamed members we might need to adjust or create dummy users.
	// For now, let's keep it simple and only link the owner.
	db.Create(&members[0])
	
	// For Ming and Mei, we'll just create dummy records or rely on Name field in splits for now.
	// Actually, LedgerMember table has UserID NOT NULL in some migrations. Let's check.
	// 000012_create_ledger_sharing.up.sql: user_id UUID NOT NULL
	
	fmt.Println("✓ Added owner to trip space")

	// 6. Create comparison stores for trip space
	store1 := &models.ComparisonStore{
		ID:       "donki1",
		LedgerID: tripSpace.ID,
		Name:     "唐吉軻德 新宿店",
		Location: "新宿",
	}
	store2 := &models.ComparisonStore{
		ID:       "bic01",
		LedgerID: tripSpace.ID,
		Name:     "Bic Camera 有楽町",
		Location: "有楽町",
	}
	db.Create(store1)
	db.Create(store2)
	fmt.Println("✓ Created 2 comparison stores for trip space")

	// 7. Create comparison products
	products := []models.ComparisonProduct{
		{ID: uuid.New(), StoreID: store1.ID, Name: "藥妝", Price: decimal.NewFromInt(2980), Currency: "JPY"},
		{ID: uuid.New(), StoreID: store2.ID, Name: "藥妝", Price: decimal.NewFromInt(3200), Currency: "JPY"},
	}
	for _, p := range products {
		db.Create(&p)
	}
	fmt.Println("✓ Created 2 comparison products")

	// 8. Create Trip Transactions with Splits
	txn1 := models.Transaction{
		ID:            "trip_txn01",
		LedgerID:      tripSpace.ID,
		Title:         "第一天晚餐",
		Payer:         "我",
		Date:          now.AddDate(0, 0, 1),
		Currency:      "JPY",
		TotalAmount:   decimal.NewFromInt(3000),
		Category:      "飲食",
		PaymentMethod: "信用卡",
	}
	db.Create(&txn1)

	// Using Name based splits for simplicity in seed
	splits1 := []models.TransactionSplit{
		{TransactionID: txn1.ID, Name: "我", Amount: decimal.NewFromInt(3000), IsPayer: true, MemberID: &members[0].ID},
		{TransactionID: txn1.ID, Name: "我", Amount: decimal.NewFromInt(1000), IsPayer: false, MemberID: &members[0].ID},
		{TransactionID: txn1.ID, Name: "小明", Amount: decimal.NewFromInt(1000), IsPayer: false},
		{TransactionID: txn1.ID, Name: "小美", Amount: decimal.NewFromInt(1000), IsPayer: false},
	}
	for _, s := range splits1 {
		db.Create(&s)
	}

	fmt.Println("✓ Created sample trip transaction with splits")

	fmt.Println("\n🎉 Seed completed!")
	fmt.Println("   Login: dev / dev123")
}

func timePtr(t time.Time) *time.Time {
	return &t
}