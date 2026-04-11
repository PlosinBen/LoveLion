<template>
  <BaseCard
    @click="router.push(`/spaces/${spaceId}/ledger/transaction/${transaction.id}`)"
    class="flex justify-between items-center hover:bg-neutral-800 transition-colors cursor-pointer group active:scale-95 shadow-sm"
  >
    <div class="flex items-center gap-3">
      <!-- Thumbnail or Type Icon -->
      <div v-if="thumbnail" class="w-10 h-10 rounded-xl overflow-hidden shrink-0">
        <img :src="thumbnail" class="w-full h-full object-cover" />
      </div>
      <div v-else-if="transaction.type === 'payment'" class="w-10 h-10 rounded-xl bg-emerald-500/10 flex items-center justify-center text-emerald-400 border border-emerald-500/20 shrink-0">
        <Icon icon="mdi:swap-horizontal" class="text-xl" />
      </div>
      <div v-else class="w-10 h-10 rounded-xl bg-neutral-800 flex items-center justify-center text-neutral-600 border border-neutral-700 shrink-0">
        <Icon icon="mdi:image-off-outline" class="text-xl" />
      </div>

      <div class="flex flex-col min-w-0">
        <div class="flex items-center gap-1.5">
          <h4 class="text-sm font-semibold text-neutral-100 truncate">{{ transaction.title || '未命名' }}</h4>
          <!-- Badge hides itself for completed / NULL statuses, so a normal
               transaction row looks identical to before. -->
          <AiStatusBadge :status="transaction.ai_status" :error="transaction.ai_error" />
        </div>
        <p class="text-xs text-neutral-500 mt-0.5">{{ formatDate(transaction.date) }}</p>
      </div>
    </div>

    <div class="text-right">
      <div class="text-xs text-neutral-500 uppercase font-bold leading-none mb-1">
        {{ transaction.currency }}
      </div>
      <div class="text-lg font-bold" :class="transaction.type === 'payment' ? 'text-emerald-400' : 'text-indigo-500'">
        {{ formatAmount(transaction.total_amount) }}
      </div>
    </div>
  </BaseCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'
import { useRouter } from 'vue-router'
import BaseCard from '~/components/BaseCard.vue'
import AiStatusBadge from '~/components/AiStatusBadge.vue'
import type { Transaction } from '~/types'

const props = defineProps<{
  transaction: Transaction
  spaceId: string | number
}>()

const router = useRouter()

const thumbnail = computed(() => {
  const images = props.transaction.images
  if (images && images.length > 0 && images[0]) {
    return images[0].file_path
  }
  return null
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', { month: 'short', day: 'numeric' })
}

const formatAmount = (amount: string | number | undefined) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : (amount || 0)
  return num.toLocaleString('zh-TW')
}
</script>
