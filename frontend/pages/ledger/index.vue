<template>
  <div class="ledger-page">
    <header class="flex justify-between items-center mb-5">
      <h1 class="text-2xl font-bold">我的帳本</h1>
    </header>

    <!-- FAB -->
    <NuxtLink 
        to="/ledger/add" 
        class="fixed bottom-24 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg shadow-indigo-500/30 rounded-full flex items-center justify-center text-white transition-transform active:scale-90 z-20 cursor-pointer border-0"
    >
        <Icon icon="mdi:plus" class="text-3xl" />
    </NuxtLink>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <div v-else-if="transactions.length === 0" class="text-center py-16 px-5 bg-neutral-900 rounded-2xl border border-neutral-800">
      <Icon icon="mdi:receipt-text-outline" class="text-6xl text-neutral-400 mb-4" />
      <p class="text-neutral-400 mb-5">還沒有任何記錄</p>
      <NuxtLink to="/ledger/add" class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors text-sm no-underline">開始記帳</NuxtLink>
    </div>

    <div v-else class="flex flex-col gap-3">
      <div
        v-for="txn in transactions"
        :key="txn.id"
        class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 cursor-pointer transition-all duration-200 hover:-translate-y-0.5 hover:border-indigo-500"
        @click="router.push(`/ledger/${txn.id}`)"
      >
        <div class="flex items-center gap-3">
          <div class="flex items-center justify-center w-11 h-11 rounded-xl bg-neutral-800">
            <Icon :icon="getCategoryIcon(txn.category)" class="text-xl text-indigo-500" />
          </div>
          <div class="flex-1">
            <h3 class="text-[15px] font-semibold mb-1">{{ txn.category || '未分類' }}</h3>
            <p class="text-xs text-neutral-400">{{ formatDate(txn.date) }} • {{ txn.items?.length || 0 }} 項目</p>
          </div>
          <div class="text-right">
            <span class="block text-xs text-neutral-400 mb-0.5">{{ txn.currency }}</span>
            <span class="text-lg font-semibold">{{ formatAmount(txn.total_amount) }}</span>
          </div>
        </div>
        <div v-if="txn.note" class="mt-3 pt-3 border-t border-neutral-800 text-[13px] text-neutral-400">{{ txn.note }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const transactions = ref<any[]>([])
const ledger = ref<any>(null)
const loading = ref(true)

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
  }
  return icons[category] || 'mdi:receipt'
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW', {
    month: 'numeric',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toLocaleString('zh-TW')
}

const fetchData = async () => {
  try {
    // Get or create default ledger
    const ledgers = await api.get<any[]>('/api/ledgers')

    if (ledgers.length === 0) {
      // Create default ledger
      ledger.value = await api.post<any>('/api/ledgers', {
        name: '我的帳本',
        type: 'personal',
        currencies: ['TWD'],
        categories: ['餐飲', '交通', '購物', '娛樂', '生活', '其他']
      })
    } else {
      ledger.value = ledgers[0]
    }

    // Fetch transactions
    const txns = await api.get<any[]>(`/api/ledgers/${ledger.value.id}/transactions`)
    transactions.value = txns
  } catch (e) {
    console.error('Failed to fetch data:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchData()
})
</script>
