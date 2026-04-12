<template>
  <span
    v-if="visible"
    :title="tooltip"
    class="inline-flex items-center justify-center rounded-full"
    :class="wrapperClass"
    :aria-label="tooltip"
  >
    <Icon :icon="icon" :class="iconClass" />
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'
import type { AiStatus } from '~/types'

const props = withDefaults(
  defineProps<{
    status?: AiStatus | null
    error?: string
    /** Size preset — 'sm' for list rows (default), 'md' for inline banners. */
    size?: 'sm' | 'md'
  }>(),
  { size: 'sm' },
)

// `completed` and null/undefined hide the badge entirely — a successful AI run
// should look identical to a manually-entered expense in the list.
const visible = computed(
  () => props.status === 'pending' || props.status === 'processing' || props.status === 'failed',
)

const wrapperClass = computed(() => {
  const base = props.size === 'md' ? 'w-8 h-8' : 'w-5 h-5'
  if (props.status === 'failed') return `${base} bg-red-500/10 border border-red-500/30`
  if (props.status === 'processing') return `${base} bg-indigo-500/10 border border-indigo-500/30`
  // pending
  return `${base} bg-neutral-700/40 border border-neutral-600`
})

const iconClass = computed(() => {
  const size = props.size === 'md' ? 'text-lg' : 'text-xs'
  if (props.status === 'failed') return `${size} text-red-400`
  if (props.status === 'processing') return `${size} text-indigo-400 animate-spin`
  return `${size} text-neutral-400`
})

const icon = computed(() => {
  if (props.status === 'failed') return 'mdi:alert-circle-outline'
  if (props.status === 'processing') return 'mdi:loading'
  // pending
  return 'mdi:robot-outline'
})

const tooltip = computed(() => {
  if (props.status === 'failed') return props.error || 'AI 辨識失敗'
  if (props.status === 'processing') return 'AI 辨識中...'
  if (props.status === 'pending') return '等待 AI 辨識'
  return ''
})
</script>
