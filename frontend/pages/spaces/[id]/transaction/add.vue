<template>
  <div class="add-transaction-page">
    <PageTitle
      title="新增交易"
      :back-to="`/spaces/${route.params.id}/ledger`"
      :breadcrumbs="[{ label: detailStore.space?.name || '空間', to: `/spaces/${route.params.id}/ledger` }]"
    />

    <div>
      <!-- Quick Scan Receipt -->
      <button
        type="button"
        @click="scanInputRef?.click()"
        class="w-full mb-6 p-4 bg-gradient-to-r from-indigo-500/10 to-violet-500/10 border border-indigo-500/30 rounded-2xl cursor-pointer flex items-center gap-3 hover:from-indigo-500/20 hover:to-violet-500/20 transition-all"
      >
        <div class="w-10 h-10 bg-indigo-500/20 rounded-xl flex items-center justify-center shrink-0">
          <Icon icon="mdi:camera" class="text-xl text-indigo-400" />
        </div>
        <div class="text-left">
          <div class="text-sm font-bold text-indigo-300">掃描收據</div>
          <div class="text-xs text-neutral-500 mt-0.5">拍照或上傳發票，AI 自動辨識明細</div>
        </div>
        <Icon icon="mdi:chevron-right" class="text-xl text-neutral-600 ml-auto" />
      </button>
      <input
        ref="scanInputRef"
        type="file"
        accept="image/*"
        multiple
        class="hidden"
        @change="handleScanReceipt"
      >

      <!-- Divider -->
      <div class="flex items-center gap-3 mb-6">
        <div class="flex-1 border-t border-neutral-800" />
        <span class="text-xs text-neutral-600">或手動輸入</span>
        <div class="flex-1 border-t border-neutral-800" />
      </div>

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

      <!-- Template Button (expense only) -->
      <div v-if="transactionType === 'expense'" class="flex justify-end mb-4">
        <button
          type="button"
          @click="showTemplatePicker = true; templateComposable.fetchTemplates()"
          class="text-xs font-bold text-indigo-400 bg-transparent border border-indigo-500/30 rounded-lg px-3 py-1.5 cursor-pointer hover:bg-indigo-500/10 transition-colors flex items-center gap-1"
        >
          <Icon icon="mdi:file-document-outline" class="text-sm" />
          套用模板
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
        :disabled="aiExtract"
        @update:debts="debts = $event"
        @submit="handleExpenseSubmit"
      >
        <template #images>
          <!-- AI extract toggle lives in the images slot so it sits right next
               to the upload area — which is the relevant interaction when the
               toggle is on (just drop photos in and submit). -->
          <div class="bg-neutral-900 rounded-xl p-4 border border-neutral-800 mb-3 flex flex-col gap-2">
            <label class="flex items-center gap-3 cursor-pointer select-none">
              <input
                type="checkbox"
                v-model="aiExtract"
                class="w-4 h-4 rounded text-indigo-500 bg-neutral-800 border-neutral-700"
              >
              <span class="flex items-center gap-1.5 text-sm font-bold text-indigo-300">
                <Icon icon="mdi:robot-outline" class="text-base" />
                使用 AI 辨識收據
              </span>
            </label>
            <p class="text-xs text-neutral-500 leading-relaxed pl-7">
              上傳的發票會傳送至 Google Gemini 進行辨識，不會用於模型訓練。系統會自動帶入金額、項目、日期等欄位。
            </p>
          </div>
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

    <ClientOnly>
      <TemplatePickerModal
        :visible="showTemplatePicker"
        :templates="templateComposable.templates.value"
        :loading="templateComposable.loading.value"
        @close="showTemplatePicker = false"
        @select="handleTemplateSelect"
        @delete="handleTemplateDelete"
      />
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useLoading } from '~/composables/useLoading'
import { useToast } from '~/composables/useToast'
import { useTransactionForm } from '~/composables/useTransactionForm'
import { useExpenseTemplates } from '~/composables/useExpenseTemplates'
import { useConfirm } from '~/composables/useConfirm'
import PageTitle from '~/components/PageTitle.vue'
import ImageManager from '~/components/ImageManager.vue'
import ExpenseForm from '~/components/ExpenseForm.vue'
import PaymentForm from '~/components/PaymentForm.vue'
import TemplatePickerModal from '~/components/TemplatePickerModal.vue'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import type { ExpenseTemplate } from '~/types'

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
const scanInputRef = ref<HTMLInputElement | null>(null)

