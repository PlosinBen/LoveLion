<template>
  <div
    v-if="broadcast && !isDismissed"
    class="bg-indigo-500/10 border-b border-indigo-500/20 px-4 py-2.5 flex items-center gap-3"
  >
    <Icon icon="mdi:bullhorn-outline" class="text-indigo-400 text-lg shrink-0" />
    <NuxtLink
      :to="`/announcements/${broadcast.id}`"
      class="flex-1 min-w-0 text-sm text-indigo-300 font-bold truncate no-underline hover:text-indigo-200 transition-colors"
    >
      {{ broadcast.title }}
    </NuxtLink>
    <button
      @click.stop="dismiss"
      class="shrink-0 w-7 h-7 flex items-center justify-center text-neutral-500 hover:text-neutral-300 transition-colors cursor-pointer rounded-full hover:bg-neutral-800"
    >
      <Icon icon="mdi:close" class="text-base" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import type { Announcement } from '~/types'

const api = useApi()
const broadcast = ref<Announcement | null>(null)

const DISMISS_KEY = 'lovelion_dismissed_broadcasts'

const dismissedIds = ref<string[]>([])

const isDismissed = computed(() => {
  if (!broadcast.value) return true
  return dismissedIds.value.includes(broadcast.value.id)
})

const loadDismissed = () => {
  if (typeof window === 'undefined') return
  try {
    const stored = localStorage.getItem(DISMISS_KEY)
    dismissedIds.value = stored ? JSON.parse(stored) : []
  } catch {
    dismissedIds.value = []
  }
}

const dismiss = () => {
  if (!broadcast.value) return
  dismissedIds.value.push(broadcast.value.id)
  if (typeof window !== 'undefined') {
    localStorage.setItem(DISMISS_KEY, JSON.stringify(dismissedIds.value))
  }
}

onMounted(async () => {
  loadDismissed()
  try {
    const result = await api.get<Announcement | null>('/api/announcements/broadcast')
    broadcast.value = result
  } catch {
    // Silently fail — no broadcast is fine
  }
})
</script>
