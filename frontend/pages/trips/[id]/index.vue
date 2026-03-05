<template>
  <div class="flex flex-col">
    <!-- Single Row Header -->
    <div class="px-2 pt-0 pb-6 flex items-center gap-3">
      <button @click="router.push('/trips')" class="w-10 h-10 rounded-full bg-neutral-800 text-white flex items-center justify-center hover:bg-neutral-700 transition-colors border-0 cursor-pointer shrink-0">
          <Icon icon="mdi:arrow-left" class="text-xl" />
      </button>

      <div class="flex-1 min-w-0">
         <h1 class="text-xl font-bold text-white tracking-tight truncate">{{ trip?.name || '載入中...' }}</h1>
         <div class="flex items-center gap-1.5 text-neutral-500 text-[10px] font-medium mt-0.5">
            <Icon icon="mdi:calendar-range" class="text-indigo-500" />
            <span>{{ formatDateRange(trip?.start_date, trip?.end_date) }}</span>
         </div>
      </div>
      
      <div class="relative shrink-0">
        <button @click="showMenu = !showMenu" class="w-10 h-10 rounded-full bg-neutral-800 text-white flex items-center justify-center hover:bg-neutral-700 transition-colors border-0 cursor-pointer">
            <Icon icon="mdi:dots-vertical" class="text-xl" />
        </button>
        
        <Transition
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="transform scale-95 opacity-0 -translate-y-2"
          enter-to-class="transform scale-100 opacity-100 translate-y-0"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="transform scale-100 opacity-100 translate-y-0"
          leave-to-class="transform scale-95 opacity-0 -translate-y-2"
        >
          <div v-if="showMenu" v-click-outside="() => showMenu = false" class="absolute top-12 right-0 bg-neutral-900 border border-neutral-800 rounded-2xl shadow-2xl z-50 overflow-hidden py-1 min-w-[140px]">
              <NuxtLink :to="`/trips/${trip?.id}/edit`" class="flex items-center gap-2 px-4 py-3 text-white hover:bg-neutral-800 transition-colors no-underline text-sm font-medium">
                <Icon icon="mdi:pencil-outline" /> 編輯旅行
              </NuxtLink>
              <div class="border-t border-neutral-800 my-1"></div>
              <button @click="handleDelete" class="w-full flex items-center gap-2 px-4 py-3 text-red-500 hover:bg-neutral-800 transition-colors border-0 bg-transparent cursor-pointer text-sm font-medium">
                <Icon icon="mdi:trash-can-outline" /> 刪除旅行
              </button>
          </div>
        </Transition>
      </div>
    </div>

    <div class="flex flex-col gap-6">
      <div v-if="loading" class="flex flex-col items-center justify-center py-20 gap-4">
          <div class="w-8 h-8 border-2 border-indigo-500/30 border-t-indigo-500 rounded-full animate-spin"></div>
          <div class="text-neutral-500 text-sm">載入中...</div>
      </div>

      <template v-else-if="trip">
          
          <!-- Tabs -->
          <div class="flex items-center gap-1 bg-neutral-900 p-1 rounded-2xl border border-neutral-800/50">
              <button 
                  v-for="tab in tabs" 
                  :key="tab.id"
                  @click="activeTab = tab.id"
                  class="flex-1 flex items-center justify-center gap-2 py-3 rounded-xl text-sm font-bold transition-all"
                  :class="activeTab === tab.id ? 'bg-neutral-800 text-white shadow-sm' : 'text-neutral-500 hover:text-neutral-300'"
              >
                  <Icon :icon="tab.icon" class="text-lg" />
                  {{ tab.label }}
              </button>
          </div>

          <!-- Tab Content -->
          
          <!-- Overview Tab -->
          <div v-if="activeTab === 'overview'" class="flex flex-col gap-6 animate-in fade-in slide-in-from-bottom-2 duration-500 pb-24">
              
              <!-- KPI Cards -->
              <div class="grid grid-cols-2 gap-3">
                  <div class="bg-neutral-900 rounded-3xl p-5 border border-neutral-800/60 flex flex-col justify-center relative overflow-hidden">
                       <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/5 to-transparent z-0"></div>
                       <span class="text-neutral-500 text-[10px] font-bold mb-1 relative z-10 uppercase tracking-wider">總支出 ({{ trip.base_currency }})</span>
                       <div class="text-2xl font-black text-white tracking-tight relative z-10">
                          {{ formatNumber(stats.totalSpent) }}
                       </div>
                  </div>
                  
                  <div class="bg-neutral-900 rounded-3xl p-5 border border-neutral-800/60 flex flex-col justify-center relative overflow-hidden">
                       <div class="absolute inset-0 bg-gradient-to-br from-purple-500/5 to-transparent z-0"></div>
                       <span class="text-neutral-500 text-[10px] font-bold mb-1 relative z-10 uppercase tracking-wider">每人平均</span>
                       <div class="text-2xl font-black text-white tracking-tight relative z-10">
                          {{ formatNumber(stats.averagePerMember) }}
                       </div>
                  </div>
              </div>

              <!-- Members Section -->
              <div class="bg-neutral-900 rounded-3xl p-6 border border-neutral-800/60">
                  <div class="flex justify-between items-center mb-5">
                    <h3 class="text-base font-bold">成員 ({{ trip.members?.length || 0 }})</h3>
                    <button @click="router.push(`/trips/${trip.id}/members/add`)" class="w-8 h-8 rounded-full bg-indigo-500/10 text-indigo-500 flex items-center justify-center border-0 cursor-pointer hover:bg-indigo-500/20 transition-colors">
                      <Icon icon="mdi:plus" class="text-xl" />
                    </button>
                  </div>
                  <div class="flex flex-wrap gap-2">
                    <div v-for="member in trip.members" :key="member.id" class="flex items-center gap-2 px-4 py-2 rounded-2xl bg-neutral-800/50 border border-neutral-800 text-sm font-medium">
                      <Icon icon="mdi:account-outline" class="text-indigo-400" />
                      <span>{{ member.name }}</span>
                      <span v-if="member.is_owner" class="text-[10px] bg-indigo-500/20 text-indigo-400 px-1.5 py-0.5 rounded-md ml-1 font-bold uppercase">主辦</span>
                    </div>
                  </div>
              </div>
              
              <!-- Quick Link to Ledger -->
              <NuxtLink :to="`/trips/${trip.id}/ledger`" class="group flex items-center justify-between bg-indigo-500 rounded-3xl p-6 text-white no-underline shadow-lg shadow-indigo-500/20 active:scale-[0.98] transition-all">
                  <div class="flex items-center gap-4">
                    <div class="w-12 h-12 rounded-2xl bg-white/20 flex items-center justify-center">
                      <Icon icon="mdi:notebook-outline" class="text-2xl" />
                    </div>
                    <div class="flex flex-col text-left">
                      <span class="font-bold text-lg">進入記帳本</span>
                      <span class="text-white/70 text-sm">查看明細與拆帳細節</span>
                    </div>
                  </div>
                  <Icon icon="mdi:chevron-right" class="text-2xl opacity-50 group-hover:translate-x-1 transition-transform" />
              </NuxtLink>

          </div>

          <!-- Stats Tab -->
          <div v-else-if="activeTab === 'stats'" class="animate-in fade-in slide-in-from-bottom-2 duration-500">
              <TripStats 
                  :transactions="transactions" 
                  :members="trip.members" 
                  :base-currency="trip.ledger?.base_currency || trip.base_currency || 'TWD'"
              />
          </div>

      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Icon } from '@iconify/vue'
