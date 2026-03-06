<template>
  <div class="flex flex-col">
    <!-- Unified Space Header -->
    <div class="px-2 pt-0 pb-6 flex items-center gap-3">
      <button @click="router.push('/')" class="w-10 h-10 rounded-full bg-neutral-800 text-white flex items-center justify-center hover:bg-neutral-700 transition-colors border-0 cursor-pointer shrink-0">
          <Icon icon="mdi:arrow-left" class="text-xl" />
      </button>

      <div class="flex-1 min-w-0">
         <h1 class="text-xl font-bold text-white tracking-tight truncate">{{ space?.name || '載入中...' }}</h1>
         <div class="flex items-center gap-1.5 text-neutral-500 text-[10px] font-medium mt-0.5">
            <Icon :icon="spaceTheme.icon" :class="spaceTheme.color" />
            <span v-if="space?.type === 'trip'">{{ formatDateRange(space?.start_date, space?.end_date) }}</span>
            <span v-else>{{ spaceTheme.label }}</span>
         </div>
      </div>
      
      <button @click="router.push(`/spaces/${route.params.id}/settings`)" class="w-10 h-10 rounded-full bg-neutral-800 text-white flex items-center justify-center hover:bg-neutral-700 transition-colors border-0 cursor-pointer shrink-0">
          <Icon icon="mdi:cog-outline" class="text-xl" />
      </button>
    </div>

    <div class="flex flex-col gap-6">
      <div v-if="loading" class="flex flex-col items-center justify-center py-20 gap-4">
          <div class="w-8 h-8 border-2 border-indigo-500/30 border-t-indigo-500 rounded-full animate-spin"></div>
          <div class="text-neutral-500 text-sm">載入空間資料...</div>
      </div>

      <template v-else-if="space">
          <!-- Main Module Tabs -->
          <div class="flex items-center gap-1 bg-neutral-900 p-1 rounded-2xl border border-neutral-800/50">
              <button 
                  v-for="tab in availableTabs" 
                  :key="tab.id"
                  @click="activeTab = tab.id"
                  class="flex-1 flex items-center justify-center gap-2 py-3 rounded-xl text-sm font-bold transition-all"
                  :class="activeTab === tab.id ? 'bg-neutral-800 text-white shadow-sm' : 'text-neutral-500 hover:text-neutral-300'"
              >
                  <Icon :icon="tab.icon" class="text-lg" />
                  {{ tab.label }}
              </button>
          </div>

          <!-- Tab Contents -->
          
          <!-- 1. Ledger (Transactions) Tab -->
          <div v-if="activeTab === 'ledger'" class="flex flex-col gap-4 animate-in fade-in slide-in-from-bottom-2 duration-500 pb-24">
              <!-- FAB for Ledger -->
              <NuxtLink 
                  :to="`/spaces/add?ledger_id=${space.id}`" 
                  class="fixed bottom-24 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg shadow-indigo-500/30 rounded-full flex items-center justify-center text-white transition-transform active:scale-90 z-20 cursor-pointer border-0"
              >
                  <Icon icon="mdi:plus" class="text-3xl" />
              </NuxtLink>

              <div v-if="transactions.length === 0" class="text-center py-20 px-6 bg-neutral-900 rounded-3xl border border-neutral-800/50 flex flex-col items-center">
                <Icon icon="mdi:receipt-text-outline" class="text-4xl text-neutral-700 mb-4" />
                <h3 class="text-lg font-bold mb-1">尚無交易</h3>
                <p class="text-neutral-500 text-sm">點擊下方按鈕開始記錄第一筆支出</p>
              </div>

              <div v-else class="flex flex-col gap-3">
                <div
                  v-for="txn in transactions"
                  :key="txn.id"
                  class="bg-neutral-900 rounded-3xl p-5 border border-neutral-800/60 cursor-pointer transition-all duration-300 hover:border-indigo-500/50 hover:bg-neutral-800/30 group"
                  @click="router.push(`/spaces/${space.id}/transaction/${txn.id}`)"
                >
                  <div class="flex items-center gap-4">
                    <div class="flex items-center justify-center w-12 h-12 rounded-2xl bg-neutral-800 group-hover:bg-indigo-500/10 transition-colors text-indigo-500">
                      <Icon :icon="getCategoryIcon(txn.category)" class="text-2xl" />
                    </div>
                    <div class="flex-1 min-w-0">
                      <h3 class="text-[15px] font-bold mb-0.5 truncate text-neutral-100">{{ txn.title || txn.category || '未分類' }}</h3>
                      <p class="text-xs text-neutral-500 flex items-center gap-1.5 font-medium">
                         <span>{{ formatDateShort(txn.date) }}</span>
                         <span class="w-0.5 h-0.5 rounded-full bg-neutral-700"></span>
                         <span>{{ txn.payer }} 支付</span>
                      </p>
                    </div>
                    <div class="text-right flex flex-col gap-0.5">
                      <span class="text-xs font-bold text-neutral-600">{{ (txn.billing_amount && Number(txn.billing_amount) > 0) ? space.base_currency : txn.currency }}</span>
                      <span class="text-lg font-black tracking-tight text-white">{{ formatAmount((txn.billing_amount && Number(txn.billing_amount) > 0) ? txn.billing_amount : txn.total_amount) }}</span>
                    </div>
                  </div>
                </div>
              </div>
          </div>

          <!-- 2. Comparison Tab -->
          <div v-else-if="activeTab === 'comparison'" class="flex flex-col gap-6 animate-in fade-in slide-in-from-bottom-2 duration-500 pb-24">
              <!-- FAB for Comparison -->
              <button 
                  @click="showAddStore = true"
                  class="fixed bottom-24 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg shadow-indigo-500/30 rounded-full flex items-center justify-center text-white transition-transform active:scale-90 z-20 cursor-pointer border-0"
              >
                  <Icon icon="mdi:store-plus-outline" class="text-3xl" />
              </button>

              <div v-if="stores.length === 0" class="text-center py-20 px-6 bg-neutral-900 rounded-3xl border border-neutral-800/50 flex flex-col items-center">
                <Icon icon="mdi:scale-balance" class="text-4xl text-neutral-700 mb-4" />
                <h3 class="text-lg font-bold mb-1">尚無比價商店</h3>
                <p class="text-neutral-500 text-sm">在此記錄不同店家的商品價格</p>
                <button @click="showAddStore = true" class="mt-6 px-6 py-2 bg-indigo-500 text-white rounded-full font-bold text-sm border-0 cursor-pointer">新增商店</button>
              </div>
              
              <div v-else class="grid grid-cols-1 gap-4">
                  <div v-for="store in stores" :key="store.id" class="bg-neutral-900 rounded-3xl p-6 border border-neutral-800/60">
                      <div class="flex justify-between items-start mb-4">
                          <div>
                              <h3 class="font-bold text-lg text-white">{{ store.name }}</h3>
                              <p v-if="store.location" class="text-xs text-neutral-500">{{ store.location }}</p>
                          </div>
                          <div class="flex gap-2">
                              <button v-if="store.google_map_url" @click="openMap(store.google_map_url)" class="w-8 h-8 rounded-full bg-neutral-800 flex items-center justify-center text-indigo-400 border-0 cursor-pointer hover:bg-neutral-700">
                                  <Icon icon="mdi:map-marker-outline" />
                              </button>
                              <button @click="openStoreDetails(store.id)" class="w-8 h-8 rounded-full bg-neutral-800 flex items-center justify-center text-neutral-400 border-0 cursor-pointer hover:bg-neutral-700">
                                  <Icon icon="mdi:chevron-right" />
                              </button>
                          </div>
                      </div>
                      
                      <!-- Product Quick List -->
                      <div class="flex flex-col gap-2">
                          <div v-for="p in store.products.slice(0, 3)" :key="p.id" class="flex justify-between text-sm py-1 border-b border-neutral-800/30 last:border-0">
                              <span class="text-neutral-400">{{ p.name }}</span>
                              <span class="text-white font-medium">{{ p.currency }} {{ formatAmount(p.price) }}</span>
                          </div>
                          <p v-if="store.products.length > 3" class="text-[10px] text-center text-neutral-600 font-bold uppercase tracking-widest mt-1">還有 {{ store.products.length - 3 }} 個項目</p>
                      </div>
                  </div>
                  
                  <button @click="showAddStore = true" class="p-5 rounded-3xl border-2 border-dashed border-neutral-800 text-neutral-500 font-bold hover:border-indigo-500/50 hover:text-indigo-400 transition-all bg-transparent cursor-pointer">
                      + 新增商店
                  </button>
              </div>
          </div>

          <!-- 3. Stats Tab -->
          <div v-else-if="activeTab === 'stats'" class="animate-in fade-in slide-in-from-bottom-2 duration-500 pb-24">
              <TripStats 
                  :transactions="transactions" 
                  :members="space.members" 
                  :base-currency="space.base_currency"
              />
          </div>

      </template>
    </div>

    <!-- Modals (Add Store, etc.) -->
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
        <div v-if="showAddStore" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-4 bg-black/60 backdrop-blur-sm">
            <div class="w-full max-w-lg bg-neutral-900 rounded-t-[32px] sm:rounded-[32px] p-8 border border-neutral-800 animate-in slide-in-from-bottom duration-300">
                <div class="flex justify-between items-center mb-6">
                    <h2 class="text-xl font-bold">新增商店</h2>
                    <button @click="showAddStore = false" class="w-10 h-10 rounded-full bg-neutral-800 flex items-center justify-center text-neutral-400 border-0 cursor-pointer">
                        <Icon icon="mdi:close" class="text-xl" />
                    </button>
                </div>

                <div class="flex flex-col gap-6">
                    <BaseInput 
                        v-model="storeForm.name" 
                        label="商店名稱" 
                        placeholder="例如：唐吉訶德、全家" 
                        required
                    />
                    <BaseInput 
                        v-model="storeForm.location" 
                        label="分店或地點" 
                        placeholder="例如：澀谷店、新宿" 
                    />
                    <BaseInput 
                        v-model="storeForm.google_map_url" 
                        label="Google Map 連結" 
                        placeholder="選填" 
                    />

                    <button 
                        @click="handleCreateStore" 
                        :disabled="submittingStore || !storeForm.name"
                        class="w-full py-4 rounded-2xl bg-indigo-500 text-white font-bold hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 mt-2"
                    >
                        {{ submittingStore ? '建立中...' : '建立商店' }}
                    </button>
                </div>
            </div>
        </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Icon } from '@iconify/vue'
