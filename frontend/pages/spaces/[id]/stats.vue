<template>
  <div class="space-stats-page">
    <PageTitle
      title="消費統計"
      :show-back="true"
      :settings-to="`/spaces/${route.params.id}/settings`"
      :breadcrumbs="[{ label: store.space?.name || '空間', to: `/spaces/${route.params.id}/ledger` }]"
    />

    <div v-if="store.loading.space || store.loading.transactions" class="flex justify-center items-center py-20 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="store.transactions.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
      <div class="flex flex-col items-center justify-center gap-4">
        <Icon icon="mdi:chart-arc" class="text-5xl opacity-20" />
        <p>目前還沒有交易紀錄，無法產生統計資訊</p>
        <BaseButton @click="router.push(`/spaces/${route.params.id}/ledger/transaction/add`)" variant="primary">
          立即新增第一筆交易
        </BaseButton>
      </div>
    </div>

    <SpaceStats
      v-else
      :transactions="store.transactions"
      :members="store.members"
      :base-currency="store.space?.base_currency || 'TWD'"
    />

    <!-- FAB for adding transaction -->
    <BaseFab @click="router.push(`/spaces/${route.params.id}/ledger/transaction/add`)" />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import PageTitle from '~/components/PageTitle.vue'
import SpaceStats from '~/components/SpaceStats.vue'
import BaseFab from '~/components/BaseFab.vue'
import BaseButton from '~/components/BaseButton.vue'

// Map both /spaces/:id and /spaces/:id/stats to this file
definePageMeta({
  alias: '/spaces/:id'
})

const route = useRoute()
const router = useRouter()
const { isAuthenticated, initAuth } = useAuth()
const store = useSpaceDetailStore()

onMounted(async () => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  store.setSpaceId(route.params.id as string)
  try {
    await Promise.all([store.fetchSpace(), store.fetchTransactions()])
  } catch (e) {
    router.push('/')
  }
})
</script>
