<template>
  <div class="space-stats flex flex-col gap-8 pb-10">
    <!-- Total Balance Card -->
    <div class="bg-neutral-900 rounded-3xl p-8 border border-neutral-800 shadow-sm relative overflow-hidden flex flex-col items-center">
        <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/5 to-emerald-500/5 pointer-events-none"></div>
        <div class="relative z-10 flex flex-col items-center">
            <h3 class="text-[10px] font-black text-neutral-500 uppercase tracking-[0.2em] mb-2">空間累積支出</h3>
            <div class="text-4xl font-black text-white tracking-tighter">
                <span class="text-lg text-neutral-600 mr-1">{{ baseCurrency }}</span>
                {{ totalAmount.toLocaleString() }}
            </div>
        </div>
    </div>

    <!-- Members / Split Summary -->
    <section class="flex flex-col gap-4">
        <h2 class="text-xs font-black text-neutral-500 uppercase tracking-widest px-1">成員分攤摘要</h2>
        <div class="flex flex-col gap-3">
            <div v-for="member in stats" :key="member.user_id" class="bg-neutral-900 rounded-3xl p-5 border border-neutral-800/60 flex items-center justify-between shadow-sm">
                <div class="flex items-center gap-4">
                    <div class="w-12 h-12 rounded-2xl bg-neutral-800 flex items-center justify-center border border-neutral-700/50">
                        <Icon icon="mdi:account-outline" class="text-2xl text-neutral-400" />
                    </div>
                    <div class="flex flex-col">
                        <span class="font-bold text-neutral-100 text-sm tracking-tight">{{ member.name || '未知使用者' }}</span>
                        <span class="text-[10px] text-neutral-500 font-black uppercase tracking-wider mt-0.5">{{ member.percentage }}% 權重</span>
                    </div>
                </div>
                <div class="text-right flex flex-col">
                    <span class="font-black text-white tracking-tight">{{ member.spent.toLocaleString() }}</span>
                    <span class="text-[10px] text-neutral-600 font-bold uppercase">{{ baseCurrency }}</span>
                </div>
            </div>
        </div>
    </section>

    <!-- Spending Categories -->
    <section class="flex flex-col gap-4">
        <h2 class="text-xs font-black text-neutral-500 uppercase tracking-widest px-1">類別分布</h2>
        <div class="flex flex-col gap-3">
             <div v-for="cat in categoryStats" :key="cat.name" class="bg-neutral-900 rounded-3xl p-5 border border-neutral-800/60 shadow-sm">
                <div class="flex justify-between items-center mb-4">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 rounded-xl bg-neutral-800 flex items-center justify-center text-xl text-indigo-400">
                             <Icon :icon="getCategoryIcon(cat.name)" />
                        </div>
                        <span class="font-bold text-neutral-100 text-sm">{{ cat.name }}</span>
                    </div>
                    <div class="text-right flex flex-col">
                         <span class="font-black text-white tracking-tight">{{ cat.amount.toLocaleString() }}</span>
                         <span class="text-[10px] text-neutral-600 font-bold uppercase">{{ baseCurrency }}</span>
                    </div>
                </div>
                <!-- Progress Bar -->
                <div class="h-1.5 bg-neutral-800 rounded-full overflow-hidden">
                    <div class="h-full bg-indigo-500 rounded-full transition-all duration-700" :style="{ width: cat.percentage + '%' }"></div>
                </div>
             </div>
        </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps<{
  transactions: any[]
  members: any[]
  baseCurrency: string
}>()

const totalAmount = computed(() => {
    return props.transactions.reduce((sum, t) => sum + (Number(t.billing_amount) || Number(t.total_amount) || 0), 0)
})

const stats = computed(() => {
    return (props.members || []).map(m => {
        const name = m.alias || m.user?.display_name || m.user?.username || '未知使用者'
        // Simplified mock logic for spent breakdown per member if not provided by backend
        return {
            user_id: m.user_id,
            name: name,
            spent: totalAmount.value * (m.weight / 100),
            percentage: m.weight || 0
        }
    })
})

const categoryStats = computed(() => {
    const cats: Record<string, number> = {}
    props.transactions.forEach(t => {
        const name = t.category || '其他'
        const amount = Number(t.billing_amount) || Number(t.total_amount) || 0
        cats[name] = (cats[name] || 0) + amount
    })

    return Object.entries(cats).map(([name, amount]) => ({
        name,
        amount,
        percentage: totalAmount.value > 0 ? Math.round((amount / totalAmount.value) * 100) : 0
    })).sort((a,b) => b.amount - a.amount)
})

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
</script>
