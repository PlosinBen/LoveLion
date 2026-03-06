<template>
  <BaseLayout>
    <template #header>
      <Header class="px-4" />
    </template>
    
    <slot />

    <template #footer>
      <!-- Priority: 1. Slot from page, 2. Global BottomNav (if not hidden) -->
      <slot name="footer">
        <BottomNav v-if="shouldShowGlobalNav" />
      </slot>
    </template>
  </BaseLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuth } from '~/composables/useAuth'
import Header from '~/components/Header.vue'
import BottomNav from '~/components/BottomNav.vue'

const { isAuthenticated, initAuth } = useAuth()
const route = useRoute()

onMounted(() => {
  initAuth()
})

const shouldShowGlobalNav = computed(() => {
  // Hide if not authenticated or if page meta explicitly says so
  if (!isAuthenticated.value) return false
  if (route.meta.hideGlobalNav) return false
  return true
})
</script>
