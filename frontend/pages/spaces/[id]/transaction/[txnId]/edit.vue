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
      <!-- AI processing banner — covers pending and processing identically.
           While the worker owns the row the whole form is locked; the only
           escape is to cancel, which flips ai_status back to NULL. -->
      <div
        v-if="aiStatus === 'pending' || aiStatus === 'processing'"
        class="bg-indigo-500/10 border border-indigo-500/30 rounded-xl p-4 mb-4 flex items-center gap-3"
      >
        <AiStatusBadge :status="aiStatus" size="md" />
        <div class="flex-1 min-w-0">
          <div class="text-sm font-bold text-indigo-300">
            {{ aiStatus === 'processing' ? 'AI 辨識中...' : 'AI 等待辨識中...' }}
          </div>
          <div class="text-xs text-neutral-400 mt-0.5">
            背景 worker 正在處理，完成後會自動填入日期與項目
          </div>
        </div>
        <button
          type="button"
          @click="handleCancelAI"
          :disabled="cancelling"
          class="shrink-0 text-xs font-bold text-red-400 bg-red-500/10 border border-red-500/30 rounded-lg px-3 py-2 cursor-pointer hover:bg-red-500/20 disabled:opacity-50 transition-colors"
        >
          取消辨識
        </button>
      </div>

      <!-- AI failure banner — pure info, no action button.
           The user recovers by editing the form and either keeping AI
           off (becomes a normal transaction) or ticking the retry box
           (backend bumps ai_status back to pending on save). -->
      <div
        v-else-if="aiStatus === 'failed'"
        class="bg-red-500/10 border border-red-500/30 rounded-xl p-4 mb-4 flex items-start gap-3"
      >
        <AiStatusBadge status="failed" size="md" />
        <div class="flex-1 min-w-0">
          <div class="text-sm font-bold text-red-300">AI 辨識失敗</div>
          <div class="text-xs text-neutral-400 mt-0.5 break-words">
            {{ transaction?.ai_error || '未知錯誤' }}
          </div>
        </div>
      </div>

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
        :disabled="formDisabled"
        @update:debts="debts = $event"
        @submit="handleExpenseSubmit"
      >
        <template #images>
          <!-- AI retry checkbox only appears in the failed state. -->
          <div
            v-if="aiStatus === 'failed'"
            class="bg-neutral-900 rounded-xl p-4 border border-neutral-800 mb-3 flex flex-col gap-2"
          >
            <label class="flex items-center gap-3 cursor-pointer select-none">
              <input
                type="checkbox"
                v-model="aiExtract"
                class="w-4 h-4 rounded text-indigo-500 bg-neutral-800 border-neutral-700"
              >
              <span class="flex items-center gap-1.5 text-sm font-bold text-indigo-300">
                <Icon icon="mdi:robot-outline" class="text-base" />
                使用 AI 重新辨識
              </span>
            </label>
            <p class="text-xs text-neutral-500 leading-relaxed pl-7">
              勾選後儲存，worker 會使用現有的發票圖片重新嘗試辨識，先前填入的項目會被清除。
            </p>
          </div>
          <ImageManager
            :entity-id="route.params.txnId as string"
            entity-type="transaction"
          />
        </template>

        <!-- Hide the save button while the worker owns the row — nothing can
             be saved anyway, and the cancel button in the banner is the only
             meaningful action. -->
        <template v-if="formDisabled" #actions>
          <div></div>
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
import { computed, ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useLoading } from '~/composables/useLoading'
import { useToast } from '~/composables/useToast'
import { useTransactionForm } from '~/composables/useTransactionForm'
import PageTitle from '~/components/PageTitle.vue'
import ImageManager from '~/components/ImageManager.vue'
import ExpenseForm from '~/components/ExpenseForm.vue'
import PaymentForm from '~/components/PaymentForm.vue'
import AiStatusBadge from '~/components/AiStatusBadge.vue'
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
const cancelling = ref(false)

const {
  transactionType, baseCurrency, categories, availableCurrencies,
  paymentMethods, memberOptions, debts, expenseForm, paymentForm, aiExtract,
  fetchSpaceConfig, populateFromTransaction,
  buildExpensePayload, buildPaymentPayload,
  validateExpense, validatePayment,
} = useTransactionForm(route.params.id as string)

const aiStatus = computed(() => transaction.value?.ai_status ?? null)
const formDisabled = computed(
  () => aiStatus.value === 'pending' || aiStatus.value === 'processing',
)

const fetchData = async () => {
  try {
    const [txn] = await Promise.all([
      api.get<Transaction>(`/api/spaces/${route.params.id}/transactions/${route.params.txnId}`),
      fetchSpaceConfig(),
    ])
    transaction.value = txn
    populateFromTransaction(txn)
    // Reset the retry checkbox on every (re)load so stale state from a
    // previous cancel/save round doesn't leak into the next submission.
    aiExtract.value = false
  } catch (e) {
    console.error('Failed to fetch transaction:', e)
  } finally {
    loading.value = false
  }
}

const handleCancelAI = async () => {
  cancelling.value = true
  try {
    await api.post(
      `/api/spaces/${route.params.id}/transactions/${route.params.txnId}/ai-cancel`,
      {},
    )
    toast.success('已取消 AI 辨識')
    // Refetch so the banner disappears and the form becomes editable.
    await fetchData()
  } catch (e: any) {
    toast.error(e.message || '取消失敗')
  } finally {
    cancelling.value = false
  }
}

const handleExpenseSubmit = async () => {
  if (!validateExpense()) return

  showLoading()
  try {
    await api.put(`/api/spaces/${route.params.id}/expenses/${route.params.txnId}`, buildExpensePayload())
    detailStore.invalidate('transactions')
    router.replace(`/spaces/${route.params.id}/ledger/transaction/${route.params.txnId}`)
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
    router.replace(`/spaces/${route.params.id}/ledger/transaction/${route.params.txnId}`)
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
