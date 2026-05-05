package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InvMember struct {
	ID            string     `gorm:"type:varchar(21);primary_key" json:"id"`
	Name          string     `gorm:"type:varchar(50);not null" json:"name"`
	UserID        *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	IsOwner       bool       `gorm:"not null;default:false" json:"is_owner"`
	Active        bool       `gorm:"not null;default:true" json:"active"`
	SortOrder     int        `gorm:"not null;default:0" json:"sort_order"`
	NetInvestment int        `gorm:"not null;default:0" json:"net_investment"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (InvMember) TableName() string {
	return "inv_members"
}

type InvSettlement struct {
	YearMonth           string    `gorm:"type:varchar(7);primary_key" json:"year_month"`
	Status              string    `gorm:"type:varchar(10);not null;default:'draft'" json:"status"`
	TotalProfitLoss     int       `gorm:"not null;default:0" json:"total_profit_loss"`
	TotalWeight         int       `gorm:"not null;default:0" json:"total_weight"`
	ProfitLossPerWeight int       `gorm:"not null;default:0" json:"profit_loss_per_weight"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Allocations      []InvSettlementAllocation `gorm:"foreignKey:YearMonth" json:"allocations,omitempty"`
	FuturesStatement *InvFuturesStatement      `gorm:"foreignKey:YearMonth" json:"futures_statement,omitempty"`
	StockStatement   *InvStockStatement        `gorm:"foreignKey:YearMonth" json:"stock_statement,omitempty"`
}

func (InvSettlement) TableName() string {
	return "inv_settlements"
}

type InvMemberTransaction struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	MemberID string    `gorm:"type:varchar(21);not null" json:"member_id"`
	Date     time.Time `gorm:"type:date;not null" json:"date"`
	Type     string    `gorm:"type:varchar(20);not null" json:"type"`
	Amount   int       `gorm:"not null;default:0" json:"amount"`
	Note     string    `gorm:"type:text" json:"note"`

	Member *InvMember `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}

func (InvMemberTransaction) TableName() string {
	return "inv_member_transactions"
}

type InvSettlementAllocation struct {
	YearMonth  string `gorm:"type:varchar(7);primary_key" json:"year_month"`
	MemberID   string `gorm:"type:varchar(21);primary_key" json:"member_id"`
	Weight     int    `gorm:"not null;default:0" json:"weight"`
	Amount     int    `gorm:"not null;default:0" json:"amount"`
	Deposit    int    `gorm:"not null;default:0" json:"deposit"`
	Withdrawal int    `gorm:"not null;default:0" json:"withdrawal"`
	Balance    int    `gorm:"not null;default:0" json:"balance"`

	Member *InvMember `gorm:"foreignKey:MemberID" json:"member,omitempty"`
}

func (InvSettlementAllocation) TableName() string {
	return "inv_settlement_allocations"
}

type InvFuturesStatement struct {
	YearMonth          string `gorm:"type:varchar(7);primary_key" json:"year_month"`
	EndingEquity       int    `gorm:"not null;default:0" json:"ending_equity"`
	FloatingProfitLoss int    `gorm:"not null;default:0" json:"floating_profit_loss"`
	RealizedProfitLoss int    `gorm:"not null;default:0" json:"realized_profit_loss"`
	Deposit            int    `gorm:"not null;default:0" json:"deposit"`
	Withdrawal         int    `gorm:"not null;default:0" json:"withdrawal"`
	ProfitLoss         int    `gorm:"not null;default:0" json:"profit_loss"`
}

func (InvFuturesStatement) TableName() string {
	return "inv_futures_statements"
}

type InvStockStatement struct {
	YearMonth      string `gorm:"type:varchar(7);primary_key" json:"year_month"`
	AccountBalance int    `gorm:"not null;default:0" json:"account_balance"`
	MarketValue    int    `gorm:"not null;default:0" json:"market_value"`
	Deposit        int    `gorm:"not null;default:0" json:"deposit"`
	Withdrawal     int    `gorm:"not null;default:0" json:"withdrawal"`
	ProfitLoss     int    `gorm:"not null;default:0" json:"profit_loss"`

	Holdings []InvStockHolding `gorm:"foreignKey:YearMonth" json:"holdings,omitempty"`
}

func (InvStockStatement) TableName() string {
	return "inv_stock_statements"
}

type InvStockHolding struct {
	ID           uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	YearMonth    string          `gorm:"type:varchar(7);not null" json:"year_month"`
	Symbol       string          `gorm:"type:varchar(20);not null" json:"symbol"`
	Shares       int             `gorm:"not null;default:0" json:"shares"`
	ClosingPrice decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"closing_price"`
	MarketValue  int             `gorm:"not null;default:0" json:"market_value"`
}

func (InvStockHolding) TableName() string {
	return "inv_stock_holdings"
}

type InvStockTrade struct {
	ID        uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TradeDate time.Time       `gorm:"type:date;not null" json:"trade_date"`
	Symbol    string          `gorm:"type:varchar(20);not null" json:"symbol"`
	Shares    int             `gorm:"not null;default:0" json:"shares"`
	Price     decimal.Decimal `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	Fee       int             `gorm:"not null;default:0" json:"fee"`
	Tax       int             `gorm:"not null;default:0" json:"tax"`
}

func (InvStockTrade) TableName() string {
	return "inv_stock_trades"
}
