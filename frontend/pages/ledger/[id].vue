<template>
  <div class="transaction-detail">
    <header class="page-header">
      <button @click="router.back()" class="back-btn">
        <Icon icon="mdi:arrow-left" />
      </button>
      <h1>交易詳情</h1>
      <button @click="handleDelete" class="delete-btn">
        <Icon icon="mdi:delete" />
      </button>
    </header>

    <div v-if="loading" class="loading">載入中...</div>

    <div v-else-if="transaction" class="content">
      <div class="header-card card">
        <div class="category-badge">
          <Icon :icon="getCategoryIcon(transaction.category)" />
          <span>{{ transaction.category || '未分類' }}</span>
        </div>
        <div class="total-amount">
          {{ transaction.currency }} {{ formatAmount(transaction.total_amount) }}
        </div>
        <div class="date">{{ formatDate(transaction.date) }}</div>
      </div>

      <div class="section">
        <h2>項目明細</h2>
        <div class="items-list">
          <div v-for="item in transaction.items" :key="item.id" class="item-row card">
            <div class="item-info">
              <span class="item-name">{{ item.name }}</span>
              <span class="item-qty">× {{ item.quantity }}</span>
            </div>
            <div class="item-price">{{ transaction.currency }} {{ formatAmount(item.amount) }}</div>
          </div>
        </div>
      </div>

      <div v-if="transaction.note" class="section">
        <h2>備註</h2>
        <div class="note-card card">{{ transaction.note }}</div>
      </div>

      <div class="section meta">
        <p>建立時間：{{ formatDateTime(transaction.created_at) }}</p>
        <p>ID：{{ transaction.id }}</p>
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
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const transaction = ref<any>(null)
const loading = ref(true)
const ledgerId = ref('')

const getCategoryIcon = (category: string) => {
  const icons: Record<string, string> = {
    '餐飲': 'mdi:food',
    '交通': 'mdi:train-car',
    '購物': 'mdi:shopping',
    '娛樂': 'mdi:movie',
    '生活': 'mdi:home',
  }
  return icons[category] || 'mdi:receipt'
}

const formatAmount = (amount: string | number) => {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount
  return num.toLocaleString('zh-TW')
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
}

const formatDateTime = (dateStr: string) => {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-TW')
}

const fetchData = async () => {
  try {
    const ledgers = await api.get<any[]>('/api/ledgers')
    if (ledgers.length > 0) {
      ledgerId.value = ledgers[0].id
      transaction.value = await api.get<any>(
        `/api/ledgers/${ledgerId.value}/transactions/${route.params.id}`
      )
    }
  } catch (e) {
    console.error('Failed to fetch transaction:', e)
  } finally {
    loading.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('確定要刪除這筆交易嗎？')) return

  try {
    await api.del(`/api/ledgers/${ledgerId.value}/transactions/${route.params.id}`)
    router.push('/ledger')
  } catch (e: any) {
    alert(e.message || '刪除失敗')
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
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-header h1 {
  font-size: 20px;
}

.back-btn,
.delete-btn {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  border: none;
  background: var(--bg-secondary);
  color: var(--text-primary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-btn svg,
.delete-btn svg {
  font-size: 24px;
}

.delete-btn {
  color: var(--danger);
}

.loading {
  text-align: center;
  color: var(--text-secondary);
  padding: 40px;
}

.header-card {
  text-align: center;
  padding: 32px 20px;
  margin-bottom: 24px;
}

.category-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--bg-tertiary);
  border-radius: 20px;
  margin-bottom: 16px;
}

.category-badge svg {
  font-size: 20px;
  color: var(--primary);
}

.total-amount {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 8px;
}

.date {
  color: var(--text-secondary);
}

.section {
  margin-bottom: 24px;
}

.section h2 {
  font-size: 16px;
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.items-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.item-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
}

.item-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.item-name {
  font-weight: 500;
}

.item-qty {
  color: var(--text-secondary);
  font-size: 14px;
}

.item-price {
  font-weight: 600;
}

.note-card {
  padding: 16px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.meta {
  color: var(--text-secondary);
  font-size: 12px;
}

.meta p {
  margin-bottom: 4px;
}
</style>
