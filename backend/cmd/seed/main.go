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

	// 3. Create personal space
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
			Title:         "星巴克",
			Payer:         "Antigravity",
			Date:          now.Add(-2 * time.Hour),
			Currency:      "TWD",
			ExchangeRate:  decimal.NewFromInt(1),
			TotalAmount:   decimal.NewFromInt(310),
			BillingAmount: decimal.NewFromInt(310),
			Category:      "餐飲",
			PaymentMethod: "信用卡",
			Note:          "跟同事下午茶",
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

	// 4.1 Add items to personal transactions
	db.Create(&models.TransactionItem{
		ID:            uuid.New(),
		TransactionID: "txn_p_01",
		Name:          "特大杯拿鐵",
		UnitPrice:     decimal.NewFromInt(155),
		Quantity:      decimal.NewFromInt(2),
		Amount:        decimal.NewFromInt(310),
	})
	fmt.Println("✓ Created sample transactions and items for personal space")

	// 5. Create a sample space (unified Ledger + Trip)
	sampleSpace := &models.Ledger{
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
	db.Create(sampleSpace)
	fmt.Println("✓ Created space: 2024 東京春櫻季 (pinned)")

	// 6. Add ledger members for sample space
	memberDev := models.LedgerMember{ID: uuid.New(), LedgerID: sampleSpace.ID, UserID: user.ID, Role: "owner", Alias: "Antigravity"}
	memberMing := models.LedgerMember{ID: uuid.New(), LedgerID: sampleSpace.ID, UserID: userMing.ID, Role: "member", Alias: "小明"}
	memberMei := models.LedgerMember{ID: uuid.New(), LedgerID: sampleSpace.ID, UserID: userMei.ID, Role: "member", Alias: "小美"}
	db.Create(&memberDev)
	db.Create(&memberMing)
	db.Create(&memberMei)
	fmt.Println("✓ Added members to sample space")

	// 7. Create comparison stores for sample space
	store1 := &models.ComparisonStore{
		ID:       "store_s_01",
		LedgerID: sampleSpace.ID,
		Name:     "唐吉軻德 澀谷店",
		Location: "澀谷",
		GoogleMapURL: "https://maps.app.goo.gl/ShibuyaDonki",
	}
	store2 := &models.ComparisonStore{
		ID:       "store_s_02",
		LedgerID: sampleSpace.ID,
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

	// 9. Create sample transactions with splits
	txnSampleID := "txn_s_01"
	txnSample := models.Transaction{
		ID:            txnSampleID,
		LedgerID:      sampleSpace.ID,
		Title:         "利木津巴士",
		Payer:         "Antigravity",
		Date:          now.AddDate(0, 1, 0),
		Currency:      "JPY",
		TotalAmount:   decimal.NewFromInt(9000),
		ExchangeRate:  decimal.NewFromFloat(0.216),
		BillingAmount: decimal.NewFromFloat(1944),
		HandlingFee:   decimal.NewFromFloat(29.16), // 1.5% fee
		Category:      "交通",
		PaymentMethod: "信用卡",
	}
	db.Create(&txnSample)

	// Add items to trip transaction
	db.Create(&models.TransactionItem{
		ID:            uuid.New(),
		TransactionID: txnSampleID,
		Name:          "成人票",
		UnitPrice:     decimal.NewFromInt(3000),
		Quantity:      decimal.NewFromInt(3),
		Amount:        decimal.NewFromInt(9000),
	})

	// Add another complex transaction with items
	txnComplexID := "txn_s_02"
	txnComplex := models.Transaction{
		ID:            txnComplexID,
		LedgerID:      sampleSpace.ID,
		Title:         "一蘭拉麵",
		Payer:         "Antigravity",
		Date:          now.AddDate(0, 1, 1),
		Currency:      "JPY",
		TotalAmount:   decimal.NewFromInt(5800),
		ExchangeRate:  decimal.NewFromFloat(0.216),
		BillingAmount: decimal.NewFromFloat(1253),
		HandlingFee:   decimal.NewFromInt(0), // Cash payment
		Category:      "飲食",
		PaymentMethod: "現金",
	}
	db.Create(&txnComplex)

	items := []models.TransactionItem{
		{ID: uuid.New(), TransactionID: txnComplexID, Name: "天然豚骨拉麵", UnitPrice: decimal.NewFromInt(980), Quantity: decimal.NewFromInt(3), Amount: decimal.NewFromInt(2940)},
		{ID: uuid.New(), TransactionID: txnComplexID, Name: "加麵", UnitPrice: decimal.NewFromInt(210), Quantity: decimal.NewFromInt(2), Amount: decimal.NewFromInt(420)},
		{ID: uuid.New(), TransactionID: txnComplexID, Name: "生啤酒", UnitPrice: decimal.NewFromInt(580), Quantity: decimal.NewFromInt(3), Amount: decimal.NewFromInt(1740)},
		{ID: uuid.New(), TransactionID: txnComplexID, Name: "半熟鹽味蛋", UnitPrice: decimal.NewFromInt(140), Quantity: decimal.NewFromInt(5), Amount: decimal.NewFromInt(700)},
	}
	for _, item := range items {
		db.Create(&item)
	}

	// Splits: Antigravity paid for everyone, split 3 ways
	splits := []models.TransactionSplit{
		{ID: uuid.New(), TransactionID: txnSampleID, Name: "Antigravity", Amount: decimal.NewFromInt(9000), IsPayer: true, MemberID: &memberDev.ID},
		{ID: uuid.New(), TransactionID: txnSampleID, Name: "Antigravity", Amount: decimal.NewFromInt(3000), IsPayer: false, MemberID: &memberDev.ID},
		{ID: uuid.New(), TransactionID: txnSampleID, Name: "小明", Amount: decimal.NewFromInt(3000), IsPayer: false, MemberID: &memberMing.ID},
		{ID: uuid.New(), TransactionID: txnSampleID, Name: "小美", Amount: decimal.NewFromInt(3000), IsPayer: false, MemberID: &memberMei.ID},
	}
	for _, s := range splits {
		db.Create(&s)
	}
	fmt.Println("✓ Created sample transaction with 3-way split")

	fmt.Println("\n🎉 Seed completed successfully!")
	fmt.Println("   Login User: dev")
	fmt.Println("   Password:   dev123")
}

func timePtr(t time.Time) *time.Time {
	return &t
}