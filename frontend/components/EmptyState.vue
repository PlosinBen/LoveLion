<template>
  <div class="text-center py-16 px-5 bg-neutral-900 rounded-2xl border border-neutral-800">
    <Icon :icon="icon" class="text-6xl text-indigo-500 mb-4" />
    <h2 class="text-xl font-bold mb-2">{{ title }}</h2>
    <p class="text-neutral-400 mb-5">{{ description }}</p>
    <NuxtLink v-if="actionLink" :to="actionLink" class="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors text-sm no-underline">
      {{ actionLabel }}
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

defineProps<{
  icon: string
  title: string
  description: string
  actionLabel?: string
  actionLink?: string
}>()
</script>
