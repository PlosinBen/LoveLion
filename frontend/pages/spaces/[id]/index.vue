<template>
  <div class="flex flex-col">
    <!-- Space Navigation -->
    <SpaceHeader 
      v-if="space"
      :title="space.name"
      :show-back="true"
      :show-switcher="false"
      class="pt-0 px-2"
    >
      <template #right>
        <button 
            @click="router.push(`/spaces/${route.params.id}/settings`)" 
            class="w-10 h-10 rounded-full bg-neutral-800 text-white flex items-center justify-center hover:bg-neutral-700 transition-colors border-0 cursor-pointer shrink-0"
        >
          <Icon icon="mdi:cog-outline" class="text-xl" />
        </button>
      </template>
    </SpaceHeader>

    <div v-if="space" class="p-4 pt-0">
      <!-- 1. Ledger Tab -->
      <div v-if="activeTab === 'ledger'" class="flex flex-col gap-6 animate-in fade-in slide-in-from-bottom-2 duration-500">
          <div class="space-summary bg-neutral-900 rounded-3xl border border-neutral-800/60 p-6 flex flex-col gap-5 shadow-sm">
              <div class="flex items-start justify-between">
                <div class="flex flex-col gap-1.5">
                   <div class="flex items-center gap-2 mb-1">
                      <div class="px-3 py-1 bg-indigo-500/10 text-indigo-400 rounded-full text-[10px] font-black uppercase tracking-widest border border-indigo-500/20">
                         {{ space.type === 'trip' ? '旅行專案' : '個人空間' }}
                      </div>
                   </div>
                   <h2 class="text-neutral-500 text-xs font-bold uppercase tracking-widest px-0.5">總支出 ({{ space.base_currency }})</h2>
                   <div class="text-4xl font-black text-white tracking-tighter">{{ formatAmount(totalExpenses) }}</div>
                </div>
                <div class="w-12 h-12 rounded-2xl bg-indigo-500/10 flex items-center justify-center text-indigo-500 border border-indigo-500/20">
                  <Icon icon="mdi:finance" class="text-2xl" />
                </div>
              </div>

              <div class="grid grid-cols-2 gap-3">
                 <div class="bg-neutral-800/40 rounded-2xl p-4 border border-white/5 flex flex-col gap-1">
                    <span class="text-[10px] font-black text-neutral-500 uppercase tracking-widest">交易數量</span>
                    <span class="text-lg font-bold">{{ transactions.length }} 筆</span>
                 </div>
                 <div class="bg-neutral-800/40 rounded-2xl p-4 border border-white/5 flex flex-col gap-1">
                    <span class="text-[10px] font-black text-neutral-500 uppercase tracking-widest">成員人數</span>
                    <span class="text-lg font-bold">{{ members.length }} 位</span>
                 </div>
              </div>
          </div>

          <div class="flex flex-col gap-4">
              <div class="flex items-center justify-between px-1">
                <h2 class="text-xs font-black text-neutral-500 uppercase tracking-widest">最近交易</h2>
                <div class="flex gap-1.5">
                   <button class="w-8 h-8 rounded-full bg-neutral-800/50 flex items-center justify-center text-neutral-400 hover:text-white transition-colors border-0 cursor-pointer">
                      <Icon icon="mdi:filter-variant" class="text-lg" />
                   </button>
                   <button class="w-8 h-8 rounded-full bg-neutral-800/50 flex items-center justify-center text-neutral-400 hover:text-white transition-colors border-0 cursor-pointer">
                      <Icon icon="mdi:magnify" class="text-lg" />
                   </button>
                </div>
              </div>

              <div v-if="transactions.length === 0" class="text-center py-20 px-6 bg-neutral-900 rounded-3xl border border-neutral-800/50 flex flex-col items-center">
                <div class="w-16 h-16 rounded-3xl bg-neutral-800 flex items-center justify-center text-neutral-600 mb-4 border border-white/5">
                  <Icon icon="mdi:receipt-text-outline" class="text-3xl" />
                </div>
                <p class="text-neutral-500 font-bold mb-1">尚無交易明細</p>
                <p class="text-neutral-600 text-xs">點擊下方的 + 按鈕開始記帳吧！</p>
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
                      <h3 class="text-sm font-bold mb-0.5 truncate text-neutral-100">{{ txn.title || txn.category || '未分類' }}</h3>
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
      </div>

      <!-- 2. Comparison Tab -->
      <div v-else-if="activeTab === 'comparison'" class="flex flex-col gap-6 animate-in fade-in slide-in-from-bottom-2 duration-500">
          <button 
              @click="router.push(`/spaces/${route.params.id}/stores/add`)"
              class="fixed bottom-24 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg shadow-indigo-500/30 rounded-full flex items-center justify-center text-white transition-transform active:scale-90 z-20 cursor-pointer border-0"
          >
              <Icon icon="mdi:store-plus-outline" class="text-3xl" />
          </button>

          <div v-if="stores.length === 0" class="text-center py-20 px-6 bg-neutral-900 rounded-3xl border border-neutral-800/50 flex flex-col items-center">
            <div class="w-16 h-16 rounded-3xl bg-neutral-800 flex items-center justify-center text-neutral-600 mb-4 border border-white/5">
              <Icon icon="mdi:scale-balance" class="text-3xl" />
            </div>
            <p class="text-neutral-500 font-bold mb-1">還沒有店家資料</p>
            <p class="text-neutral-600 text-xs">新增店家開始比較商品價格</p>
          </div>

          <div v-else class="flex flex-col gap-5">
            <div v-for="store in stores" :key="store.id" class="bg-neutral-900 rounded-3xl border border-neutral-800/60 overflow-hidden shadow-sm">
              <div class="p-5 border-b border-neutral-800/60 flex items-center justify-between bg-neutral-800/10">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-xl bg-indigo-500/10 flex items-center justify-center text-indigo-500 border border-indigo-500/10">
                    <Icon icon="mdi:store-outline" class="text-xl" />
                  </div>
                  <h3 class="font-bold text-white">{{ store.name }}</h3>
                </div>
                <button 
                  @click="router.push(`/spaces/${space.id}/stores/${store.id}`)"
                  class="px-4 py-2 rounded-xl bg-neutral-800 text-neutral-100 text-xs font-bold hover:bg-neutral-700 transition-colors border-0 cursor-pointer"
                >
                  查看明細
                </button>
              </div>
              
              <div class="p-5">
                <div v-if="!store.products || store.products.length === 0" class="text-center py-6 text-neutral-600 text-xs italic">
                  目前沒有商品紀錄
                </div>
                <div v-else class="flex flex-col gap-3">
                  <div v-for="product in (store.products || []).slice(0, 3)" :key="product.id" class="flex justify-between items-center text-sm">
                    <span class="text-neutral-300 font-medium">{{ product.name }}</span>
                    <span class="text-white font-black">{{ space.base_currency }} {{ formatAmount(product.price) }}</span>
                  </div>
                  <div v-if="store.products.length > 3" class="text-center pt-2 text-[10px] font-black text-neutral-600 uppercase tracking-widest">
                    及其他 {{ store.products.length - 3 }} 項商品
                  </div>
                </div>
              </div>
            </div>
          </div>
      </div>

      <!-- 3. Stats Tab -->
      <div v-else-if="activeTab === 'stats'" class="flex flex-col gap-6 animate-in fade-in slide-in-from-bottom-2 duration-500">
         <SpaceStats 
            v-if="space"
            :transactions="transactions" 
            :members="members" 
            :base-currency="space.base_currency" 
         />
      </div>
    </div>

    <!-- Floating Action Button -->
    <button 
        v-if="activeTab === 'ledger' && space"
        @click="router.push(`/spaces/${route.params.id}/transaction/add`)"
        class="fixed bottom-24 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg shadow-indigo-500/30 rounded-full flex items-center justify-center text-white transition-transform active:scale-90 z-20 cursor-pointer border-0"
    >
        <Icon icon="mdi:plus" class="text-3xl" />
    </button>

    <!-- Custom Footer for Tab Navigation -->
    <NuxtLayout name="default">
      <template #footer>
        <BottomNav v-model="activeTab" :items="navItems" />
      </template>
    </NuxtLayout>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import SpaceHeader from '~/components/SpaceHeader.vue'
