<template>
  <div class="add-transaction-page">
    <header class="page-header">
      <button @click="router.back()" class="back-btn">
        <Icon icon="mdi:arrow-left" />
      </button>
      <h1>{{ isEdit ? '編輯交易' : '新增記帳' }}</h1>
      <div style="width: 40px;"></div>
    </header>

    <form @submit.prevent="handleSubmit" class="form">
      <!-- Date -->
      <div class="form-group">
        <label class="label">日期</label>
        <input v-model="form.date" type="date" class="input" />
      </div>

      <!-- Category -->
      <div class="form-group">
        <label class="label">類別</label>
        <div class="category-grid">
          <button
            v-for="cat in categories"
            :key="cat"
            type="button"
            class="category-btn"
            :class="{ active: form.category === cat }"
            @click="form.category = cat"
          >
            {{ cat }}
          </button>
        </div>
      </div>

      <!-- Items -->
      <div class="form-group">
        <label class="label">項目明細</label>
        <div class="items-list">
          <div v-for="(item, index) in form.items" :key="index" class="item-row card">
            <div class="item-inputs">
              <input
                v-model="item.name"
                type="text"
                class="input item-name"
                placeholder="項目名稱"
              />
              <div class="item-price-row">
                <input
                  v-model.number="item.unit_price"
                  type="number"
                  class="input item-price"
                  placeholder="單價"
                />
                <span class="multiply">×</span>
                <input
                  v-model.number="item.quantity"
                  type="number"
                  class="input item-qty"
                  placeholder="數量"
                  min="1"
                />
                <button type="button" @click="removeItem(index)" class="remove-btn" v-if="form.items.length > 1">
                  <Icon icon="mdi:close" />
                </button>
              </div>
            </div>
            <div class="item-subtotal">
              {{ currency }} {{ (item.unit_price * item.quantity).toLocaleString() }}
            </div>
          </div>
        </div>
        <button type="button" @click="addItem" class="add-item-btn">
          <Icon icon="mdi:plus" /> 新增項目
        </button>
      </div>

      <!-- Total -->
      <div class="total-section card">
        <span>總計</span>
        <span class="total-amount">{{ currency }} {{ totalAmount.toLocaleString() }}</span>
      </div>

      <!-- Note -->
      <div class="form-group">
        <label class="label">備註</label>
        <textarea v-model="form.note" class="input textarea" rows="2" placeholder="選填"></textarea>
      </div>

      <!-- Submit -->
      <button type="submit" class="btn btn-primary btn-block" :disabled="submitting">
        {{ submitting ? '儲存中...' : '儲存' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const isEdit = computed(() => !!route.params.id)
const currency = 'TWD'
const categories = ['餐飲', '交通', '購物', '娛樂', '生活', '其他']
const ledgerId = ref('')
const submitting = ref(false)

const form = ref({
  date: new Date().toISOString().split('T')[0],
  category: '',
  note: '',
  items: [{ name: '', unit_price: 0, quantity: 1 }]
})

const totalAmount = computed(() => {
  return form.value.items.reduce((sum, item) => {
    return sum + (item.unit_price * item.quantity)
  }, 0)
})

const addItem = () => {
  form.value.items.push({ name: '', unit_price: 0, quantity: 1 })
}

const removeItem = (index: number) => {
  form.value.items.splice(index, 1)
}

const handleSubmit = async () => {
  if (!form.value.items.some(item => item.name && item.unit_price > 0)) {
    alert('請至少填寫一個項目')
    return
  }

  submitting.value = true

  try {
    const payload = {
      date: new Date(form.value.date).toISOString(),
      category: form.value.category,
      note: form.value.note,
      items: form.value.items.filter(item => item.name && item.unit_price > 0)
    }

    await api.post(`/api/ledgers/${ledgerId.value}/transactions`, payload)
    router.push('/ledger')
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    submitting.value = false
  }
}

const fetchLedger = async () => {
  try {
    const ledgers = await api.get<any[]>('/api/ledgers')
    if (ledgers.length === 0) {
      const newLedger = await api.post<any>('/api/ledgers', {
        name: '我的帳本',
        type: 'personal',
        currencies: ['TWD'],
        categories
      })
      ledgerId.value = newLedger.id
    } else {
      ledgerId.value = ledgers[0].id
    }
  } catch (e) {
    console.error('Failed to fetch ledger:', e)
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchLedger()
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

.back-btn {
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

.back-btn svg {
  font-size: 24px;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.category-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.category-btn {
  padding: 10px 16px;
  border-radius: 20px;
  border: 1px solid var(--border-color);
  background: var(--bg-secondary);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.2s;
}

.category-btn.active {
  background: var(--primary);
  border-color: var(--primary);
}

.items-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.item-row {
  padding: 16px;
}

.item-inputs {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.item-name {
  font-weight: 500;
}

.item-price-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.item-price {
  flex: 2;
}

.item-qty {
  flex: 1;
  width: 60px;
}

.multiply {
  color: var(--text-secondary);
}

.remove-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: none;
  background: var(--danger);
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.item-subtotal {
  margin-top: 12px;
  text-align: right;
  color: var(--text-secondary);
  font-size: 14px;
}

.add-item-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px;
  border: 2px dashed var(--border-color);
  border-radius: 12px;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  margin-top: 12px;
}

.add-item-btn:hover {
  border-color: var(--primary);
  color: var(--primary);
}

.total-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 18px;
}

.total-amount {
  font-size: 24px;
  font-weight: 700;
  color: var(--primary);
}

.textarea {
  resize: none;
}

.btn-block {
  width: 100%;
  margin-top: 12px;
}
</style>
