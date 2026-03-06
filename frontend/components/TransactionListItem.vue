<template>
  <div class="flex justify-between items-center p-4 rounded-xl border border-neutral-800 bg-neutral-900">
    <div class="flex flex-col">
      <h4 class="text-sm font-medium">{{ transaction.category || '未分類' }}</h4>
      <p class="text-xs text-neutral-400 mt-1">{{ formatDate(transaction.date) }}</p>
    </div>
    <div class="font-semibold text-indigo-500 text-right">
      <div class="text-xs text-neutral-500 leading-none mb-1">
        {{ (transaction.billing_amount && Number(transaction.billing_amount) > 0) ? (baseCurrency || 'TWD') : transaction.currency }}
      </div>
      <div>
        {{ formatAmount((transaction.billing_amount && Number(transaction.billing_amount) > 0) ? transaction.billing_amount : transaction.total_amount) }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  transaction: {
    category?: string
    date: string
    currency: string
    total_amount: number | string
    billing_amount?: number | string
  },
  baseCurrency?: string
}>()

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', { month: 'short', day: 'numeric' })
}

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toLocaleString('zh-TW')
}
</script>
