<template>
  <div class="fixed bottom-0 left-0 right-0 h-16 bg-neutral-900 border-t border-neutral-800 flex justify-around items-center z-50 shadow-lg px-2">
    <template v-for="item in navItems" :key="item.id || item.label">
      <!-- Link Mode (for Global Nav) -->
      <NuxtLink 
          v-if="item.to"
          :to="item.to" 
          class="flex-1 flex flex-col items-center gap-1 text-neutral-500 no-underline text-xs py-2 px-1 transition-colors hover:text-indigo-400 active:text-indigo-400 group" 
          active-class="text-indigo-400"
      >
        <div class="relative">
            <Icon :icon="item.icon" class="text-2xl transition-transform group-active:scale-95" />
        </div>
        <span class="font-bold tracking-tight">{{ item.label }}</span>
      </NuxtLink>

      <!-- Tab Mode (for Space Internal Nav) -->
      <button 
          v-else
          @click="$emit('update:modelValue', item.id)"
          class="flex-1 flex flex-col items-center gap-1 bg-transparent border-0 cursor-pointer py-2 px-1 transition-colors group"
          :class="modelValue === item.id ? 'text-indigo-400' : 'text-neutral-500 hover:text-neutral-300'"
      >
        <div class="relative">
            <Icon :icon="item.icon" class="text-2xl transition-transform group-active:scale-95" />
        </div>
        <span class="font-bold tracking-tight text-xs">{{ item.label }}</span>
      </button>
    </template>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  items: {
    type: Array as () => any[],
    default: null
  }
})

defineEmits(['update:modelValue'])

const defaultItems = [
    { label: '擐?', icon: 'mdi:view-dashboard-outline', to: '/' },
    { label: '蝯梯???', icon: 'mdi:chart-bar', to: '/spaces/stats' },
    { label: '?末閮剖?', icon: 'mdi:cog-outline', to: '/settings' },
]

const navItems = computed(() => props.items || defaultItems)
</script>
