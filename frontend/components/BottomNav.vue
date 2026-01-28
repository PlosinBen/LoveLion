<template>
  <div v-if="navItems.length > 0" class="fixed bottom-0 left-0 right-0 h-16 bg-neutral-900 border-t border-neutral-800 flex justify-around items-center z-50 shadow-lg px-2">
    <div 
        v-for="item in navItems" 
        :key="item.label"
        class="flex-1"
    >
        <NuxtLink 
            :to="item.to" 
            class="flex flex-col items-center gap-1 text-neutral-500 no-underline text-xs py-2 px-1 transition-colors hover:text-indigo-400 active:text-indigo-400 group" 
            :class="{ 'text-indigo-400': isActive(item.to) }"
        >
          <div class="relative">
              <Icon :icon="item.icon" class="text-2xl transition-transform group-active:scale-95" />
              <div v-if="isActive(item.to)" class="absolute -bottom-2 left-1/2 -translate-x-1/2 w-1 h-1 bg-indigo-500 rounded-full"></div>
          </div>
          <span class="font-medium">{{ item.label }}</span>
        </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const navItems = computed(() => {
    const path = route.path

    // Trip Context
    if (path.startsWith('/trips/') && route.params.id) {
        const tripId = route.params.id
        // Ensure we don't show this nav if we are creating a trip or something not specific to an ID
        // But the route pattern matches /trips/[id], so it should be fine.
        return [
            { label: '統計分析', icon: 'mdi:chart-pie', to: `/trips/${tripId}/stats` },
            { label: '商品比價', icon: 'mdi:shopping-search', to: `/trips/${tripId}/stores` },
            { label: '行程記帳', icon: 'mdi:notebook-edit', to: `/trips/${tripId}` },
        ]
    }

    // Daily Ledger Context
    if (path.startsWith('/ledger')) {
         return [
            { label: '統計分析', icon: 'mdi:chart-bar', to: `/ledger/stats` },
            { label: '日常記帳', icon: 'mdi:wallet', to: `/ledger` },
        ]
    }

    // Home / Settings / Other -> No Bottom Nav
    return []
})

const isActive = (to: string) => {
    if (to === route.path) return true
    
    // Handle nested routes activation (e.g. stores products under stores tab)
    if (to.endsWith('/stores') && route.path.includes('/products')) return true
    
    return false
}
</script>
