<template>
  <div class="flex flex-col gap-6 min-h-screen bg-neutral-900 text-neutral-50 p-4">
    <header class="flex justify-between items-center">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">新增旅行</h1>
      <div class="w-10"></div>
    </header>

    <form @submit.prevent="handleSubmit" class="flex flex-col gap-5">
      <!-- Name -->
      <BaseInput v-model="form.name" label="旅行名稱 *" placeholder="例如：日本東京五日遊" required />

      <!-- Description -->
      <div class="flex flex-col gap-2">
        <BaseTextarea 
          v-model="form.description" 
          label="描述" 
          rows="2" 
          placeholder="旅行簡介（選填）" 
        />
      </div>

      <!-- Dates -->
      <div class="grid grid-cols-2 gap-3">
        <BaseInput v-model="form.start_date" type="date" label="開始日期" />
        <BaseInput v-model="form.end_date" type="date" label="結束日期" />
      </div>

      <div class="flex flex-col gap-3">
        <label class="text-sm font-medium text-neutral-400">幣別設定</label>
        
        <!-- Input Row -->
        <div class="grid grid-cols-2 gap-3">
          <!-- Base Currency -->
          <BaseSelect 
            v-model="form.base_currency" 
            label="基準貨幣 *"
            :options="currencyOptions"
            :disabled="form.currencies.length === 0"
            :placeholder="form.currencies.length === 0 ? '需先新增幣別' : '選擇基準'"
          />
          
          <!-- Add Currency Input -->
           <div class="flex flex-col gap-1">
             <span class="text-sm text-neutral-400">新增幣別</span>
             <div class="relative">
              <BaseInput
                v-model="newCurrency" 
                placeholder="輸入代碼 (如 USD)"
                class="pr-10"
                @keydown.enter.prevent="addCurrency"
                @blur="addCurrency"
              />
              <button 
                type="button" 
                @click="addCurrency" 
                class="absolute right-3 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-white transition-all"
                :class="{ 'opacity-0 scale-75 pointer-events-none': !newCurrency, 'opacity-100 scale-100': newCurrency }"
              >
                <Icon icon="mdi:plus" class="text-xl" />
              </button>
            </div>
           </div>
        </div>

        <!-- Currencies List (Full Width) -->
        <div class="flex flex-wrap gap-2">
          <div v-for="(item, index) in form.currencies" :key="index" class="flex items-center gap-1 bg-neutral-700 px-3 py-1 rounded text-sm text-white">
            <span>{{ item }}</span>
            <button type="button" @click="removeCurrency(index)" class="text-neutral-400 hover:text-white transition-colors">
              <Icon icon="mdi:close" class="text-base" />
            </button>
          </div>
          <span v-if="form.currencies.length === 0" class="text-sm text-neutral-500 italic">尚未新增任何幣別</span>
        </div>
      </div>

      <!-- Members -->
      <ListEditor v-model="form.members" label="旅行成員" placeholder="新增成員 (如 Kevin)" />

      <div class="border-t border-neutral-800 my-2 pt-4"></div>
      <h2 class="text-lg font-semibold text-neutral-200">帳本設定</h2>

      <!-- Moved Currencies List up -->

      <!-- Categories -->
      <ListEditor v-model="form.categories" label="消費分類" placeholder="新增分類 (如 交通)" />

      <!-- Payment Methods -->
      <ListEditor v-model="form.payment_methods" label="支付方式" placeholder="新增支付方式 (如 信用卡)" />

      <!-- Submit -->
      <button type="submit" class="w-full mt-3 px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed" :disabled="submitting">
        {{ submitting ? '建立中...' : '建立旅行' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'

definePageMeta({
  layout: 'empty'
})

const router = useRouter()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const submitting = ref(false)
const form = ref({
  name: '',
  description: '',
  start_date: '',
  end_date: '',
  base_currency: 'TWD',
  members: [] as string[],
  currencies: ['TWD', 'JPY', 'USD'],
  categories: [],
  payment_methods: []
})

const newCurrency = ref('')

const addCurrency = () => {
  const val = newCurrency.value.trim().toUpperCase()
  if (val && !form.value.currencies.includes(val)) {
    form.value.currencies.push(val)
    newCurrency.value = ''
  } else {
    newCurrency.value = ''
  }
}

const removeCurrency = (index: number) => {
  form.value.currencies.splice(index, 1)
}

const currencyOptions = computed(() => {
  return form.value.currencies.map(c => ({ label: c, value: c }))
})

// Ensure base_currency is valid
watch(() => form.value.currencies, (newVal) => {
  if (newVal.length > 0 && !newVal.includes(form.value.base_currency)) {
    form.value.base_currency = newVal[0]
  } else if (newVal.length === 0) {
      form.value.base_currency = ''
  }
}, { deep: true })

const handleSubmit = async () => {
  if (!form.value.name.trim()) return
  if (!form.value.base_currency) {
      alert('請選擇基準貨幣')
      return
  }

  submitting.value = true
  try {
    const payload: any = {
      name: form.value.name,
      description: form.value.description,
      base_currency: form.value.base_currency,
      members: form.value.members,
      currencies: form.value.currencies,
      categories: form.value.categories,
      payment_methods: form.value.payment_methods
    }
    if (form.value.start_date) {
      payload.start_date = new Date(form.value.start_date).toISOString()
    }
    if (form.value.end_date) {
      payload.end_date = new Date(form.value.end_date).toISOString()
    }

    const trip = await api.post<any>('/api/trips', payload)
    router.push(`/trips/${trip.id}`)
  } catch (e: any) {
    alert(e.message || '建立失敗')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
  }
})
</script>
