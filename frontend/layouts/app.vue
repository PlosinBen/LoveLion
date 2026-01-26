<template>
  <div class="flex flex-col h-screen bg-black text-neutral-50 font-sans">
    <header class="px-4 shrink-0">
      <NuxtLink to="/">
        <h1 class="text-3xl font-bold bg-gradient-to-br from-indigo-500 to-purple-500 bg-clip-text text-transparent">
          <img src="/dark-logo.png" alt="LoveLion" class="w-56">
        </h1>
      </NuxtLink>
    </header>
    <main class="flex-1 w-full max-w-lg mx-auto relative overflow-hidden flex flex-col">
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
