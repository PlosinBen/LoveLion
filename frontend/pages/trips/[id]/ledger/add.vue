<template>
  <div class="add-transaction-page flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">新增旅行記帳</h1>
      <div style="width: 40px;"></div>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <form v-else @submit.prevent="handleSubmit" class="flex flex-col gap-5">
      <!-- Date -->
      <!-- Date & Time -->
      <!-- Date & Time -->
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

      <!-- Currency -->
      <div class="flex flex-col gap-2">
        <label class="block mb-1 text-sm text-neutral-400">幣別</label>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="cur in currencies"
            :key="cur"
            type="button"
            class="px-4 py-2.5 rounded-3xl border border-neutral-800 bg-neutral-900 text-white cursor-pointer transition-all duration-200 hover:border-indigo-500"
            :class="{ 'bg-indigo-500 border-indigo-500': form.currency === cur }"
            @click="form.currency = cur"
          >
            {{ cur }}
          </button>
        </div>
      </div>

      <!-- Category -->
      <div class="flex flex-col gap-2">
        <label class="block mb-1 text-sm text-neutral-400">類別</label>
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
        <label class="block mb-1 text-sm text-neutral-400">項目明細</label>
        <div class="flex flex-col gap-3">
          <div v-for="(item, index) in form.items" :key="index" class="bg-neutral-900 rounded-2xl p-4 border border-neutral-800">
            <div class="flex flex-col gap-3">
              <div class="flex items-end gap-2">
                <div class="flex-1">
                  <BaseInput
                    v-model="item.name"
                    placeholder="項目名稱"
                    label="項目名稱"
                    input-class="font-medium"
                  />
                </div>
                <button type="button" @click="removeItem(index)" class="flex justify-center items-center w-9 h-9 rounded-lg text-neutral-400 hover:text-red-500 hover:bg-neutral-800 transition-colors" v-if="form.items.length > 1">
                  <Icon icon="mdi:close" class="text-xl" />
                </button>
              </div>
              <div class="grid grid-cols-3 gap-2">
                <div class="flex flex-col gap-1">
                  <label class="text-xs text-neutral-400">單價</label>
                  <BaseInput
                    v-model.number="item.unit_price"
                    type="number"
                    placeholder="0"
                  />
                </div>
                <!-- - -->
                <div class="flex flex-col gap-1">
                  <label class="text-xs text-neutral-400">折扣</label>
                   <BaseInput
                    v-model.number="item.discount"
                    type="number"
                    placeholder="0"
                  />
                </div>
                <!-- * -->
                <div class="flex flex-col gap-1">
                  <label class="text-xs text-neutral-400">數量</label>
                  <BaseInput
                    v-model.number="item.quantity"
                    type="number"
                    placeholder="1"
                    min="1"
                  />
                </div>
              </div>
            </div>
            <div class="flex justify-end items-center gap-2 mt-2 pt-2 border-t border-neutral-800">
               <span class="text-neutral-500 text-xs">( {{ item.unit_price || 0 }} - {{ item.discount || 0 }} ) × {{ item.quantity || 1 }} = </span>
               <span class="text-neutral-300 font-medium">{{ currency }} {{ ((item.unit_price - (item.discount || 0)) * item.quantity).toLocaleString() }}</span>
            </div>
          </div>
        </div>
        <button type="button" @click="addItem" class="flex justify-center items-center gap-1.5 p-3 border-2 border-dashed border-neutral-800 rounded-xl bg-transparent text-neutral-400 cursor-pointer mt-3 hover:border-indigo-500 hover:text-indigo-500 transition-colors">
          <Icon icon="mdi:plus" /> 新增項目
        </button>
      </div>

      <!-- Total -->
      <div class="flex justify-between items-center text-lg bg-neutral-900 rounded-2xl p-5 border border-neutral-800">
        <span>總計</span>
        <span class="text-2xl font-bold text-indigo-500">{{ currency }} {{ totalAmount.toLocaleString() }}</span>
      </div>

      <!-- Payer & Split -->
      <div class="bg-neutral-900 rounded-2xl p-5 border border-neutral-800 flex flex-col gap-4">
        <!-- Payer -->
        <div class="flex flex-col gap-2">
          <label class="text-sm text-neutral-400">誰付錢？</label>
          <BaseSelect 
            v-model="form.payer" 
            :options="trip.members.map((m: any) => ({ label: m.name, value: m.id }))"
          />
        </div>

        <hr class="border-t border-neutral-800 my-1" />

        <!-- Split Toggle -->
        <div class="flex flex-col gap-3">
            <div class="flex items-center justify-between">
                <label class="flex items-center gap-2 cursor-pointer">
                    <input type="checkbox" v-model="isSplitEnabled" class="w-5 h-5 accent-indigo-500 rounded" @change="handleSplitToggle" />
                    <span class="text-sm text-neutral-200">多人分攤 (均分/自訂)</span>
                </label>
                <div v-if="isSplitEnabled" class="text-xs text-neutral-400">
                    剩餘: <span :class="remainingAmount !== 0 ? 'text-red-500 font-bold' : 'text-green-500'">{{ formatCurrency(remainingAmount) }}</span>
                </div>
            </div>

             <!-- Members Split List -->
            <div v-if="isSplitEnabled" class="flex flex-col gap-3 mt-1 bg-neutral-800/30 p-3 rounded-xl">
                 <div class="flex justify-end mb-2">
                    <button type="button" @click="resetToEven" class="text-xs text-indigo-400 hover:text-indigo-300">
                        重設為均分
                    </button>
                 </div>
                 <div v-for="m in splitList" :key="m.name" class="flex items-center gap-3">
                     <input type="checkbox" v-model="m.involved" class="w-5 h-5 accent-indigo-500" @change="recalculateSplits(false)" />
                     <div class="flex-1">
                        <div class="text-base">{{ m.name }}</div>
                     </div>
                     <div class="w-28 relative">
                         <input 
                           v-model.number="m.customAmount" 
                           type="number" 
                           class="w-full px-2 py-1.5 rounded-lg border border-neutral-700 bg-neutral-800 text-white text-right focus:outline-none focus:border-indigo-500 text-sm"
                           :disabled="!m.involved"
                           @input="handleAmountInput"
                         />
                     </div>
                 </div>
            </div>
        </div>
      </div>

      <!-- Payment Method -->
      <div v-if="paymentMethods.length > 0" class="flex flex-col gap-2">
        <label class="block mb-1 text-sm text-neutral-400">付款方式</label>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="method in paymentMethods"
            :key="method"
            type="button"
            class="px-4 py-2.5 rounded-3xl border border-neutral-800 bg-neutral-900 text-white cursor-pointer transition-all duration-200 hover:border-indigo-500"
            :class="{ 'bg-indigo-500 border-indigo-500': form.payment_method === method }"
            @click="form.payment_method = method"
          >
            {{ method }}
          </button>
        </div>
      </div>

      <!-- Note -->
      <div class="flex flex-col gap-2">
        <label class="block mb-1 text-sm text-neutral-400">備註</label>
        <BaseTextarea 
          v-model="form.note" 
          rows="2" 
          placeholder="選填" 
        />
      </div>

      <!-- Images -->
      <div class="flex flex-col gap-2">
        <label class="block mb-1 text-sm text-neutral-400">照片 / 收據</label>
        <ImageManager 
          ref="imageManager"
          entity-id="temp-new-txn" 
          entity-type="transaction" 
          :allow-reorder="true"
          :instant-delete="false"
          :instant-upload="false"
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
import BaseSelect from '~/components/BaseSelect.vue'
import BaseTextarea from '~/components/BaseTextarea.vue'
import BaseInput from '~/components/BaseInput.vue'
import ImageManager from '~/components/ImageManager.vue'
import { VueDatePicker } from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const trip = ref<any>(null)
const loading = ref(true)
const submitting = ref(false)
const imageManager = ref<any>(null)
const isSplitEnabled = ref(false)

