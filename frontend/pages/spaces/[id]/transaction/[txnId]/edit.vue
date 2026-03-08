<template>
  <div class="edit-transaction-page">
    <SpaceHeader 
      title="編輯交易" 
      :show-back="true"
      class="px-2"
    />

    <div v-if="loading" class="p-20 flex justify-center items-center text-neutral-500">
        <Icon icon="mdi:loading" class="text-4xl animate-spin" />
    </div>

    <div v-else>
      <form @submit.prevent="handleSubmit" class="flex flex-col gap-6">
        <!-- Date & Time -->
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">日期與時間</label>
          <VueDatePicker 
            v-model="form.date" 
            :dark="true"
            :formats="{input: 'yyyy-MM-dd HH:mm'}"
            :enable-seconds="false"
            time-picker-inline
            cancel-text="取消"
            select-text="確定"
            placeholder="點擊選擇時間"
          />
        </div>

        <!-- Currency & Category -->
        <div class="grid grid-cols-2 gap-4">
          <div class="flex flex-col gap-2">
            <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">幣別</label>
            <BaseSelect
              v-model="form.currency"
              :options="availableCurrencies"
            />
          </div>
          <div class="flex flex-col gap-2">
            <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">類別</label>
            <BaseSelect
              v-model="form.category"
              :options="categories"
            />
          </div>
        </div>

        <!-- Items Section -->
        <div class="flex flex-col gap-3">
          <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1 flex justify-between">
            <span>項目明細 ({{ form.currency }})</span>
            <span class="text-indigo-400">總計 {{ totalAmount.toLocaleString() }}</span>
          </label>
          
          <div class="flex flex-col gap-3">
            <div v-for="(item, index) in form.items" :key="index" class="bg-neutral-900 rounded-xl p-4 border border-neutral-800 flex flex-col gap-4 relative">
              <BaseInput
                v-model="item.name"
                placeholder="項目名稱"
                input-class="font-bold"
              />
              <div class="flex items-center gap-3">
                <div class="flex-1 flex flex-col gap-1">
                  <label class="text-[10px] font-bold text-neutral-500 uppercase px-1">單價</label>
                  <BaseInput
                    v-model.number="item.unit_price"
                    type="number"
                    placeholder="單價"
                  />
                </div>
                <div class="w-20 flex flex-col gap-1">
                  <label class="text-[10px] font-bold text-neutral-500 uppercase px-1">數量</label>
                  <BaseInput
                    v-model.number="item.quantity"
                    type="number"
                    placeholder="數量"
                    min="1"
                  />
                </div>
              </div>
              
              <div class="flex items-end gap-3">
                <div class="flex-1 flex flex-col gap-1">
                  <label class="text-[10px] font-bold text-neutral-500 uppercase px-1">折扣</label>
                  <BaseInput
                    v-model.number="item.discount"
                    type="number"
                    placeholder="折扣"
                  />
                </div>
                <div class="flex-1 text-right pb-3">
                  <div class="text-[10px] font-bold text-neutral-500 uppercase px-1 mb-1">小計</div>
                  <div class="text-lg font-bold text-white">
                    {{ ((Number(item.unit_price) - Number(item.discount || 0)) * Number(item.quantity)).toLocaleString() }}
                  </div>
                </div>
                <button 
                  v-if="form.items.length > 1"
                  type="button" 
                  @click="removeItem(index)" 
                  class="mb-2 w-10 h-10 rounded-lg bg-red-500/10 text-red-500 flex items-center justify-center border-0 cursor-pointer hover:bg-red-500/20 active:scale-95 transition-transform"
                >
                  <Icon icon="mdi:close" />
                </button>
              </div>
            </div>
          </div>
          
          <button type="button" @click="addItem" class="flex justify-center items-center gap-2 p-4 border-2 border-dashed border-neutral-700 rounded-xl bg-transparent text-neutral-500 font-bold cursor-pointer hover:border-indigo-500/50 hover:text-indigo-400 transition-colors">
            <Icon icon="mdi:plus" class="text-xl" /> 新增項目
          </button>
        </div>

        <!-- Foreign Currency Settlement (Conditional) -->
        <div v-if="form.currency !== baseCurrency" class="bg-indigo-500/5 rounded-xl p-5 border border-indigo-500/10 flex flex-col gap-5">
          <div class="flex justify-between items-center">
            <h3 class="font-bold text-xs text-indigo-400 uppercase tracking-wider">外幣結算</h3>
            <label class="flex items-center gap-2 cursor-pointer select-none">
              <input type="checkbox" v-model="form.manual_rate" class="w-4 h-4 rounded text-indigo-500 bg-neutral-800 border-neutral-700">
              <span class="text-xs font-bold text-neutral-400">手動輸入匯率</span>
            </label>
          </div>

          <div v-if="form.manual_rate" class="flex flex-col gap-4">
             <div class="flex flex-col gap-2">
               <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">匯率 (1 TWD = ? 外幣)</label>
               <BaseInput
                  v-model.number="form.exchange_rate"
                  type="number"
                  step="0.0001"
               />
             </div>
             <div class="flex flex-col gap-2">
               <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">銀行手續費 (TWD)</label>
               <BaseInput
                  v-model.number="form.handling_fee"
                  type="number"
               />
             </div>
             <div class="flex flex-col gap-2">
                <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">折合台幣 (含手續費)</label>
                <div class="w-full bg-neutral-800/50 border border-neutral-700 border-dashed text-white py-3 px-4 rounded-xl text-xl font-bold">
                   TWD {{ (calculatedBillingAmount + Number(form.handling_fee)).toLocaleString() }}
                </div>
             </div>
          </div>

          <div v-else class="flex flex-col gap-4">
             <div class="flex flex-col gap-2">
               <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">銀行入帳金額 (TWD，不含手續費)</label>
               <BaseInput
                  v-model.number="form.billing_amount"
                  type="number"
               />
             </div>
             <div class="flex flex-col gap-2">
               <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">銀行手續費 (TWD)</label>
               <BaseInput
                  v-model.number="form.handling_fee"
                  type="number"
               />
             </div>
             <div class="flex flex-col gap-2">
                <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">換算匯率</label>
                <div class="w-full bg-neutral-800/50 border border-neutral-700 border-dashed text-indigo-400 py-3 px-4 rounded-xl text-xl font-bold">
                   {{ calculatedExchangeRate }}
                </div>
             </div>
          </div>
        </div>

        <!-- Note -->
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">備註</label>
          <BaseTextarea 
            v-model="form.note" 
            placeholder="選填" 
            rows="3"
          />
        </div>

        <!-- Actions -->
        <div class="flex flex-col gap-3 mt-4">
          <button 
            type="submit" 
            class="w-full py-4 rounded-xl font-bold bg-indigo-500 text-white hover:bg-indigo-600 transition-all active:scale-95 disabled:opacity-50 border-0 cursor-pointer shadow-lg" 
            :disabled="submitting"
          >
            {{ submitting ? '儲存中...' : '儲存變更' }}
          </button>
          
          <button 
            type="button"
            @click="handleDelete"
            class="w-full py-4 rounded-xl font-bold bg-red-500/10 text-red-500 hover:bg-red-500/20 transition-all active:scale-95 border-0 cursor-pointer"
          >
            刪除此筆交易
          </button>
        </div>
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

