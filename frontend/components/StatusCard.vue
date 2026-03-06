<template>
  <div class="flex items-center gap-4 p-4 rounded-xl border border-neutral-800 bg-neutral-900 transition-all shadow-sm" :class="statusBorderClass">
    <div class="text-2xl" :class="statusIconColor">
      <Icon :icon="statusIcon" />
    </div>
    <div class="flex flex-col">
      <h3 class="text-xs font-bold text-neutral-500 uppercase tracking-wider">系統狀態</h3>
      <p class="text-sm font-medium text-white">{{ message }}</p>
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
  if (props.status === 'online') return 'border-emerald-500/20'
  if (props.status === 'offline') return 'border-red-500/20'
  return 'border-neutral-800'
})

const statusIconColor = computed(() => {
  if (props.status === 'online') return 'text-emerald-500'
  if (props.status === 'offline') return 'text-red-500'
  return 'text-neutral-500'
})

const statusIcon = computed(() => {
  switch (props.status) {
    case 'online': return 'mdi:check-circle'
    case 'offline': return 'mdi:alert-circle'
    default: return 'mdi:loading'
  }
})
</script>
