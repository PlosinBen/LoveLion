<template>
  <div class="flex items-center gap-5 p-6 rounded-3xl border border-neutral-800 bg-neutral-900/50 backdrop-blur-sm transition-all duration-500 shadow-sm" :class="statusBorderClass">
    <div class="w-12 h-12 rounded-2xl flex items-center justify-center text-2xl" :class="statusIconBg">
      <Icon :icon="statusIcon" :class="statusIconColor" />
    </div>
    <div class="flex flex-col">
      <h3 class="text-[10px] font-black text-neutral-500 uppercase tracking-[0.2em] mb-1">系統服務狀態</h3>
      <p class="text-sm font-bold text-white leading-none">{{ message }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps<{
  status: 'checking' | 'online' | 'offline'
  message: string
}>()

const statusBorderClass = computed(() => {
  if (props.status === 'online') return 'border-emerald-500/20 shadow-emerald-500/5'
  if (props.status === 'offline') return 'border-red-500/20 shadow-red-500/5'
  return 'border-neutral-800'
})

const statusIconBg = computed(() => {
  if (props.status === 'online') return 'bg-emerald-500/10 border border-emerald-500/20'
  if (props.status === 'offline') return 'bg-red-500/10 border border-red-500/20'
  return 'bg-neutral-800 border border-neutral-700/50'
})

const statusIconColor = computed(() => {
  if (props.status === 'online') return 'text-emerald-500'
  if (props.status === 'offline') return 'text-red-500'
  return 'text-neutral-500 animate-pulse'
})

const statusIcon = computed(() => {
  switch (props.status) {
    case 'online': return 'mdi:check-circle-outline'
    case 'offline': return 'mdi:alert-circle-outline'
    default: return 'mdi:progress-clock'
  }
})
</script>
