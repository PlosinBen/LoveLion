<template>
  <div class="flex flex-col min-h-screen bg-neutral-900 text-neutral-50 font-sans transition-colors duration-200">
    <header class="px-4">
      <NuxtLink to="/">
        <h1 class="text-3xl font-bold bg-gradient-to-br from-indigo-500 to-purple-500 bg-clip-text text-transparent">
          <img src="/dark-logo.png" alt="LoveLion" class="w-56">
        </h1>
      </NuxtLink>
    </header>
    <main class="flex-1 w-full max-w-lg mx-auto p-4 relative" :class="{'pb-20': showBottomNav, 'pb-4': !showBottomNav}">

      <slot />
    </main>
    <BottomNav v-if="isAuthenticated" />
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'

const { isAuthenticated, initAuth } = useAuth()
const route = useRoute()

const showBottomNav = computed(() => {
    const path = route.path
    // Trip context or Daily Ledger context
    return (path.startsWith('/trips/') && !!route.params.id) || path.startsWith('/ledger')
})

onMounted(() => {
  initAuth()
})
</script>
