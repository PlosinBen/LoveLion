<template>
  <div class="flex flex-col">
    <SpaceHeader
      title="日常統計"
      icon="mdi:chart-bar"
    />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-400">
        <Icon icon="eos-icons:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else class="flex flex-col gap-6 px-4">
        
        <!-- Total Spending Card -->
        <div class="bg-neutral-800 rounded-2xl p-6 flex flex-col items-center justify-center relative overflow-hidden">
            <div class="absolute inset-0 bg-gradient-to-br from-green-500/10 to-teal-500/10 z-0"></div>
            <div class="relative z-10 text-center">
                <span class="text-neutral-400 text-sm mb-1 block">本月總支出 ({{ currentLedger?.base_currency || 'TWD' }})</span>
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
            <div class="flex flex-col gap-3 pb-24">
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
import { ref, onMounted, watch } from 'vue'
import { Icon } from '@iconify/vue'
import { useSpace } from '~/composables/useSpace'
import { useApi } from '~/composables/useApi'

definePageMeta({
  layout: 'main'
})

const router = useRouter()
const api = useApi()
const { currentSpace: currentLedger, currentSpaceId: currentLedgerId, fetchSpaces: fetchLedgers } = useSpace()

const loading = ref(true)
const totalSpent = ref(0)
const categories = ref<any[]>([])
const weeklyData = ref<number[]>([])

const formatNumber = (num: number) => {
    return num.toLocaleString()
}

// Watch for ledger selection changes to reload stats
watch(currentLedgerId, () => {
  if (currentLedgerId.value) {
    fetchData()
  }
})

// Generator for visualization based on transactions
const processStats = (transactions: any[]) => {
    // Basic calculation for now
    let total = 0
    const catMap: Record<string, number> = {}
    
    transactions.forEach(t => {
        const amount = Number(t.billing_amount || t.total_amount)
        total += amount
        catMap[t.category || '未分類'] = (catMap[t.category || '未分類'] || 0) + amount
    })

    totalSpent.value = total
    
    const catIcons: Record<string, string> = {
        '餐飲': '🍔', '交通': '🚆', '購物': '🛍️', '生活': '🏠', '娛樂': '🎮', '食物': '🍔'
    }

    categories.value = Object.entries(catMap).map(([name, amount]) => ({
        name,
        amount,
        percentage: total > 0 ? Math.round((amount / total) * 100) : 0,
        icon: catIcons[name] || '📦'
    })).sort((a,b) => b.amount - a.amount)

    // Mock weekly spending relative percentage
    weeklyData.value = [30, 60, 45, 80]
}

const fetchData = async () => {
    if (!currentLedger.value) return

    loading.value = true
    try {
        const txns = await api.get<any[]>(`/api/spaces/${currentLedger.value.id}/transactions`)
        processStats(txns)
    } catch (e) {
        console.error(e)
    } finally {
        loading.value = false
    }
}

onMounted(async () => {
    await fetchLedgers()
    fetchData()
})
</script>
