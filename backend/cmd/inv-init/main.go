package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"lovelion/internal/models"
	"lovelion/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var scanner *bufio.Scanner

func ask(prompt, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultVal)
	} else {
		fmt.Printf("%s: ", prompt)
	}
	scanner.Scan()
	val := strings.TrimSpace(scanner.Text())
	if val == "" {
		return defaultVal
	}
	return val
}

func parseAmount(s string) (int, error) {
	s = strings.ReplaceAll(s, ",", "")
	return strconv.Atoi(s)
}

func askInt(prompt string, defaultVal int) int {
	for {
		s := ask(prompt, formatNumber(defaultVal))
		n, err := parseAmount(s)
		if err != nil {
			fmt.Printf("  ⚠ 請輸入有效數字（可含千分位逗號）: %q\n", s)
			continue
		}
		return n
	}
}

func askYN(prompt string) bool {
	for {
		fmt.Printf("%s (y/n): ", prompt)
		scanner.Scan()
		val := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if val == "y" {
			return true
		}
		if val == "n" {
			return false
		}
		fmt.Println("  ⚠ 請輸入 y 或 n")
	}
}

func lastDayOfMonth(ym string) time.Time {
	parts := strings.Split(ym, "-")
	year, _ := strconv.Atoi(parts[0])
	month, _ := strconv.Atoi(parts[1])
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC)
}

func main() {
	scanner = bufio.NewScanner(os.Stdin)

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}

	fmt.Println("=== 投資初始化 ===")
	fmt.Println()

	// [1] Owner
	fmt.Println("[1] 設定 Owner")
	username := ask("Username", "")
	if username == "" {
		log.Fatal("username is required")
	}

	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Fatalf("user %q not found", username)
	}

	var ownerMember models.InvMember
	ownerExists := db.Where("user_id = ?", user.ID).First(&ownerMember).Error == nil

	if ownerExists {
		fmt.Printf("  ✓ Owner 已存在: %s (id=%s)\n", ownerMember.Name, ownerMember.ID)
	} else {
		displayName := ask("Display name", user.DisplayName)
		id := utils.MustNewShortID(db, "inv_members", "id")
		ownerMember = models.InvMember{
			ID:      id,
			Name:    displayName,
			UserID:  &user.ID,
			IsOwner: true,
			Active:  true,
		}
		if err := db.Create(&ownerMember).Error; err != nil {
			log.Fatalf("create owner: %v", err)
		}
		fmt.Printf("  ✓ Owner 已建立: %s (id=%s)\n", displayName, id)
	}
	fmt.Println()

	// [2] Base period
	fmt.Println("[2] 基準期")
	baseYM := ask("基準年月", "2023-12")
	baseDate := lastDayOfMonth(baseYM)
	fmt.Printf("  入金日期: %s\n", baseDate.Format("2006-01-02"))
	fmt.Println()

	// Ensure base settlement exists (needed for FK on futures/stocks statements)
	var baseSettlement models.InvSettlement
	if db.First(&baseSettlement, "year_month = ?", baseYM).Error != nil {
		baseSettlement = models.InvSettlement{
			YearMonth: baseYM,
			Status:    "baseline",
		}
		if err := db.Create(&baseSettlement).Error; err != nil {
			log.Fatalf("create base settlement: %v", err)
		}
	}

	// [3] Futures base
	fmt.Println("[3] 期貨實質權益基準值")
	futuresEquity := askInt("實質權益", 0)

	var existingFutures models.InvFuturesStatement
	if db.First(&existingFutures, "year_month = ?", baseYM).Error == nil {
		fmt.Printf("  ✓ %s 期貨基準已存在 (ending_equity=%d)\n", baseYM, existingFutures.EndingEquity)
	} else {
		stmt := models.InvFuturesStatement{
			YearMonth:    baseYM,
			EndingEquity: futuresEquity,
		}
		if err := db.Create(&stmt).Error; err != nil {
			log.Fatalf("create futures statement: %v", err)
		}
		fmt.Printf("  ✓ %s 期貨基準已建立 (實質權益=%d)\n", baseYM, futuresEquity)
	}
	fmt.Println()

	// [4] Members & initial balances
	fmt.Println("[4] 成員初始值")

	// Owner's initial balance
	fmt.Printf("  --- %s (owner) ---\n", ownerMember.Name)
	ownerBalance := askInt("  初始結餘", 0)
	if ownerBalance > 0 {
		createDeposit(db, ownerMember.ID, baseDate, ownerBalance)
		fmt.Printf("  ✓ 入金 %s @ %s (繼承舊檔)\n", formatNumber(ownerBalance), baseDate.Format("2006-01-02"))
	}
	fmt.Println()

	// Additional members
	for {
		if !askYN("新增成員？") {
			break
		}
		var name string
		for {
			name = ask("  成員名稱", "")
			if name != "" {
				break
			}
			fmt.Println("  ⚠ 成員名稱不可為空")
		}

		id := utils.MustNewShortID(db, "inv_members", "id")
		member := models.InvMember{
			ID:     id,
			Name:   name,
			Active: true,
		}
		if err := db.Create(&member).Error; err != nil {
			log.Fatalf("create member: %v", err)
		}
		fmt.Printf("  ✓ 成員已建立: %s (id=%s)\n", name, id)

		balance := askInt("  初始結餘", 0)
		if balance > 0 {
			createDeposit(db, id, baseDate, balance)
			fmt.Printf("  ✓ 入金 %s @ %s (繼承舊檔)\n", formatNumber(balance), baseDate.Format("2006-01-02"))
		}
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("✓ 初始化完成")
}

func createDeposit(db *gorm.DB, memberID string, date time.Time, amount int) {
	txn := models.InvMemberTransaction{
		MemberID: memberID,
		Date:     date,
		Type:     "deposit",
		Amount:   amount,
		Note:     "繼承舊檔",
	}
	if err := db.Create(&txn).Error; err != nil {
		log.Fatalf("create deposit: %v", err)
	}

	db.Model(&models.InvMember{}).Where("id = ?", memberID).
		Update("net_investment", gorm.Expr("net_investment + ?", amount))
}

func formatNumber(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}
	var parts []string
	for i := len(s); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{s[start:i]}, parts...)
	}
	return strings.Join(parts, ",")
}