import TripStats from '~/components/TripStats.vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useImages } from '~/composables/useImages'

definePageMeta({
  layout: 'main'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const { getImages, getImageUrl } = useImages()

const trip = ref<any>(null)
const transactions = ref<any[]>([])
const loading = ref(true)
const showMenu = ref(false)
const coverImage = ref<string | null>(null)

// Tabs
const activeTab = ref('overview')
const tabs = [
    { id: 'overview', label: '概覽', icon: 'mdi:view-dashboard-outline' },
    { id: 'stats', label: '統計', icon: 'mdi:chart-pie' }
]

// Stats Data (Kept for Overview KPI)
const stats = ref({
    totalSpent: 0,
    averagePerMember: 0
})

const formatNumber = (num: number) => {
    return num.toLocaleString(undefined, { maximumFractionDigits: 0 })
}

const formatDateRange = (start: string | null, end: string | null) => {
  if (!start && !end) return '日期未設定'
  const opts: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'short', day: 'numeric' }
  const startStr = start ? new Date(start).toLocaleDateString('zh-TW', opts) : '?'
  const endStr = end ? new Date(end).toLocaleDateString('zh-TW', opts) : '?'
  return `${startStr} - ${endStr}`
}

const calculateStats = () => {
    if (!trip.value || !transactions.value) return

    let total = 0
    const baseC = trip.value.ledger?.base_currency || trip.value.base_currency || 'TWD'

    transactions.value.forEach(txn => {
        if (txn.currency === baseC) {
            total += parseFloat(txn.total_amount)
        } else {
            const billing = parseFloat(txn.billing_amount)
            if (billing > 0) total += billing
        }
    })

    stats.value.totalSpent = total
    stats.value.averagePerMember = trip.value.members.length > 0 ? total / trip.value.members.length : 0
}

const fetchTripData = async () => {
  try {
    trip.value = await api.get<any>(`/api/trips/${route.params.id}`)
    
    if (trip.value.ledger_id) {
         const txns = await api.get<any[]>(`/api/ledgers/${trip.value.ledger_id}/transactions`)
         transactions.value = txns
         calculateStats()
    }
  } catch (e) {
    console.error('Failed to fetch trip data:', e)
    router.push('/trips')
  } finally {
    loading.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('確定要刪除此旅行？')) return
  try {
    await api.del(`/api/trips/${route.params.id}`)
    router.push('/trips')
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
  fetchTripData()
})
</script>
