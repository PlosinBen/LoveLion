<template>
  <div class="flex flex-col pb-24">
    <SpaceHeader
      title="統計分析"
      :show-back="false"
      class="pt-0 px-2"
    />

    <div v-if="loading" class="flex justify-center items-center py-20 text-neutral-400">
        <Icon icon="mdi:loading" class="text-4xl animate-spin" />
    </div>

    <div v-else class="flex flex-col gap-8 pt-0 animate-in fade-in duration-500">
        
        <!-- Summary Header -->
        <div class="bg-neutral-900 rounded-2xl border border-neutral-800/60 p-8 flex flex-col items-center justify-center relative overflow-hidden shadow-sm">
            <div class="relative z-10 text-center">
                <span class="text-xs font-bold text-neutral-500 uppercase tracking-widest mb-2 block">當前空間總支出</span>
                <div class="text-4xl font-bold text-white tracking-tighter">
                    <span class="text-lg text-neutral-600 mr-1">{{ currentSpace?.base_currency || 'TWD' }}</span>
                    {{ formatNumber(totalSpent) }}
                </div>
            </div>
        </div>

        <!-- Category Breakdown -->
        <section class="flex flex-col gap-4">
            <div class="flex items-center justify-between px-1">
                <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest">支出類別佔比</h2>
                <Icon icon="mdi:chart-donut" class="text-indigo-500 text-xl" />
            </div>
            
            <div v-if="categories.length === 0" class="text-center py-16 bg-neutral-900 rounded-2xl border border-neutral-800/50 border-dashed">
                <p class="text-neutral-600 text-sm font-medium">尚無交易數據</p>
            </div>

            <div v-else class="flex flex-col gap-3">
                <div v-for="cat in categories" :key="cat.name" class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800/60 hover:border-indigo-500/30 transition-colors shadow-sm">
                     <div class="flex justify-between items-center mb-4">
                         <div class="flex items-center gap-3">
                             <div class="w-10 h-10 rounded-xl bg-neutral-800 flex items-center justify-center text-xl text-indigo-400">
                                 <Icon :icon="cat.icon" />
                             </div>
                             <div class="flex flex-col">
                                 <span class="font-bold text-neutral-100 text-sm">{{ cat.name }}</span>
                                 <span class="text-xs text-neutral-500 font-bold uppercase tracking-wider">{{ cat.percentage }}%</span>
                             </div>
                         </div>
                         <div class="text-right flex flex-col">
                             <span class="font-bold text-white tracking-tight">{{ formatNumber(cat.amount) }}</span>
                             <span class="text-xs text-neutral-600 font-bold uppercase">{{ currentSpace?.base_currency }}</span>
                         </div>
                     </div>
                     <!-- Progress Bar -->
                     <div class="h-1.5 bg-neutral-800 rounded-full overflow-hidden">
                         <div class="h-full bg-indigo-500 rounded-full transition-all duration-1000" :style="{ width: cat.percentage + '%' }"></div>
                     </div>
                </div>
            </div>
        </section>

        <!-- Weekly Activity (Bar Chart Mock) -->
        <section class="flex flex-col gap-4">
            <h2 class="text-xs font-bold text-neutral-500 uppercase tracking-widest px-1">支出週趨勢</h2>
            <div class="bg-neutral-900 rounded-2xl border border-neutral-800/60 p-6 shadow-sm">
                <div class="flex justify-between items-end h-32 px-2 gap-3">
                    <div v-for="(val, index) in weeklyData" :key="index" class="flex-1 flex flex-col items-center gap-3 group">
                        <div class="w-full bg-neutral-800/50 rounded-xl relative group-hover:bg-neutral-800 transition-colors overflow-hidden h-full">
                            <div class="absolute bottom-0 w-full bg-indigo-500/40 rounded-xl group-hover:bg-indigo-500 transition-all" :style="{ height: val + '%' }"></div>
                        </div>
                        <span class="text-xs font-bold text-neutral-600 uppercase tracking-widest">W{{ index + 1 }}</span>
                    </div>
                </div>
            </div>
        </section>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Icon } from '@iconify/vue'
import { useSpace } from '~/composables/useSpace'
import { useApi } from '~/composables/useApi'
import SpaceHeader from '~/components/SpaceHeader.vue'

definePageMeta({
  layout: 'default'
})

const api = useApi()
const { currentSpace, currentSpaceId, fetchSpaces } = useSpace()

const loading = ref(true)
const totalSpent = ref(0)
const categories = ref<any[]>([])
const weeklyData = ref<number[]>([30, 60, 45, 80])

const formatNumber = (num: number) => {
    return num.toLocaleString('zh-TW', { maximumFractionDigits: 0 })
}

const getCategoryIcon = (category: string) => {
  const icons: Record<string, string> = {
    '餐飲': 'mdi:food',
    '交通': 'mdi:train-car',
    '購物': 'mdi:shopping',
    '娛樂': 'mdi:movie',
    '生活': 'mdi:home',
    '其他': 'mdi:receipt'
  }
  return icons[category] || 'mdi:receipt'
}

watch(currentSpaceId, () => {
  if (currentSpaceId.value) {
    fetchData()
  }
})

const processStats = (transactions: any[]) => {
    let total = 0
    const catMap: Record<string, number> = {}
    
    transactions.forEach(t => {
        const amount = Number(t.billing_amount && Number(t.billing_amount) > 0 ? t.billing_amount : t.total_amount)
        total += amount
        const cat = t.category || '其他'
        catMap[cat] = (catMap[cat] || 0) + amount
    })

    totalSpent.value = total
    
    categories.value = Object.entries(catMap).map(([name, amount]) => ({
        name,
        amount,
        percentage: total > 0 ? Math.round((amount / total) * 100) : 0,
        icon: getCategoryIcon(name)
    })).sort((a,b) => b.amount - a.amount)
}

const fetchData = async () => {
    if (!currentSpace.value) return

    loading.value = true
    try {
        const txns = await api.get<any[]>(`/api/spaces/${currentSpace.value.id}/transactions`)
        processStats(txns || [])
    } catch (e) {
        console.error('Failed to fetch stats:', e)
    } finally {
        loading.value = false
    }
}

onMounted(async () => {
    await fetchSpaces()
    fetchData()
})
</script>
