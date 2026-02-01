<template>
  <BaseLayout content-class="max-w-lg p-4 pb-20">
    <template #header>
      <BrandLogo class="px-4" />
    </template>
    
    <slot />

    <template #footer>
        <BottomNav v-if="isAuthenticated" />
    </template>
  </BaseLayout>
</template>

<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'

const { isAuthenticated, initAuth } = useAuth()

onMounted(() => {
  initAuth()
})
</script>
