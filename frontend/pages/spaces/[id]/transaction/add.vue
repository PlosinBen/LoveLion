<template>
  <div class="add-transaction-page min-h-screen bg-neutral-900 text-neutral-50">
    <SpaceHeader 
      title="新增交易" 
      :show-back="true"
      class="pt-0 px-2"
    />

    <div class="p-4 pt-0">
      <form @submit.prevent="handleSubmit" class="flex flex-col gap-6">
        <!-- Date & Time -->
        <div class="flex flex-col gap-2">
          <label class="text-xs font-black text-neutral-500 uppercase tracking-widest px-1">日期與時間</label>
          <VueDatePicker 
            v-model="form.date" 
            :dark="true"
            :formats="{input: 'yyyy-MM-dd HH:mm'}"
            :enable-seconds="false"
            time-picker-inline
            cancel-text="取消"
            select-text="確定"
            placeholder="點擊選擇時間"
            class="date-picker-dark"
          />
        </div>

        <!-- Currency & Category -->
        <div class="grid grid-cols-2 gap-4">
          <div class="flex flex-col gap-2">
            <label class="text-xs font-black text-neutral-500 uppercase tracking-widest px-1">幣別</label>
            <BaseSelect
              v-model="form.currency"
              :options="availableCurrencies"
            />
          </div>
          <div class="flex flex-col gap-2">
            <label class="text-xs font-black text-neutral-500 uppercase tracking-widest px-1">類別</label>
            <BaseSelect
              v-model="form.category"
              :options="categories"
            />
          </div>
        </div>

        <!-- Items Section -->
        <div class="flex flex-col gap-3">
          <label class="text-xs font-black text-neutral-500 uppercase tracking-widest px-1 flex justify-between">
            <span>項目明細 ({{ form.currency }})</span>
            <span class="text-indigo-400">總計 {{ totalAmount.toLocaleString() }}</span>
          </label>
          
          <div class="flex flex-col gap-3">
            <div v-for="(item, index) in form.items" :key="index" class="bg-neutral-900 rounded-3xl p-5 border border-neutral-800 flex flex-col gap-4 relative">
              <BaseInput
                v-model="item.name"
                placeholder="項目名稱"
                input-class="font-bold"
              />
              <div class="flex items-center gap-3">
                <BaseInput
                  v-model.number="item.unit_price"
                  type="number"
                  placeholder="單價"
                  class="flex-1"
                />
                <span class="text-neutral-600 font-black">×</span>
                <BaseInput
                  v-model.number="item.quantity"
                  type="number"
                  placeholder="數量"
                  min="1"
                  class="w-20"
                />
                <button 
                  v-if="form.items.length > 1"
                  type="button" 
                  @click="removeItem(index)" 
                  class="w-10 h-10 rounded-xl bg-red-500/10 text-red-500 flex items-center justify-center border-0 cursor-pointer hover:bg-red-500/20"
                >
                  <Icon icon="mdi:close" />
                </button>
              </div>
            </div>
          </div>
          
          <button type="button" @click="addItem" class="flex justify-center items-center gap-2 p-4 border-2 border-dashed border-neutral-800 rounded-2xl bg-transparent text-neutral-500 font-bold cursor-pointer hover:border-indigo-500/50 hover:text-indigo-400 transition-colors">
            <Icon icon="mdi:plus" class="text-xl" /> 新增項目
          </button>
        </div>

        <!-- Foreign Currency (Conditional) -->
        <div v-if="form.currency !== baseCurrency" class="bg-indigo-500/5 rounded-3xl p-6 border border-indigo-500/10 flex flex-col gap-5">
          <div class="flex justify-between items-center">
            <h3 class="font-black text-sm text-indigo-400 uppercase tracking-widest">外幣結算</h3>
            <label class="flex items-center gap-2 cursor-pointer select-none">
              <input type="checkbox" v-model="form.manual_rate" class="w-4 h-4 rounded text-indigo-500 bg-neutral-800 border-neutral-700">
              <span class="text-xs font-bold text-neutral-400">手動輸入匯率</span>
            </label>
          </div>

          <div v-if="form.manual_rate" class="flex flex-col gap-4">
             <BaseInput
                v-model.number="form.exchange_rate"
                type="number"
                label="匯率 (1 TWD = ? 外幣)"
                step="0.0001"
             />
             <div class="flex justify-between items-center p-4 bg-neutral-900 rounded-2xl border border-neutral-800">
                <span class="text-neutral-500 text-xs font-bold">折合台幣</span>
                <span class="text-xl font-black text-white">TWD {{ calculatedBillingAmount.toLocaleString() }}</span>
             </div>
          </div>

          <div v-else class="flex flex-col gap-4">
             <BaseInput
                v-model.number="form.billing_amount"
                type="number"
                label="銀行入帳金額 (TWD)"
             />
             <div class="flex justify-between items-center p-4 bg-neutral-900 rounded-2xl border border-neutral-800">
                <span class="text-neutral-500 text-xs font-bold">換算匯率</span>
                <span class="text-xl font-black text-indigo-400">{{ calculatedExchangeRate }}</span>
             </div>
          </div>
        </div>

        <!-- Note -->
        <div class="flex flex-col gap-2">
          <label class="text-xs font-black text-neutral-500 uppercase tracking-widest px-1">備註</label>
          <BaseTextarea 
            v-model="form.note" 
            placeholder="選填" 
            rows="3"
          />
        </div>

        <!-- Submit -->
        <button 
          type="submit" 
          class="w-full mt-4 py-4 rounded-2xl font-black bg-indigo-500 text-white hover:bg-indigo-600 transition-all active:scale-[0.98] disabled:opacity-50 border-0 cursor-pointer shadow-lg shadow-indigo-500/20" 
          :disabled="submitting"
        >
          {{ submitting ? '儲存中...' : '儲存交易' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { VueDatePicker } from '@vuepic/vue-datepicker'
import '@vuepic/vue-datepicker/dist/main.css'
import SpaceHeader from '~/components/SpaceHeader.vue'

definePageMeta({
  layout: 'default',
  hideGlobalNav: true
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const categories = [
  { label: '餐飲', value: '餐飲' },
  { label: '交通', value: '交通' },
  { label: '購物', value: '購物' },
  { label: '娛樂', value: '娛樂' },
  { label: '生活', value: '生活' },
  { label: '其他', value: '其他' }
]

const availableCurrencies = [
  { label: 'TWD', value: 'TWD' },
  { label: 'JPY', value: 'JPY' },
  { label: 'USD', value: 'USD' },
  { label: 'EUR', value: 'EUR' }
]

const baseCurrency = ref('TWD')
const submitting = ref(false)

const form = ref({
  date: new Date(),
  category: '其他',
  note: '',
  items: [{ name: '', unit_price: 0, quantity: 1 }],
  currency: 'TWD',
  manual_rate: true, 
  exchange_rate: 0,
  billing_amount: 0,
  handling_fee: 0
})

const totalAmount = computed(() => {
  return form.value.items.reduce((sum, item) => {
    return sum + (Number(item.unit_price) * Number(item.quantity))
  }, 0)
})

const calculatedBillingAmount = computed(() => {
    if (form.value.manual_rate && form.value.exchange_rate > 0) {
        return Math.round(totalAmount.value * form.value.exchange_rate)
    }
    return 0
})

const calculatedExchangeRate = computed(() => {
    if (!form.value.manual_rate && totalAmount.value > 0 && form.value.billing_amount > 0) {
        return (form.value.billing_amount / totalAmount.value).toFixed(4)
    }
    return 0
})

watch(calculatedBillingAmount, (newVal) => {
    if (form.value.manual_rate) form.value.billing_amount = newVal
})

watch(calculatedExchangeRate, (newVal) => {
    if (!form.value.manual_rate) form.value.exchange_rate = Number(newVal)
})

const addItem = () => {
  form.value.items.push({ name: '', unit_price: 0, quantity: 1 })
}

const removeItem = (index: number) => {
  form.value.items.splice(index, 1)
}

const handleSubmit = async () => {
  if (!form.value.items.some(item => item.name && Number(item.unit_price) > 0)) {
    alert('請至少填寫一個項目名稱與金額')
    return
  }
  
  submitting.value = true
  try {
    const payload = {
      date: form.value.date.toISOString(),
      category: form.value.category,
      note: form.value.note,
      currency: form.value.currency,
      items: form.value.items.filter(item => item.name && Number(item.unit_price) > 0),
      exchange_rate: form.value.currency === baseCurrency.value ? 1 : form.value.exchange_rate,
      billing_amount: form.value.currency === baseCurrency.value ? totalAmount.value : form.value.billing_amount,
      total_amount: totalAmount.value
    }

    await api.post(`/api/spaces/${route.params.id}/transactions`, payload)
    router.push(`/spaces/${route.params.id}`)
  } catch (e: any) {
    alert(e.message || '儲存失敗')
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  initAuth()
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  
  try {
    const space = await api.get<any>(`/api/spaces/${route.params.id}`)
    baseCurrency.value = space.base_currency || 'TWD'
    form.value.currency = baseCurrency.value
  } catch (e) {
    console.error('Failed to fetch space currency', e)
  }
})
</script>
