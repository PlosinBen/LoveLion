<template>
  <div class="space-stats-page">
    <SpaceHeader title="統計" />

    <div v-if="store.loading.space || store.loading.transactions" class="text-center py-10 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="store.transactions.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
      目前還沒有交易紀錄
    </div>

    <SpaceStats
      v-else
      :transactions="store.transactions"
      :members="store.members"
      :base-currency="store.space?.base_currency || 'TWD'"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import { useSpaceDetailStore } from '~/stores/spaceDetail'
import SpaceHeader from '~/components/SpaceHeader.vue'
import SpaceStats from '~/components/SpaceStats.vue'

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
