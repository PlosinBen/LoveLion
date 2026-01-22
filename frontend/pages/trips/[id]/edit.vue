<template>
  <div class="flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">編輯旅行</h1>
      <div class="w-10"></div>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <form v-else @submit.prevent="handleSubmit" class="flex flex-col gap-5">
      <!-- Name -->
      <div class="flex flex-col gap-2">
        <label class="text-sm text-neutral-400">旅行名稱 *</label>
        <input v-model="form.name" type="text" class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white focus:outline-none focus:border-indigo-500 placeholder-neutral-400" placeholder="例如：日本東京五日遊" required />
      </div>

      <!-- Description -->
      <div class="flex flex-col gap-2">
        <label class="text-sm text-neutral-400">描述</label>
        <textarea v-model="form.description" class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white focus:outline-none focus:border-indigo-500 placeholder-neutral-400 resize-none" rows="2" placeholder="旅行簡介（選填）"></textarea>
      </div>

      <!-- Dates -->
      <div class="grid grid-cols-2 gap-3">
        <div class="flex flex-col gap-2">
          <label class="text-sm text-neutral-400">開始日期</label>
          <input v-model="form.start_date" type="date" class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white focus:outline-none focus:border-indigo-500" />
        </div>
        <div class="flex flex-col gap-2">
          <label class="text-sm text-neutral-400">結束日期</label>
          <input v-model="form.end_date" type="date" class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white focus:outline-none focus:border-indigo-500" />
        </div>
      </div>

      <!-- Base Currency -->
      <div class="flex flex-col gap-2">
        <label class="text-sm text-neutral-400">基準貨幣</label>
        <select v-model="form.base_currency" class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white focus:outline-none focus:border-indigo-500">
          <option value="TWD">TWD - 新台幣</option>
          <option value="JPY">JPY - 日圓</option>
          <option value="USD">USD - 美元</option>
          <option value="EUR">EUR - 歐元</option>
          <option value="KRW">KRW - 韓元</option>
        </select>
      </div>

      <div class="border-t border-neutral-800 my-2 pt-4"></div>
      <h2 class="text-lg font-semibold text-neutral-200">帳本設定</h2>

      <!-- Active Currencies -->
      <ListEditor v-model="form.currencies" label="使用幣別" placeholder="新增幣別 (如 USD)" />

      <!-- Categories -->
      <ListEditor v-model="form.categories" label="消費分類" placeholder="新增分類 (如 交通)" />

      <!-- Payment Methods -->
      <ListEditor v-model="form.payment_methods" label="支付方式" placeholder="新增支付方式 (如 信用卡)" />

      <!-- Ledger Members -->
      <ListEditor v-model="form.ledger_members" label="旅伴 (Ledger Members)" placeholder="新增旅伴 (如 Kevin)" />

      <!-- Submit -->
      <button type="submit" class="w-full mt-3 px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed" :disabled="submitting">
        {{ submitting ? '儲存中...' : '儲存變更' }}
      </button>
    </form>
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

const loading = ref(true)
const submitting = ref(false)
interface TripForm {
  name: string
  description: string
  start_date: string
  end_date: string
  base_currency: string
  currencies: string[]
  categories: string[]
  payment_methods: string[]
  ledger_members: string[]
}

const form = ref<TripForm>({
  name: '',
  description: '',
  start_date: '',
  end_date: '',
  base_currency: 'TWD',
  currencies: ['TWD'],
  categories: [],
  payment_methods: [],
  ledger_members: []
})

const fetchTrip = async () => {
  try {
    const trip = await api.get<any>(`/api/trips/${route.params.id}`)
    form.value.name = trip.name
    form.value.description = trip.description
    form.value.base_currency = trip.base_currency || 'TWD'
    
    if (trip.start_date) {
      form.value.start_date = new Date(trip.start_date).toISOString().split('T')[0]
    }
    if (trip.end_date) {
      form.value.end_date = new Date(trip.end_date).toISOString().split('T')[0]
    }

    if (trip.ledger) {
      form.value.currencies = trip.ledger.currencies || ['TWD']
      form.value.categories = trip.ledger.categories || []
      form.value.payment_methods = trip.ledger.payment_methods || []
      form.value.ledger_members = trip.ledger.members || []
    }
  } catch (e) {
    console.error('Failed to fetch trip:', e)
    router.push(`/trips/${route.params.id}`)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!form.value.name.trim()) return

  submitting.value = true
  try {
    const payload: any = {
      name: form.value.name,
      description: form.value.description,
      base_currency: form.value.base_currency,
      currencies: form.value.currencies,
      categories: form.value.categories,
      payment_methods: form.value.payment_methods,
      ledger_members: form.value.ledger_members
    }
    if (form.value.start_date) {
      payload.start_date = new Date(form.value.start_date).toISOString()
    }
    if (form.value.end_date) {
      payload.end_date = new Date(form.value.end_date).toISOString()
    }

    await api.put(`/api/trips/${route.params.id}`, payload)
    router.push(`/trips/${route.params.id}`)
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  fetchTrip()
})
</script>
