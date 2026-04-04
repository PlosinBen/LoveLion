<template>
  <BaseLayout>
    <template #header>
      <Header class="px-4" />
    </template>
    
    <slot />

    <template #footer>
      <slot name="footer">
        <BottomNav v-if="shouldShowGlobalNav" :items="spaceNavItems" />
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

const { isAuthenticated } = useAuth()
const route = useRoute()

const shouldShowGlobalNav = computed(() => {
  if (!isAuthenticated.value) return false
  if (route.meta.hideGlobalNav) return false
  return true
})

const spaceNavItems = computed(() => {
  const spaceId = route.params.id
  if (!spaceId || typeof spaceId !== 'string') return undefined

  // Inside space: only show space-related tools
  return [
    { label: '統計', icon: 'mdi:chart-bar', to: `/spaces/${spaceId}/stats` },
    { label: '比價', icon: 'mdi:scale-balance', to: `/spaces/${spaceId}/stores`, alternateTo: `/spaces/${spaceId}/products` },
    { label: '記帳', icon: 'mdi:wallet-outline', to: `/spaces/${spaceId}/ledger` }
  ]
})
</script>
