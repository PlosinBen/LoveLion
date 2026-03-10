<template>
  <NuxtLink 
    v-if="to" 
    :to="to" 
    class="flex flex-col items-center gap-2 p-5 rounded-2xl border border-neutral-800 bg-neutral-900 text-white no-underline transition-all duration-200 cursor-pointer hover:-translate-y-0.5 hover:border-indigo-500/50 group active:scale-95 shadow-sm"
  >
    <Icon :icon="icon" class="text-3xl" :class="[iconColor || 'text-indigo-500']" />
    <span class="text-xs font-black uppercase tracking-widest">{{ label }}</span>
  </NuxtLink>
  
  <button
    v-else
    @click="$emit('click')"
    class="flex flex-col items-center gap-2 p-5 rounded-2xl border border-neutral-800 bg-neutral-800 text-neutral-400 hover:text-white hover:bg-neutral-700 hover:border-red-500/50 no-underline transition-all duration-200 cursor-pointer active:scale-95 shadow-sm"
  >
    <Icon :icon="icon" class="text-3xl" :class="[iconColor || 'text-red-500']" />
    <span class="text-xs font-black uppercase tracking-widest">{{ label }}</span>
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
