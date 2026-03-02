<template>
  <div class="transaction-detail min-h-screen bg-neutral-900 text-neutral-50 p-4">
    <header class="flex justify-between items-center mb-6">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">交易詳情</h1>
      <div class="flex gap-2">
        <button @click="router.push(`/trips/${route.params.id}/ledger/${route.params.txnId}/edit`)" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-indigo-400 border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
          <Icon icon="mdi:pencil" class="text-2xl" />
        </button>
        <button @click="handleDelete" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-red-500 border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
          <Icon icon="mdi:delete" class="text-2xl" />
        </button>
      </div>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <div v-else-if="transaction" class="content flex flex-col gap-6">
      <!-- Amount Card -->
      <div class="text-center py-8 px-5 bg-neutral-900 rounded-2xl border border-neutral-800">
        <div class="inline-flex items-center gap-2 px-4 py-2 bg-neutral-800 rounded-3xl mb-4">
          <Icon :icon="getCategoryIcon(transaction.category)" class="text-xl text-indigo-500" />
          <span>{{ transaction.category || '未分類' }}</span>
        </div>
        <h2 class="text-lg font-bold mb-1">{{ transaction.title || '無標題' }}</h2>
        <div class="text-4xl font-bold mb-2">
          {{ transaction.currency }} {{ formatAmount(transaction.total_amount) }}
        </div>
        <div class="text-neutral-400">{{ formatDate(transaction.date) }}</div>

        <!-- Foreign Currency Details -->
        <div v-if="transaction.currency !== baseCurrency" class="mt-6 pt-6 border-t border-neutral-800/50 flex flex-col gap-2 text-sm">
            <div class="flex justify-between text-neutral-400">
                <span>匯率</span>
                <span>{{ transaction.exchange_rate }}</span>
            </div>
            <div class="flex justify-between text-neutral-400">
                <span>折合{{ baseCurrency }}</span>
                <span class="text-white font-bold">{{ baseCurrency }} {{ formatAmount(transaction.billing_amount) }}</span>
            </div>
            <div v-if="parseFloat(transaction.handling_fee) > 0" class="flex justify-between text-neutral-400">
                <span>手續費</span>
                <span>{{ baseCurrency }} {{ formatAmount(transaction.handling_fee) }}</span>
            </div>
        </div>
      </div>

      <!-- Split Details (Unique to Trips) -->
      <div v-if="transaction.splits && transaction.splits.length > 0">
        <h2 class="text-base text-neutral-400 mb-3 font-normal">分帳與付款</h2>
        <div class="flex flex-col gap-2">
          <!-- Payer -->
          <div class="flex justify-between items-center p-4 bg-neutral-800/50 rounded-2xl border border-neutral-800">
             <div class="flex items-center gap-2">
                <Icon icon="mdi:wallet-outline" class="text-indigo-400" />
                <span class="text-sm">由 <span class="font-bold text-white">{{ transaction.payer }}</span> 付款</span>
             </div>
             <span class="font-mono text-sm">{{ transaction.currency }} {{ formatAmount(transaction.total_amount) }}</span>
          </div>
          
          <!-- Consumers -->
          <div class="flex flex-col gap-1 mt-1">
             <div v-for="split in consumers" :key="split.id" class="flex justify-between items-center px-4 py-3 bg-neutral-900 rounded-xl border border-neutral-800">
                <span class="text-sm font-medium">{{ split.name }}</span>
                <span class="font-mono text-sm text-neutral-400">{{ transaction.currency }} {{ formatAmount(split.amount) }}</span>
             </div>
          </div>
        </div>
      </div>

      <!-- Items -->
      <div v-if="transaction.items && transaction.items.length > 0">
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

      <!-- Images -->
      <div v-if="transaction.images && transaction.images.length > 0">
        <h2 class="text-base text-neutral-400 mb-3 font-normal">照片</h2>
        <div class="grid grid-cols-2 gap-2">
           <div 
            v-for="img in transaction.images" 
            :key="img.id" 
            class="aspect-square rounded-xl bg-neutral-800 overflow-hidden cursor-pointer"
            @click="previewImage(img)"
           >
              <img :src="getImageUrl(img.file_path)" class="w-full h-full object-cover" alt="Transaction image" />
           </div>
        </div>
      </div>

      <!-- Note -->
      <div v-if="transaction.note" class="mb-6">
        <h2 class="text-base text-neutral-400 mb-3 font-normal">備註</h2>
        <div class="p-4 bg-neutral-900 rounded-2xl border border-neutral-800 text-neutral-400 leading-relaxed">{{ transaction.note }}</div>
      </div>

      <div class="text-neutral-400 text-xs pb-10">
        <p class="mb-1">建立時間：{{ formatDateTime(transaction.created_at) }}</p>
        <p>ID：{{ transaction.id }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useImages } from '~/composables/useImages'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const { getImageUrl } = useImages()

definePageMeta({
  layout: 'empty'
})

const trip = ref<any>(null)
const transaction = ref<any>(null)
const loading = ref(true)

const baseCurrency = computed(() => trip.value?.ledger?.base_currency || trip.value?.base_currency || 'TWD')

const consumers = computed(() => {
  if (!transaction.value?.splits) return []
  return transaction.value.splits.filter((s: any) => !s.is_payer)
})

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

const previewImage = (img: any) => {
  router.push({
    path: '/image-preview',
    query: {
      id: transaction.value.id,
      type: 'transaction'
    },
    state: {
      images: transaction.value.images.map((i: any) => ({
        url: getImageUrl(i.file_path),
        hash: i.blur_hash
      }))
    }
  })
}

const fetchData = async () => {
  try {
    // 1. Fetch Trip
    const tripData = await api.get<any>(`/api/trips/${route.params.id}`)
    trip.value = tripData
    
    // 2. Fetch Transaction
    const txn = await api.get<any>(`/api/ledgers/${tripData.ledger_id}/transactions/${route.params.txnId}`)
    transaction.value = txn
  } catch (e) {
    console.error('Failed to fetch data:', e)
    router.push(`/trips/${route.params.id}/ledger`)
  } finally {
    loading.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('確定要刪除這筆交易嗎？')) return

  try {
    await api.del(`/api/ledgers/${trip.value.ledger_id}/transactions/${route.params.txnId}`)
    router.push(`/trips/${route.params.id}/ledger`)
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
