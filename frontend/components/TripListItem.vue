<template>
  <div
    class="relative w-full h-44 rounded-3xl overflow-hidden cursor-pointer transition-all duration-300 hover:scale-[1.02] active:scale-[0.98] group shadow-lg mb-4"
    @click="$emit('click')"
  >
    <!-- Background Image -->
    <div class="absolute inset-0 bg-neutral-800">
      <img 
        v-if="coverImage" 
        :src="getImageUrl(coverImage)" 
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110" 
        :alt="trip.name"
      />
      <div v-else class="w-full h-full bg-gradient-to-br from-indigo-900 to-slate-900 flex items-center justify-center">
        <Icon icon="mdi:airplane-takeoff" class="text-6xl text-white/10" />
      </div>
      
      <!-- Overlays -->
      <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent"></div>
    </div>

    <!-- Content -->
    <div class="absolute inset-0 p-6 flex flex-col justify-end">
      <div class="flex flex-col gap-1">
        <h3 class="text-xl font-bold text-white tracking-tight drop-shadow-md">{{ trip.name }}</h3>
        
        <div class="flex items-center gap-3 text-white/80 text-xs font-medium">
          <div class="flex items-center gap-1">
            <Icon icon="mdi:calendar-range" class="text-sm" />
            <span>{{ formatDateRange(trip.start_date, trip.end_date) }}</span>
          </div>
          
          <div v-if="trip.members?.length" class="flex items-center gap-1">
            <Icon icon="mdi:account-group" class="text-sm" />
            <span>{{ trip.members.length }} 人</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Status Badge (Optional, if you want to show it's active) -->
    <div v-if="isCurrentTrip(trip.start_date, trip.end_date)" class="absolute top-4 right-4 px-3 py-1 bg-green-500/20 backdrop-blur-md border border-green-500/30 rounded-full">
      <span class="text-[10px] font-bold text-green-400 uppercase tracking-wider">進行中</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { useImages } from '~/composables/useImages'

const props = defineProps<{
  trip: any
  coverImage?: string
}>()

defineEmits(['click'])

const { getImageUrl } = useImages()

const isCurrentTrip = (start: string | null, end: string | null) => {
  if (!start || !end) return false
  const now = new Date()
  const s = new Date(start)
  const e = new Date(end)
  return now >= s && now <= e
}

const formatDateRange = (start: string | null, end: string | null) => {
  if (!start && !end) return '日期未設定'
  const opts: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric' }
  const startStr = start ? new Date(start).toLocaleDateString('zh-TW', opts) : '?'
  const endStr = end ? new Date(end).toLocaleDateString('zh-TW', opts) : '?'
  return `${startStr} - ${endStr}`
}
</script>
