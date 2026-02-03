<template>
  <div class="edit-transaction-page flex flex-col gap-6">
    <header class="flex justify-between items-center">
      <button @click="router.back()" class="flex justify-center items-center w-10 h-10 rounded-xl bg-neutral-900 text-white border-0 cursor-pointer hover:bg-neutral-800 transition-colors">
        <Icon icon="mdi:arrow-left" class="text-2xl" />
      </button>
      <h1 class="text-xl font-bold">編輯交易</h1>
      <button @click="handleDelete" class="flex justify-center items-center w-10 h-10 rounded-xl bg-red-500/10 text-red-500 border-0 cursor-pointer hover:bg-red-500/20 transition-colors">
        <Icon icon="mdi:trash-can-outline" class="text-xl" />
      </button>
    </header>

    <div v-if="loading" class="text-center text-neutral-400 p-10">載入中...</div>

    <form v-else @submit.prevent="handleSubmit" class="flex flex-col gap-5">
      <!-- Date & Time -->
      <!-- Core Info Card -->
      <div class="bg-neutral-900 rounded-2xl p-4 border border-neutral-800 flex flex-col gap-4">
         <div>
            <label class="text-xs text-neutral-400">日期</label>
            <VueDatePicker 
                v-model="form.date" 
                :dark="true"
                :formats="{input: 'yyyy-MM-dd HH:mm'}"
                :enable-seconds="false"
                time-picker-inline
                cancel-text="取消"
                select-text="確定"
                placeholder="日期與時間"
                :ui="{input: '!border-neutral-800 !hover:border-gray-400'}"
            />
         </div>

         <div class="flex gap-3">
             <div class="flex-1">
                <BaseInput
                    v-model="form.title"
                    placeholder="例如：第一天晚餐 (必填)"
                    label="標題"
                />
             </div>
         </div>

         <div class="flex items-end gap-3">
             <div class="flex-1">
                <BaseSelect
                  v-model="form.category"
                  :options="categories"
                  label="類別"
                />
             </div>
             <div class="flex-1">
                <BaseSelect
                  v-model="form.currency"
                  :options="currencies"
                  label="幣別"
                />
             </div>
             <div class="flex-1">
                <BaseInput
                    v-model.number="form.manualAmount"
                    type="number"
                    placeholder="0"
                    label="總金額"
                    input-class="font-mono text-right"
                />
             </div>
         </div>

         <div class="flex gap-3">
             <div class="flex-1">
                <BaseSelect 
                  v-model="form.payer" 
                  :options="trip.members.map((m: any) => ({ label: m.name, value: m.id }))"
                  label="付款人"
                />
             </div>
             <div v-if="paymentMethods.length > 0" class="w-1/3 min-w-28">
                <BaseSelect
                  v-model="form.payment_method"
                  :options="paymentMethods"
                  placeholder="付款方式"
                  label="付款方式"
                />
             </div>
         </div>
      </div>


      <!-- Items -->
      <div class="flex flex-col gap-2">
        <label class="text-xs text-neutral-400">項目明細</label>
        <div class="flex flex-col rounded-2xl bg-neutral-900 border border-neutral-800 overflow-hidden">
             <div v-for="(item, index) in form.items" :key="index" class="border-b border-neutral-800 last:border-0 transition-all duration-300">
                <BaseCollapsible v-model="item.expanded">
                    <!-- Header (Summary) -->
                    <template #header="{ open }">
                        <div class="flex items-center gap-3 flex-1 min-w-0">
                             <Icon :icon="open ? 'mdi:chevron-down' : 'mdi:chevron-right'" class="text-xl text-neutral-500 shrink-0" />
                             <span class="font-medium truncate" :class="item.name ? 'text-white' : 'text-neutral-500'">{{ item.name || '未命名項目' }}</span>
                        </div>
                        <div class="flex items-center gap-3">
                             <div class="flex flex-col items-end">
                                <span class="text-xs text-neutral-500 font-mono">
                                     (${{ item.unit_price }} - {{ item.discount || 0 }}) x {{ item.quantity }}
                                </span>
                                <span class="font-mono text-indigo-400">{{ currency }} {{ ((item.unit_price - (item.discount || 0)) * item.quantity).toLocaleString() }}</span>
                             </div>
                             <button type="button" @click.stop="removeItem(index)" class="text-neutral-500 hover:text-red-500 transition-colors p-1" v-if="form.items.length > 1">
                                <Icon icon="mdi:close" class="text-lg" />
                             </button>
                        </div>
                    </template>

                    <!-- Body (Inputs) -->
                    <div class="p-4 gap-3 flex flex-col">
                        <BaseInput
                            v-model="item.name"
                            placeholder="項目名稱"
                            label="項目名稱"
                        />
                        
                        <div class="flex gap-3">
                            <div class="flex-1">
                                <BaseInput
                                    v-model.number="item.unit_price"
                                    type="number"
                                    placeholder="0"
                                    label="單價"
                                    input-class="text-right"
                                />
                            </div>
                            <div class="flex-1">
                                <BaseInput
                                    v-model.number="item.discount"
                                    type="number"
                                    placeholder="0"
                                    label="折扣"
                                    input-class="text-right text-red-400"
                                />
                            </div>
                            <div class="w-24">
                                <BaseInput
                                    v-model.number="item.quantity"
                                    type="number"
                                    placeholder="1"
                                    label="數量"
                                    input-class="text-right"
                                />
                            </div>
                        </div>
                    </div>
                </BaseCollapsible>
             </div>
             
             <!-- Add Item Button inside the list container at the bottom -->
             <button type="button" @click="addItem" class="flex justify-center items-center gap-2 p-3 text-sm text-neutral-400 hover:text-indigo-400 hover:bg-neutral-800/50 transition-colors">
                <Icon icon="mdi:plus" /> 新增項目
            </button>
        </div>
      </div>
      
      <!-- Split -->
      <div class="bg-neutral-900 rounded-2xl border border-neutral-800 overflow-hidden transition-all duration-300">
        <BaseCollapsible :default-open="isSplitEnabled">
             <template #header="{ open }">
                <div class="flex items-center gap-2">
                     <Icon :icon="open ? 'mdi:chevron-down' : 'mdi:chevron-right'" class="text-xl text-neutral-400" />
                     <h3 class="font-bold">分帳設定</h3>
                </div>
                
                <div v-if="!open" class="text-sm text-neutral-400">
                    <span v-if="isSplitEnabled" class="text-indigo-400">已啟用拆帳</span>
                    <span v-else>平均分攤</span>
                </div>
            </template>

            <div class="p-4 flex flex-col gap-4">
                 <div class="flex items-center justify-between">
                    <label class="flex items-center gap-2 cursor-pointer select-none">
                        <input type="checkbox" v-model="isSplitEnabled" class="w-4 h-4 rounded text-indigo-500 focus:ring-indigo-500 bg-neutral-800 border-gray-600" @change="handleSplitToggle" />
                        <span class="text-sm">啟用進階拆帳</span>
                    </label>
                    <div v-if="isSplitEnabled" class="text-xs text-neutral-400">
                        剩餘: <span :class="remainingAmount !== 0 ? 'text-red-500 font-bold' : 'text-green-500'">{{ formatCurrency(remainingAmount) }}</span>
                    </div>
                 </div>

                 <!-- Members Split List -->
                 <div v-if="isSplitEnabled" class="flex flex-col gap-3 bg-neutral-950/50 p-4 rounded-xl border border-neutral-800 animate-fade-in">
                     <div class="flex justify-end mb-2">
                        <button type="button" @click="resetToEven" class="text-xs text-indigo-400 hover:text-indigo-300">
                            重設為均分
                        </button>
                     </div>
                     <div v-for="m in splitList" :key="m.name" class="flex items-center gap-3">
                         <input type="checkbox" v-model="m.involved" class="w-5 h-5 accent-indigo-500 rounded" @change="recalculateSplits(false)" />
                         <div class="flex-1">
                            <div class="text-base font-medium overflow-hidden text-ellipsis whitespace-nowrap">{{ m.name }}</div>
                         </div>
                         <div class="w-32">
                             <BaseInput 
                               v-model.number="m.customAmount" 
                               type="number" 
                               input-class="text-right font-mono"
                               :disabled="!m.involved"
                               placeholder="0"
                               @input="handleAmountInput"
                             />
                         </div>
                     </div>
                 </div>
            </div>
        </BaseCollapsible>
      </div>







      <!-- Foreign Currency Settlement -->
      <div v-if="form.currency && form.currency !== (trip?.base_currency || 'TWD')" class="bg-neutral-900 rounded-2xl border border-neutral-800 overflow-hidden transition-all duration-300">
        <BaseCollapsible v-model="isSettlementExpanded">
            <template #header="{ open }">
                <div class="flex items-center gap-2">
                     <Icon :icon="open ? 'mdi:chevron-down' : 'mdi:chevron-right'" class="text-xl text-neutral-400" />
                     <h3 class="font-bold">外幣結算</h3>
                </div>
                
                 <!-- Show summary if collapsed and has value -->
                 <div v-if="!open" class="text-sm text-neutral-400">
                    <span v-if="form.manual_rate && form.exchange_rate > 0 && form.exchange_rate !== 1">匯率: {{ form.exchange_rate }}</span>
                    <span v-else-if="!form.manual_rate && form.billing_amount > 0">已填寫</span>
                    <span v-else>未設定</span>
                 </div>
            </template>
            
            <div class="p-4 flex flex-col gap-4">
                 <div class="flex justify-end">
                    <label class="flex items-center gap-2 cursor-pointer select-none">
                        <input type="checkbox" v-model="form.manual_rate" class="w-4 h-4 rounded text-indigo-500 focus:ring-indigo-500 bg-neutral-800 border-gray-600">
                        <span class="text-sm">自行輸入匯率 (現金)</span>
                    </label>
                 </div>

                 <div v-if="form.manual_rate" class="flex flex-col gap-4">
                    <!-- Manual Rate Mode -->
                    <BaseInput
                       v-model.number="form.exchange_rate"
                       type="number"
                       label="匯率"
                       step="0.0001"
                       :placeholder="`1 ${trip?.ledger?.base_currency || trip?.base_currency || 'TWD'} = ? ${form.currency}`"
                    />
                    <div class="flex justify-between items-center p-2 bg-neutral-800 rounded">
                       <span class="text-neutral-400">折合 {{ trip?.ledger?.base_currency || trip?.base_currency || 'TWD' }}</span>
                       <span class="font-bold">{{ trip?.ledger?.base_currency || trip?.base_currency || 'TWD' }} {{ calculatedBillingAmount.toLocaleString() }}</span>
                    </div>
                 </div>

            <div v-else class="flex flex-col gap-4">
               <!-- Auto Rate Mode (Credit Card) -->
               <BaseInput
                  v-model.number="form.billing_amount"
                  type="number"
                  :label="`銀行入帳金額 (${trip?.ledger?.base_currency || trip?.base_currency || 'TWD'})`"
                  placeholder="信用卡帳單金額"
               />
               <BaseInput
                  v-model.number="form.handling_fee"
                  type="number"
                  :label="`海外手續費 (${trip?.ledger?.base_currency || trip?.base_currency || 'TWD'})`"
                  placeholder="選填"
               />
               <div class="flex justify-between items-center p-2 bg-neutral-800 rounded">
                  <span class="text-neutral-400">換算匯率</span>
                  <span class="text-xl font-bold text-indigo-400">{{ calculatedExchangeRate }}</span>
               </div>
            </div>
            </div>
        </BaseCollapsible>
      </div>



      <!-- Images -->
      <div class="bg-neutral-900 rounded-2xl p-4 border border-neutral-800 flex flex-col gap-2">
        <label class="block text-xs text-neutral-400">照片 / 收據</label>
        <ImageManager 
          ref="imageManager"
          :entity-id="route.params.txnId as string" 
          entity-type="transaction" 
          :allow-reorder="true"
          :instant-delete="false"
          :instant-upload="false"
        />
      </div>

      <!-- Submit -->
      <button type="submit" class="w-full mt-3 px-6 py-3 rounded-xl font-semibold bg-indigo-500 text-white hover:bg-indigo-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed" :disabled="submitting">
        {{ submitting ? '更新' : '更新' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import BaseSelect from '~/components/BaseSelect.vue'
import BaseTextarea from '~/components/BaseTextarea.vue'
import BaseInput from '~/components/BaseInput.vue'
import ImageManager from '~/components/ImageManager.vue'
import BaseCollapsible from '~/components/BaseCollapsible.vue'
import { VueDatePicker } from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css'

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()

const trip = ref<any>(null)
const transaction = ref<any>(null)
const loading = ref(true)
const submitting = ref(false)
const imageManager = ref<any>(null)
const currencies = computed(() => {
  const ledgerCurrencies = trip.value?.ledger?.currencies
  return ledgerCurrencies?.length > 0 ? ledgerCurrencies : ['TWD']
})
const paymentMethods = computed(() => {
  return trip.value?.ledger?.payment_methods || []
})
const categories = ['餐飲', '交通', '購物', '娛樂', '住宿', '生活', '其他']
const isSplitEnabled = ref(false)
const isSettlementExpanded = ref(false)

// hourOptions and minuteOptions removed



const form = ref({
  date: new Date(),
  title: '',
  category: '餐飲',
  currency: '',
  payment_method: '', // Add payment_method if missing in original
  note: '',
  manualAmount: 0,
  items: [{ name: '', unit_price: 0, quantity: 1, discount: 0, expanded: true }],
  payer: '', // Member ID
  manual_rate: true,
  exchange_rate: 0,
  billing_amount: 0,
  handling_fee: 0
})

const currency = computed(() => form.value.currency || trip.value?.ledger?.base_currency || trip.value?.base_currency || 'TWD')

const hasItems = computed(() => {
    return form.value.items.some(item => item.name && item.unit_price > 0)
})

const totalAmount = computed(() => {
    return form.value.manualAmount
})

const itemsTotal = computed(() => {
    return form.value.items.reduce((sum, item) => {
        return sum + ((item.unit_price - (item.discount || 0)) * item.quantity)
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

// Watchers
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

const splitList = ref<{ name: string, involved: boolean, customAmount: number }[]>([])

const remainingAmount = computed(() => {
    if (!isSplitEnabled.value) return 0
    const allocated = splitList.value.reduce((sum, m) => sum + (m.involved ? (m.customAmount || 0) : 0), 0)
    return totalAmount.value - allocated
})

const addItem = () => {
  form.value.items.push({ name: '', unit_price: 0, quantity: 1, discount: 0, expanded: true })
}

const removeItem = (index: number) => {
  form.value.items.splice(index, 1)
}

const resetToEven = () => {
    recalculateSplits(true)
}

const handleSplitToggle = () => {
    if (isSplitEnabled.value) {
        splitList.value.forEach(m => m.involved = true)
        recalculateSplits(true)
    }
}

const handleAmountInput = () => {
    // Just allow typing
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
  
  // Distribute remainder
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

watch([totalAmount, itemsTotal], ([newTotal, newItemsTotal]) => {
    if (isSplitEnabled.value) {
        recalculateSplits(true)
    }
})

const fetchData = async () => {
  try {
    // 1. Fetch Trip
    const tripData = await api.get<any>(`/api/trips/${route.params.id}`)
    trip.value = tripData
    
    // 2. Fetch Transaction
    const txn = await api.get<any>(`/api/ledgers/${tripData.ledger_id}/transactions/${route.params.txnId}`)
    transaction.value = txn
    
    // 3. Populate Form
    form.value.date = new Date(txn.date)

    
    form.value.title = txn.title || ''
    form.value.category = txn.category
    form.value.currency = txn.currency || tripData.ledger?.base_currency || tripData.base_currency
    form.value.note = txn.note || ''
    form.value.payment_method = txn.payment_method || ''
    
    // Populate Foreign Currency logic
    form.value.exchange_rate = Number(txn.exchange_rate || 1)
    form.value.billing_amount = Number(txn.billing_amount || 0)
    form.value.handling_fee = Number(txn.handling_fee || 0)
    
    // Determine mode: If handling fee exists, it's likely Auto Rate (Credit Card)
    if (form.value.handling_fee > 0) {
        form.value.manual_rate = false
    } else {
        form.value.manual_rate = true
    }

    if (txn.items && txn.items.length > 0) {
      form.value.items = txn.items.map((i: any) => ({
        name: i.name,
        unit_price: Number(i.unit_price),
        quantity: Number(i.quantity),
        discount: Number(i.discount || 0),
        expanded: false
      }))
    } 
    
    // Always set manualAmount from txn.total_amount
    form.value.manualAmount = Number(txn.total_amount)
    
    // 4. Parse Splits
    // Initialize bases from Ledger Settings
    const baseNames = new Set<string>(tripData.ledger?.members || tripData.members.map((m: any) => m.name))
    
    // Add names from existing splits (History preservation)
    const expenseSplits = txn.splits?.filter((s: any) => !s.is_payer) || []
    expenseSplits.forEach((s: any) => {
        if (s.name) baseNames.add(s.name)
        // Fallback for old data with member_id but no name
        else if (s.member_id) {
           const tripMem = tripData.members.find((m:any) => m.id === s.member_id)
           if (tripMem) baseNames.add(tripMem.name)
        }
    })

    // Init list
    splitList.value = Array.from(baseNames).map(name => ({
        name,
        involved: false,
        customAmount: 0
    }))
    
    // Find Payer
    const payerSplit = txn.splits?.find((s: any) => s.is_payer)
    if (payerSplit) {
      // If payer has name, great. If only ID, map it.
      // Current UI binds payer value to Member ID (User). This is tricky.
      // Ideally Payer is also just a Name string? 
      // User requirement: "Ledger.Members is ... list".
      // But Payer selection is usually "Who paid?" -> User.
      // Use fallback: If payer ID is in TripMember, select it.
      if (payerSplit.member_id) {
          form.value.payer = payerSplit.member_id
      }
    }
    
    // Determine if split is enabled
    if (expenseSplits.length > 0) {
        let isSimple = false
        // Logic for simple split detection (One consumer, same person, same amount) - Name based now?
        // Or simplified: Just check expense count.
        // If 1 expense = payer?
        
        // Strict mapping
        expenseSplits.forEach((s: any) => {
            const row = splitList.value.find(m => m.name === s.name || (s.member_id && tripData.members.find((tm:any) => tm.id === s.member_id && tm.name === m.name)))
            if (row) {
                row.involved = true
                row.customAmount = Number(s.amount)
            }
        })
        
        // Check simple toggle logic
        // If only 1 involved and it matches Payer name and amount matches Payer amount...
        const involvedCount = splitList.value.filter(m => m.involved).length
        const payerName = tripData.members.find((m:any) => m.id === form.value.payer)?.name
        
        if (involvedCount === 1 && payerName) {
             const soleConsumer = splitList.value.find(m => m.involved)
             if (soleConsumer && soleConsumer.name === payerName && Math.abs(soleConsumer.customAmount - totalAmount.value) < 1) {
                 isSimple = true
             }
        }
        
        if (isSimple) {
             isSplitEnabled.value = false
             // Reset defaults just in case
             splitList.value.forEach(m => m.involved = true)
        } else {
             isSplitEnabled.value = true
        }

    } else {
        isSplitEnabled.value = false
        splitList.value.forEach(m => m.involved = true)
    }

    // Determine if settlement should be expanded
    if (form.value.currency !== (tripData.ledger?.base_currency || tripData.base_currency || 'TWD')) {
        if (form.value.manual_rate) {
            if (form.value.exchange_rate !== 1) {
                isSettlementExpanded.value = true
            }
        } else {
             if (form.value.billing_amount > 0) {
                 isSettlementExpanded.value = true
             }
        }
    }

  } catch (e) {
    console.error('Failed to fetch data:', e)
    router.push(`/trips/${route.params.id}/ledger`)
  } finally {
    loading.value = false
  }
}

const handleDelete = async () => {
  if (!confirm('確定要刪除此紀錄？')) return
  try {
    await api.del(`/api/ledgers/${trip.value.ledger_id}/transactions/${route.params.txnId}`)
    router.push(`/trips/${trip.value.id}/ledger`)
  } catch (e: any) {
    alert(e.message || '刪除失敗')
  }
}

const handleSubmit = async () => {
  if (totalAmount.value <= 0) {
    alert('總金額必須大於 0')
    return
  }
  if (!form.value.title.trim()) {
    alert('請輸入標題')
    return
  }
  if (!form.value.payer) {
    alert('請選擇付款人')
    return
  }

  // Validate Total vs Items
  if (hasItems.value) {
      if (Math.abs(itemsTotal.value - form.value.manualAmount) > 1) {
          if (!confirm(`總金額 (${form.value.manualAmount}) 與 項目明細總和 (${itemsTotal.value}) 不符，確定要儲存？`)) {
              return
          }
      }
  }

  const baseCurrency = trip.value?.ledger?.base_currency || trip.value?.base_currency || 'TWD'
  if (form.value.currency !== baseCurrency) {
      if (form.value.manual_rate && form.value.exchange_rate <= 0) {
          alert('請輸入有效的匯率')
          return
      }
      if (!form.value.manual_rate && form.value.billing_amount <= 0) {
          alert('請輸入銀行入帳金額')
          return
      }
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
      payer: payerMember?.name || 'Unknown',
      date: form.value.date.toISOString(),
      title: form.value.title,
      category: form.value.category,
      note: form.value.note,
      currency: form.value.currency || baseCurrency,
      payment_method: form.value.payment_method, // Include if present
      items: hasItems.value ? form.value.items
        .filter(item => item.name && item.unit_price > 0)
        .map(({ expanded, ...item }) => item) : [],
      
      total_amount: totalAmount.value,
      exchange_rate: form.value.currency === baseCurrency ? 1 : form.value.exchange_rate,
      billing_amount: form.value.currency === baseCurrency ? Math.round(totalAmount.value) : form.value.billing_amount,
      handling_fee: form.value.currency === baseCurrency ? 0 : form.value.handling_fee,

      splits: splits
    }

    await api.put(`/api/ledgers/${trip.value.ledger_id}/transactions/${transaction.value.id}`, payload)
    
    // Commit Image Deletes
    if (imageManager.value) {
        await imageManager.value.commit()
    }

    router.push(`/trips/${trip.value.id}/ledger`)
  } catch (e: any) {
    alert(e.message || '更新失敗')
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
  fetchData()
})
</script>