import TripStats from '~/components/TripStats.vue'
import BaseInput from '~/components/BaseInput.vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

definePageMeta({
  layout: 'main'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const space = ref<any>(null)
const transactions = ref<any[]>([])
const stores = ref<any[]>([])
const loading = ref(true)
const activeTab = ref('ledger')

// Add Store Modal Logic
const showAddStore = ref(false)
const submittingStore = ref(false)
const storeForm = ref({
    name: '',
    location: '',
    google_map_url: ''
})

const handleCreateStore = async () => {
    if (!storeForm.value.name.trim()) return
    submittingStore.value = true
    try {
        const newStore = await api.post<any>(`/api/spaces/${route.params.id}/stores`, storeForm.value)
        stores.value.push(newStore)
        showAddStore.value = false
        storeForm.value = { name: '', location: '', google_map_url: '' }
    } catch (e: any) {
        alert(e.message || '建立商店失敗')
    } finally {
        submittingStore.value = false
    }
}

// Atmosphere Engine: Styles based on Space Type
const spaceTheme = computed(() => {
    const type = space.value?.type || 'personal'
    const themes: Record<string, { icon: string, color: string, bg: string, label: string }> = {
        'trip': { icon: 'mdi:airplane', color: 'text-blue-400', bg: 'bg-blue-500/10', label: '旅遊計畫' },
        'personal': { icon: 'mdi:wallet-outline', color: 'text-green-400', bg: 'bg-green-500/10', label: '個人空間' },
        'group': { icon: 'mdi:account-group-outline', color: 'text-indigo-400', bg: 'bg-indigo-500/10', label: '團隊空間' }
    }
    return themes[type] || themes.personal
})

const availableTabs = [
    { id: 'ledger', label: '明細', icon: 'mdi:receipt-text-outline' },
    { id: 'stats', label: '統計', icon: 'mdi:chart-pie-outline' },
    { id: 'comparison', label: '比價', icon: 'mdi:scale-balance' }
]

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toLocaleString('zh-TW')
}

const formatDateShort = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW', { month: 'numeric', day: 'numeric', weekday: 'short' })
}

