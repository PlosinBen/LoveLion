<template>
  <NuxtLink v-if="to" :to="to" class="dashboard-card group">
    <Icon :icon="icon" class="text-3xl" :class="[iconColor || 'text-indigo-500']" />
    <span>{{ label }}</span>
  </NuxtLink>
  <button v-else @click="$emit('click')" class="dashboard-card group">
    <Icon :icon="icon" class="text-3xl" :class="[iconColor || 'text-red-500']" />
    <span>{{ label }}</span>
  </button>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

defineProps<{
  to?: string
  icon: string
  label: string
  iconColor?: string
}>()

defineEmits(['click'])
</script>

<style scoped>
.dashboard-card {
  @apply flex flex-col items-center gap-2 p-5 rounded-2xl border border-neutral-800 bg-neutral-900 text-white no-underline transition-all duration-200 cursor-pointer hover:-translate-y-0.5 hover:border-indigo-500;
}
</style>
