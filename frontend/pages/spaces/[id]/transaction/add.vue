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
            付款
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
import { useToast } from '~/composables/useToast'
import { useTransactionForm } from '~/composables/useTransactionForm'
import PageTitle from '~/components/PageTitle.vue'
import ImageManager from '~/components/ImageManager.vue'
import ExpenseForm from '~/components/ExpenseForm.vue'
import PaymentForm from '~/components/PaymentForm.vue'
import { useSpaceDetailStore } from '~/stores/spaceDetail'

definePageMeta({
  path: '/spaces/:id/ledger/transaction/add',
  layout: 'default'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const detailStore = useSpaceDetailStore()
const { showLoading, hideLoading } = useLoading()
const toast = useToast()
const imageManagerRef = ref<InstanceType<typeof ImageManager> | null>(null)

const {
  transactionType, baseCurrency, categories, availableCurrencies,
  paymentMethods, memberOptions, debts, expenseForm, paymentForm,
  fetchSpaceConfig, buildExpensePayload, buildPaymentPayload,
  validateExpense, validatePayment,
} = useTransactionForm(route.params.id as string)

const handleExpenseSubmit = async () => {
  if (!validateExpense()) return

  showLoading()
  try {
    const created = await api.post<any>(`/api/spaces/${route.params.id}/expenses`, buildExpensePayload())
    if (imageManagerRef.value) {
      await imageManagerRef.value.commit(created.id)
    }
    detailStore.invalidate('transactions')
    router.push(`/spaces/${route.params.id}/ledger`)
  } catch (e: any) {
    toast.error(e.message || '儲存失敗')
  } finally {
    hideLoading()
  }
}

const handlePaymentSubmit = async () => {
  if (!validatePayment()) return

  showLoading()
  try {
    await api.post(`/api/spaces/${route.params.id}/payments`, buildPaymentPayload())
    detailStore.invalidate('transactions')
    router.push(`/spaces/${route.params.id}/ledger`)
  } catch (e: any) {
    toast.error(e.message || '儲存失敗')
  } finally {
    hideLoading()
  }
}

onMounted(async () => {
  detailStore.setSpaceId(route.params.id as string)
  detailStore.fetchSpace()
  await fetchSpaceConfig()
})
</script>