// hourOptions/minuteOptions/manual time logic removed

const defaultCategories = ['餐飲', '交通', '購物', '娛樂', '住宿', '生活', '其他']

const categories = computed(() => {
  const ledgerCategories = trip.value?.ledger?.categories
  return ledgerCategories?.length > 0 ? ledgerCategories : defaultCategories
})

const currencies = computed(() => {
  const ledgerCurrencies = trip.value?.ledger?.currencies
  return ledgerCurrencies?.length > 0 ? ledgerCurrencies : ['TWD']
})

const paymentMethods = computed(() => {
  return trip.value?.ledger?.payment_methods || []
})

const form = ref({
  date: new Date(),
  category: '餐飲',
  currency: '',
  payment_method: '',
  note: '',
  items: [{ name: '', unit_price: 0, quantity: 1, discount: 0 }],
  payer: '' // Member ID
})

const currency = computed(() => trip.value?.base_currency || 'TWD')

const totalAmount = computed(() => {
  return form.value.items.reduce((sum, item) => {
    return sum + ((item.unit_price - (item.discount || 0)) * item.quantity)
  }, 0)
})

const splitList = ref<{ name: string, involved: boolean, customAmount: number }[]>([])

const remainingAmount = computed(() => {
    if (!isSplitEnabled.value) return 0
    const allocated = splitList.value.reduce((sum, m) => sum + (m.involved ? (m.customAmount || 0) : 0), 0)
    return totalAmount.value - allocated
})

const addItem = () => {
  form.value.items.push({ name: '', unit_price: 0, quantity: 1, discount: 0 })
}

const removeItem = (index: number) => {
  form.value.items.splice(index, 1)
}

const resetToEven = () => {
    recalculateSplits(true)
}

const handleSplitToggle = () => {
    if (isSplitEnabled.value) {
        // Enable: Reset to even
        splitList.value.forEach(m => m.involved = true)
        recalculateSplits(true)
    }
}

