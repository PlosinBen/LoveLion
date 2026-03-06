<template>
  <div class="transaction-detail min-h-screen bg-neutral-900 text-neutral-50">
    <SpaceHeader 
      title="交易詳情" 
      class="pt-0 px-2"
    >
      <template #right>
        <div class="flex gap-1">
          <button @click="router.push(`/spaces/${route.params.id}/transaction/${route.params.txnId}/edit`)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-indigo-400 border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
            <Icon icon="mdi:pencil-outline" class="text-xl" />
          </button>
          <button @click="handleDelete" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-red-500 border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
            <Icon icon="mdi:trash-can-outline" class="text-xl" />
          </button>
        </div>
      </template>
    </SpaceHeader>

    <div class="p-4 pt-0">
      <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

      <div v-else-if="transaction" class="content">
        <div class="text-center py-8 px-5 mb-6 bg-neutral-900 rounded-3xl border border-neutral-800">
          <div class="inline-flex items-center gap-2 px-4 py-2 bg-neutral-800 rounded-3xl mb-4">
            <Icon :icon="getCategoryIcon(transaction.category)" class="text-xl text-indigo-500" />
            <span class="font-bold text-sm">{{ transaction.category || '未分類' }}</span>
          </div>
          <div class="text-4xl font-black mb-2 tracking-tighter">
            <span class="text-neutral-500 text-lg mr-1 font-bold">{{ transaction.currency }}</span>
            {{ formatAmount(transaction.total_amount) }}
          </div>
          <div class="text-neutral-400 text-sm font-medium">{{ formatDate(transaction.date) }}</div>

          <!-- Foreign Currency Details -->
          <div v-if="transaction.currency !== 'TWD'" class="mt-6 pt-6 border-t border-neutral-800/50 flex flex-col gap-3 text-sm">
              <div class="flex justify-between items-center text-neutral-400">
                  <span>匯率</span>
                  <span class="font-mono bg-neutral-800 px-2 py-0.5 rounded text-neutral-200">{{ transaction.exchange_rate }}</span>
              </div>
              <div class="flex justify-between items-center text-neutral-400">
                  <span>折合台幣</span>
                  <span class="text-white font-bold text-base">TWD {{ formatAmount(transaction.billing_amount) }}</span>
              </div>
              <div v-if="parseFloat(transaction.handling_fee) > 0" class="flex justify-between items-center text-neutral-400">
                  <span>手續費</span>
                  <span class="text-neutral-200">TWD {{ formatAmount(transaction.handling_fee) }}</span>
              </div>
          </div>
        </div>

        <div class="mb-6">
          <h2 class="text-xs font-black text-neutral-500 uppercase tracking-widest mb-3 px-1">項目明細</h2>
          <div class="flex flex-col gap-2">
            <div v-for="item in transaction.items" :key="item.id" class="flex justify-between items-center p-4 bg-neutral-900 rounded-2xl border border-neutral-800/60">
              <div class="flex flex-col gap-0.5">
                <span class="font-bold text-neutral-100">{{ item.name }}</span>
                <span class="text-neutral-500 text-xs font-medium">單價 {{ transaction.currency }} {{ formatAmount(item.unit_price) }} × {{ item.quantity }}</span>
              </div>
              <div class="font-black text-white">{{ transaction.currency }} {{ formatAmount(item.amount) }}</div>
            </div>
          </div>
        </div>

        <div v-if="transaction.note" class="mb-6">
          <h2 class="text-xs font-black text-neutral-500 uppercase tracking-widest mb-3 px-1">備註</h2>
          <div class="p-4 bg-neutral-900 rounded-2xl border border-neutral-800/60 text-neutral-400 text-sm leading-relaxed">{{ transaction.note }}</div>
        </div>

        <div class="px-1 text-neutral-600 text-[10px] font-medium space-y-1">
          <p>建立於 {{ formatDateTime(transaction.created_at) }}</p>
          <p class="font-mono">TXN: {{ transaction.id }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  hideGlobalNav: true
})
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import SpaceHeader from '~/components/SpaceHeader.vue'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const transaction = ref<any>(null)
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
    transaction.value = await api.get<any>(
      `/api/spaces/${route.params.id}/transactions/${route.params.txnId}`
    )
  } catch (e) {
    console.error('Failed to fetch transaction:', e)
    router.push(`/spaces/${route.params.id}`)
  } finally {
    loading.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('確定要刪除這筆交易嗎？')) return

  try {
    await api.del(`/api/spaces/${route.params.id}/transactions/${route.params.txnId}`)
    router.push(`/spaces/${route.params.id}`)
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
