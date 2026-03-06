<template>
  <nav class="fixed bottom-0 left-0 right-0 bg-neutral-900 border-t border-neutral-800 flex justify-around items-center px-4 py-3 z-50">
    <template v-if="!items">
      <NuxtLink 
        v-for="item in defaultNavItems" 
        :key="item.to"
        :to="item.to"
        class="flex flex-col items-center gap-1 no-underline"
        :class="isCurrent(item.to) ? 'text-indigo-500' : 'text-neutral-500'"
      >
        <Icon 
          :icon="isCurrent(item.to) ? item.activeIcon : item.icon" 
          class="text-2xl"
        />
        <span class="text-xs font-medium">
          {{ item.label }}
        </span>
      </NuxtLink>
    </template>

    <template v-else>
      <button 
        v-for="item in items" 
        :key="item.id"
        @click="$emit('update:modelValue', item.id)"
        class="flex flex-col items-center gap-1 bg-transparent border-0 cursor-pointer"
        :class="modelValue === item.id ? 'text-indigo-500' : 'text-neutral-500'"
      >
        <Icon 
          :icon="modelValue === item.id ? (item.activeIcon || item.icon) : item.icon" 
          class="text-2xl"
        />
        <span class="text-xs font-medium">
          {{ item.label }}
        </span>
      </button>
    </template>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'
import { useRoute } from 'vue-router'
import { useSpace } from '~/composables/useSpace'

interface NavItem {
  id?: string
  label: string
  icon: string
  activeIcon?: string
  to?: string
}

const props = defineProps<{
  modelValue?: string
  items?: NavItem[]
}>()

defineEmits(['update:modelValue'])

const route = useRoute()
const { currentSpaceId } = useSpace()

const defaultNavItems = computed(() => {
  const base = currentSpaceId.value ? `/spaces/${currentSpaceId.value}` : '/'
  return [
    { label: '我的空間', icon: 'mdi:view-grid-outline', activeIcon: 'mdi:view-grid', to: '/' },
    { label: '目前活動', icon: 'mdi:wallet-outline', activeIcon: 'mdi:wallet', to: base },
    { label: '數據統計', icon: 'mdi:chart-bar', activeIcon: 'mdi:chart-bar', to: '/spaces/stats' }
  ]
})

const isCurrent = (path: string) => {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}
</script>
