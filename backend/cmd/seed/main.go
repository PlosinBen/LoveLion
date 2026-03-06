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

	// 1. Create main test user
	user := &models.User{
		ID:          uuid.New(),
		Username:    "dev",
		DisplayName: "Antigravity",
	}
	user.SetPassword("dev123")
	db.Create(user)
	fmt.Println("✓ Created user: dev / dev123")

	// 2. Create other test users for collaboration
	userMing := &models.User{ID: uuid.New(), Username: "ming", DisplayName: "小明"}
	userMei := &models.User{ID: uuid.New(), Username: "mei", DisplayName: "小美"}
	db.Create(userMing)
	db.Create(userMei)
	fmt.Println("✓ Created collaborator users: ming, mei")

	// 3. Create personal space (ledger)
	personalSpace := &models.Ledger{
		ID:             uuid.New(),
		UserID:         user.ID,
		Name:           "日常開銷",
		Type:           "personal",
		BaseCurrency:   "TWD",
		Currencies:     datatypes.JSON([]byte(`["TWD", "JPY", "USD"]`)),
		Categories:     datatypes.JSON([]byte(`["餐飲", "交通", "購物", "娛樂", "生活"]`)),
		PaymentMethods: datatypes.JSON([]byte(`["現金", "信用卡", "Line Pay"]`)),
		IsPinned:       true,
	}
	db.Create(personalSpace)
	db.Create(&models.LedgerMember{ID: uuid.New(), LedgerID: personalSpace.ID, UserID: user.ID, Role: "owner"})
	fmt.Println("✓ Created space: 日常開銷 (pinned)")

	// 4. Create sample transactions for personal space
	now := time.Now()
	transactions := []models.Transaction{
		{
			ID:            "txn_p_01",
			LedgerID:      personalSpace.ID,
			Title:         "星巴克拿鐵",
			Payer:         "Antigravity",
			Date:          now.Add(-2 * time.Hour),
			Currency:      "TWD",
			ExchangeRate:  decimal.NewFromInt(1),
			TotalAmount:   decimal.NewFromInt(155),
			BillingAmount: decimal.NewFromInt(155),
			Category:      "餐飲",
			PaymentMethod: "信用卡",
		},
		{
			ID:            "txn_p_02",
			LedgerID:      personalSpace.ID,
			Title:         "捷運定期票",
			Payer:         "Antigravity",
			Date:          now.Add(-24 * time.Hour),
			Currency:      "TWD",
			ExchangeRate:  decimal.NewFromInt(1),
			TotalAmount:   decimal.NewFromInt(1200),
			BillingAmount: decimal.NewFromInt(1200),
			Category:      "交通",
			PaymentMethod: "現金",
		},
	}
	for _, txn := range transactions {
		db.Create(&txn)
	}
	fmt.Println("✓ Created sample transactions for personal space")

	// 5. Create a trip space (unified Ledger + Trip)
	tripSpace := &models.Ledger{
		ID:             uuid.New(),
		UserID:         user.ID,
		Name:           "2024 東京春櫻季",
		Description:    "5 天 4 夜 東京賞櫻團",
		Type:           "trip",
		BaseCurrency:   "TWD",
		Currencies:     datatypes.JSON([]byte(`["TWD", "JPY"]`)),
		Categories:     datatypes.JSON([]byte(`["住宿", "交通", "飲食", "購物", "娛樂"]`)),
		PaymentMethods: datatypes.JSON([]byte(`["現金", "信用卡"]`)),
		StartDate:      timePtr(now.AddDate(0, 1, 0)),
		EndDate:        timePtr(now.AddDate(0, 1, 5)),
		IsPinned:       true,
		CoverImage:     "https://images.unsplash.com/photo-1493976040374-85c8e12f0c0e?q=80&w=800&auto=format&fit=crop",
	}
	db.Create(tripSpace)
	fmt.Println("✓ Created space: 2024 東京春櫻季 (pinned trip)")

	// 6. Add ledger members for trip space
	memberDev := models.LedgerMember{ID: uuid.New(), LedgerID: tripSpace.ID, UserID: user.ID, Role: "owner", Alias: "Antigravity"}
	memberMing := models.LedgerMember{ID: uuid.New(), LedgerID: tripSpace.ID, UserID: userMing.ID, Role: "member", Alias: "小明"}
	memberMei := models.LedgerMember{ID: uuid.New(), LedgerID: tripSpace.ID, UserID: userMei.ID, Role: "member", Alias: "小美"}
	db.Create(&memberDev)
	db.Create(&memberMing)
	db.Create(&memberMei)
	fmt.Println("✓ Added members to trip space")

	// 7. Create comparison stores for trip space
	store1 := &models.ComparisonStore{
		ID:       "store_s_01",
		LedgerID: tripSpace.ID,
		Name:     "唐吉軻德 澀谷店",
		Location: "澀谷",
		GoogleMapURL: "https://maps.app.goo.gl/ShibuyaDonki",
	}
	store2 := &models.ComparisonStore{
		ID:       "store_s_02",
		LedgerID: tripSpace.ID,
		Name:     "Bic Camera 新宿",
		Location: "新宿",
	}
	db.Create(store1)
	db.Create(store2)
	fmt.Println("✓ Created comparison stores")

	// 8. Create comparison products
	products := []models.ComparisonProduct{
		{ID: uuid.New(), StoreID: store1.ID, Name: "一蘭拉麵泡麵", Price: decimal.NewFromInt(1850), Currency: "JPY"},
		{ID: uuid.New(), StoreID: store2.ID, Name: "一蘭拉麵泡麵", Price: decimal.NewFromInt(1980), Currency: "JPY"},
		{ID: uuid.New(), StoreID: store1.ID, Name: "Dyson 吹風機", Price: decimal.NewFromInt(45000), Currency: "JPY"},
	}
	for _, p := range products {
		db.Create(&p)
	}
	fmt.Println("✓ Created comparison products")

	// 9. Create Trip Transactions with Splits
	txnTripID := "txn_t_01"
	txnTrip := models.Transaction{
		ID:            txnTripID,
		LedgerID:      tripSpace.ID,
		Title:         "成田機場利木津巴士",
		Payer:         "Antigravity",
		Date:          now.AddDate(0, 1, 0),
		Currency:      "JPY",
		TotalAmount:   decimal.NewFromInt(9000),
		BillingAmount: decimal.NewFromInt(1950),
		ExchangeRate:  decimal.NewFromFloat(0.216),
		Category:      "交通",
		PaymentMethod: "信用卡",
	}
	db.Create(&txnTrip)

	// Splits: Antigravity paid for everyone, split 3 ways
	splits := []models.TransactionSplit{
		{ID: uuid.New(), TransactionID: txnTripID, Name: "Antigravity", Amount: decimal.NewFromInt(9000), IsPayer: true, MemberID: &memberDev.ID},
		{ID: uuid.New(), TransactionID: txnTripID, Name: "Antigravity", Amount: decimal.NewFromInt(3000), IsPayer: false, MemberID: &memberDev.ID},
		{ID: uuid.New(), TransactionID: txnTripID, Name: "小明", Amount: decimal.NewFromInt(3000), IsPayer: false, MemberID: &memberMing.ID},
		{ID: uuid.New(), TransactionID: txnTripID, Name: "小美", Amount: decimal.NewFromInt(3000), IsPayer: false, MemberID: &memberMei.ID},
	}
	for _, s := range splits {
		db.Create(&s)
	}
	fmt.Println("✓ Created trip transaction with 3-way split")

	fmt.Println("\n🎉 Seed completed successfully!")
	fmt.Println("   Login User: dev")
	fmt.Println("   Password:   dev123")
}

func timePtr(t time.Time) *time.Time {
	return &t
}