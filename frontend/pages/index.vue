<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <h1>LoveLion</h1>
      <p>嗨，{{ user?.display_name || '使用者' }}！</p>
    </header>

    <div class="status-card card" :class="statusClass">
      <div class="status-icon">
        <Icon :icon="statusIcon" />
      </div>
      <div class="status-text">
        <h3>系統狀態</h3>
        <p>{{ statusMessage }}</p>
      </div>
    </div>

    <div class="quick-actions">
      <NuxtLink to="/ledger/add" class="action-card card">
        <Icon icon="mdi:plus-circle" />
        <span>新增記帳</span>
      </NuxtLink>
      <NuxtLink to="/ledger" class="action-card card">
        <Icon icon="mdi:wallet" />
        <span>我的帳本</span>
      </NuxtLink>
      <NuxtLink to="/trips" class="action-card card">
        <Icon icon="mdi:airplane" />
        <span>旅行</span>
      </NuxtLink>
      <button @click="handleLogout" class="action-card card logout">
        <Icon icon="mdi:logout" />
        <span>登出</span>
      </button>
    </div>

    <div v-if="recentTransactions.length > 0" class="recent-section">
      <h2>最近交易</h2>
      <div class="transaction-list">
        <div
          v-for="txn in recentTransactions"
          :key="txn.id"
          class="transaction-item card"
        >
          <div class="txn-info">
            <h4>{{ txn.category || '未分類' }}</h4>
            <p>{{ formatDate(txn.date) }}</p>
          </div>
          <div class="txn-amount">
            {{ txn.currency }} {{ txn.total_amount }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAuth } from '~/composables/useAuth'
import { useApi } from '~/composables/useApi'

const router = useRouter()
const { user, isAuthenticated, initAuth, logout } = useAuth()
const api = useApi()

const backendStatus = ref<'checking' | 'online' | 'offline'>('checking')
const recentTransactions = ref<any[]>([])
const ledgers = ref<any[]>([])

const statusClass = computed(() => ({
  online: backendStatus.value === 'online',
  offline: backendStatus.value === 'offline',
}))

const statusIcon = computed(() => {
  switch (backendStatus.value) {
    case 'online': return 'mdi:check-circle'
    case 'offline': return 'mdi:close-circle'
    default: return 'mdi:loading'
  }
})

const statusMessage = computed(() => {
  switch (backendStatus.value) {
    case 'online': return '後端服務正常運行'
    case 'offline': return '無法連接後端服務'
    default: return '檢查中...'
  }
})

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', { month: 'short', day: 'numeric' })
}

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

<style scoped>
.dashboard-header {
  margin-bottom: 24px;
}

.dashboard-header h1 {
  font-size: 28px;
  background: linear-gradient(135deg, var(--primary), #a855f7);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.dashboard-header p {
  color: var(--text-secondary);
  margin-top: 4px;
}

.status-card {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.status-card.online {
  border-color: var(--success);
}

.status-card.offline {
  border-color: var(--danger);
}

.status-icon svg {
  font-size: 32px;
}

.status-card.online .status-icon {
  color: var(--success);
}

.status-card.offline .status-icon {
  color: var(--danger);
}

.status-text h3 {
  font-size: 14px;
  color: var(--text-secondary);
}

.status-text p {
  font-size: 16px;
  margin-top: 4px;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.action-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 20px;
  text-decoration: none;
  color: var(--text-primary);
  transition: all 0.2s;
  cursor: pointer;
  border: none;
  background: var(--bg-secondary);
}

.action-card:hover {
  transform: translateY(-2px);
  border-color: var(--primary);
}

.action-card svg {
  font-size: 28px;
  color: var(--primary);
}

.action-card.logout svg {
  color: var(--danger);
}

.recent-section h2 {
  font-size: 18px;
  margin-bottom: 12px;
}

.transaction-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.transaction-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
}

.txn-info h4 {
  font-size: 14px;
}

.txn-info p {
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 4px;
}

.txn-amount {
  font-weight: 600;
  color: var(--primary);
}
</style>
