<template>
  <nav class="fixed bottom-0 left-0 right-0 bg-neutral-900 border-t border-neutral-800 flex justify-around items-center px-4 py-3 z-50">
    <!-- Navigation Items (Static Style) -->
    <template v-if="!items">
      <NuxtLink 
        v-for="item in defaultNavItems" 
        :key="item.to"
        :to="item.to"
        class="flex flex-col items-center gap-1 no-underline text-neutral-500 transition-transform active:scale-95"
      >
        <Icon :icon="item.icon" class="text-2xl" />
        <span class="text-xs font-bold uppercase tracking-widest scale-90">
          {{ item.label }}
        </span>
      </NuxtLink>
    </template>

    <!-- Tab Selection Items (Space Local - Static Style) -->
    <template v-else>
      <button 
        v-for="item in items" 
        :key="item.id"
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
import { useSpace } from '~/composables/useSpace'

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

const { currentSpaceId } = useSpace()

const defaultNavItems = computed(() => {
  const base = currentSpaceId.value ? `/spaces/${currentSpaceId.value}` : '/'
  return [
    { label: '空間', icon: 'mdi:view-grid-outline', to: '/' },
    { label: '活動', icon: 'mdi:wallet-outline', to: base },
    { label: '數據', icon: 'mdi:chart-bar', to: '/spaces/stats' },
    { label: '設定', icon: 'mdi:cog-outline', to: '/settings' }
  ]
})
</script>
