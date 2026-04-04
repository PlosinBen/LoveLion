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
import { useTransactionForm } from '~/composables/useTransactionForm'
import PageTitle from '~/components/PageTitle.vue'
import ImageManager from '~/components/ImageManager.vue'
import ExpenseForm from '~/components/ExpenseForm.vue'
import PaymentForm from '~/components/PaymentForm.vue'
import type { Transaction } from '~/types'
import { useSpaceDetailStore } from '~/stores/spaceDetail'

definePageMeta({
  path: '/spaces/:id/ledger/transaction/:txnId/edit',
  layout: 'default'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const detailStore = useSpaceDetailStore()
const { showLoading, hideLoading } = useLoading()
const toast = useToast()

const transaction = ref<Transaction | null>(null)
const loading = ref(true)

const {
  transactionType, baseCurrency, categories, availableCurrencies,
  paymentMethods, memberOptions, debts, expenseForm, paymentForm,
  fetchSpaceConfig, populateFromTransaction,
  buildExpensePayload, buildPaymentPayload,
  validateExpense, validatePayment,
} = useTransactionForm(route.params.id as string)

const fetchData = async () => {
  try {
    const [txn] = await Promise.all([
      api.get<Transaction>(`/api/spaces/${route.params.id}/transactions/${route.params.txnId}`),
      fetchSpaceConfig(),
    ])
    transaction.value = txn
    populateFromTransaction(txn)
  } catch (e) {
    console.error('Failed to fetch transaction:', e)
  } finally {
    loading.value = false
  }
}

const handleExpenseSubmit = async () => {
  if (!validateExpense()) return

  showLoading()
  try {
    await api.put(`/api/spaces/${route.params.id}/expenses/${route.params.txnId}`, buildExpensePayload())
    detailStore.invalidate('transactions')
    router.push(`/spaces/${route.params.id}/ledger/transaction/${route.params.txnId}`)
  } catch (e: any) {
    toast.error(e.message || '更新失敗')
  } finally {
    hideLoading()
  }
}

const handlePaymentSubmit = async () => {
  if (!validatePayment()) return

  showLoading()
  try {
    await api.put(`/api/spaces/${route.params.id}/payments/${route.params.txnId}`, buildPaymentPayload())
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
