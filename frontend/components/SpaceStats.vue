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

    <!-- Settlement Summary (應收應付) -->
    <section v-if="settlement.items.length > 0" class="flex flex-col gap-4">
        <h2 class="text-sm font-bold text-neutral-400 uppercase tracking-wider px-1">應收應付</h2>
        <BaseCard padding="p-0" class="divide-y divide-neutral-800">
            <div v-for="person in settlement.items" :key="person.name" class="flex justify-between items-center px-5 py-4">
                <span class="font-bold text-neutral-100 text-sm">{{ person.name }}</span>
                <div class="text-right">
                    <div class="font-bold text-sm" :class="person.baseAmount >= 0 ? 'text-indigo-400' : 'text-red-500'">
                        {{ person.baseAmount >= 0 ? '+' : '-' }}{{ baseCurrency }} {{ Math.abs(person.baseAmount).toLocaleString() }}
                    </div>
                    <div v-for="(amount, currency) in person.foreignAmounts" :key="currency" class="text-xs text-neutral-500 font-medium">
                        {{ amount >= 0 ? '+' : '-' }}{{ currency }} {{ Math.abs(amount).toLocaleString() }}
                    </div>
                </div>
            </div>
            <div v-if="settlement.unsettledCount > 0" class="px-5 py-3 flex items-center gap-1 text-xs text-neutral-500 font-medium">
                <Icon icon="mdi:information-outline" class="text-sm" />
                有 {{ settlement.unsettledCount }} 筆外幣拆帳交易尚未結算
            </div>
        </BaseCard>
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
import type { Transaction } from '~/types'

const props = defineProps<{
  transactions: Transaction[]
  baseCurrency: string
}>()

// Only count expenses for total spending (exclude payments)
const expenseTransactions = computed(() => props.transactions.filter(t => t.type === 'expense'))

const totalAmount = computed(() => {
    return expenseTransactions.value.reduce((sum, t) => {
        const billing = t.expense?.billing_amount ? Number(t.expense.billing_amount) : 0
        return sum + (billing || Number(t.total_amount) || 0)
    }, 0)
})

const settlement = computed(() => {
    const baseNet: Record<string, number> = {}
    const foreignNet: Record<string, Record<string, number>> = {}
    let unsettledCount = 0

    for (const t of props.transactions) {
        if (!t.debts || t.debts.length === 0) continue

        for (const debt of t.debts) {
            const settledAmount = Number(debt.settled_amount) || 0

            if (settledAmount > 0) {
                // Settled debt: use settled_amount (base currency)
                baseNet[debt.payer_name] = (baseNet[debt.payer_name] || 0) - settledAmount
                baseNet[debt.payee_name] = (baseNet[debt.payee_name] || 0) + settledAmount
            } else if (!debt.is_spot_paid) {
                // Unsettled foreign currency debt
                unsettledCount++
                const payerForeign = foreignNet[debt.payer_name] ??= {}
                const payeeForeign = foreignNet[debt.payee_name] ??= {}
                payerForeign[t.currency] = (payerForeign[t.currency] || 0) - Number(debt.amount)
                payeeForeign[t.currency] = (payeeForeign[t.currency] || 0) + Number(debt.amount)
            }
            // is_spot_paid debts are excluded from receivables/payables
        }
    }

    const allNames = new Set([...Object.keys(baseNet), ...Object.keys(foreignNet)])
    const items = Array.from(allNames).map(name => ({
        name,
        baseAmount: baseNet[name] || 0,
        foreignAmounts: Object.fromEntries(
            Object.entries(foreignNet[name] || {}).filter(([, v]) => v !== 0)
        )
    })).sort((a, b) => b.baseAmount - a.baseAmount)

    return { items, unsettledCount }
})

const categoryStats = computed(() => {
    const cats: Record<string, number> = {}
    expenseTransactions.value.forEach(t => {
        const name = t.expense?.category || '其他'
        const billing = t.expense?.billing_amount ? Number(t.expense.billing_amount) : 0
        const amount = billing || Number(t.total_amount) || 0
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
