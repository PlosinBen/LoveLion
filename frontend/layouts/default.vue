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
import { computed, onMounted } from 'vue'
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
    { label: '記帳', icon: 'mdi:wallet-outline', to: `/spaces/${spaceId}` },
    { label: '比價', icon: 'mdi:scale-balance', to: `/spaces/${spaceId}/stores` },
    { label: '設定', icon: 'mdi:cog-outline', to: `/spaces/${spaceId}/settings` }
  ]
})
</script>
