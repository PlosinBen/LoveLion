<template>
  <div class="flex items-center gap-4 mb-6 p-5 rounded-2xl border bg-neutral-900 duration-200" :class="statusClasses">
    <div class="text-3xl" :class="statusIconColor">
      <Icon :icon="statusIcon" />
    </div>
    <div class="flex flex-col">
      <h3 class="text-sm text-neutral-400">系統狀態</h3>
      <p class="text-base mt-1">{{ message }}</p>
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

const statusClasses = computed(() => {
  if (props.status === 'online') return 'border-green-500'
  if (props.status === 'offline') return 'border-red-500'
  return 'border-neutral-800'
})

const statusIconColor = computed(() => {
  if (props.status === 'online') return 'text-green-500'
  if (props.status === 'offline') return 'text-red-500'
  return 'text-neutral-500'
})

const statusIcon = computed(() => {
  switch (props.status) {
    case 'online': return 'mdi:check-circle'
    case 'offline': return 'mdi:close-circle'
    default: return 'mdi:loading'
  }
})
</script>