const formatDateRange = (start: string | null, end: string | null) => {
  if (!start && !end) return '日期未設定'
  const opts: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric' }
  const startStr = start ? new Date(start).toLocaleDateString('zh-TW', opts) : '?'
  const endStr = end ? new Date(end).toLocaleDateString('zh-TW', opts) : '?'
  return `${startStr} - ${endStr}`
}

const getCategoryIcon = (category: string) => {
  const icons: Record<string, string> = {
    '餐飲': 'mdi:food', '交通': 'mdi:train-car', '購物': 'mdi:shopping', '生活': 'mdi:home-outline', '娛樂': 'mdi:movie-open-outline',
  }
  return icons[category] || 'mdi:receipt-outline'
}

const fetchData = async () => {
  try {
    // Fetch Space Details
    space.value = await api.get<any>(`/api/spaces/${route.params.id}`)
    
    // Fetch Transactions
    const txns = await api.get<any[]>(`/api/spaces/${route.params.id}/transactions`)
    transactions.value = txns

    // Fetch Comparison Stores
    const storesData = await api.get<any[]>(`/api/spaces/${route.params.id}/stores`)
    stores.value = storesData
  } catch (e) {
    console.error('Failed to fetch space data:', e)
    router.push('/')
  } finally {
    loading.value = false
  }
}

const openMap = (url: string) => window.open(url, '_blank')
const openStoreDetails = (storeId: string) => {
    // Navigate to store details (to be implemented/refactored)
    router.push(`/spaces/${route.params.id}/stores/${storeId}`)
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
