<template>
  <div class="space-stats-page">
    <SpaceHeader title="統計" />

    <div v-if="loading" class="text-center py-10 text-neutral-500">
      <Icon icon="mdi:loading" class="text-3xl animate-spin" />
    </div>

    <div v-else-if="transactions.length === 0" class="bg-neutral-900/50 rounded-2xl border border-neutral-800 border-dashed p-10 flex flex-col items-center justify-center text-neutral-500 text-sm italic">
      目前還沒有交易紀錄
    </div>

    <SpaceStats
      v-else
      :transactions="transactions"
      :members="members"
      :base-currency="baseCurrency"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import SpaceHeader from '~/components/SpaceHeader.vue'
import SpaceStats from '~/components/SpaceStats.vue'

const route = useRoute()
const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const loading = ref(true)
const transactions = ref<any[]>([])
const members = ref<any[]>([])
const baseCurrency = ref('TWD')

const fetchData = async () => {
  try {
    const spaceId = route.params.id
    const [spaceData, txnData] = await Promise.all([
      api.get<any>(`/api/spaces/${spaceId}`),
      api.get<any>(`/api/spaces/${spaceId}/transactions`)
    ])
    members.value = spaceData.members || []
    baseCurrency.value = spaceData.base_currency || 'TWD'
    transactions.value = txnData || []
  } catch (e) {
    console.error('Failed to fetch stats data:', e)
    router.push('/')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchData()
})
</script>
