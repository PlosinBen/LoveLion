<template>
  <BaseCard padding="py-16 px-6" class="text-center">
    <Icon :icon="icon" class="text-5xl text-indigo-500 mb-4" />
    <h2 class="text-xl font-bold mb-2 text-white">{{ title }}</h2>
    <p class="text-neutral-400 text-sm mb-8">{{ description }}</p>
    
    <NuxtLink v-if="actionLink" :to="actionLink" class="inline-flex items-center gap-2 px-6 py-2.5 rounded-xl font-bold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors text-sm no-underline border-0">
      <Icon icon="mdi:plus" />
      {{ actionLabel }}
    </NuxtLink>
  </BaseCard>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import BaseCard from '~/components/BaseCard.vue'

defineProps<{
  icon: string
  title: string
  description: string
  actionLabel?: string
  actionLink?: string
}>()
</script>
