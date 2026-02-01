<template>
  <div class="flex flex-col gap-6">
    <ImmersiveHeader
        fallback-icon="mdi:chart-bar"
        class="rounded-2xl"
    >
        <template #top-left>
             <button @click="router.push('/ledger')" class="w-10 h-10 rounded-full bg-black/30 backdrop-blur-md text-white flex items-center justify-center hover:bg-black/50 transition-colors border-0 cursor-pointer">
                <Icon icon="mdi:arrow-left" class="text-xl" />
             </button>
        </template>
        
        <template #bottom>
             <h1 class="text-2xl font-bold text-white shadow-sm mb-1">日常統計</h1>
             <p class="text-sm text-neutral-300">本月支出概況</p>
        </template>
    </ImmersiveHeader>

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-400">
        <Icon icon="eos-icons:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else class="flex flex-col gap-6 px-4">
        
        <!-- Total Spending Card -->
        <div class="bg-neutral-800 rounded-2xl p-6 flex flex-col items-center justify-center relative overflow-hidden">
            <div class="absolute inset-0 bg-gradient-to-br from-green-500/10 to-teal-500/10 z-0"></div>
            <div class="relative z-10 text-center">
                <span class="text-neutral-400 text-sm mb-1 block">本月總支出</span>
                <div class="text-4xl font-bold text-white tracking-tight">
                    <span class="text-2xl text-neutral-500 mr-1">$</span>
                    {{ formatNumber(totalSpent) }}
                </div>
            </div>
        </div>

        <!-- Weekly Activity (Bar Chart Mock) -->
        <div>
            <h2 class="text-lg font-bold mb-4 flex items-center gap-2">
                <Icon icon="mdi:calendar-week" class="text-green-400" />
                每週趨勢
            </h2>
            <div class="flex justify-between items-end h-32 px-2 gap-2">
                <div v-for="(val, index) in weeklyData" :key="index" class="flex-1 flex flex-col items-center gap-2 group">
                    <div class="w-full bg-neutral-800 rounded-t-lg relative group-hover:bg-neutral-700 transition-colors overflow-hidden">
                        <div class="absolute bottom-0 w-full bg-green-500/50 rounded-t-lg group-hover:bg-green-500 transition-all" :style="{ height: val + '%' }"></div>
                    </div>
                    <span class="text-xs text-neutral-500">W{{ index + 1 }}</span>
                </div>
            </div>
        </div>

        <!-- Category Breakdown -->
        <div>
            <h2 class="text-lg font-bold mb-4 flex items-center gap-2">
                <Icon icon="mdi:shape-outline" class="text-green-400" />
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
                         <div class="h-full bg-green-500 rounded-full transition-all duration-1000" :style="{ width: cat.percentage + '%' }"></div>
                     </div>
                </div>
            </div>
        </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import ImmersiveHeader from '~/components/ImmersiveHeader.vue'

definePageMeta({
  layout: 'main'
})

const router = useRouter()
const loading = ref(true)
const totalSpent = ref(0)
const categories = ref<any[]>([])
const weeklyData = ref<number[]>([])

const formatNumber = (num: number) => {
    return num.toLocaleString()
}

// Mock data generator for visualization before backend integration
const processMockData = () => {
    const total = 15600
    totalSpent.value = total
    
    categories.value = [
        { name: '餐飲', amount: 6240, percentage: 40, icon: '🍔' },
        { name: '交通', amount: 3120, percentage: 20, icon: '🚆' },
        { name: '租金', amount: 4680, percentage: 30, icon: '🏠' },
        { name: '其他', amount: 1560, percentage: 10, icon: '📦' },
    ].sort((a,b) => b.amount - a.amount)

    // Mock weekly spending relative percentage (0-100)
    weeklyData.value = [30, 60, 45, 80]
}

const fetchData = async () => {
    try {
        // Here we would fetch real ledger transactions
        await new Promise(resolve => setTimeout(resolve, 500)) // Fake delay
        processMockData()

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
