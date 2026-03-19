import type { Image } from './image'

export interface Transaction {
  id: string
  space_id: string
  title: string
  payer: string
  date: string
  currency: string
  total_amount: string
  exchange_rate: string
  billing_amount: string
  handling_fee: string
  category: string
  payment_method: string
  note: string
  created_at: string
  updated_at: string
  items?: TransactionItem[]
  splits?: TransactionSplit[]
  images?: Image[]
}

export interface TransactionItem {
  id: string
  transaction_id: string
  name: string
  unit_price: string
  quantity: string
  discount: string
  amount: string
  created_at: string
  updated_at: string
}

export interface TransactionSplit {
  id: string
  transaction_id: string
  name: string
  amount: string
  is_payer: boolean
  created_at: string
}
