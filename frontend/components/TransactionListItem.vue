<template>
  <BaseCard
    @click="router.push(`/spaces/${spaceId}/ledger/transaction/${transaction.id}`)"
    class="flex justify-between items-center hover:bg-neutral-800 transition-colors cursor-pointer group active:scale-95 shadow-sm"
  >
    <div class="flex items-center gap-3">
      <!-- Thumbnail or Category Icon -->
      <div v-if="thumbnail" class="w-10 h-10 rounded-xl overflow-hidden shrink-0">
        <img :src="thumbnail" class="w-full h-full object-cover" />
      </div>
      <div v-else class="w-10 h-10 rounded-xl bg-neutral-800 flex items-center justify-center text-indigo-500 border border-neutral-700 shrink-0">
        <Icon :icon="getCategoryIcon(transaction.category)" class="text-xl" />
      </div>
      
      <div class="flex flex-col">
        <h4 class="text-sm font-semibold text-neutral-100">{{ transaction.title || transaction.category || '未分類' }}</h4>
        <p class="text-xs text-neutral-500 mt-0.5">{{ formatDate(transaction.date) }}</p>
      </div>
    </div>

    <div class="text-right">
      <div class="text-xs text-neutral-500 uppercase font-bold leading-none mb-1">
        {{ isBaseCurrency ? (baseCurrency || 'TWD') : transaction.currency }}
      </div>
      <div class="text-lg font-bold text-indigo-500">
        {{ formatAmount(displayAmount) }}
      </div>
    </div>
  </BaseCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'
import { useRouter } from 'vue-router'
import BaseCard from '~/components/BaseCard.vue'

const props = defineProps<{
  transaction: {
    id: string | number
    title?: string
    category?: string
    date: string
    currency: string
    total_amount: number | string
    billing_amount?: number | string
    images?: { id: string; file_path: string }[]
  },
  spaceId: string | number,
  baseCurrency?: string
}>()

const router = useRouter()

const thumbnail = computed(() => {
  const images = props.transaction.images
  if (images && images.length > 0 && images[0]) {
    return images[0].file_path
  }
  return null
})

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
  return num.toLocaleString('zh-TW')
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
