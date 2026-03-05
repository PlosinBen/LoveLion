<template>
  <div class="ledger-page pb-24">
    <header class="flex justify-between items-center mb-6">
      <div class="flex flex-col">
        <h1 class="text-2xl font-bold tracking-tight">個人記帳</h1>
        <LedgerSwitcher />
      </div>

      <div class="flex gap-2">
        <button v-if="isOwner" @click="router.push(`/ledger/${currentLedger?.id}/settings`)" class="w-10 h-10 rounded-xl bg-neutral-900 border border-neutral-800 flex items-center justify-center text-neutral-400 hover:text-white transition-colors cursor-pointer">
          <Icon icon="mdi:share-variant" class="text-xl" />
        </button>
      </div>
    </header>

    <!-- FAB -->
    <NuxtLink 
        v-if="currentLedger"
        to="/ledger/add" 
        class="fixed bottom-24 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg shadow-indigo-500/30 rounded-full flex items-center justify-center text-white transition-transform active:scale-90 z-20 cursor-pointer border-0"
    >
        <Icon icon="mdi:plus" class="text-3xl" />
    </NuxtLink>

    <div v-if="loadingTransactions" class="flex flex-col items-center justify-center py-20 gap-4">
      <div class="w-8 h-8 border-2 border-indigo-500/30 border-t-indigo-500 rounded-full animate-spin"></div>
      <div class="text-neutral-500 text-sm">載入交易中...</div>
    </div>

    <div v-else-if="transactions.length === 0" class="text-center py-20 px-6 bg-neutral-900 rounded-3xl border border-neutral-800/50 flex flex-col items-center">
      <div class="w-20 h-20 bg-neutral-800 rounded-full flex items-center justify-center mb-6">
        <Icon icon="mdi:receipt-text-outline" class="text-4xl text-neutral-600" />
      </div>
      <h3 class="text-lg font-bold mb-2">空空如也</h3>
      <p class="text-neutral-500 text-sm mb-8 leading-relaxed">這個帳本還沒有任何記錄，<br/>點擊下方按鈕開始記下第一筆支出吧！</p>
      <NuxtLink to="/ledger/add" class="px-6 py-3 rounded-2xl font-bold bg-indigo-500 text-white hover:bg-indigo-600 transition-all shadow-lg shadow-indigo-500/20 active:scale-95 no-underline">開始記帳</NuxtLink>
    </div>

    <div v-else class="flex flex-col gap-3">
      <div
        v-for="txn in transactions"
        :key="txn.id"
        class="bg-neutral-900 rounded-3xl p-5 border border-neutral-800/60 cursor-pointer transition-all duration-300 hover:border-indigo-500/50 hover:bg-neutral-800/30 group"
        @click="router.push(`/ledger/${txn.id}`)"
      >
        <div class="flex items-center gap-4">
          <div class="flex items-center justify-center w-12 h-12 rounded-2xl bg-neutral-800 group-hover:bg-indigo-500/10 transition-colors text-indigo-500">
            <Icon :icon="getCategoryIcon(txn.category)" class="text-2xl" />
          </div>
          <div class="flex-1 min-w-0">
            <h3 class="text-[15px] font-bold mb-0.5 truncate text-neutral-100">{{ txn.title || txn.category || '未分類' }}</h3>
            <p class="text-xs text-neutral-500 flex items-center gap-1.5 font-medium">
               <span>{{ formatDate(txn.date) }}</span>
               <span class="w-0.5 h-0.5 rounded-full bg-neutral-700"></span>
               <span>由 {{ txn.user?.display_name || '某人' }} 記錄</span>
            </p>
          </div>
          <div class="text-right flex flex-col gap-0.5">
            <span class="text-xs font-bold text-neutral-600">{{ (txn.billing_amount && Number(txn.billing_amount) > 0) ? (currentLedger?.base_currency || 'TWD') : txn.currency }}</span>
            <span class="text-lg font-black tracking-tight text-white">{{ formatAmount((txn.billing_amount && Number(txn.billing_amount) > 0) ? txn.billing_amount : txn.total_amount) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useLedger } from '~/composables/useLedger'
import LedgerSwitcher from '~/components/LedgerSwitcher.vue'

definePageMeta({
  layout: 'main'
})

const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth, user } = useAuth()
const { currentLedger, currentLedgerId, fetchLedgers } = useLedger()

const transactions = ref<any[]>([])
const loadingTransactions = ref(true)

const isOwner = computed(() => {
  return currentLedger.value && currentLedger.value.user_id === user.value?.id
})

// Watch for ledger selection changes to reload data
watch(currentLedgerId, () => {
  if (currentLedgerId.value) {
    fetchTransactions()
  }
})

const getCategoryIcon = (category: string) => {
  const icons: Record<string, string> = {
    'Food': 'mdi:food',
    'Transport': 'mdi:train-car',
    'Shopping': 'mdi:shopping',
    'Entertainment': 'mdi:movie',
    'Health': 'mdi:hospital',
    '餐飲': 'mdi:food',
    '交通': 'mdi:train-car',
    '購物': 'mdi:shopping',
    '生活': 'mdi:home-outline',
    '娛樂': 'mdi:movie-open-outline',
  }
  return icons[category] || 'mdi:receipt-outline'
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW', {
    month: 'numeric',
    day: 'numeric',
    weekday: 'short'
  })
}

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toLocaleString('zh-TW')
}

const fetchTransactions = async () => {
  if (!currentLedger.value) return
  
  loadingTransactions.value = true
  try {
    const txns = await api.get<any[]>(`/api/ledgers/${currentLedger.value.id}/transactions`)
    transactions.value = txns
  } catch (e) {
    console.error('Failed to fetch transactions:', e)
  } finally {
    loadingTransactions.value = false
  }
}

onMounted(async () => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  
  await fetchLedgers()
  fetchTransactions()
})
</script>
