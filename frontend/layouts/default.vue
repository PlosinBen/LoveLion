<template>
  <BaseLayout>
    <template #header>
      <Header class="px-4" />
      <BroadcastBar />
    </template>
    
    <slot />

    <template #footer>
      <slot name="footer">
        <BottomNav v-if="shouldShowGlobalNav" :items="spaceNavItems || globalNavItems" />
      </slot>
    </template>
  </BaseLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuth } from '~/composables/useAuth'
import Header from '~/components/Header.vue'
import BroadcastBar from '~/components/BroadcastBar.vue'
import BottomNav from '~/components/BottomNav.vue'

const { isAuthenticated, user } = useAuth()
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

const globalNavItems = computed(() => {
  if (spaceNavItems.value) return undefined
  const items: any[] = [
    { label: '空間', icon: 'mdi:view-grid-outline', to: '/' },
  ]
  if (user.value?.inv_access) {
    items.push({ label: '投資', icon: 'mdi:chart-line', to: '/investments' })
  }
  items.push({ label: '設定', icon: 'mdi:cog-outline', to: '/settings' })
  return items
})
</script>
