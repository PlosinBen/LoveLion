<template>
  <div class="ledger-page">
    <header class="page-header">
      <h1>我的帳本</h1>
      <NuxtLink to="/ledger/add" class="btn btn-primary">
        <Icon icon="mdi:plus" /> 新增
      </NuxtLink>
    </header>

    <div v-if="loading" class="loading">載入中...</div>

    <div v-else-if="transactions.length === 0" class="empty-state card">
      <Icon icon="mdi:receipt-text-outline" />
      <p>還沒有任何記錄</p>
      <NuxtLink to="/ledger/add" class="btn btn-primary">開始記帳</NuxtLink>
    </div>

    <div v-else class="transactions">
      <div
        v-for="txn in transactions"
        :key="txn.id"
        class="transaction-card card"
        @click="router.push(`/ledger/${txn.id}`)"
      >
        <div class="txn-main">
          <div class="txn-category">
            <Icon :icon="getCategoryIcon(txn.category)" />
          </div>
          <div class="txn-info">
            <h3>{{ txn.category || '未分類' }}</h3>
            <p>{{ formatDate(txn.date) }} • {{ txn.items?.length || 0 }} 項目</p>
          </div>
          <div class="txn-amount">
            <span class="currency">{{ txn.currency }}</span>
            <span class="amount">{{ formatAmount(txn.total_amount) }}</span>
          </div>
        </div>
        <div v-if="txn.note" class="txn-note">{{ txn.note }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const transactions = ref<any[]>([])
const ledger = ref<any>(null)
const loading = ref(true)

const getCategoryIcon = (category: string) => {
  const icons: Record<string, string> = {
    'Food': 'mdi:food',
    'Transport': 'mdi:train-car',
    'Shopping': 'mdi:shopping',
    'Entertainment': 'mdi:movie',
    'Health': 'mdi:hospital',
    '餐飲': 'mdi:food',
    '交通': 'mdi:train-car',
    '購物': 'mdi:shopping',
  }
  return icons[category] || 'mdi:receipt'
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    month: 'short',
    day: 'numeric',
    weekday: 'short'
  })
}

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toLocaleString('zh-TW')
}

const fetchData = async () => {
  try {
    // Get or create default ledger
    const ledgers = await api.get<any[]>('/api/ledgers')

    if (ledgers.length === 0) {
      // Create default ledger
      ledger.value = await api.post<any>('/api/ledgers', {
        name: '我的帳本',
        type: 'personal',
        currencies: ['TWD'],
        categories: ['餐飲', '交通', '購物', '娛樂', '生活', '其他']
      })
    } else {
      ledger.value = ledgers[0]
    }

    // Fetch transactions
    const txns = await api.get<any[]>(`/api/ledgers/${ledger.value.id}/transactions`)
    transactions.value = txns
  } catch (e) {
    console.error('Failed to fetch data:', e)
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

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h1 {
  font-size: 24px;
}

.page-header .btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
}

.loading {
  text-align: center;
  color: var(--text-secondary);
  padding: 40px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-state svg {
  font-size: 64px;
  color: var(--text-secondary);
  margin-bottom: 16px;
}

.empty-state p {
  color: var(--text-secondary);
  margin-bottom: 20px;
}

.transactions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.transaction-card {
  cursor: pointer;
  transition: all 0.2s;
}

.transaction-card:hover {
  transform: translateY(-2px);
  border-color: var(--primary);
}

.txn-main {
  display: flex;
  align-items: center;
  gap: 12px;
}

.txn-category {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: var(--bg-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.txn-category svg {
  font-size: 22px;
  color: var(--primary);
}

.txn-info {
  flex: 1;
}

.txn-info h3 {
  font-size: 15px;
  margin-bottom: 4px;
}

.txn-info p {
  font-size: 12px;
  color: var(--text-secondary);
}

.txn-amount {
  text-align: right;
}

.txn-amount .currency {
  font-size: 12px;
  color: var(--text-secondary);
  display: block;
}

.txn-amount .amount {
  font-size: 18px;
  font-weight: 600;
}

.txn-note {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
  font-size: 13px;
  color: var(--text-secondary);
}
</style>
