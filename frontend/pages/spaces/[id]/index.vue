<template>
  <div class="flex flex-col min-h-screen">
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

    <div v-if="space" class="p-4 pt-0 pb-32">
      <!-- 1. Ledger Tab -->
      <div v-if="activeTab === 'ledger'" class="flex flex-col gap-6">
          <div class="bg-neutral-900 rounded-2xl border border-neutral-800 p-6 flex flex-col gap-5 shadow-sm">
              <div class="flex items-start justify-between">
                <div class="flex flex-col gap-1">
                   <div class="flex items-center gap-2 mb-1">
                      <div class="px-2.5 py-0.5 bg-indigo-500/10 text-indigo-400 rounded-full text-[10px] font-bold uppercase tracking-wider border border-indigo-500/20">
                         {{ space.type === 'trip' ? '專案' : '空間' }}
                      </div>
                   </div>
                   <h2 class="text-neutral-500 text-xs font-medium uppercase tracking-wider">總支出 ({{ space.base_currency }})</h2>
                   <div class="text-3xl font-bold text-white">{{ formatAmount(totalExpenses) }}</div>
                </div>
                <div class="w-12 h-12 rounded-xl bg-neutral-800 flex items-center justify-center text-indigo-500 border border-neutral-700">
                  <Icon icon="mdi:finance" class="text-2xl" />
                </div>
              </div>

              <div class="grid grid-cols-2 gap-3">
                 <div class="bg-neutral-800/50 rounded-xl p-4 border border-neutral-700/30 flex flex-col gap-1">
                    <span class="text-[10px] font-bold text-neutral-500 uppercase tracking-wider">交易</span>
                    <span class="text-lg font-bold">{{ transactions.length }} 筆</span>
                 </div>
                 <div class="bg-neutral-800/50 rounded-xl p-4 border border-neutral-700/30 flex flex-col gap-1">
                    <span class="text-[10px] font-bold text-neutral-500 uppercase tracking-wider">成員</span>
                    <span class="text-lg font-bold">{{ members.length }} 位</span>
                 </div>
              </div>
          </div>

          <div class="flex flex-col gap-4">
              <div class="flex items-center justify-between px-1">
                <h2 class="text-sm font-bold text-neutral-400 uppercase tracking-wider">最近交易</h2>
                <button class="w-8 h-8 rounded-full bg-neutral-800 flex items-center justify-center text-neutral-400 hover:text-white transition-colors border-0 cursor-pointer">
                  <Icon icon="mdi:magnify" class="text-lg" />
                </button>
              </div>

              <div v-if="transactions.length === 0" class="text-center py-16 px-6 bg-neutral-900 rounded-2xl border border-neutral-800 flex flex-col items-center">
                <div class="w-12 h-12 rounded-2xl bg-neutral-800 flex items-center justify-center text-neutral-600 mb-4">
                  <Icon icon="mdi:receipt-text-outline" class="text-2xl" />
                </div>
                <p class="text-neutral-500 font-medium">尚無交易明細</p>
              </div>

              <div v-else class="flex flex-col gap-3">
                <TransactionListItem
                  v-for="txn in transactions"
                  :key="txn.id"
                  :transaction="txn"
                  :base-currency="space.base_currency"
                  @click="router.push(`/spaces/${space.id}/transaction/${txn.id}`)"
                />
              </div>
          </div>
      </div>

      <!-- 2. Comparison Tab -->
      <div v-else-if="activeTab === 'comparison'" class="flex flex-col gap-6">
          <div v-if="stores.length === 0" class="text-center py-16 px-6 bg-neutral-900 rounded-2xl border border-neutral-800 flex flex-col items-center">
            <div class="w-12 h-12 rounded-2xl bg-neutral-800 flex items-center justify-center text-neutral-600 mb-4">
              <Icon icon="mdi:scale-balance" class="text-2xl" />
            </div>
            <p class="text-neutral-500 font-medium">還沒有店家資料</p>
          </div>

          <div v-else class="flex flex-col gap-4">
            <div v-for="store in stores" :key="store.id" class="bg-neutral-900 rounded-2xl border border-neutral-800 overflow-hidden shadow-sm">
              <div class="p-4 border-b border-neutral-800 flex items-center justify-between bg-neutral-800/20">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-xl bg-neutral-800 flex items-center justify-center text-indigo-500 border border-neutral-700">
                    <Icon icon="mdi:store-outline" class="text-xl" />
                  </div>
                  <h3 class="font-bold text-white">{{ store.name }}</h3>
                </div>
                <button 
                  @click="router.push(`/spaces/${space.id}/stores/${store.id}`)"
                  class="px-4 py-1.5 rounded-lg bg-neutral-800 text-neutral-100 text-xs font-bold hover:bg-neutral-700 transition-colors border-0 cursor-pointer"
                >
                  明細
                </button>
              </div>
              
              <div class="p-4">
                <div v-if="!store.products || store.products.length === 0" class="text-center py-2 text-neutral-600 text-xs italic">
                  目前沒有商品紀錄
                </div>
                <div v-else class="flex flex-col gap-2.5">
                  <div v-for="product in (store.products || []).slice(0, 3)" :key="product.id" class="flex justify-between items-center text-sm">
                    <span class="text-neutral-400 font-medium">{{ product.name }}</span>
                    <span class="text-white font-bold">{{ space.base_currency }} {{ formatAmount(product.price) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
      </div>

      <!-- 3. Stats Tab -->
      <div v-else-if="activeTab === 'stats'" class="flex flex-col gap-6">
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
        class="fixed bottom-24 right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg rounded-full flex items-center justify-center text-white transition-all active:scale-90 z-20 cursor-pointer border-0"
    >
        <Icon icon="mdi:plus" class="text-3xl" />
    </button>
    
    <button 
        v-if="activeTab === 'comparison' && space"
        @click="router.push(`/spaces/${route.params.id}/stores/add`)"
        class="fixed bottom-24 right-6 w-14 h-14 bg-green-500 hover:bg-green-600 shadow-lg rounded-full flex items-center justify-center text-white transition-all active:scale-90 z-20 cursor-pointer border-0"
    >
        <Icon icon="mdi:store-plus-outline" class="text-3xl" />
    </button>

    <!-- Standard Navigation Tabs -->
    <BottomNav v-model="activeTab" :items="navItems" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import SpaceHeader from '~/components/SpaceHeader.vue'
import SpaceStats from '~/components/SpaceStats.vue'
import TransactionListItem from '~/components/TransactionListItem.vue'
import BottomNav from '~/components/BottomNav.vue'

definePageMeta({
  layout: 'default',
  hideGlobalNav: true
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
  { id: 'ledger', label: '帳本', icon: 'mdi:receipt-text-outline', activeIcon: 'mdi:receipt-text' },
  { id: 'comparison', label: '比價', icon: 'mdi:scale-balance', activeIcon: 'mdi:scale-balance' },
  { id: 'stats', label: '統計', icon: 'mdi:chart-pie-outline', activeIcon: 'mdi:chart-pie' }
]

const totalExpenses = computed(() => {
  return transactions.value.reduce((sum, txn) => {
    const amount = (txn.billing_amount && Number(txn.billing_amount) > 0) ? txn.billing_amount : txn.total_amount
    return sum + Number(amount || 0)
  }, 0)
})

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toLocaleString('zh-TW')
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
