import { ref } from 'vue'
import { useApi } from '~/composables/useApi'
import { useToast } from '~/composables/useToast'
import { parseNaiveDate } from '~/utils/date'

function toLocalISOString(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}Z`
}
import type { ExpenseFormData } from '~/components/ExpenseForm.vue'
import type { PaymentFormData } from '~/components/PaymentForm.vue'
import type { DebtItem } from '~/components/DebtEditor.vue'
import type { Transaction, TransactionType, ExpenseTemplateData } from '~/types'

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
  // When true, the worker will fill in items/date for us, so we skip those
  // fields from validation and let the backend pick them up from the LLM.
  const aiExtract = ref(false)
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
    location_url: '',
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
        date: parseNaiveDate(txn.date),
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
        location_url: txn.expense.location_url || '',
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
        date: parseNaiveDate(txn.date),
        title: txn.title,
        payer_name: debt?.payer_name || '',
        payee_name: debt?.payee_name || '',
        total_amount: Number(txn.total_amount),
        note: txn.note,
      }
    }
  }

  // Populate forms from template data
  const populateFromTemplate = (data: ExpenseTemplateData) => {
    transactionType.value = 'expense'
    expenseForm.value = {
      date: new Date(),
      title: data.title,
      total_amount: Number(data.total_amount),
      category: data.category,
      payment_method: data.payment_method,
      note: data.note,
      items: data.items && data.items.length > 0
        ? data.items.map(i => ({
            name: i.name,
            unit_price: Number(i.unit_price),
            quantity: Number(i.quantity),
            discount: Number(i.discount),
          }))
        : [{ name: '', unit_price: 0, quantity: 1, discount: 0 }],
      currency: data.currency || baseCurrency.value,
      manual_rate: true,
      exchange_rate: 0,
      billing_amount: 0,
      handling_fee: 0,
      location_url: data.location_url || '',
    }
    if (data.debts && data.debts.length > 0) {
      debts.value = data.debts.map(d => ({
        payer_name: d.payer_name,
        payee_name: d.payee_name,
        amount: Number(d.amount),
        is_spot_paid: d.is_spot_paid,
      }))
    }
  }

  // Build expense API payload
  const buildExpensePayload = () => {
    const form = expenseForm.value
    const payload: Record<string, any> = {
      title: form.title,
      total_amount: Number(form.total_amount),
      date: toLocalISOString(form.date),
      currency: form.currency,
      note: form.note,
      expense: {
        category: form.category,
        exchange_rate: form.currency === baseCurrency.value ? 1 : Number(form.exchange_rate),
        billing_amount: form.currency === baseCurrency.value ? Number(form.total_amount) : Number(form.billing_amount),
        handling_fee: form.currency === baseCurrency.value ? 0 : Number(form.handling_fee),
        payment_method: form.payment_method,
        location_url: form.location_url || '',
        items: form.items
          .filter(item => item.name)
          .map(item => ({
            name: item.name,
            unit_price: Number(item.unit_price),
            quantity: Number(item.quantity),
            discount: Number(item.discount) || 0,
          })),
      },
      debts: debts.value.length > 0
        ? debts.value
            .filter(d => d.payer_name && d.payee_name && d.amount > 0)
            .map(d => ({ payer_name: d.payer_name, payee_name: d.payee_name, amount: Number(d.amount), is_spot_paid: d.is_spot_paid }))
        : undefined,
    }
    if (aiExtract.value) {
      payload.ai_extract = true
    }
    return payload
  }

  // Build a multipart/form-data body the backend Create handler accepts:
  //   data   — JSON string matching buildExpensePayload()
  //   images — N files attached under the same field name
  const buildExpenseMultipart = (files: File[]): FormData => {
    const fd = new FormData()
    fd.append('data', JSON.stringify(buildExpensePayload()))
    for (const f of files) {
      fd.append('images', f, f.name)
    }
    return fd
  }

  // Build payment API payload
  const buildPaymentPayload = () => {
    const form = paymentForm.value
    return {
      date: toLocalISOString(form.date),
      title: form.title,
      note: form.note,
      total_amount: Number(form.total_amount),
      payer_name: form.payer_name,
      payee_name: form.payee_name,
    }
  }

  // Validate expense form, returns true if valid.
  // When ai_extract is on we trust the worker to fill in total/items, so the
  // only requirement is that the caller attached at least one image (checked
  // by the caller — this composable has no visibility into the ImageManager).
  const validateExpense = (): boolean => {
    if (aiExtract.value) {
      return true
    }
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
    aiExtract,
    fetchSpaceConfig,
    populateFromTransaction,
    populateFromTemplate,
    buildExpensePayload,
    buildExpenseMultipart,
    buildPaymentPayload,
    validateExpense,
    validatePayment,
  }
}
