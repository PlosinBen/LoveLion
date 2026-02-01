<template>
  <div class="transaction-detail min-h-screen bg-neutral-900 text-neutral-50 p-4">
    <header class="flex justify-between items-center mb-6">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">交易詳情</h1>
      <button @click="handleDelete" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-red-500 border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:delete" class="text-2xl" />
      </button>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <div v-else-if="transaction" class="content">
      <div class="text-center py-8 px-5 mb-6 bg-neutral-900 rounded-2xl border border-neutral-800">
        <div class="inline-flex items-center gap-2 px-4 py-2 bg-neutral-800 rounded-3xl mb-4">
          <Icon :icon="getCategoryIcon(transaction.category)" class="text-xl text-indigo-500" />
          <span>{{ transaction.category || '未分類' }}</span>
        </div>
        <div class="text-4xl font-bold mb-2">
          {{ transaction.currency }} {{ formatAmount(transaction.total_amount) }}
        </div>
        <div class="text-neutral-400">{{ formatDate(transaction.date) }}</div>
      </div>

      <div class="mb-6">
        <h2 class="text-base text-neutral-400 mb-3 font-normal">項目明細</h2>
        <div class="flex flex-col gap-2">
          <div v-for="item in transaction.items" :key="item.id" class="flex justify-between items-center p-4 bg-neutral-900 rounded-2xl border border-neutral-800">
            <div class="flex items-center gap-2">
              <span class="font-medium">{{ item.name }}</span>
              <span class="text-neutral-400 text-sm">× {{ item.quantity }}</span>
            </div>
            <div class="font-semibold">{{ transaction.currency }} {{ formatAmount(item.amount) }}</div>
          </div>
        </div>
      </div>

      <div v-if="transaction.note" class="mb-6">
        <h2 class="text-base text-neutral-400 mb-3 font-normal">備註</h2>
        <div class="p-4 bg-neutral-900 rounded-2xl border border-neutral-800 text-neutral-400 leading-relaxed">{{ transaction.note }}</div>
      </div>

      <div class="text-neutral-400 text-xs">
        <p class="mb-1">建立時間：{{ formatDateTime(transaction.created_at) }}</p>
        <p>ID：{{ transaction.id }}</p>
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
definePageMeta({
  layout: 'empty'
})
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const transaction = ref<any>(null)
const loading = ref(true)
const ledgerId = ref('')

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
    const ledgers = await api.get<any[]>('/api/ledgers')
    if (ledgers.length > 0) {
      ledgerId.value = ledgers[0].id
      transaction.value = await api.get<any>(
        `/api/ledgers/${ledgerId.value}/transactions/${route.params.id}`
      )
    }
  } catch (e) {
    console.error('Failed to fetch transaction:', e)
  } finally {
    loading.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('確定要刪除這筆交易嗎？')) return

  try {
    await api.del(`/api/ledgers/${ledgerId.value}/transactions/${route.params.id}`)
    router.push('/ledger')
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
  fetchData()
})
</script>
