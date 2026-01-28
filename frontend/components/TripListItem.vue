<template>
  <div
    class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 cursor-pointer transition-all duration-200 hover:-translate-y-0.5 hover:border-indigo-500"
    @click="$emit('click')"
  >
    <div class="flex items-center gap-3">
      <div class="relative flex items-center justify-center w-12 h-12 rounded-xl bg-indigo-500/20 overflow-hidden">
         <img v-if="coverImage" :src="getImageUrl(coverImage)" class="w-full h-full object-cover" />
         <Icon v-else icon="mdi:airplane" class="text-2xl text-indigo-500" />
      </div>
      <div class="flex-1">
        <h3 class="text-base font-bold mb-1">{{ trip.name }}</h3>
        <p class="text-xs text-neutral-400">
          {{ formatDateRange(trip.start_date, trip.end_date) }}
          <span v-if="trip.members?.length" class="ml-2">• {{ trip.members.length }} 位成員</span>
        </p>
      </div>
      <Icon icon="mdi:chevron-right" class="text-xl text-neutral-500" />
    </div>
    <p v-if="trip.description" class="mt-3 pt-3 border-t border-neutral-800 text-sm text-neutral-400 line-clamp-2">
      {{ trip.description }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { useImages } from '~/composables/useImages'

defineProps<{
  trip: any
  coverImage?: string
}>()

defineEmits(['click'])

const { getImageUrl } = useImages()

const formatDateRange = (start: string | null, end: string | null) => {
  if (!start && !end) return '日期未設定'
  const opts: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric' }
  const startStr = start ? new Date(start).toLocaleDateString('zh-TW', opts) : '?'
  const endStr = end ? new Date(end).toLocaleDateString('zh-TW', opts) : '?'
  return `${startStr} - ${endStr}`
}
</script>
