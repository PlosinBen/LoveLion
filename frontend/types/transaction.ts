import type { Image } from './image'

export type TransactionType = 'expense' | 'payment'

export interface Transaction {
  id: string
  space_id: string
  type: TransactionType
  title: string
  date: string
  currency: string
  total_amount: string
  note: string
  created_at: string
  updated_at: string
  expense?: TransactionExpense
  debts?: TransactionDebt[]
  images?: Image[]
}

export interface TransactionExpense {
  id: string
  transaction_id: string
  category: string
  exchange_rate: string
  billing_amount: string
  handling_fee: string
  payment_method: string
  location_url: string
  items?: TransactionExpenseItem[]
}

export interface TransactionExpenseItem {
  id: string
  expense_id: string
  name: string
  unit_price: string
  quantity: string
  discount: string
  amount: string
}

export interface TransactionDebt {
  id: string
  transaction_id: string
  payer_name: string
  payee_name: string
  amount: string
  settled_amount: string
  is_spot_paid: boolean
}
