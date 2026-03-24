<template>
  <div class="add-transaction-page">
    <PageTitle
      title="新增交易"
      :back-to="`/spaces/${route.params.id}/ledger`"
      :breadcrumbs="[{ label: detailStore.space?.name || '空間', to: `/spaces/${route.params.id}/ledger` }]"
    />

    <div>
      <!-- Type Selector -->
      <div class="flex bg-neutral-900 rounded-xl p-1 border border-neutral-800 mb-6">
        <button
          type="button"
          @click="transactionType = 'expense'"
          class="flex-1 py-2.5 text-sm font-bold rounded-lg border-0 cursor-pointer transition-all"
          :class="transactionType === 'expense' ? 'bg-indigo-500 text-white' : 'bg-transparent text-neutral-500 hover:text-neutral-300'"
        >
          消費
        </button>
        <button
          type="button"
          @click="transactionType = 'payment'"
          class="flex-1 py-2.5 text-sm font-bold rounded-lg border-0 cursor-pointer transition-all"
          :class="transactionType === 'payment' ? 'bg-indigo-500 text-white' : 'bg-transparent text-neutral-500 hover:text-neutral-300'"
        >
          補款
        </button>
      </div>

      <!-- Expense Form -->
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
            entity-id="pending"
            entity-type="transaction"
            :instant-upload="false"
            :instant-delete="false"
            ref="imageManagerRef"
          />
        </template>
      </ExpenseForm>

      <!-- Payment Form -->
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
import { useLoading } from '~/composables/useLoading'
import PageTitle from '~/components/PageTitle.vue'
import ImageManager from '~/components/ImageManager.vue'
import ExpenseForm from '~/components/ExpenseForm.vue'
import PaymentForm from '~/components/PaymentForm.vue'
import type { ExpenseFormData } from '~/components/ExpenseForm.vue'
import type { PaymentFormData } from '~/components/PaymentForm.vue'
import type { DebtItem } from '~/components/DebtEditor.vue'
import type { TransactionType } from '~/types'
import { useSpaceDetailStore } from '~/stores/spaceDetail'

definePageMeta({
  path: '/spaces/:id/ledger/transaction/add',
  layout: 'default'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const detailStore = useSpaceDetailStore()

const transactionType = ref<TransactionType>('expense')
const categories = ref<{label: string, value: string}[]>([])
const availableCurrencies = ref<{label: string, value: string}[]>([])
const paymentMethods = ref<{label: string, value: string}[]>([])
const baseCurrency = ref('TWD')
const memberOptions = ref<string[]>([])
const { showLoading, hideLoading } = useLoading()
const imageManagerRef = ref<InstanceType<typeof ImageManager> | null>(null)
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

const handleExpenseSubmit = async () => {
  const form = expenseForm.value

  if (form.total_amount <= 0) {
    alert('請填寫總額')
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

    const created = await api.post<any>(`/api/spaces/${route.params.id}/expenses`, payload)
    if (imageManagerRef.value) {
      await imageManagerRef.value.commit(created.id)
    }
    detailStore.invalidate('transactions')
    router.push(`/spaces/${route.params.id}/ledger`)
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    hideLoading()
  }
}

const handlePaymentSubmit = async () => {
  const form = paymentForm.value
  if (!form.payer_name || !form.payee_name) {
    alert('請填寫付款人與收款人')
    return
  }
  if (form.total_amount <= 0) {
    alert('請填寫金額')
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

    await api.post(`/api/spaces/${route.params.id}/payments`, payload)
    detailStore.invalidate('transactions')
    router.push(`/spaces/${route.params.id}/ledger`)
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    hideLoading()
  }
}

onMounted(async () => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }

  detailStore.setSpaceId(route.params.id as string)
  detailStore.fetchSpace()

  try {
    const space = await api.get<any>(`/api/spaces/${route.params.id}`)
    baseCurrency.value = space.base_currency || 'TWD'
    expenseForm.value.currency = baseCurrency.value
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

    if (categories.value.length > 0 && !expenseForm.value.category) {
      const firstCat = categories.value[0]
      if (firstCat) expenseForm.value.category = firstCat.value
    }
  } catch (e) {
    console.error('Failed to fetch space data', e)
  }
})
</script>
