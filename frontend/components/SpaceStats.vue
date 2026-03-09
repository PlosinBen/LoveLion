<template>
  <div class="space-stats flex flex-col gap-6 pb-10">
    <!-- Total Balance Card -->
    <BaseCard padding="p-6" class="flex flex-col items-center">
        <h3 class="text-xs font-bold text-neutral-500 uppercase tracking-wider mb-2">空間累積支出</h3>
        <div class="text-3xl font-bold text-white tracking-tight">
            <span class="text-base text-neutral-500 mr-1">{{ baseCurrency }}</span>
            {{ totalAmount.toLocaleString() }}
        </div>
    </BaseCard>

    <!-- Members / Split Summary -->
    <section class="flex flex-col gap-4">
        <h2 class="text-sm font-bold text-neutral-400 uppercase tracking-wider px-1">成員分攤摘要</h2>
        <div class="flex flex-col gap-3">
            <BaseCard v-for="member in stats" :key="member.user_id" class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-xl bg-neutral-800 flex items-center justify-center border border-neutral-700">
                        <Icon icon="mdi:account-outline" class="text-xl text-neutral-400" />
                    </div>
                    <div class="flex flex-col">
                        <span class="font-bold text-neutral-100 text-sm">{{ member.name || '未知使用者' }}</span>
                        <span class="text-xs text-neutral-500 font-medium mt-0.5">{{ member.percentage }}% 權重</span>
                    </div>
                </div>
                <div class="text-right">
                    <div class="font-bold text-white">{{ member.spent.toLocaleString() }}</div>
                    <div class="text-xs text-neutral-500 uppercase">{{ baseCurrency }}</div>
                </div>
            </BaseCard>
        </div>
    </section>

    <!-- Spending Categories -->
    <section class="flex flex-col gap-4">
        <h2 class="text-sm font-bold text-neutral-400 uppercase tracking-wider px-1">類別分布</h2>
        <div class="flex flex-col gap-3">
             <BaseCard v-for="cat in categoryStats" :key="cat.name">
                <div class="flex justify-between items-center mb-3">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 rounded-xl bg-neutral-800 flex items-center justify-center text-indigo-400 border border-neutral-700">
                             <Icon :icon="getCategoryIcon(cat.name)" class="text-xl" />
                        </div>
                        <span class="font-bold text-neutral-100 text-sm">{{ cat.name }}</span>
                    </div>
                    <div class="text-right">
                         <div class="font-bold text-white">{{ cat.amount.toLocaleString() }}</div>
                         <div class="text-xs text-neutral-500 uppercase">{{ baseCurrency }}</div>
                    </div>
                </div>
                <!-- Progress Bar -->
                <div class="h-1.5 bg-neutral-800 rounded-full overflow-hidden">
                    <div class="h-full bg-indigo-500 rounded-full transition-all duration-700" :style="{ width: cat.percentage + '%' }"></div>
                </div>
             </BaseCard>
        </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'
import BaseCard from '~/components/BaseCard.vue'

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
        return {
            user_id: m.user_id,
            name: name,
            spent: totalAmount.value * ((m.weight || 0) / 100),
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
