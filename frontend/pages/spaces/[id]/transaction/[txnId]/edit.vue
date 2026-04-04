<template>
  <div class="edit-transaction-page">
    <PageTitle
      title="編輯交易"
      :back-to="`/spaces/${route.params.id}/ledger/transaction/${route.params.txnId}`"
      :breadcrumbs="[{ label: '我的空間', to: '/' }, { label: detailStore.space?.name || '空間', to: `/spaces/${route.params.id}/stats` }, { label: transaction?.title || '交易詳情', to: `/spaces/${route.params.id}/ledger/transaction/${route.params.txnId}` }]"
    />

    <div v-if="loading" class="p-20 flex justify-center items-center text-neutral-500">
      <Icon icon="mdi:loading" class="text-4xl animate-spin" />
    </div>

    <div v-else>
      <!-- Expense Edit -->
      <ExpenseForm
        v-if="transactionType === 'expense'"
        v-model:form="expenseForm"
        :base-currency="baseCurrency"
        :categories="categories"
        :available-currencies="availableCurrencies"
        :payment-methods="paymentMethods"
        :debts="debts"
        :member-options="memberOptions"
        @update:debts="debts = $event"
        @submit="handleExpenseSubmit"
      >
        <template #images>
          <ImageManager
            :entity-id="route.params.txnId as string"
            entity-type="transaction"
          />
        </template>
      </ExpenseForm>

      <!-- Payment Edit -->
      <PaymentForm
        v-else
        v-model:form="paymentForm"
        :base-currency="baseCurrency"
        :member-options="memberOptions"
        @submit="handlePaymentSubmit"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useLoading } from '~/composables/useLoading'
import { useToast } from '~/composables/useToast'
import PageTitle from '~/components/PageTitle.vue'
import ImageManager from '~/components/ImageManager.vue'
import ExpenseForm from '~/components/ExpenseForm.vue'
import PaymentForm from '~/components/PaymentForm.vue'
import type { ExpenseFormData } from '~/components/ExpenseForm.vue'
import type { PaymentFormData } from '~/components/PaymentForm.vue'
import type { DebtItem } from '~/components/DebtEditor.vue'
import type { Transaction, TransactionType } from '~/types'
import { useSpaceDetailStore } from '~/stores/spaceDetail'

definePageMeta({
  path: '/spaces/:id/ledger/transaction/:txnId/edit',
  layout: 'default'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const detailStore = useSpaceDetailStore()

const transaction = ref<Transaction | null>(null)
const transactionType = ref<TransactionType>('expense')
const categories = ref<{label: string, value: string}[]>([])
const availableCurrencies = ref<{label: string, value: string}[]>([])
const paymentMethods = ref<{label: string, value: string}[]>([])
const { showLoading, hideLoading } = useLoading()
const toast = useToast()
const loading = ref(true)
const baseCurrency = ref('TWD')
const memberOptions = ref<string[]>([])
const debts = ref<DebtItem[]>([])

const expenseForm = ref<ExpenseFormData>({
  date: new Date(),
  title: '',
  total_amount: 0,
  category: '',
  payment_method: '',
  note: '',
  items: [{ name: '', unit_price: 0, quantity: 1, discount: 0 }],
  currency: 'TWD',
  manual_rate: true,
  exchange_rate: 0,
  billing_amount: 0,
  handling_fee: 0
})

const paymentForm = ref<PaymentFormData>({
  date: new Date(),
  title: '',
  payer_name: '',
  payee_name: '',
  total_amount: 0,
  note: ''
})

const fetchData = async () => {
  try {
    const [txn, space] = await Promise.all([
      api.get<Transaction>(`/api/spaces/${route.params.id}/transactions/${route.params.txnId}`),
      api.get<any>(`/api/spaces/${route.params.id}`),
    ])

    transaction.value = txn
    transactionType.value = txn.type
    baseCurrency.value = space.base_currency || 'TWD'
    memberOptions.value = space.split_members || []

    if (space.categories) {
      categories.value = space.categories.map((c: string) => ({ label: c, value: c }))
    } else {
      categories.value = [
        { label: '餐飲', value: '餐飲' },
        { label: '交通', value: '交通' },
        { label: '購物', value: '購物' },
        { label: '娛樂', value: '娛樂' },
        { label: '生活', value: '生活' },
        { label: '其他', value: '其他' }
      ]
    }

    if (space.currencies) {
      availableCurrencies.value = space.currencies.map((c: string) => ({ label: c, value: c }))
    } else {
      availableCurrencies.value = [
        { label: 'TWD', value: 'TWD' },
        { label: 'JPY', value: 'JPY' },
        { label: 'USD', value: 'USD' }
      ]
    }

    if (space.payment_methods) {
      paymentMethods.value = space.payment_methods.map((m: string) => ({ label: m, value: m }))
    }

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
              discount: Number(i.discount)
            }))
          : [{ name: '', unit_price: 0, quantity: 1, discount: 0 }],
        currency: txn.currency,
        manual_rate: Number(txn.expense.handling_fee) === 0 || !txn.expense.handling_fee,
        exchange_rate: Number(txn.expense.exchange_rate),
        billing_amount: Number(txn.expense.billing_amount),
        handling_fee: Number(txn.expense.handling_fee) || 0
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
        note: txn.note
      }
    }
  } catch (e) {
    console.error('Failed to fetch transaction:', e)
  } finally {
    loading.value = false
  }
}

const handleExpenseSubmit = async () => {
  const form = expenseForm.value

  if (form.total_amount <= 0) {
    toast.error('請填寫總額')
    return
  }

  showLoading()
  try {
    const payload = {
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

    await api.put(`/api/spaces/${route.params.id}/expenses/${route.params.txnId}`, payload)
    detailStore.invalidate('transactions')
    router.push(`/spaces/${route.params.id}/ledger/transaction/${route.params.txnId}`)
  } catch (e: any) {
    toast.error(e.message || '更新失敗')
  } finally {
    hideLoading()
  }
}

const handlePaymentSubmit = async () => {
  const form = paymentForm.value
  if (!form.payer_name || !form.payee_name) {
    toast.error('請填寫付款人與收款人')
    return
  }

  showLoading()
  try {
    const payload = {
      date: form.date.toISOString(),
      title: form.title,
      note: form.note,
      total_amount: form.total_amount,
      payer_name: form.payer_name,
      payee_name: form.payee_name,
    }

    await api.put(`/api/spaces/${route.params.id}/payments/${route.params.txnId}`, payload)
    detailStore.invalidate('transactions')
    router.push(`/spaces/${route.params.id}/ledger/transaction/${route.params.txnId}`)
  } catch (e: any) {
    toast.error(e.message || '更新失敗')
  } finally {
    hideLoading()
  }
}


onMounted(() => {
  detailStore.setSpaceId(route.params.id as string)
  detailStore.fetchSpace()
  fetchData()
})
</script>