import SpaceStats from '~/components/SpaceStats.vue'
import BottomNav from '~/components/BottomNav.vue'

// Back to default layout which provides Header
definePageMeta({
  layout: 'default'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const space = ref<any>(null)
const transactions = ref<any[]>([])
const stores = ref<any[]>([])
const members = ref<any[]>([])
const activeTab = ref(route.query.tab?.toString() || 'ledger')

const navItems = [
  { id: 'ledger', label: '帳本', icon: 'mdi:receipt-text-outline' },
  { id: 'comparison', label: '比價', icon: 'mdi:scale-balance' },
  { id: 'stats', label: '分帳', icon: 'mdi:chart-pie-outline' }
]

const totalExpenses = computed(() => {
  return transactions.value.reduce((sum, txn) => {
    const amount = (txn.billing_amount && Number(txn.billing_amount) > 0) ? txn.billing_amount : txn.total_amount
    return sum + Number(amount || 0)
  }, 0)
})

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toLocaleString('zh-TW', { maximumFractionDigits: 0 })
}

const formatDateShort = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', { month: '2-digit', day: '2-digit' })
}

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

const fetchData = async () => {
  try {
    const [spaceData, txnData, storeData, membersData] = await Promise.all([
      api.get<any>(`/api/spaces/${route.params.id}`),
      api.get<any[]>(`/api/spaces/${route.params.id}/transactions`),
      api.get<any[]>(`/api/spaces/${route.params.id}/stores`),
      api.get<any[]>(`/api/spaces/${route.params.id}/members`)
    ])
    
    space.value = spaceData
    transactions.value = txnData || []
    stores.value = storeData || []
    members.value = membersData || []
  } catch (e) {
    console.error('Failed to fetch space data:', e)
    router.push('/')
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
