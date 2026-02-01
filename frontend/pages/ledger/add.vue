<template>
  <div class="add-transaction-page min-h-screen bg-neutral-900 text-neutral-50 p-4">
    <header class="flex justify-between items-center mb-6">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">{{ isEdit ? '編輯交易' : '新增記帳' }}</h1>
      <div style="width: 40px;"></div>
    </header>

    <form @submit.prevent="handleSubmit" class="flex flex-col gap-5">
      <!-- Date -->
      <VueDatePicker 
        v-model="form.date" 
        :dark="true"
        :formats="{input: 'yyyy-MM-dd HH:mm'}"
        :enable-seconds="false"
        time-picker-inline
        cancel-text="取消"
        select-text="確定"
        placeholder="日期與時間"
        class="date-picker-dark"
      />

      <!-- Currency Selection -->
      <div class="flex flex-col gap-2">
        <label class="block mb-2 text-sm text-neutral-400">幣別</label>
        <BaseSelect
          v-model="form.currency"
          :options="availableCurrencies"
          placeholder="選擇幣別"
        />
      </div>

      <!-- Category -->
      <div class="flex flex-col gap-2">
        <label class="block mb-2 text-sm text-neutral-400">類別</label>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="cat in categories"
            :key="cat"
            type="button"
            class="px-4 py-2.5 rounded-3xl border border-neutral-800 bg-neutral-900 text-white cursor-pointer transition-all duration-200 hover:border-indigo-500"
            :class="{ 'bg-indigo-500 border-indigo-500': form.category === cat }"
            @click="form.category = cat"
          >
            {{ cat }}
          </button>
        </div>
      </div>

      <!-- Items -->
      <div class="flex flex-col gap-2">
        <label class="block mb-2 text-sm text-neutral-400">項目明細 ({{ form.currency }})</label>
        <div class="flex flex-col gap-3">
          <div v-for="(item, index) in form.items" :key="index" class="bg-neutral-900 rounded-2xl p-4 border border-neutral-800">
            <div class="flex flex-col gap-2">
              <BaseInput
                v-model="item.name"
                placeholder="項目名稱"
                input-class="font-medium"
              />
              <div class="flex items-center gap-2">
                <BaseInput
                  v-model.number="item.unit_price"
                  type="number"
                  placeholder="單價"
                  input-class="flex-[2]"
                />
                <span class="text-neutral-400">×</span>
                <BaseInput
                  v-model.number="item.quantity"
                  type="number"
                  placeholder="數量"
                  min="1"
                  input-class="flex-1 w-[60px]"
                />
                <button type="button" @click="removeItem(index)" class="flex justify-center items-center w-9 h-9 rounded-lg bg-red-500 text-white border-0 cursor-pointer hover:bg-red-600 transition-colors" v-if="form.items.length > 1">
                  <Icon icon="mdi:close" />
                </button>
              </div>
            </div>
            <div class="text-right text-neutral-400 text-sm mt-3">
              {{ form.currency }} {{ (item.unit_price * item.quantity).toLocaleString() }}
            </div>
          </div>
        </div>
        <button type="button" @click="addItem" class="flex justify-center items-center gap-1.5 p-3 border-2 border-dashed border-neutral-800 rounded-xl bg-transparent text-neutral-400 cursor-pointer mt-3 hover:border-indigo-500 hover:text-indigo-500 transition-colors">
          <Icon icon="mdi:plus" /> 新增項目
        </button>
      </div>

      <!-- Total -->
      <div class="flex justify-between items-center text-lg bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
        <span>總計 ({{ form.currency }})</span>
        <span class="text-2xl font-bold text-indigo-500">{{ form.currency }} {{ totalAmount.toLocaleString() }}</span>
      </div>

      <!-- Foreign Currency Settlement -->
      <div v-if="form.currency !== 'TWD'" class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 flex flex-col gap-4">
        <div class="flex justify-between items-center">
          <h3 class="font-bold text-lg">外幣結算</h3>
          <label class="flex items-center gap-2 cursor-pointer select-none">
            <input type="checkbox" v-model="form.manual_rate" class="w-4 h-4 rounded text-indigo-500 focus:ring-indigo-500 bg-neutral-800 border-gray-600">
            <span class="text-sm">自行輸入匯率 (現金)</span>
          </label>
        </div>

        <div v-if="form.manual_rate" class="flex flex-col gap-4 animate-fade-in">
           <!-- Manual Rate Mode -->
           <BaseInput
              v-model.number="form.exchange_rate"
              type="number"
              label="匯率"
              step="0.0001"
              placeholder="1 TWD = ? Foreign"
           />
           <div class="flex justify-between items-center p-3 bg-neutral-800 rounded-xl">
              <span class="text-neutral-400">折合台幣</span>
              <span class="text-xl font-bold">TWD {{ calculatedBillingAmount.toLocaleString() }}</span>
           </div>
        </div>

        <div v-else class="flex flex-col gap-4 animate-fade-in">
           <!-- Auto Rate Mode (Credit Card) -->
           <BaseInput
              v-model.number="form.billing_amount"
              type="number"
              label="銀行入帳金額 (TWD)"
              placeholder="信用卡帳單上的台幣金額"
           />
           <BaseInput
              v-model.number="form.handling_fee"
              type="number"
              label="海外手續費 (TWD)"
              placeholder="選填"
           />
           <div class="flex justify-between items-center p-3 bg-neutral-800 rounded-xl">
              <span class="text-neutral-400">換算匯率</span>
              <span class="text-xl font-bold text-indigo-400">{{ calculatedExchangeRate }}</span>
           </div>
        </div>
      </div>

      <!-- Note -->
      <div class="flex flex-col gap-2">
        <label class="block mb-2 text-sm text-neutral-400">備註</label>
        <BaseTextarea 
          v-model="form.note" 
          rows="2" 
          placeholder="選填" 
        />
      </div>

      <!-- Submit -->
      <button type="submit" class="w-full mt-3 px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed" :disabled="submitting">
        {{ submitting ? '儲存中...' : '儲存' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { VueDatePicker } from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css'
import BaseSelect from '~/components/BaseSelect.vue'

definePageMeta({
  layout: 'empty'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const isEdit = computed(() => !!route.params.id)
const categories = ['餐飲', '交通', '購物', '娛樂', '生活', '其他']
const availableCurrencies = ref(['TWD', 'JPY', 'USD', 'EUR', 'KRW', 'THB', 'CNY', 'GBP']) // Default list, could be dynamic
const ledgerId = ref('')
const submitting = ref(false)

const form = ref({
  date: new Date(),
  category: '',
  note: '',
  items: [{ name: '', unit_price: 0, quantity: 1 }],
  currency: 'TWD',
  manual_rate: true, // true = input rate, false = input billing
  exchange_rate: 0,
  billing_amount: 0,
  handling_fee: 0
})

const totalAmount = computed(() => {
  return form.value.items.reduce((sum, item) => {
    return sum + (item.unit_price * item.quantity)
  }, 0)
})

// Calculated Billing Amount (For Manual Rate Mode)
const calculatedBillingAmount = computed(() => {
    if (form.value.manual_rate && form.value.exchange_rate > 0) {
        return Math.round(totalAmount.value * form.value.exchange_rate)
    }
    return 0
})

// Calculated Exchange Rate (For Auto Rate Mode)
const calculatedExchangeRate = computed(() => {
    if (!form.value.manual_rate && totalAmount.value > 0 && form.value.billing_amount > 0) {
        const netBilling = form.value.billing_amount - form.value.handling_fee
        if (netBilling <= 0) return 0
        return (netBilling / totalAmount.value).toFixed(4)
    }
    return 0
})

// Watchers to keep internal state consistent
watch(calculatedBillingAmount, (newVal) => {
    if (form.value.manual_rate) {
        form.value.billing_amount = newVal
        form.value.handling_fee = 0 // Reset fee in cash mode? Or optional? Usually cash has no separate fee, but exchange rate implies it.
    }
})

watch(calculatedExchangeRate, (newVal) => {
    if (!form.value.manual_rate) {
        form.value.exchange_rate = Number(newVal)
    }
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
  
  // Validation for foreign currency
  if (form.value.currency !== 'TWD') {
      if (form.value.manual_rate && form.value.exchange_rate <= 0) {
          alert('請輸入有效的匯率')
          return
      }
      if (!form.value.manual_rate && form.value.billing_amount <= 0) {
          alert('請輸入銀行入帳金額')
          return
      }
  }

  submitting.value = true

  try {
    const payload = {
      date: form.value.date.toISOString(),
      category: form.value.category,
      note: form.value.note,
      currency: form.value.currency,
      items: form.value.items.filter(item => item.name && item.unit_price > 0),
      // Foreign currency fields
      exchange_rate: form.value.currency === 'TWD' ? 1 : form.value.exchange_rate,
      billing_amount: form.value.currency === 'TWD' ? Math.round(totalAmount.value) : form.value.billing_amount,
      handling_fee: form.value.currency === 'TWD' ? 0 : form.value.handling_fee,
      total_amount: totalAmount.value
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
      // Optional: Load supported currencies from ledger settings if implemented
      // if (ledgers[0].currencies) availableCurrencies.value = ledgers[0].currencies
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
.animate-fade-in {
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-5px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
