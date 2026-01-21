<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.push(`/trips/${route.params.id}`)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold truncate">旅行帳本</h1>
      <div v-if="trip" class="flex gap-2">
        <NuxtLink :to="`/trips/${trip.id}/ledger/add`" class="flex items-center justify-center w-10 h-10 rounded-xl bg-indigo-500 text-white hover:bg-indigo-600 transition-colors no-underline">
          <Icon icon="mdi:plus" class="text-2xl" />
        </NuxtLink>
      </div>
      <div v-else class="w-10"></div>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <template v-else-if="trip">
      <!-- Summary Card -->
      <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
        <h2 class="text-sm text-neutral-400 mb-2">總支出</h2>
        <div class="text-3xl font-bold text-white mb-1">
          <span class="text-lg text-neutral-500 mr-1">{{ trip.base_currency }}</span>
          {{ formatAmount(totalExpense) }}
        </div>
      </div>

      <!-- Transactions List -->
      <div v-if="transactions.length === 0" class="text-center py-10 px-5">
        <p class="text-neutral-400 mb-4">還沒有任何記帳</p>
        <NuxtLink :to="`/trips/${trip.id}/ledger/add`" class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg bg-neutral-800 text-white hover:bg-neutral-700 transition-colors text-sm no-underline">
          <Icon icon="mdi:plus" /> 新增第一筆
        </NuxtLink>
      </div>

      <div v-else class="flex flex-col gap-3">
        <div
          v-for="txn in transactions"
          :key="txn.id"
          class="bg-neutral-900 rounded-2xl p-4 border border-neutral-800 cursor-pointer transition-all duration-200 hover:-translate-y-0.5 hover:border-indigo-500"
          @click="router.push(`/trips/${trip.id}/ledger/${txn.id}`)"
        >
          <div class="flex items-center gap-3">
            <div class="flex items-center justify-center w-10 h-10 rounded-xl bg-neutral-800">
              <Icon :icon="getCategoryIcon(txn.category)" class="text-xl text-indigo-500" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex justify-between items-start mb-0.5">
                <h3 class="font-semibold truncate pr-2">{{ txn.items?.[0]?.name || txn.category || '未分類' }}</h3>
                <span class="font-bold whitespace-nowrap">{{ formatAmount(txn.total_amount) }}</span>
              </div>
              <div class="flex justify-between items-center text-xs text-neutral-400">
                <span>{{ formatDate(txn.date) }} • {{ txn.payer }} 付款</span>
                <span>{{ txn.currency }}</span>
              </div>
            </div>
          </div>
          <!-- Splits Preview -->
          <div v-if="txn.splits?.length" class="mt-3 pt-2 border-t border-neutral-800 flex items-center gap-2 overflow-x-auto pb-1 scrollbar-hide">
            <span class="text-xs text-neutral-500 whitespace-nowrap">分攤成員:</span>
            <div v-for="split in txn.splits" :key="split.id" class="flex items-center gap-1 bg-neutral-800 px-2 py-0.5 rounded text-xs text-neutral-300 whitespace-nowrap">
              <span>{{ getMemberName(split.member_id) }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const trip = ref<any>(null)
const transactions = ref<any[]>([])
const loading = ref(true)

const totalExpense = computed(() => {
  return transactions.value.reduce((sum, txn) => {
    // Simplified total calculation (assuming base currency for now or ignoring exchange rates for simple sum)
    // To do it properly, we should convert everything to base currency.
    // For now, let's just sum up ignoring currency mixing or assume specific logic.
    // Actually better to just show raw number if multiple currencies is complex.
    // Or just sum everything as number.
    return sum + Number(txn.total_amount)
  }, 0)
})

const getMemberName = (memberId: string) => {
  const member = trip.value?.members?.find((m: any) => m.id === memberId)
  return member ? member.name : '未知'
}

const getCategoryIcon = (category: string) => {
  const icons: Record<string, string> = {
    'Food': 'mdi:food',
    'Transport': 'mdi:train-car',
    'Shopping': 'mdi:shopping',
    'Entertainment': 'mdi:movie',
    'Health': 'mdi:hospital',
    'Accomodation': 'mdi:bed',
    '餐飲': 'mdi:food',
    '交通': 'mdi:train-car',
    '購物': 'mdi:shopping',
    '娛樂': 'mdi:movie',
    '住宿': 'mdi:bed',
    '生活': 'mdi:home',
  }
  return icons[category] || 'mdi:receipt'
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    month: 'short',
    day: 'numeric'
  })
}

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toLocaleString('zh-TW')
}

const fetchData = async () => {
  try {
    // 1. Fetch Trip info to get ledger_id
    trip.value = await api.get<any>(`/api/trips/${route.params.id}`)
    
    if (trip.value.ledger_id) {
      // 2. Fetch transactions
      transactions.value = await api.get<any[]>(`/api/ledgers/${trip.value.ledger_id}/transactions`)
    }
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
