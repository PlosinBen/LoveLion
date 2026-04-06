import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

// Mock useApi
const mockGet = vi.fn()
vi.mock('~/composables/useApi', () => ({
  useApi: () => ({
    get: mockGet,
    post: vi.fn(),
    put: vi.fn(),
    loading: { value: false },
    error: { value: null },
  }),
}))

import { useTransactionForm } from '~/composables/useTransactionForm'
import type { Transaction } from '~/types'

describe('useTransactionForm', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  // --- loadSpaceConfig / fetchSpaceConfig ---

  it('fetchSpaceConfig loads space and sets categories/currencies from space data', async () => {
    mockGet.mockResolvedValueOnce({
      base_currency: 'JPY',
      split_members: ['Alice', 'Bob'],
      categories: ['食費', '交通費'],
      currencies: ['JPY', 'USD'],
      payment_methods: ['現金', 'カード'],
    })

    const form = useTransactionForm('space-1')
    await form.fetchSpaceConfig()

    expect(form.baseCurrency.value).toBe('JPY')
    expect(form.memberOptions.value).toEqual(['Alice', 'Bob'])
    expect(form.categories.value).toEqual([
      { label: '食費', value: '食費' },
      { label: '交通費', value: '交通費' },
    ])
    expect(form.availableCurrencies.value).toEqual([
      { label: 'JPY', value: 'JPY' },
      { label: 'USD', value: 'USD' },
    ])
    expect(form.paymentMethods.value).toEqual([
      { label: '現金', value: '現金' },
      { label: 'カード', value: 'カード' },
    ])
    // currency should be set to base_currency
    expect(form.expenseForm.value.currency).toBe('JPY')
  })

  it('fetchSpaceConfig uses defaults when space has no categories/currencies', async () => {
    mockGet.mockResolvedValueOnce({
      base_currency: 'TWD',
    })

    const form = useTransactionForm('space-2')
    await form.fetchSpaceConfig()

    expect(form.categories.value).toHaveLength(6) // default 6 categories
    expect(form.categories.value[0]).toEqual({ label: '餐飲', value: '餐飲' })
    expect(form.availableCurrencies.value).toHaveLength(3) // default TWD/JPY/USD
    expect(form.paymentMethods.value).toEqual([])
  })

  it('fetchSpaceConfig returns null on API error', async () => {
    mockGet.mockRejectedValueOnce(new Error('Network error'))

    const form = useTransactionForm('space-err')
    const result = await form.fetchSpaceConfig()

    expect(result).toBeNull()
  })

  // --- populateFromTransaction ---

  it('populateFromTransaction fills expense form from transaction', () => {
    const form = useTransactionForm('space-1')
    const txn: Transaction = {
      id: 'txn-1',
      space_id: 'space-1',
      type: 'expense',
      title: '午餐',
      date: '2026-04-01T12:00:00Z',
      currency: 'JPY',
      total_amount: '1500',
      note: '拉麵',
      created_at: '',
      updated_at: '',
      expense: {
        id: 'exp-1',
        transaction_id: 'txn-1',
        category: '餐飲',
        exchange_rate: '0.22',
        billing_amount: '330',
        handling_fee: '10',
        payment_method: '現金',
        items: [
          { id: 'i1', expense_id: 'exp-1', name: '拉麵', unit_price: '1500', quantity: '1', discount: '0', amount: '1500' },
        ],
      },
      debts: [
        { id: 'd1', transaction_id: 'txn-1', payer_name: 'Alice', payee_name: 'Bob', amount: '750', settled_amount: '0', is_spot_paid: false },
      ],
    }

    form.populateFromTransaction(txn)

    expect(form.transactionType.value).toBe('expense')
    expect(form.expenseForm.value.title).toBe('午餐')
    expect(form.expenseForm.value.total_amount).toBe(1500)
    expect(form.expenseForm.value.currency).toBe('JPY')
    expect(form.expenseForm.value.category).toBe('餐飲')
    expect(form.expenseForm.value.exchange_rate).toBe(0.22)
    expect(form.expenseForm.value.handling_fee).toBe(10)
    expect(form.expenseForm.value.manual_rate).toBe(false) // handling_fee > 0
    expect(form.expenseForm.value.items).toHaveLength(1)
    expect(form.expenseForm.value.items[0]?.name).toBe('拉麵')
    expect(form.debts.value).toHaveLength(1)
    expect(form.debts.value[0]?.payer_name).toBe('Alice')
  })

  it('populateFromTransaction fills payment form from transaction', () => {
    const form = useTransactionForm('space-1')
    const txn: Transaction = {
      id: 'txn-2',
      space_id: 'space-1',
      type: 'payment',
      title: '付款',
      date: '2026-04-02T10:00:00Z',
      currency: 'TWD',
      total_amount: '500',
      note: '',
      created_at: '',
      updated_at: '',
      debts: [
        { id: 'd2', transaction_id: 'txn-2', payer_name: 'Bob', payee_name: 'Alice', amount: '500', settled_amount: '0', is_spot_paid: false },
      ],
    }

    form.populateFromTransaction(txn)

    expect(form.transactionType.value).toBe('payment')
    expect(form.paymentForm.value.title).toBe('付款')
    expect(form.paymentForm.value.total_amount).toBe(500)
    expect(form.paymentForm.value.payer_name).toBe('Bob')
    expect(form.paymentForm.value.payee_name).toBe('Alice')
  })

  it('populateFromTransaction uses empty items array as default', () => {
    const form = useTransactionForm('space-1')
    const txn: Transaction = {
      id: 'txn-3',
      space_id: 'space-1',
      type: 'expense',
      title: '雜費',
      date: '2026-04-03T00:00:00Z',
      currency: 'TWD',
      total_amount: '100',
      note: '',
      created_at: '',
      updated_at: '',
      expense: {
        id: 'exp-3',
        transaction_id: 'txn-3',
        category: '其他',
        exchange_rate: '1',
        billing_amount: '100',
        handling_fee: '0',
        payment_method: '',
        items: [], // empty
      },
    }

    form.populateFromTransaction(txn)

    // Should fall back to default empty item
    expect(form.expenseForm.value.items).toHaveLength(1)
    expect(form.expenseForm.value.items[0]?.name).toBe('')
  })

  // --- buildExpensePayload ---

  it('buildExpensePayload builds correct payload for same-currency expense', () => {
    const form = useTransactionForm('space-1')
    form.baseCurrency.value = 'TWD'
    form.expenseForm.value = {
      date: new Date('2026-04-01T12:00:00Z'),
      title: '午餐',
      total_amount: 300,
      category: '餐飲',
      payment_method: '現金',
      note: '好吃',
      items: [
        { name: '便當', unit_price: 100, quantity: 2, discount: 0 },
        { name: '飲料', unit_price: 50, quantity: 2, discount: 0 },
        { name: '', unit_price: 0, quantity: 1, discount: 0 }, // empty, should be filtered
      ],
      currency: 'TWD',
      manual_rate: true,
      exchange_rate: 0,
      billing_amount: 0,
      handling_fee: 0,
    }

    const payload = form.buildExpensePayload()

    expect(payload.title).toBe('午餐')
    expect(payload.total_amount).toBe(300)
    expect(payload.currency).toBe('TWD')
    expect(payload.expense.exchange_rate).toBe(1) // same currency
    expect(payload.expense.billing_amount).toBe(300) // same currency = total_amount
    expect(payload.expense.handling_fee).toBe(0) // same currency
    expect(payload.expense.items).toHaveLength(2) // empty item filtered
    expect(payload.debts).toBeUndefined() // no debts
  })

  it('buildExpensePayload uses custom rates for foreign currency', () => {
    const form = useTransactionForm('space-1')
    form.baseCurrency.value = 'TWD'
    form.expenseForm.value = {
      date: new Date('2026-04-01T12:00:00Z'),
      title: '拉麵',
      total_amount: 1000,
      category: '餐飲',
      payment_method: '',
      note: '',
      items: [{ name: '拉麵', unit_price: 1000, quantity: 1, discount: 0 }],
      currency: 'JPY',
      manual_rate: false,
      exchange_rate: 0.22,
      billing_amount: 220,
      handling_fee: 5,
    }

    const payload = form.buildExpensePayload()

    expect(payload.expense.exchange_rate).toBe(0.22)
    expect(payload.expense.billing_amount).toBe(220)
    expect(payload.expense.handling_fee).toBe(5)
  })

  it('buildExpensePayload includes debts when present', () => {
    const form = useTransactionForm('space-1')
    form.baseCurrency.value = 'TWD'
    form.expenseForm.value.currency = 'TWD'
    form.expenseForm.value.total_amount = 600
    form.debts.value = [
      { payer_name: 'Alice', payee_name: 'Bob', amount: 300, is_spot_paid: false },
      { payer_name: '', payee_name: '', amount: 0, is_spot_paid: false }, // invalid, should be filtered
    ]

    const payload = form.buildExpensePayload()

    expect(payload.debts).toHaveLength(1)
    expect(payload.debts![0]).toEqual({ payer_name: 'Alice', payee_name: 'Bob', amount: 300, is_spot_paid: false })
  })

  // --- buildPaymentPayload ---

  it('buildPaymentPayload builds correct payload', () => {
    const form = useTransactionForm('space-1')
    form.paymentForm.value = {
      date: new Date('2026-04-02T10:00:00Z'),
      title: '付款',
      payer_name: 'Bob',
      payee_name: 'Alice',
      total_amount: 500,
      note: '上次的',
    }

    const payload = form.buildPaymentPayload()

    expect(payload.title).toBe('付款')
    expect(payload.total_amount).toBe(500)
    expect(payload.payer_name).toBe('Bob')
    expect(payload.payee_name).toBe('Alice')
    expect(payload.note).toBe('上次的')
  })

  // --- validateExpense ---

  it('validateExpense returns false when total_amount is 0', () => {
    const form = useTransactionForm('space-1')
    form.expenseForm.value.total_amount = 0

    expect(form.validateExpense()).toBe(false)
  })

  it('validateExpense returns true when total_amount > 0', () => {
    const form = useTransactionForm('space-1')
    form.expenseForm.value.total_amount = 100

    expect(form.validateExpense()).toBe(true)
  })

  // --- validatePayment ---

  it('validatePayment returns false when payer_name is empty', () => {
    const form = useTransactionForm('space-1')
    form.paymentForm.value.payer_name = ''
    form.paymentForm.value.payee_name = 'Alice'
    form.paymentForm.value.total_amount = 100

    expect(form.validatePayment()).toBe(false)
  })

  it('validatePayment returns false when total_amount is 0', () => {
    const form = useTransactionForm('space-1')
    form.paymentForm.value.payer_name = 'Bob'
    form.paymentForm.value.payee_name = 'Alice'
    form.paymentForm.value.total_amount = 0

    expect(form.validatePayment()).toBe(false)
  })

  it('validatePayment returns true when all fields are valid', () => {
    const form = useTransactionForm('space-1')
    form.paymentForm.value.payer_name = 'Bob'
    form.paymentForm.value.payee_name = 'Alice'
    form.paymentForm.value.total_amount = 500

    expect(form.validatePayment()).toBe(true)
  })
})
