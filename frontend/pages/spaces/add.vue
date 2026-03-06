<template>
  <div class="add-transaction-page min-h-screen bg-neutral-900 text-neutral-50 p-4">
    <header class="flex justify-between items-center mb-6">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">{{ isEdit ? '蝺刻摩鈭斗?' : '?啣?閮董' }}</h1>
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
        cancel-text="??"
        select-text="蝣箏?"
        placeholder="?交?????
        class="date-picker-dark"
      />

      <!-- Currency Selection -->
      <div class="flex flex-col gap-2">
        <label class="block mb-2 text-sm text-neutral-400">撟?</label>
        <BaseSelect
          v-model="form.currency"
          :options="availableCurrencies"
          placeholder="?豢?撟?"
        />
      </div>

      <!-- Category -->
      <div class="flex flex-col gap-2">
        <label class="block mb-2 text-sm text-neutral-400">憿</label>
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
        <label class="block mb-2 text-sm text-neutral-400">??敦 ({{ form.currency }})</label>
        <div class="flex flex-col gap-3">
          <div v-for="(item, index) in form.items" :key="index" class="bg-neutral-900 rounded-2xl p-4 border border-neutral-800">
            <div class="flex flex-col gap-2">
              <BaseInput
                v-model="item.name"
                placeholder="??迂"
                input-class="font-medium"
              />
              <div class="flex items-center gap-2">
                <BaseInput
                  v-model.number="item.unit_price"
                  type="number"
                  placeholder="?桀"
                  input-class="flex-auto"
                />
                <span class="text-neutral-400">?</span>
                <BaseInput
                  v-model.number="item.quantity"
                  type="number"
                  placeholder="?賊?"
                  min="1"
                  input-class="flex-1 w-14"
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
          <Icon icon="mdi:plus" /> ?啣??
        </button>
      </div>

      <!-- Total -->
      <div class="flex justify-between items-center text-lg bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
        <span>蝮質? ({{ form.currency }})</span>
        <span class="text-2xl font-bold text-indigo-500">{{ form.currency }} {{ totalAmount.toLocaleString() }}</span>
      </div>

      <!-- Foreign Currency Settlement -->
      <div v-if="form.currency !== 'TWD'" class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 flex flex-col gap-4">
        <div class="flex justify-between items-center">
          <h3 class="font-bold text-lg">憭馳蝯?</h3>
          <label class="flex items-center gap-2 cursor-pointer select-none">
            <input type="checkbox" v-model="form.manual_rate" class="w-4 h-4 rounded text-indigo-500 focus:ring-indigo-500 bg-neutral-800 border-gray-600">
            <span class="text-sm">?芾?頛詨?舐? (?暸?)</span>
          </label>
        </div>

        <div v-if="form.manual_rate" class="flex flex-col gap-4 animate-fade-in">
           <!-- Manual Rate Mode -->
           <BaseInput
              v-model.number="form.exchange_rate"
              type="number"
              label="?舐?"
              step="0.0001"
              placeholder="1 TWD = ? Foreign"
           />
           <div class="flex justify-between items-center p-3 bg-neutral-800 rounded-xl">
              <span class="text-neutral-400">???啣馳</span>
              <span class="text-xl font-bold">TWD {{ calculatedBillingAmount.toLocaleString() }}</span>
           </div>
        </div>

        <div v-else class="flex flex-col gap-4 animate-fade-in">
           <!-- Auto Rate Mode (Credit Card) -->
           <BaseInput
              v-model.number="form.billing_amount"
              type="number"
              label="?銵撣喲?憿?(TWD)"
              placeholder="靽∠?∪董?桐??撟??憿?
           />
           <BaseInput
              v-model.number="form.handling_fee"
              type="number"
              label="瘚瑕???鞎?(TWD)"
              placeholder="?詨‵"
           />
           <div class="flex justify-between items-center p-3 bg-neutral-800 rounded-xl">
              <span class="text-neutral-400">???舐?</span>
              <span class="text-xl font-bold text-indigo-400">{{ calculatedExchangeRate }}</span>
           </div>
        </div>
      </div>

      <!-- Note -->
      <div class="flex flex-col gap-2">
        <label class="block mb-2 text-sm text-neutral-400">?酉</label>
        <BaseTextarea 
          v-model="form.note" 
          rows="2" 
          placeholder="?詨‵" 
        />
      </div>

      <!-- Submit -->
      <button type="submit" class="w-full mt-3 px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed" :disabled="submitting">
        {{ submitting ? '?脣?銝?..' : '?脣?' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { useSpace } from '~/composables/useSpace'
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
const { currentSpace, fetchSpaces } = useSpace()

const isEdit = computed(() => !!route.params.id)
const categories = ['擗ㄡ', '鈭日?, '鞈潛', '憡?', '?暑', '?嗡?']
const availableCurrencies = ref(['TWD', 'JPY', 'USD', 'EUR', 'KRW', 'THB', 'CNY', 'GBP'])
const submitting = ref(false)

const ledgerId = computed(() => (route.query.ledger_id as string) || currentSpace.value?.id)

const form = ref({
  date: new Date(),
  category: '',
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
    return sum + (item.unit_price * item.quantity)
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
        const netBilling = form.value.billing_amount - form.value.handling_fee
        if (netBilling <= 0) return 0
        return (netBilling / totalAmount.value).toFixed(4)
    }
    return 0
})

watch(calculatedBillingAmount, (newVal) => {
    if (form.value.manual_rate) {
        form.value.billing_amount = newVal
        form.value.handling_fee = 0
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
  if (!ledgerId.value) {
    alert('蝛粹? ID ?箏仃')
    return
  }
  
  if (!form.value.items.some(item => item.name && item.unit_price > 0)) {
    alert('隢撠‵撖思?????)
    return
  }
  
  if (form.value.currency !== 'TWD') {
      if (form.value.manual_rate && form.value.exchange_rate <= 0) {
          alert('隢撓?交????舐?')
          return
      }
      if (!form.value.manual_rate && form.value.billing_amount <= 0) {
          alert('隢撓?仿?銵撣喲?憿?)
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
      exchange_rate: form.value.currency === 'TWD' ? 1 : form.value.exchange_rate,
      billing_amount: form.value.currency === 'TWD' ? Math.round(totalAmount.value) : form.value.billing_amount,
      handling_fee: form.value.currency === 'TWD' ? 0 : form.value.handling_fee,
      total_amount: totalAmount.value
    }

    await api.post(`/api/spaces/${ledgerId.value}/transactions`, payload)
    router.push(`/spaces/${ledgerId.value}`)
  } catch (e: any) {
    alert(e.message || '?脣?憭望?')
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
  await fetchSpaces()
  if (currentSpace.value?.base_currency) {
    form.value.currency = currentSpace.value.base_currency
  }
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