const {
  transactionType, baseCurrency, categories, availableCurrencies,
  paymentMethods, memberOptions, debts, expenseForm, paymentForm, aiExtract,
  fetchSpaceConfig, populateFromTemplate, buildExpenseMultipart, buildPaymentPayload,
  validateExpense, validatePayment,
} = useTransactionForm(route.params.id as string)

const templateComposable = useExpenseTemplates(route.params.id as string)
const showTemplatePicker = ref(false)
const confirm = useConfirm()

const handleTemplateSelect = (template: ExpenseTemplate) => {
  populateFromTemplate(template.data)
  showTemplatePicker.value = false
}

const handleTemplateDelete = async (templateId: string) => {
  if (!await confirm({ message: '確定要刪除此模板嗎？', destructive: true })) return
  try {
    await templateComposable.deleteTemplate(templateId)
  } catch (e: any) {
    toast.error(e.message || '刪除失敗')
  }
}

const handleScanReceipt = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const files = Array.from(input.files || [])
  if (files.length === 0) return

  // Reset input so same file can be re-selected
  input.value = ''

  const now = new Date()
  const mm = String(now.getMonth() + 1).padStart(2, '0')
  const dd = String(now.getDate()).padStart(2, '0')
  const hh = String(now.getHours()).padStart(2, '0')
  const min = String(now.getMinutes()).padStart(2, '0')

  // Each image becomes its own transaction, uploaded in parallel
  aiExtract.value = true
  expenseForm.value.date = now
  expenseForm.value.total_amount = 0

  // Build all FormData payloads first (synchronous), then send in parallel
  const payloads = files.map((file, i) => {
    const suffix = files.length > 1 ? ` (${i + 1})` : ''
    expenseForm.value.title = `Receipt ${mm}/${dd} ${hh}:${min}${suffix}`
    return buildExpenseMultipart([file])
  })
  const uploads = payloads.map(fd =>
    api.upload(`/api/spaces/${route.params.id}/expenses`, fd)
  )

  showLoading()
  try {
    await Promise.all(uploads)
    detailStore.invalidate('transactions')
    const msg = files.length > 1
      ? `${files.length} 筆收據已送出，AI 辨識中...`
      : '收據已送出，AI 辨識中...'
    toast.success(msg)
    router.push(`/spaces/${route.params.id}/ledger`)
  } catch (e: any) {
    toast.error(e.message || '送出失敗')
  } finally {
    hideLoading()
  }
}

const handleExpenseSubmit = async () => {
  if (!validateExpense()) return

  // The ImageManager is in buffered mode (instant-upload=false), so files live
  // in its pending queue until we pull them out here. Submitting as multipart
  // lets the backend create the transaction and attach images atomically.
  const files = imageManagerRef.value?.getBufferedFiles() ?? []
  if (aiExtract.value && files.length === 0) {
    toast.error('使用 AI 辨識時必須上傳至少一張收據圖片')
    return
  }

  showLoading()
  try {
    await api.upload(`/api/spaces/${route.params.id}/expenses`, buildExpenseMultipart(files))
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

  // Pre-fill payment form from query params (e.g. from stats suggested transfers)
  if (route.query.type === 'payment') {
    transactionType.value = 'payment'
    if (route.query.payer) paymentForm.value.payer_name = route.query.payer as string
    if (route.query.payee) paymentForm.value.payee_name = route.query.payee as string
    if (route.query.amount) paymentForm.value.total_amount = Number(route.query.amount)
  }
})
</script>
