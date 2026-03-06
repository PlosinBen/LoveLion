<template>
  <div class="text-center py-20 px-8 bg-neutral-900 rounded-[2.5rem] border border-neutral-800/60 shadow-inner">
    <div class="w-20 h-20 bg-indigo-500/5 rounded-full flex items-center justify-center mx-auto mb-6">
      <Icon :icon="icon" class="text-5xl text-indigo-500/80" />
    </div>
    <h2 class="text-xl font-black mb-2 text-white tracking-tight">{{ title }}</h2>
    <p class="text-neutral-500 text-sm font-medium mb-8 leading-relaxed max-w-xs mx-auto">{{ description }}</p>
    
    <NuxtLink v-if="actionLink" :to="actionLink" class="inline-flex items-center gap-2 px-8 py-4 rounded-2xl font-black bg-indigo-500 text-white hover:bg-indigo-600 transition-all active:scale-[0.98] shadow-lg shadow-indigo-500/20 text-sm no-underline border-0">
      <Icon icon="mdi:plus" />
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
