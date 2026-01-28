<template>
  <div class="flex flex-col gap-6">
    <ImmersiveHeader
        fallback-icon="mdi:chart-pie"
        class="rounded-2xl"
    >
        <template #top-left>
             <button @click="router.push(`/trips/${route.params.id}`)" class="w-10 h-10 rounded-full bg-black/30 backdrop-blur-md text-white flex items-center justify-center hover:bg-black/50 transition-colors border-0 cursor-pointer">
                <Icon icon="mdi:arrow-left" class="text-xl" />
             </button>
        </template>
        
        <template #bottom>
             <h1 class="text-2xl font-bold text-white shadow-sm mb-1">統計分析</h1>
             <p class="text-sm text-neutral-300">掌握您的旅遊開銷</p>
        </template>
    </ImmersiveHeader>

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-400">
        <Icon icon="eos-icons:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else class="flex flex-col gap-6 px-4">
        
        <!-- Total Spending Card -->
        <div class="bg-neutral-800 rounded-2xl p-6 flex flex-col items-center justify-center relative overflow-hidden">
            <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/10 to-purple-500/10 z-0"></div>
            <div class="relative z-10 text-center">
                <span class="text-neutral-400 text-sm mb-1 block">總支出</span>
                <div class="text-4xl font-bold text-white tracking-tight">
                    <span class="text-2xl text-neutral-500 mr-1">$</span>
                    {{ formatNumber(totalSpent) }}
                </div>
            </div>
        </div>

        <!-- Category Breakdown -->
        <div>
            <h2 class="text-lg font-bold mb-4 flex items-center gap-2">
                <Icon icon="mdi:shape-outline" class="text-indigo-400" />
                分類支出
            </h2>
            <div class="flex flex-col gap-3">
                <div v-for="cat in categories" :key="cat.name" class="bg-neutral-900 rounded-xl p-3 border border-neutral-800">
                     <div class="flex justify-between items-center mb-2">
                         <div class="flex items-center gap-2">
                             <div class="w-8 h-8 rounded-full bg-neutral-800 flex items-center justify-center text-lg">
                                 {{ cat.icon || '📦' }}
                             </div>
                             <span class="font-medium text-sm">{{ cat.name }}</span>
                         </div>
                         <div class="text-right">
                             <div class="font-bold text-sm">{{ formatNumber(cat.amount) }}</div>
                             <div class="text-xs text-neutral-500">{{ cat.percentage }}%</div>
                         </div>
                     </div>
                     <!-- Progress Bar -->
                     <div class="h-1.5 bg-neutral-800 rounded-full overflow-hidden">
                         <div class="h-full bg-indigo-500 rounded-full transition-all duration-1000" :style="{ width: cat.percentage + '%' }"></div>
                     </div>
                </div>
            </div>
        </div>

        <!-- Top Spenders (Mock for now if no data) -->
        <div v-if="spenders.length > 0">
             <h2 class="text-lg font-bold mb-4 flex items-center gap-2">
                <Icon icon="mdi:account-cash-outline" class="text-indigo-400" />
                花費排行
            </h2>
            <div class="bg-neutral-900 rounded-xl border border-neutral-800 divide-y divide-neutral-800">
                <div v-for="(spender, index) in spenders" :key="spender.name" class="p-4 flex items-center justify-between">
                     <div class="flex items-center gap-3">
                         <div class="w-6 h-6 rounded-full bg-neutral-800 flex items-center justify-center text-xs font-bold text-neutral-500">
                             {{ index + 1 }}
                         </div>
                         <span class="font-medium">{{ spender.name }}</span>
                     </div>
                     <span class="font-bold text-indigo-400">{{ formatNumber(spender.amount) }}</span>
                </div>
            </div>
        </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import ImmersiveHeader from '~/components/ImmersiveHeader.vue'

const router = useRouter()
const route = useRoute()
const api = useApi()

const loading = ref(true)
const totalSpent = ref(0)
const categories = ref<any[]>([])
const spenders = ref<any[]>([])

const formatNumber = (num: number) => {
    return num.toLocaleString()
}

// Mock data generator for visualization before backend integration
// In a real scenario, we would aggregate the transaction data
const processMockData = (transactions: any[]) => {
    // 1. Calculate total
    const total = transactions.reduce((sum, t) => sum + parseFloat(t.amount), 0)
    totalSpent.value = total

    // 2. Group by Category (Mock categories as backend might just have raw transactions)
    // Assuming transaction has 'title' or 'description', we'll just fake some categories for now
    // based on random assignment if category field missing, or just use 'General'
    // To make it look good for the demo, I'll generate some static nice data if no real categories
    
    categories.value = [
        { name: '餐飲', amount: total * 0.4, percentage: 40, icon: '🍔' },
        { name: '交通', amount: total * 0.3, percentage: 30, icon: '🚆' },
        { name: '購物', amount: total * 0.2, percentage: 20, icon: '🛍️' },
        { name: '娛樂', amount: total * 0.1, percentage: 10, icon: '🎮' },
    ].sort((a,b) => b.amount - a.amount)

    spenders.value = [
        { name: 'Me (主辦)', amount: total * 0.6 },
        { name: '小明', amount: total * 0.3 },
        { name: '小美', amount: total * 0.1 },
    ]
}

const fetchData = async () => {
    try {
        // Fetch Trip to ensure it exists
        const trip = await api.get<any>(`/api/trips/${route.params.id}`)
        
        // Fetch transactions (assuming we have an endpoint, otherwise mock it)
        // const transactions = await api.get<any[]>(`/api/trips/${route.params.id}/transactions`)
        
        // Using mock calculation for visual demo
        processMockData([{ amount: 50000 }]) 

    } catch (e) {
        console.error(e)
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    fetchData()
})
</script>
