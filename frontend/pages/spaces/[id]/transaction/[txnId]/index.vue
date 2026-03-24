<template>
  <div class="transaction-detail">
    <PageTitle
      title="交易詳情"
      :back-to="`/spaces/${route.params.id}/ledger`"
      :breadcrumbs="[{ label: '我的空間', to: '/' }, { label: detailStore.space?.name || '空間', to: `/spaces/${route.params.id}/stats` }, { label: '記帳', to: `/spaces/${route.params.id}/ledger` }]"
    >
      <template #right>
        <div class="flex gap-1">
          <NuxtLink :to="`/spaces/${route.params.id}/ledger/transaction/${route.params.txnId}/edit`" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-indigo-400 hover:bg-neutral-800 transition-colors active:scale-95 no-underline">
            <Icon icon="mdi:pencil-outline" class="text-xl" />
          </NuxtLink>
          <button @click="handleDelete" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-red-500 border-0 cursor-pointer hover:bg-neutral-800 transition-colors active:scale-95">
            <Icon icon="mdi:trash-can-outline" class="text-xl" />
          </button>
        </div>
      </template>
    </PageTitle>

    <div class="content-wrapper">
      <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

      <div v-else-if="transaction" class="content">
        <!-- Header Card -->
        <BaseCard padding="py-8 px-5" class="text-center mb-6">
          <!-- Type badge -->
          <div v-if="transaction.type === 'payment'" class="inline-flex items-center gap-2 px-4 py-2 bg-emerald-500/10 rounded-2xl mb-4">
            <Icon icon="mdi:swap-horizontal" class="text-xl text-emerald-400" />
            <span class="font-bold text-sm text-emerald-400">補款</span>
          </div>
          <div v-else-if="transaction.expense" class="inline-flex items-center gap-2 px-4 py-2 bg-neutral-800 rounded-2xl mb-4">
            <Icon :icon="getCategoryIcon(transaction.expense.category)" class="text-xl text-indigo-500" />
            <span class="font-bold text-sm">{{ transaction.expense.category || '未分類' }}</span>
          </div>

          <h2 v-if="transaction.title" class="text-lg font-bold text-white mb-1">{{ transaction.title }}</h2>

          <div class="text-4xl font-bold mb-2 tracking-tight">
            <span class="text-neutral-500 text-lg mr-1 font-bold">{{ transaction.currency }}</span>
            {{ formatAmount(transaction.total_amount) }}
          </div>
          <div class="text-neutral-400 text-sm font-medium">{{ formatDate(transaction.date) }}</div>

          <!-- Payment Method -->
          <div v-if="transaction.expense?.payment_method" class="mt-4 text-sm text-neutral-500">
            {{ transaction.expense.payment_method }}
          </div>

          <!-- Foreign Currency Details (expense only) -->
          <div v-if="transaction.expense && transaction.currency !== detailStore.space?.base_currency" class="mt-6 pt-6 border-t border-neutral-800/50 flex flex-col gap-3 text-sm">
            <div class="flex justify-between items-center text-neutral-400">
              <span>匯率</span>
              <span class="font-mono bg-neutral-800 px-2 py-0.5 rounded text-neutral-200">{{ transaction.expense.exchange_rate }}</span>
            </div>
            <div class="flex justify-between items-center text-neutral-400">
              <span>折合 {{ detailStore.space?.base_currency }}</span>
              <span class="text-white font-bold text-base">{{ detailStore.space?.base_currency }} {{ formatAmount(transaction.expense.billing_amount) }}</span>
            </div>
            <div v-if="parseFloat(transaction.expense.handling_fee) > 0" class="flex justify-between items-center text-neutral-400">
              <span>手續費</span>
              <span class="text-neutral-200">{{ detailStore.space?.base_currency }} {{ formatAmount(transaction.expense.handling_fee) }}</span>
            </div>
          </div>
        </BaseCard>

        <!-- Images -->
        <div v-if="transaction.images && transaction.images.length > 0" class="mb-6">
          <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest mb-3 px-1">收據 / 照片</h2>
          <div class="grid grid-cols-3 gap-2">
            <div
              v-for="(image, idx) in transaction.images"
              :key="image.id"
              class="aspect-square rounded-xl overflow-hidden bg-neutral-800 border border-neutral-700 cursor-pointer"
              @click="router.push(`/image-preview?id=${transaction.id}&type=transaction&index=${idx}`)"
            >
              <img :src="image.file_path" class="w-full h-full object-cover" />
            </div>
          </div>
        </div>

        <!-- Expense Items -->
        <div v-if="transaction.expense?.items && transaction.expense.items.length > 0" class="mb-6">
          <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest mb-3 px-1">項目明細</h2>
          <div class="bg-neutral-900 rounded-xl border border-neutral-800/60 divide-y divide-neutral-800">
            <div v-for="item in transaction.expense.items" :key="item.id" class="flex justify-between items-center px-5 py-4">
              <span class="font-bold text-neutral-100 text-sm">{{ item.name }}</span>
              <div class="flex items-center gap-2">
                <span class="text-xs text-neutral-500 font-medium">{{ formatAmount(item.unit_price) }}<template v-if="Number(item.discount) > 0"> - {{ formatAmount(item.discount) }}</template> × {{ item.quantity }}</span>
                <span class="font-bold text-white text-sm">{{ formatAmount(item.amount) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Debts -->
        <div v-if="transaction.debts && transaction.debts.length > 0" class="mb-6">
          <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest mb-3 px-1">
            {{ transaction.type === 'payment' ? '補款詳情' : '分帳' }}
          </h2>
          <div class="bg-neutral-900 rounded-xl border border-neutral-800/60 divide-y divide-neutral-800">
            <div v-for="debt in transaction.debts" :key="debt.id" class="flex justify-between items-center px-5 py-4">
              <div class="flex items-center gap-2">
                <span class="font-bold text-neutral-100 text-sm">{{ debt.payer_name }}</span>
                <span class="text-neutral-600 text-xs">→</span>
                <span class="font-bold text-neutral-100 text-sm">{{ debt.payee_name }}</span>
                <span v-if="debt.is_spot_paid" class="text-xs text-emerald-400 font-bold">已結清</span>
              </div>
              <div class="text-right">
                <div class="font-bold text-white text-sm">{{ transaction.currency }} {{ formatAmount(debt.amount) }}</div>
                <div v-if="Number(debt.settled_amount) > 0 && Number(debt.settled_amount) !== Number(debt.amount)" class="text-xs text-neutral-500">
                  ≈ {{ detailStore.space?.base_currency }} {{ formatAmount(debt.settled_amount) }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="transaction.note" class="mb-6">
          <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest mb-3 px-1">備註</h2>
          <div class="p-4 bg-neutral-900 rounded-xl border border-neutral-800/60 text-neutral-400 text-sm leading-relaxed">{{ transaction.note }}</div>
        </div>

        <div class="px-1 text-neutral-600 text-xs font-medium space-y-1">
          <p>建立於 {{ formatDateTime(transaction.created_at) }}</p>
          <p class="font-mono">TXN: {{ transaction.id }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'
import BaseCard from '~/components/BaseCard.vue'
import type { Transaction } from '~/types'

definePageMeta({
  path: '/spaces/:id/ledger/transaction/:txnId'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const detailStore = useSpaceDetailStore()

const transaction = ref<Transaction | null>(null)
const loading = ref(true)

const getCategoryIcon = (category: string) => {
  const icons: Record<string, string> = {
    '餐飲': 'mdi:food',
    '交通': 'mdi:train-car',
    '購物': 'mdi:shopping',
    '娛樂': 'mdi:movie',
    '生活': 'mdi:home',
  }
  return icons[category] || 'mdi:receipt'
}

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toLocaleString('zh-TW')
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'short',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

const formatDateTime = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW', { hour12: false })
}

const fetchData = async () => {
  try {
    transaction.value = await api.get<Transaction>(
      `/api/spaces/${route.params.id}/transactions/${route.params.txnId}`
    )
  } catch (e) {
    console.error('Failed to fetch transaction:', e)
    router.push(`/spaces/${route.params.id}/ledger`)
  } finally {
    loading.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('確定要刪除此交易嗎？')) return

  try {
    await api.del(`/api/spaces/${route.params.id}/transactions/${route.params.txnId}`)
    detailStore.invalidate('transactions')
    router.push(`/spaces/${route.params.id}/ledger`)
  } catch (e: any) {
    alert(e.message || '刪除失敗')
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  detailStore.setSpaceId(route.params.id as string)
  detailStore.fetchSpace()
  fetchData()
})
</script>