const loading = ref(true)
const submitting = ref(false)
const baseCurrency = ref('TWD')

const form = ref({
  date: new Date(),
  category: '',
  note: '',
  items: [{ name: '', unit_price: 0, quantity: 1, discount: 0 }],
  currency: 'TWD',
  manual_rate: true, 
  exchange_rate: 0,
  billing_amount: 0,
  handling_fee: 0
})

const totalAmount = computed(() => {
  return form.value.items.reduce((sum, item) => {
    const subtotal = (Number(item.unit_price) - Number(item.discount || 0)) * Number(item.quantity)
    return sum + subtotal
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
  form.value.items.push({ name: '', unit_price: 0, quantity: 1, discount: 0 })
}

const removeItem = (index: number) => {
  form.value.items.splice(index, 1)
}

const fetchData = async () => {
  try {
    const [txn, space] = await Promise.all([
      api.get<any>(`/api/spaces/${route.params.id}/transactions/${route.params.txnId}`),
      api.get<any>(`/api/spaces/${route.params.id}`)
    ])
    
    baseCurrency.value = space.base_currency || 'TWD'
    
    const itemsData = txn.items || []
    
    form.value = {
      date: new Date(txn.date),
      category: txn.category,
      note: txn.note,
      items: itemsData.length > 0 ? itemsData : [{ name: '', unit_price: 0, quantity: 1, discount: 0 }],
      currency: txn.currency,
      manual_rate: txn.handling_fee === 0 || !txn.handling_fee,
      exchange_rate: txn.exchange_rate,
      billing_amount: txn.billing_amount,
      handling_fee: txn.handling_fee || 0
    }
  } catch (e) {
    console.error('Failed to fetch transaction:', e)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    const payload = {
      date: form.value.date.toISOString(),
      category: form.value.category,
      note: form.value.note,
      currency: form.value.currency,
      items: form.value.items
        .filter(item => item.name && Number(item.unit_price) > 0)
        .map(item => ({
          ...item,
          amount: (Number(item.unit_price) - Number(item.discount || 0)) * Number(item.quantity)
        })),
      exchange_rate: form.value.currency === baseCurrency.value ? 1 : form.value.exchange_rate,
      billing_amount: form.value.currency === baseCurrency.value ? totalAmount.value : form.value.billing_amount,
      handling_fee: form.value.currency === baseCurrency.value ? 0 : form.value.handling_fee,
      total_amount: totalAmount.value
    }

    await api.put(`/api/spaces/${route.params.id}/transactions/${route.params.txnId}`, payload)
    router.push(`/spaces/${route.params.id}`)
  } catch (e: any) {
    alert(e.message || '更新失敗')
  } finally {
    submitting.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('確定要刪除這筆交易嗎？')) return
  
  try {
    await api.delete(`/api/spaces/${route.params.id}/transactions/${route.params.txnId}`)
    router.push(`/spaces/${route.params.id}`)
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
