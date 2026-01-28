<template>
  <div class="flex justify-between items-center p-4 rounded-xl border border-neutral-800 bg-neutral-900">
    <div class="flex flex-col">
      <h4 class="text-sm font-medium">{{ transaction.category || '未分類' }}</h4>
      <p class="text-xs text-neutral-400 mt-1">{{ formatDate(transaction.date) }}</p>
    </div>
    <div class="font-semibold text-indigo-500">
      {{ transaction.currency }} {{ transaction.total_amount }}
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
  }
}>()

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', { month: 'short', day: 'numeric' })
}
</script>
