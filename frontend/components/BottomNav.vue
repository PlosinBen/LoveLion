<template>
  <nav class="fixed bottom-0 left-0 right-0 bg-neutral-900/90 backdrop-blur-xl border-t border-neutral-800 px-6 pt-3 pb-10 z-50 flex justify-around items-center shadow-2xl">
    <NuxtLink 
      v-for="item in navItems" 
      :key="item.to"
      :to="item.to"
      class="flex flex-col items-center gap-1.5 no-underline transition-all active:scale-90 relative"
      :class="isCurrent(item.to) ? 'text-indigo-400' : 'text-neutral-500'"
    >
      <!-- Active Glow Background -->
      <div 
        v-if="isCurrent(item.to)"
        class="absolute -top-1 w-12 h-8 bg-indigo-500/10 blur-xl rounded-full"
      ></div>

      <Icon 
        :icon="isCurrent(item.to) ? item.activeIcon : item.icon" 
        class="text-2xl transition-transform"
        :class="{ 'scale-110': isCurrent(item.to) }"
      />
      <span class="text-[10px] font-black uppercase tracking-[0.2em]">
        {{ item.label }}
      </span>

      <!-- Indicator Dot -->
      <div 
        v-if="isCurrent(item.to)"
        class="w-1 h-1 rounded-full bg-indigo-500 mt-0.5"
      ></div>
    </NuxtLink>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'
import { useRoute } from 'vue-router'
import { useSpace } from '~/composables/useSpace'

const route = useRoute()
const { currentSpaceId } = useSpace()

const navItems = computed(() => {
  const base = currentSpaceId.value ? `/spaces/${currentSpaceId.value}` : '/spaces'
  
  return [
    { 
      label: '我的空間', 
      icon: 'mdi:view-grid-outline', 
      activeIcon: 'mdi:view-grid',
      to: '/' 
    },
    { 
      label: '目前活動', 
      icon: 'mdi:wallet-outline', 
      activeIcon: 'mdi:wallet',
      to: base 
    },
    { 
      label: '數據統計', 
      icon: 'mdi:chart-bar', 
      activeIcon: 'mdi:chart-bar',
      to: '/spaces/stats' 
    }
  ]
})

const isCurrent = (path: string) => {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}
</script>
