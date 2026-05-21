<template>
  <div class="investments-page">
    <NuxtPage />

    <BottomNav
      v-if="!route.meta.hideInvNav"
      :items="navItems"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuth } from '~/composables/useAuth'
import BottomNav from '~/components/BottomNav.vue'

definePageMeta({
  hideGlobalNav: true,
})

const { user } = useAuth()
const route = useRoute()

const memberQuery = computed(() => route.query.member as string | undefined)

const navItems = computed(() => {
  const isOwnerView = user.value?.inv_is_owner && !memberQuery.value
  if (isOwnerView) {
    return [
      { label: '損益', icon: 'mdi:chart-line', to: '/investments', exact: true },
      { label: '結算', icon: 'mdi:calculator-variant', to: '/investments/settlements' },
      { label: '交易', icon: 'mdi:swap-horizontal', to: '/investments/trades' },
      { label: '出入金', icon: 'mdi:transfer', to: '/investments/transactions' },
    ]
  }
  const suffix = memberQuery.value ? `?member=${memberQuery.value}` : ''
  return [
    { label: '損益', icon: 'mdi:chart-line', to: `/investments${suffix}`, exact: true },
    { label: '出入金', icon: 'mdi:transfer', to: `/investments/transactions${suffix}` },
  ]
})
</script>
