<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center gap-4 mb-6 p-5 rounded-2xl border bg-neutral-900 duration-200" :class="statusClasses">
      <div class="text-3xl" :class="statusIconColor">
        <Icon :icon="statusIcon" />
      </div>
      <div class="flex flex-col">
        <h3 class="text-sm text-neutral-400">系統狀態</h3>
        <p class="text-base mt-1">{{ statusMessage }}</p>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-3 mb-6">
      <NuxtLink to="/ledger/add" class="flex flex-col items-center gap-2 p-5 rounded-2xl border border-neutral-800 bg-neutral-900 text-white no-underline transition-all duration-200 cursor-pointer hover:-translate-y-0.5 hover:border-indigo-500">
        <Icon icon="mdi:plus-circle" class="text-3xl text-indigo-500" />
        <span>新增記帳</span>
      </NuxtLink>
      <NuxtLink to="/ledger" class="flex flex-col items-center gap-2 p-5 rounded-2xl border border-neutral-800 bg-neutral-900 text-white no-underline transition-all duration-200 cursor-pointer hover:-translate-y-0.5 hover:border-indigo-500">
        <Icon icon="mdi:wallet" class="text-3xl text-indigo-500" />
        <span>我的帳本</span>
      </NuxtLink>
      <NuxtLink to="/trips" class="flex flex-col items-center gap-2 p-5 rounded-2xl border border-neutral-800 bg-neutral-900 text-white no-underline transition-all duration-200 cursor-pointer hover:-translate-y-0.5 hover:border-indigo-500">
        <Icon icon="mdi:airplane" class="text-3xl text-indigo-500" />
        <span>旅行</span>
      </NuxtLink>
      <button @click="handleLogout" class="flex flex-col items-center gap-2 p-5 rounded-2xl border border-neutral-800 bg-neutral-900 text-white no-underline transition-all duration-200 cursor-pointer hover:-translate-y-0.5 hover:border-indigo-500">
        <Icon icon="mdi:logout" class="text-3xl text-red-500" />
        <span>登出</span>
      </button>
    </div>

    <div v-if="recentTransactions.length > 0" class="flex flex-col gap-2">
      <h2 class="text-lg font-semibold mb-3">最近交易</h2>
      <div class="flex flex-col gap-2">
        <div
          v-for="txn in recentTransactions"
          :key="txn.id"
          class="flex justify-between items-center p-4 rounded-xl border border-neutral-800 bg-neutral-900"
        >
          <div class="flex flex-col">
            <h4 class="text-sm font-medium">{{ txn.category || '未分類' }}</h4>
            <p class="text-xs text-neutral-400 mt-1">{{ formatDate(txn.date) }}</p>
          </div>
          <div class="font-semibold text-indigo-500">
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

const statusClasses = computed(() => {
  if (backendStatus.value === 'online') return 'border-green-500'
  if (backendStatus.value === 'offline') return 'border-red-500'
  return 'border-neutral-800'
})

const statusIconColor = computed(() => {
  if (backendStatus.value === 'online') return 'text-green-500'
  if (backendStatus.value === 'offline') return 'text-red-500'
  return 'text-neutral-500'
})

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
