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

	// 2. Create personal ledger
	ledger := &models.Ledger{
		ID:             uuid.New(),
		UserID:         user.ID,
		Name:           "我的帳本",
		Type:           "personal",
		Currencies:     datatypes.JSON([]byte(`["TWD", "JPY", "USD"]`)),
		Categories:     datatypes.JSON([]byte(`["食物", "交通", "購物", "娛樂"]`)),
		PaymentMethods: datatypes.JSON([]byte(`["現金", "信用卡", "Line Pay"]`)),
		Members:        datatypes.JSON([]byte(`[]`)),
	}
	db.Create(ledger)
	fmt.Println("✓ Created ledger: 我的帳本")

	// 3. Create sample transactions
	now := time.Now()
	transactions := []models.Transaction{
		{
			ID:            "txn01",
			LedgerID:      ledger.ID,
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
			LedgerID:      ledger.ID,
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
		{
			ID:            "txn03",
			LedgerID:      ledger.ID,
			Title:         "生活用品採購",
			Payer:         "Me",
			Date:          now,
			Currency:      "TWD",
			ExchangeRate:  decimal.NewFromInt(1),
			TotalAmount:   decimal.NewFromInt(299),
			Category:      "購物",
			PaymentMethod: "信用卡",
			Note:          "生活用品",
		},
	}
	for _, txn := range transactions {
		db.Create(&txn)
	}
	fmt.Println("✓ Created 3 sample transactions")

	// 4. Create a trip
	tripLedger := &models.Ledger{
		ID:             uuid.New(),
		UserID:         user.ID,
		Name:           "2024 日本旅行 帳本",
		Type:           "trip",
		Currencies:     datatypes.JSON([]byte(`["JPY"]`)),
		Categories:     datatypes.JSON([]byte(`["住宿", "交通", "飲食", "購物", "門票"]`)),
		PaymentMethods: datatypes.JSON([]byte(`["現金", "信用卡"]`)),
		Members:        datatypes.JSON([]byte(`["我", "小明", "小美"]`)),
	}
	db.Create(tripLedger)
	fmt.Println("✓ Created ledger for trip: 2024 日本旅行 帳本")

	trip := &models.Trip{
		ID:           "japan24",
		Name:         "2024 日本旅行",
		Description:  "東京 + 大阪 5天4夜",
		BaseCurrency: "JPY",
		StartDate:    timePtr(now.AddDate(0, 1, 0)),
		EndDate:      timePtr(now.AddDate(0, 1, 5)),
		CreatedBy:    user.ID,
		LedgerID:     &tripLedger.ID,
	}
	db.Create(trip)
	fmt.Println("✓ Created trip: 2024 日本旅行")

	// 5. Add trip members
	members := []models.TripMember{
		{ID: uuid.New(), TripID: trip.ID, UserID: &user.ID, Name: "我", IsOwner: true},
		{ID: uuid.New(), TripID: trip.ID, Name: "小明"},
		{ID: uuid.New(), TripID: trip.ID, Name: "小美"},
	}
	for _, m := range members {
		db.Create(&m)
	}
	fmt.Println("✓ Added 3 trip members")

	// 6. Create comparison stores
	store1 := &models.ComparisonStore{
		ID:       "donki1",
		TripID:   trip.ID,
		Name:     "唐吉軻德 新宿店",
		Location: "新宿",
	}
	store2 := &models.ComparisonStore{
		ID:       "bic01",
		TripID:   trip.ID,
		Name:     "Bic Camera 有楽町",
		Location: "有楽町",
	}
	db.Create(store1)
	db.Create(store2)
	fmt.Println("✓ Created 2 comparison stores")

	// 7. Create comparison products
	products := []models.ComparisonProduct{
		{ID: uuid.New(), StoreID: store1.ID, Name: "藥妝", Price: decimal.NewFromInt(2980), Currency: "JPY"},
		{ID: uuid.New(), StoreID: store1.ID, Name: "零食組合", Price: decimal.NewFromInt(1500), Currency: "JPY"},
		{ID: uuid.New(), StoreID: store2.ID, Name: "藥妝", Price: decimal.NewFromInt(3200), Currency: "JPY"},
	}
	for _, p := range products {
		db.Create(&p)
	}
	fmt.Println("✓ Created 3 comparison products")

	// 8. Create Trip Transactions with Splits
	// Members: "我" (Owner), "小明", "小美"
	// Txn 1: Dinner (Food), Paid by "我", Split: Even (1000 each)
	txn1 := models.Transaction{
		ID:            "trip_txn01",
		LedgerID:      tripLedger.ID,
		Title:         "第一天晚餐",
		Payer:         "我",
		Date:          now.AddDate(0, 0, 1),
		Currency:      "JPY",
		TotalAmount:   decimal.NewFromInt(3000),
		Category:      "飲食",
		PaymentMethod: "信用卡",
		Note:          "第一天晚餐",
	}
	db.Create(&txn1)

	splits1 := []models.TransactionSplit{
		{TransactionID: txn1.ID, Name: "我", Amount: decimal.NewFromInt(3000), IsPayer: true, MemberID: &members[0].ID}, // Payer record
		{TransactionID: txn1.ID, Name: "我", Amount: decimal.NewFromInt(1000), IsPayer: false, MemberID: &members[0].ID},
		{TransactionID: txn1.ID, Name: "小明", Amount: decimal.NewFromInt(1000), IsPayer: false, MemberID: &members[1].ID},
		{TransactionID: txn1.ID, Name: "小美", Amount: decimal.NewFromInt(1000), IsPayer: false, MemberID: &members[2].ID},
	}
	for _, s := range splits1 {
		db.Create(&s)
	}

	// Txn 2: Transport, Paid by "小明", Split: Even (500 each)
	txn2 := models.Transaction{
		ID:            "trip_txn02",
		LedgerID:      tripLedger.ID,
		Title:         "機場巴士",
		Payer:         "小明",
		Date:          now.AddDate(0, 0, 1),
		Currency:      "JPY",
		TotalAmount:   decimal.NewFromInt(1500),
		Category:      "交通",
		PaymentMethod: "現金",
		Note:          "機場巴士",
	}
	db.Create(&txn2)

	splits2 := []models.TransactionSplit{
		{TransactionID: txn2.ID, Name: "小明", Amount: decimal.NewFromInt(1500), IsPayer: true, MemberID: &members[1].ID}, // Payer record
		{TransactionID: txn2.ID, Name: "我", Amount: decimal.NewFromInt(500), IsPayer: false, MemberID: &members[0].ID},
		{TransactionID: txn2.ID, Name: "小明", Amount: decimal.NewFromInt(500), IsPayer: false, MemberID: &members[1].ID},
		{TransactionID: txn2.ID, Name: "小美", Amount: decimal.NewFromInt(500), IsPayer: false, MemberID: &members[2].ID},
	}
	for _, s := range splits2 {
		db.Create(&s)
	}

	// Txn 3: Souvenirs, Paid by "小美", Living (Personal), No Split (or Split to Self)
	txn3 := models.Transaction{
		ID:            "trip_txn03",
		LedgerID:      tripLedger.ID,
		Title:         "個人伴手禮",
		Payer:         "小美",
		Date:          now.AddDate(0, 0, 2),
		Currency:      "JPY",
		TotalAmount:   decimal.NewFromInt(5000),
		Category:      "購物",
		PaymentMethod: "信用卡",
		Note:          "個人伴手禮",
	}
	db.Create(&txn3)

	splits3 := []models.TransactionSplit{
		{TransactionID: txn3.ID, Name: "小美", Amount: decimal.NewFromInt(5000), IsPayer: true, MemberID: &members[2].ID},  // Payer
		{TransactionID: txn3.ID, Name: "小美", Amount: decimal.NewFromInt(5000), IsPayer: false, MemberID: &members[2].ID}, // Consumer
	}
	for _, s := range splits3 {
		db.Create(&s)
	}

	fmt.Println("✓ Created 3 trip transactions with splits")

	fmt.Println("\n🎉 Seed completed!")
	fmt.Println("   Login: dev / dev123")
}

func timePtr(t time.Time) *time.Time {
	return &t
}
