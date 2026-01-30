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

        <!-- Recent Activity (New) -->
        <div v-if="recentTransactions.length > 0" class="px-1">
             <div class="flex justify-between items-center mb-3">
                <h3 class="text-base font-bold flex items-center gap-2">
                    <Icon icon="mdi:history" class="text-indigo-400" />
                    最新動態
                </h3>
                <NuxtLink :to="`/trips/${trip.id}/ledger`" class="text-xs text-neutral-400 hover:text-white no-underline">查看全部</NuxtLink>
             </div>
             
             <div class="flex flex-col gap-2">
                <div v-for="txn in recentTransactions" :key="txn.id" class="bg-neutral-900 rounded-xl p-3 border border-neutral-800 flex items-center justify-between" @click="router.push(`/trips/${trip.id}/ledger`)">
                    <div class="flex items-center gap-3">
                         <div class="w-10 h-10 rounded-lg bg-neutral-800 flex items-center justify-center text-xl">
                             <Icon :icon="getCategoryIcon(txn.category)" class="text-neutral-400" />
                         </div>
                         <div>
                             <div class="font-medium text-sm">{{ txn.items?.[0]?.name || txn.category }}</div>
                             <div class="text-xs text-neutral-500">{{ formatDate(txn.date) }} • {{ txn.payer }} 付款</div>
                         </div>
                    </div>
                    <div class="font-bold text-neutral-200">{{ formatNumber(Number(txn.total_amount)) }}</div>
                </div>
             </div>
        </div>

        <!-- Charts: Spend by Category -->
        <div class="px-1" v-if="stats.categories.length > 0">
             <h3 class="text-base font-bold mb-3 flex items-center gap-2">
                <Icon icon="mdi:shape-outline" class="text-indigo-400" />
                分類統計
             </h3>
             <div class="flex flex-col gap-3">
                 <div v-for="cat in stats.categories" :key="cat.name" class="bg-neutral-900 rounded-xl p-3 border border-neutral-800">
                      <div class="flex justify-between items-center mb-2">
                          <span class="font-medium text-sm">{{ cat.name }}</span>
                          <div class="text-right">
                              <span class="font-bold text-sm">{{ formatNumber(cat.amount) }}</span>
                              <span class="text-xs text-neutral-500 ml-1">{{ cat.percentage }}%</span>
                          </div>
                      </div>
                      <div class="h-1.5 bg-neutral-800 rounded-full overflow-hidden">
                          <div class="h-full bg-indigo-500 rounded-full" :style="{ width: cat.percentage + '%' }"></div>
                      </div>
                 </div>
             </div>
        </div>

        <!-- Charts: Top Spenders (Net Expense) -->
        <div class="px-1" v-if="stats.memberSpending.length > 0">
             <h3 class="text-base font-bold mb-3 flex items-center gap-2">
                <Icon icon="mdi:account-cash-outline" class="text-indigo-400" />
                個人支出 (分攤後)
             </h3>
             <div class="bg-neutral-900 rounded-xl border border-neutral-800 divide-y divide-neutral-800">
                <div v-for="(spender, index) in stats.memberSpending" :key="spender.name" class="p-3 flex items-center justify-between">
                     <div class="flex items-center gap-3">
                         <div class="w-6 h-6 rounded-full bg-neutral-800 flex items-center justify-center text-xs font-bold text-neutral-500">
                             {{ index + 1 }}
                         </div>
                         <span class="font-medium text-sm">{{ spender.name }}</span>
                     </div>
                     <span class="font-bold text-indigo-400">{{ formatNumber(spender.amount) }}</span>
                </div>
            </div>
        </div>

      <!-- Members Section (Simplified) -->
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

    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { Icon } from '@iconify/vue'
import ImmersiveHeader from '~/components/ImmersiveHeader.vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useImages } from '~/composables/useImages'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const trip = ref<any>(null)
const transactions = ref<any[]>([])
const loading = ref(true)
const showMenu = ref(false)
const coverImage = ref<string | null>(null)
const { getImages, getImageUrl } = useImages()

// Stats Data
const stats = ref({
    totalSpent: 0,
    averagePerMember: 0,
    categories: [] as any[],
    memberSpending: [] as any[]
})

const recentTransactions = computed(() => {
    return transactions.value.slice(0, 3)
})

const formatNumber = (num: number) => {
    return num.toLocaleString(undefined, { maximumFractionDigits: 0 })
}

const formatDate = (dateStr: string) => {
    const d = new Date(dateStr)
    return `${d.getMonth() + 1}/${d.getDate()}`
}

const formatDateRange = (start: string | null, end: string | null) => {
  if (!start && !end) return '日期未設定'
  const opts: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'short', day: 'numeric' }
  const startStr = start ? new Date(start).toLocaleDateString('zh-TW', opts) : '?'
  const endStr = end ? new Date(end).toLocaleDateString('zh-TW', opts) : '?'
  return `${startStr} - ${endStr}`
}

const getCategoryIcon = (category: string) => {
  const icons: Record<string, string> = {
    'Food': 'mdi:food',
    'Transport': 'mdi:train-car',
    'Shopping': 'mdi:shopping',
    'Entertainment': 'mdi:movie',
    '住宿': 'mdi:bed',
    '餐飲': 'mdi:food',
    '交通': 'mdi:train-car',
    '購物': 'mdi:shopping',
    '娛樂': 'mdi:movie',
  }
  return icons[category] || 'mdi:receipt'
}

const calculateStats = () => {
    if (!trip.value || !transactions.value) return

    let total = 0
    const catMap: Record<string, number> = {}
    const memberExpenseMap: Record<string, number> = {}

    // Init member map
    trip.value.members.forEach((m: any) => memberExpenseMap[m.name] = 0)

    transactions.value.forEach(txn => {
        const amount = parseFloat(txn.total_amount)
        total += amount

        // Category
        const cat = txn.category || '未分類'
        catMap[cat] = (catMap[cat] || 0) + amount

        // Net Expense (Splits)
        // If splits exist, use them. If not, default to Payer = Consumer (Not split) ?? 
        // Actually if no split, usually means "Payer paid for everyone" or "Payer paid for self"?
        // In this app context, if "Split" is not enabled, we usually assume Payer paid for THEMSELVES (Personal expense in trip?) OR Payer paid for GROUP (Shared)?
        // Looking at add.vue: if !isSplitEnabled, we create a split { is_payer: false, amount: total, member_id: payer }. 
        // So Payer is Consumer. Correct.
        if (txn.splits && txn.splits.length > 0) {
            txn.splits.forEach((split: any) => {
                if (!split.is_payer) {
                    // This is a consumption split
                    const name = split.name
                    memberExpenseMap[name] = (memberExpenseMap[name] || 0) + parseFloat(split.amount)
                }
            })
        }
    })

    stats.value.totalSpent = total
    stats.value.averagePerMember = trip.value.members.length > 0 ? total / trip.value.members.length : 0

    // Process Categories
    stats.value.categories = Object.entries(catMap)
        .map(([name, amount]) => ({
            name,
            amount,
            percentage: total > 0 ? Math.round((amount / total) * 100) : 0
        }))
        .sort((a, b) => b.amount - a.amount)

    // Process Members
    stats.value.memberSpending = Object.entries(memberExpenseMap)
        .map(([name, amount]) => ({
            name,
            amount
        }))
        .sort((a, b) => b.amount - a.amount)
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