const handleAmountInput = () => {
    // Just allow typing. Validation happens on submit/remaining check
}

const recalculateSplits = (forceEven: boolean) => {
  if (splitList.value.length === 0) return
  const involved = splitList.value.filter(m => m.involved)
  const count = involved.length
  if (count === 0) return

  const total = totalAmount.value
  const base = Math.floor(total / count)
  const remainder = total - (base * count)

  splitList.value.forEach(m => {
      if (m.involved) {
          m.customAmount = base
      } else {
          m.customAmount = 0
      }
  })
  
  // Distribute remainder to first involved
  let distributed = 0
  for (const m of splitList.value) {
      if (m.involved && distributed < remainder) {
          m.customAmount += 1
          distributed++
      }
  }
}

const formatCurrency = (val: number) => {
  return val.toLocaleString()
}

// Watch total amount to auto-update splits if they haven't been manually touched too much?
// Simple rule: If split is enabled, always re-run even split when total changes? 
// Or better: If total changes, we should probably warn or re-even. 
// For now let's watch total and re-even if it matches the sum of previous. 
// Actually, simple UX: create/add items -> total changes. Split amounts become stale.
// Let's re-run even split whenever total changes for simplicity in this version, unless user locked it?
// User said "Default even".
watch(totalAmount, () => {
    if (isSplitEnabled.value) {
        recalculateSplits(true)
    }
})


const fetchTrip = async () => {
  try {
    const data = await api.get<any>(`/api/trips/${route.params.id}`)
    trip.value = data

    // Initialize Split List (Ledger Members)
    const ledgerMembers: string[] = data.ledger?.members || []
    if (ledgerMembers.length > 0) {
        splitList.value = ledgerMembers.map(name => ({ name, involved: true, customAmount: 0 }))
    } else {
        // Fallback to Trip Members
        splitList.value = data.members.map((m: any) => ({ name: m.name, involved: true, customAmount: 0 }))
    }
    
    // Set default payer to current user
    const currentUser = await api.get<any>('/api/users/me')
    const me = data.members.find((m: any) => m.user_id === currentUser.id)
    if (me) {
      form.value.payer = me.id
    } else if (data.members.length > 0) {
      form.value.payer = data.members[0].id
    }
    
    // Set default currency from trip
    form.value.currency = data.base_currency || 'TWD'
    
    // Set default category from ledger categories
    if (data.ledger?.categories?.length > 0) {
      form.value.category = data.ledger.categories[0]
    }
    
    // Set default payment method from ledger
    if (data.ledger?.payment_methods?.length > 0) {
      form.value.payment_method = data.ledger.payment_methods[0]
    }
  } catch (e) {
    console.error('Failed to fetch trip:', e)
    router.push('/trips')
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  if (!form.value.items.some(item => item.name && item.unit_price > 0)) {
    alert('請至少填寫一個項目')
    return
  }
  if (!form.value.payer) {
    alert('請選擇付款人')
    return
  }
  
  // Validate total split if custom
  if (isSplitEnabled.value) {
    const customTotal = splitList.value
      .filter((m) => m.involved)
      .reduce((sum, m) => sum + (m.customAmount || 0), 0)
    
    if (Math.abs(customTotal - totalAmount.value) > 1) {
      alert(`分攤金額總和 (${customTotal}) 與 總金額 (${totalAmount.value}) 不符`)
      return
    }
  }

  submitting.value = true

  try {
    const payerMember = trip.value.members.find((m: any) => m.id === form.value.payer)
    
    // Construct Splits
    const splits = []
    
    // 1. Payer Split (Payment)
    // Payer is always a Trip Member (User) in this UI, but we record Name.
    splits.push({
      member_id: form.value.payer,
      name: payerMember?.name || 'Unknown',
      amount: totalAmount.value,
      is_payer: true
    })
    
    // 2. Consumer Splits (Expenses)
    if (isSplitEnabled.value) {
        splitList.value.forEach((m) => {
            if (m.involved && m.customAmount > 0) {
                splits.push({
                    name: m.name,
                    amount: m.customAmount,
                    is_payer: false
                })
            }
        })
    } else {
        // No split -> Payer is sole consumer
        splits.push({
            member_id: form.value.payer,
            name: payerMember?.name || 'Unknown',
            amount: totalAmount.value,
            is_payer: false
        })
    }
    

    const payload = {
      date: form.value.date.toISOString(),
      category: form.value.category,
      note: form.value.note,
      payer: payerMember?.name || 'Unknown', // Legacy string
      currency: form.value.currency || trip.value.base_currency,
      payment_method: form.value.payment_method,
      items: form.value.items.filter(item => item.name && item.unit_price > 0),
      splits: splits
    }

    const res = await api.post<{ id: string }>(`/api/ledgers/${trip.value.ledger_id}/transactions`, payload)
    
    // Commit Images with new ID
    if (imageManager.value) {
        await imageManager.value.commit(res.id)
    }

    router.push(`/trips/${trip.value.id}/ledger`)
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
