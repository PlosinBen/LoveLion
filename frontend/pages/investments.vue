<template>
  <div class="investments-page">
    <NuxtPage />

    <BottomNav
      v-if="user?.inv_is_owner && !route.meta.hideInvNav"
      :items="invNavItems"
    />
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'
import BottomNav from '~/components/BottomNav.vue'

definePageMeta({
  hideGlobalNav: true,
})

const { user } = useAuth()
const route = useRoute()

const invNavItems = [
  { label: '損益', icon: 'mdi:chart-line', to: '/investments' },
  { label: '結算', icon: 'mdi:calculator-variant', to: '/investments/settlements' },
  { label: '交易', icon: 'mdi:swap-horizontal', to: '/investments/trades' },
  { label: '異動', icon: 'mdi:transfer', to: '/investments/transactions' },
]
</script>
