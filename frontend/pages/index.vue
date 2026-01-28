<template>
  <div class="flex flex-col gap-6">
    <StatusCard
      :status="backendStatus"
      :message="statusMessage"
    />

    <div class="grid grid-cols-2 gap-3 mb-6">
      <DashboardCard
        to="/ledger/add"
        icon="mdi:plus-circle"
        label="新增記帳"
        icon-color="text-indigo-500"
      />
      <DashboardCard
        to="/ledger"
        icon="mdi:wallet"
        label="我的帳本"
        icon-color="text-indigo-500"
      />
      <DashboardCard
        to="/trips"
        icon="mdi:airplane"
        label="旅行"
        icon-color="text-indigo-500"
      />
      <DashboardCard
        @click="handleLogout"
        icon="mdi:logout"
        label="登出"
        icon-color="text-red-500"
      />
    </div>

    <div v-if="recentTransactions.length > 0" class="flex flex-col gap-2">
      <h2 class="text-lg font-semibold mb-3">最近交易</h2>
      <div class="flex flex-col gap-2">
        <TransactionListItem
          v-for="txn in recentTransactions"
          :key="txn.id"
          :transaction="txn"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import DashboardCard from '~/components/DashboardCard.vue'
import StatusCard from '~/components/StatusCard.vue'
import TransactionListItem from '~/components/TransactionListItem.vue'
import { useAuth } from '~/composables/useAuth'
import { useApi } from '~/composables/useApi'

const router = useRouter()
const { isAuthenticated, initAuth, logout } = useAuth()
const api = useApi()

const backendStatus = ref<'checking' | 'online' | 'offline'>('checking')
const recentTransactions = ref<any[]>([])
const ledgers = ref<any[]>([])

const statusMessage = computed(() => {
  switch (backendStatus.value) {
    case 'online': return '後端服務正常運行'
    case 'offline': return '無法連接後端服務'
    default: return '檢查中...'
  }
})

const checkBackend = async () => {
  try {
    const response = await fetch('/health')
    if (response.ok) {
      backendStatus.value = 'online'
    } else {
      backendStatus.value = 'offline'
    }
  } catch {
    backendStatus.value = 'offline'
  }
}

const fetchData = async () => {
  try {
    // Fetch ledgers
    const ledgerData = await api.get<any[]>('/api/ledgers')
    ledgers.value = ledgerData

    // Fetch recent transactions from first ledger
    if (ledgerData.length > 0) {
      const txns = await api.get<any[]>(`/api/ledgers/${ledgerData[0].id}/transactions`)
      recentTransactions.value = txns.slice(0, 5)
    }
  } catch (e) {
    console.error('Failed to fetch data:', e)
  }
}

const handleLogout = () => {
  logout()
  router.push('/login')
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  checkBackend()
  fetchData()
})
</script>
