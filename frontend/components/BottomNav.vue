<template>
  <nav class="fixed bottom-0 left-0 right-0 bg-neutral-900 border-t border-neutral-800 flex justify-around items-center px-4 py-3 z-50">
    <template v-for="item in navItems" :key="item.to || item.id">
      <!-- Link Mode -->
      <NuxtLink 
        v-if="item.to"
        :to="item.to"
        class="flex flex-col items-center gap-1 no-underline text-neutral-500 transition-transform active:scale-95"
      >
        <Icon :icon="item.icon" class="text-2xl" />
        <span class="text-xs font-bold uppercase tracking-widest scale-90">
          {{ item.label }}
        </span>
      </NuxtLink>

      <!-- Action Mode (Button) -->
      <button 
        v-else
        @click="$emit('update:modelValue', item.id)"
        class="flex flex-col items-center gap-1 bg-transparent border-0 cursor-pointer text-neutral-500 transition-transform active:scale-95"
      >
        <Icon :icon="item.icon" class="text-2xl" />
        <span class="text-xs font-bold uppercase tracking-widest scale-90">
          {{ item.label }}
        </span>
      </button>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

interface NavItem {
  id?: string
  label: string
  icon: string
  to?: string
}

const props = defineProps<{
  modelValue?: string
  items?: NavItem[]
}>()

defineEmits(['update:modelValue'])

const navItems = computed(() => {
  if (props.items) return props.items
  return [
    { label: '空間', icon: 'mdi:view-grid-outline', to: '/' },
    { label: '設定', icon: 'mdi:cog-outline', to: '/settings' }
  ]
})
</script>
