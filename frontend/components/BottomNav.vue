<template>
  <div class="fixed bottom-0 left-0 right-0 h-16 bg-white dark:bg-neutral-800 border-t border-neutral-200 dark:border-neutral-700 flex justify-around items-center z-50 shadow-lg">
    <NuxtLink to="/" class="flex flex-col items-center gap-1 text-neutral-500 dark:text-neutral-400 no-underline text-xs py-2 px-4 transition-colors hover:text-primary active:text-primary" :class="{ 'text-primary dark:text-primary': route.path === '/' }">
      <Icon icon="mdi:home" class="text-2xl" />
      <span>首頁</span>
    </NuxtLink>
    <NuxtLink to="/ledger" class="flex flex-col items-center gap-1 text-neutral-500 dark:text-neutral-400 no-underline text-xs py-2 px-4 transition-colors hover:text-primary active:text-primary" :class="{ 'text-primary dark:text-primary': route.path.startsWith('/ledger') }">
      <Icon icon="mdi:wallet" class="text-2xl" />
      <span>記帳</span>
    </NuxtLink>
    <NuxtLink to="/trips" class="flex flex-col items-center gap-1 text-neutral-500 dark:text-neutral-400 no-underline text-xs py-2 px-4 transition-colors hover:text-primary active:text-primary" :class="{ 'text-primary dark:text-primary': route.path.startsWith('/trips') }">
      <Icon icon="mdi:airplane" class="text-2xl" />
      <span>旅行</span>
    </NuxtLink>
    <NuxtLink to="/settings" class="flex flex-col items-center gap-1 text-neutral-500 dark:text-neutral-400 no-underline text-xs py-2 px-4 transition-colors hover:text-primary active:text-primary" :class="{ 'text-primary dark:text-primary': route.path === '/settings' }">
      <Icon icon="mdi:cog" class="text-2xl" />
      <span>設定</span>
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

const route = useRoute()
</script>
