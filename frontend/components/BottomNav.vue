<template>
  <div class="fixed bottom-0 left-0 right-0 h-16 bg-neutral-900 border-t border-neutral-800 flex justify-around items-center z-50 shadow-lg px-2">
    <NuxtLink 
        v-for="item in navItems" 
        :key="item.label"
        :to="item.to" 
        class="flex-1 flex flex-col items-center gap-1 text-neutral-500 no-underline text-xs py-2 px-1 transition-colors hover:text-indigo-400 active:text-indigo-400 group" 
        active-class="text-indigo-400"
    >
      <div class="relative">
          <Icon :icon="item.icon" class="text-2xl transition-transform group-active:scale-95" />
      </div>
      <span class="font-bold tracking-tight">{{ item.label }}</span>
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

const navItems = [
    { label: '首頁', icon: 'mdi:view-dashboard-outline', to: '/' },
    { label: '統計分析', icon: 'mdi:chart-bar', to: '/spaces/stats' },
    { label: '偏好設定', icon: 'mdi:cog-outline', to: '/settings' },
]
</script>
