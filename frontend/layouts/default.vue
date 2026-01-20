<template>
  <div class="flex flex-col min-h-screen bg-neutral-50 dark:bg-neutral-900 text-neutral-900 dark:text-neutral-50 font-sans transition-colors duration-200">
    <main class="flex-1 w-full max-w-lg mx-auto p-4 pb-20 relative">
      <slot />
    </main>
    <BottomNav v-if="isAuthenticated" />
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'

const { isAuthenticated, initAuth } = useAuth()

onMounted(() => {
  initAuth()
})
</script>
