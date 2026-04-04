import { ref } from 'vue'
import { useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import type { ExpenseFormData } from '~/components/ExpenseForm.vue'
import type { PaymentFormData } from '~/components/PaymentForm.vue'
import type { DebtItem } from '~/components/DebtEditor.vue'
import type { Transaction, TransactionType } from '~/types'

const DEFAULT_CATEGORIES = [
  { label: '餐飲', value: '餐飲' },
  { label: '交通', value: '交通' },
  { label: '購物', value: '購物' },
  { label: '娛樂', value: '娛樂' },
  { label: '生活', value: '生活' },
  { label: '其他', value: '其他' },
]

const DEFAULT_CURRENCIES = [
  { label: 'TWD', value: 'TWD' },
  { label: 'JPY', value: 'JPY' },
  { label: 'USD', value: 'USD' },
]

function toOptions(arr: string[]) {
  return arr.map(v => ({ label: v, value: v }))
}

export function useTransactionForm(spaceId: string) {
  const api = useApi()
  const toast = useToast()

  const transactionType = ref<TransactionType>('expense')
  const baseCurrency = ref('TWD')
  const categories = ref<{ label: string; value: string }[]>([])
  const availableCurrencies = ref<{ label: string; value: string }[]>([])
  const paymentMethods = ref<{ label: string; value: string }[]>([])
  const memberOptions = ref<string[]>([])
  const debts = ref<DebtItem[]>([])

  const expenseForm = ref<ExpenseFormData>({
    date: new Date(),
    title: '',
    total_amount: 0,
    category: '其他',
    payment_method: '',
    note: '',
    items: [{ name: '', unit_price: 0, quantity: 1, discount: 0 }],
    currency: 'TWD',
    manual_rate: true,
    exchange_rate: 0,
    billing_amount: 0,
    handling_fee: 0,
  })

  const paymentForm = ref<PaymentFormData>({
    date: new Date(),
    title: '',
    payer_name: '',
    payee_name: '',
    total_amount: 0,
    note: '',
  })

  // Load space config (categories, currencies, payment_methods, members)
  const loadSpaceConfig = (space: any) => {
    baseCurrency.value = space.base_currency || 'TWD'
    memberOptions.value = space.split_members || []
    categories.value = space.categories ? toOptions(space.categories) : DEFAULT_CATEGORIES
    availableCurrencies.value = space.currencies ? toOptions(space.currencies) : DEFAULT_CURRENCIES
    paymentMethods.value = space.payment_methods ? toOptions(space.payment_methods) : []
  }

  // Fetch space and apply config
  const fetchSpaceConfig = async () => {
    try {
      const space = await api.get<any>(`/api/spaces/${spaceId}`)
      loadSpaceConfig(space)
      expenseForm.value.currency = baseCurrency.value
      if (categories.value.length > 0 && !expenseForm.value.category) {
        const first = categories.value[0]
        if (first) expenseForm.value.category = first.value
      }
      return space
    } catch (e) {
      console.error('Failed to fetch space data', e)
      return null
    }
  }

  // Populate forms from existing transaction (edit mode)
  const populateFromTransaction = (txn: Transaction) => {
    transactionType.value = txn.type

    if (txn.type === 'expense' && txn.expense) {
      const itemsData = txn.expense.items || []
      expenseForm.value = {
        date: new Date(txn.date),
        title: txn.title,
        total_amount: Number(txn.total_amount),
        category: txn.expense.category,
        payment_method: txn.expense.payment_method,
        note: txn.note,
        items: itemsData.length > 0
          ? itemsData.map(i => ({
              name: i.name,
              unit_price: Number(i.unit_price),
              quantity: Number(i.quantity),
              discount: Number(i.discount),
            }))
          : [{ name: '', unit_price: 0, quantity: 1, discount: 0 }],
        currency: txn.currency,
        manual_rate: Number(txn.expense.handling_fee) === 0 || !txn.expense.handling_fee,
        exchange_rate: Number(txn.expense.exchange_rate),
        billing_amount: Number(txn.expense.billing_amount),
        handling_fee: Number(txn.expense.handling_fee) || 0,
      }

      if (txn.debts && txn.debts.length > 0) {
        debts.value = txn.debts.map(d => ({
          payer_name: d.payer_name,
          payee_name: d.payee_name,
          amount: Number(d.amount),
          is_spot_paid: d.is_spot_paid,
        }))
      }
    } else if (txn.type === 'payment') {
      const debt = txn.debts?.[0]
      paymentForm.value = {
        date: new Date(txn.date),
        title: txn.title,
        payer_name: debt?.payer_name || '',
        payee_name: debt?.payee_name || '',
        total_amount: Number(txn.total_amount),
        note: txn.note,
      }
    }
  }

  // Build expense API payload
  const buildExpensePayload = () => {
    const form = expenseForm.value
    return {
      title: form.title,
      total_amount: form.total_amount,
      date: form.date.toISOString(),
      currency: form.currency,
      note: form.note,
      expense: {
        category: form.category,
        exchange_rate: form.currency === baseCurrency.value ? 1 : form.exchange_rate,
        billing_amount: form.currency === baseCurrency.value ? form.total_amount : form.billing_amount,
        handling_fee: form.currency === baseCurrency.value ? 0 : form.handling_fee,
        payment_method: form.payment_method,
        items: form.items
          .filter(item => item.name && Number(item.unit_price) > 0)
          .map(item => ({
            name: item.name,
            unit_price: item.unit_price,
            quantity: item.quantity,
            discount: item.discount || 0,
          })),
      },
      debts: debts.value.length > 0
        ? debts.value
            .filter(d => d.payer_name && d.payee_name && d.amount > 0)
            .map(d => ({ payer_name: d.payer_name, payee_name: d.payee_name, amount: d.amount, is_spot_paid: d.is_spot_paid }))
        : undefined,
    }
  }

  // Build payment API payload
  const buildPaymentPayload = () => {
    const form = paymentForm.value
    return {
      date: form.date.toISOString(),
      title: form.title,
      note: form.note,
      total_amount: form.total_amount,
      payer_name: form.payer_name,
      payee_name: form.payee_name,
    }
  }

  // Validate expense form, returns true if valid
  const validateExpense = (): boolean => {
    if (expenseForm.value.total_amount <= 0) {
      toast.error('請填寫總額')
      return false
    }
    return true
  }

  // Validate payment form, returns true if valid
  const validatePayment = (): boolean => {
    const form = paymentForm.value
    if (!form.payer_name || !form.payee_name) {
      toast.error('請填寫付款人與收款人')
      return false
    }
    if (form.total_amount <= 0) {
      toast.error('請填寫金額')
      return false
    }
    return true
  }

  return {
    transactionType,
    baseCurrency,
    categories,
    availableCurrencies,
    paymentMethods,
    memberOptions,
    debts,
    expenseForm,
    paymentForm,
    fetchSpaceConfig,
    populateFromTransaction,
    buildExpensePayload,
    buildPaymentPayload,
    validateExpense,
    validatePayment,
  }
}
