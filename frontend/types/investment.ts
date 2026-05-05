export interface InvMember {
  id: string
  name: string
  user_id: string | null
  is_owner: boolean
  active: boolean
  sort_order: number
  net_investment: number
  created_at: string
}

export interface InvSettlement {
  year_month: string
  status: 'draft' | 'completed'
  total_profit_loss: number
  total_weight: number
  profit_loss_per_weight: number
  created_at: string
  updated_at: string
}

export interface InvMemberTransaction {
  id: string
  member_id: string
  member_name?: string
  date: string
  type: 'deposit' | 'withdrawal' | 'profit_loss'
  amount: number
  note: string
}

export interface InvSettlementAllocation {
  year_month: string
  member_id: string
  member_name?: string
  weight: number
  amount: number
  deposit: number
  withdrawal: number
  balance: number
}

export interface InvFuturesStatement {
  year_month: string
  ending_equity: number
  floating_profit_loss: number
  realized_profit_loss: number
  deposit: number
  withdrawal: number
  profit_loss: number
}

export interface InvStockStatement {
  year_month: string
  account_balance: number
  market_value: number
  deposit: number
  withdrawal: number
  profit_loss: number
}

export interface InvStockHolding {
  id: string
  year_month: string
  symbol: string
  shares: number
  closing_price: number
}

export interface InvStockTrade {
  id: string
  trade_date: string
  symbol: string
  shares: number
  price: number
  fee: number
  tax: number
  note: string
}

export interface InvSettlementDetail {
  year_month: string
  status: string
  total_profit_loss: number
  total_weight: number
  profit_loss_per_weight: number
  futures_statement: InvFuturesStatement | null
  stock_statement: (InvStockStatement & { holdings?: InvStockHolding[] }) | null
  allocations: AllocationPreview[]
}

export interface AllocationPreview {
  member_id: string
  member_name: string
  is_owner: boolean
  weight: number
  amount: number
  deposit: number
  withdrawal: number
  balance: number
}
