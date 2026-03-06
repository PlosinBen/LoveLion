<template>
  <div class="flex justify-between items-center p-5 rounded-3xl border border-neutral-800/60 bg-neutral-900 hover:bg-neutral-800/50 transition-all cursor-pointer group active:scale-[0.98] shadow-sm">
    <div class="flex items-center gap-4">
      <!-- Category Icon -->
      <div class="w-12 h-12 rounded-2xl bg-indigo-500/10 flex items-center justify-center border border-indigo-500/10">
        <Icon :icon="getCategoryIcon(transaction.category)" class="text-2xl text-indigo-400" />
      </div>
      
      <div class="flex flex-col">
        <h4 class="text-sm font-bold text-neutral-100 tracking-tight">{{ transaction.category || '其他' }}</h4>
        <p class="text-[10px] text-neutral-500 font-black uppercase tracking-widest mt-0.5">{{ formatDate(transaction.date) }}</p>
      </div>
    </div>

    <div class="text-right flex flex-col items-end">
      <div class="text-[10px] text-neutral-600 font-black uppercase tracking-widest mb-0.5 leading-none">
        {{ isBaseCurrency ? (baseCurrency || 'TWD') : transaction.currency }}
      </div>
      <div class="text-lg font-black text-indigo-400 tracking-tighter">
        {{ formatAmount(displayAmount) }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps<{
  transaction: {
    category?: string
    date: string
    currency: string
    total_amount: number | string
    billing_amount?: number | string
  },
  baseCurrency?: string
}>()

const isBaseCurrency = computed(() => {
  return props.transaction.billing_amount && Number(props.transaction.billing_amount) > 0
})

const displayAmount = computed(() => {
  return isBaseCurrency.value ? props.transaction.billing_amount : props.transaction.total_amount
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', { month: 'short', day: 'numeric' })
}

const formatAmount = (amount: string | number | undefined) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : (amount || 0)
  return num.toLocaleString('zh-TW', { maximumFractionDigits: 0 })
}

const getCategoryIcon = (category?: string) => {
  const icons: Record<string, string> = {
    '餐飲': 'mdi:food',
    '交通': 'mdi:train-car',
    '購物': 'mdi:shopping',
    '娛樂': 'mdi:movie',
    '生活': 'mdi:home',
    '其他': 'mdi:receipt'
  }
  return icons[category || ''] || 'mdi:receipt'
}
</script>
