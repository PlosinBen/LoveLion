<template>
  <div class="flex flex-col gap-6">
    <!-- Immersive Header Area -->
    <ImmersiveHeader
      :image="coverImage"
      fallback-icon="mdi:airplane"
      class="rounded-2xl"
    >
      <template #top-left>
        <button @click="router.push('/trips')" class="flex justify-center items-center w-10 h-10 rounded-xl bg-black/20 text-white backdrop-blur-md border-0 cursor-pointer hover:bg-black/40 transition-colors">
            <Icon icon="mdi:arrow-left" class="text-2xl" />
        </button>
      </template>

      <template #top-right>
        <button @click="showMenu = !showMenu" class="flex justify-center items-center w-10 h-10 rounded-xl bg-black/20 text-white backdrop-blur-md border-0 cursor-pointer hover:bg-black/40 transition-colors relative">
            <Icon icon="mdi:dots-vertical" class="text-2xl" />
            <div v-if="showMenu" class="absolute top-12 right-0 bg-neutral-800 rounded-xl border border-neutral-700 shadow-lg z-20 overflow-hidden min-w-32 py-1">
              <NuxtLink :to="`/trips/${trip?.id}/edit`" class="block w-full px-4 py-3 text-left text-white hover:bg-neutral-700 transition-colors no-underline">編輯旅行</NuxtLink>
              <button @click="handleDelete" class="w-full px-4 py-3 text-left text-red-500 hover:bg-neutral-700 transition-colors border-0 bg-transparent cursor-pointer">刪除旅行</button>
            </div>
        </button>
      </template>

      <template #bottom>
         <h1 class="text-2xl font-bold text-white mb-1 shadow-sm">{{ trip?.name || '載入中...' }}</h1>
         <div class="flex items-center gap-2 text-neutral-300 text-sm">
            <Icon icon="mdi:calendar-range" class="text-indigo-400" />
            <span>{{ formatDateRange(trip?.start_date, trip?.end_date) }}</span>
         </div>
      </template>
    </ImmersiveHeader>

    <div v-if="loading" class="text-center text-neutral-400 p-10">
        <Icon icon="eos-icons:loading" class="text-3xl animate-spin" />
    </div>

    <template v-else-if="trip">
        
        <!-- Tabs -->
        <div class="flex items-center gap-1 bg-neutral-900 mx-1 p-1 rounded-xl border border-neutral-800">
            <button 
                v-for="tab in tabs" 
                :key="tab.id"
                @click="activeTab = tab.id"
                class="flex-1 flex items-center justify-center gap-2 py-2.5 rounded-lg text-sm font-semibold transition-all"
                :class="activeTab === tab.id ? 'bg-neutral-800 text-white shadow-sm' : 'text-neutral-400 hover:text-neutral-200'"
            >
                <Icon :icon="tab.icon" />
                {{ tab.label }}
            </button>
        </div>

        <!-- Tab Content -->
        
        <!-- Overview Tab -->
        <div v-if="activeTab === 'overview'" class="flex flex-col gap-6 animate-fade-in">
            
            <!-- KPI Cards -->
            <div class="grid grid-cols-2 gap-3 px-1">
                <div class="bg-neutral-800 rounded-2xl p-4 flex flex-col justify-center relative overflow-hidden">
                     <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/10 to-transparent z-0"></div>
                     <span class="text-neutral-400 text-xs mb-1 relative z-10">總支出 ({{ trip.base_currency }})</span>
                     <div class="text-2xl font-bold text-white tracking-tight relative z-10">
                        {{ formatNumber(stats.totalSpent) }}
                     </div>
                </div>
                
                <div class="bg-neutral-800 rounded-2xl p-4 flex flex-col justify-center relative overflow-hidden">
                     <div class="absolute inset-0 bg-gradient-to-br from-purple-500/10 to-transparent z-0"></div>
                     <span class="text-neutral-400 text-xs mb-1 relative z-10">每人平均</span>
                     <div class="text-2xl font-bold text-white tracking-tight relative z-10">
                        {{ formatNumber(stats.averagePerMember) }}
                     </div>
                </div>
            </div>

            <!-- Members Section -->
            <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
                <div class="flex justify-between items-center mb-4">
                  <h3 class="text-base font-semibold">成員 ({{ trip.members?.length || 0 }})</h3>
                  <button @click="router.push(`/trips/${trip.id}/members/add`)" class="flex items-center gap-1 px-3 py-1.5 rounded-lg bg-indigo-500/20 text-indigo-500 text-sm border-0 cursor-pointer hover:bg-indigo-500/30 transition-colors">
                    <Icon icon="mdi:plus" /> 新增
                  </button>
                </div>
                <div class="flex flex-wrap gap-2">
                  <div v-for="member in trip.members" :key="member.id" class="flex items-center gap-2 px-3 py-2 rounded-full bg-neutral-800 text-sm">
                    <Icon icon="mdi:account" class="text-indigo-500" />
                    <span>{{ member.name }}</span>
                    <span v-if="member.is_owner" class="text-xs text-indigo-500">(主辦)</span>
                  </div>
                </div>
            </div>
            
            <!-- Quick Link to Ledger (Optional, but user removed detailed list so maybe useful?) -->
            <div class="px-1">
                 <NuxtLink :to="`/trips/${trip.id}/ledger`" class="block w-full bg-neutral-900 border border-neutral-800 rounded-2xl p-4 text-center text-white font-semibold hover:bg-neutral-800 transition-colors no-underline flex items-center justify-center gap-2">
                    <Icon icon="mdi:notebook-outline" class="text-xl" />
                    進入記帳本
                 </NuxtLink>
            </div>

        </div>

        <!-- Stats Tab -->
        <div v-else-if="activeTab === 'stats'" class="animate-fade-in">
            <TripStats 
                :transactions="transactions" 
                :members="trip.members" 
                :base-currency="trip.ledger?.base_currency || trip.base_currency || 'TWD'"
            />
        </div>

    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Icon } from '@iconify/vue'
import ImmersiveHeader from '~/components/ImmersiveHeader.vue'
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
    { id: 'stats', label: '帳務分析', icon: 'mdi:chart-pie' }
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
    // Simplified KPI calc for Overview (Total Base)
    // We can reuse the logic from TripStats logic or keep it simple
    // Let's assume transactions converted to base for KPI
    // Base currency
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
    // 1. Fetch Trip
    trip.value = await api.get<any>(`/api/trips/${route.params.id}`)
    
    // 2. Fetch Images
    const images = await getImages(trip.value.id, 'trip')
    if (images.length > 0) coverImage.value = images[0].file_path

    // 3. Fetch Transactions (All, for stats)
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
