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

        <!-- Split Mode -->
        <div class="flex flex-col gap-2">
          <label class="text-sm text-neutral-400">分攤方式</label>
          <div class="flex gap-2 p-1 bg-neutral-800 rounded-xl">
             <button
               type="button" 
               class="flex-1 py-1.5 rounded-lg text-sm font-medium transition-all"
               :class="splitMode === 'even' ? 'bg-indigo-500 text-white shadow' : 'text-neutral-400 hover:text-white'"
               @click="splitMode = 'even'"
             >
               平分
             </button>
             <button
               type="button" 
               class="flex-1 py-1.5 rounded-lg text-sm font-medium transition-all"
               :class="splitMode === 'custom' ? 'bg-indigo-500 text-white shadow' : 'text-neutral-400 hover:text-white'"
               @click="splitMode = 'custom'"
             >
               自訂
             </button>
          </div>
        </div>

        <!-- Members Split -->
        <div class="flex flex-col gap-3">
           <div v-for="m in trip.members" :key="m.id" class="flex items-center gap-3">
             <input type="checkbox" v-model="m.involved" class="w-5 h-5 accent-indigo-500" @change="recalculateSplits" />
             <div class="flex-1">
               <div class="text-base">{{ m.name }}</div>
             </div>
             <div v-if="splitMode === 'even'" class="text-neutral-400 text-sm">
               {{ m.involved ? formatCurrency(evenSplitAmount) : '-' }}
             </div>
             <div v-else class="w-24">
               <input 
                 v-model.number="m.customAmount" 
                 type="number" 
                 class="w-full px-2 py-1.5 rounded-lg border border-neutral-700 bg-neutral-800 text-white text-right focus:outline-none focus:border-indigo-500 text-sm"
                 :disabled="!m.involved"
               />
             </div>
           </div>
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
const categories = ['餐飲', '交通', '購物', '娛樂', '住宿', '生活', '其他']
const splitMode = ref<'even' | 'custom'>('even')

// hourOptions and minuteOptions removed



const form = ref({
  date: new Date(),
  category: '餐飲',
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

const addItem = () => {
  form.value.items.push({ name: '', unit_price: 0, quantity: 1, discount: 0 })
}

const removeItem = (index: number) => {
  form.value.items.splice(index, 1)
}

const evenSplitAmount = computed(() => {
  if (!trip.value?.members) return 0
  const count = trip.value.members.filter((m: any) => m.involved).length
  if (count === 0) return 0
  return Math.round(totalAmount.value / count)
})

const recalculateSplits = () => {
}

const formatCurrency = (val: number) => {
  return val.toLocaleString()
}

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

    
    form.value.category = txn.category
    form.value.note = txn.note || ''
    form.value.items = txn.items.map((i: any) => ({
      name: i.name,
      unit_price: Number(i.unit_price),
      quantity: Number(i.quantity),
      discount: Number(i.discount || 0)
    }))
    
    // 4. Parse Splits
    // Initialize members
    tripData.members = tripData.members.map((m: any) => ({
      ...m,
      involved: false,
      customAmount: 0
    }))
    
    // Find Payer
    const payerSplit = txn.splits?.find((s: any) => s.is_payer)
    if (payerSplit) {
      form.value.payer = payerSplit.member_id
    }
    
    // Map Expenses
    const expenseSplits = txn.splits?.filter((s: any) => !s.is_payer) || []
    if (expenseSplits.length > 0) {
      // Check if even
      const amounts = expenseSplits.map((s: any) => Number(s.amount))
      const allSame = amounts.every((val: number, i: number, arr: number[]) => Math.abs(val - (arr[0] || 0)) < 1)
      splitMode.value = allSame ? 'even' : 'custom'
      
      expenseSplits.forEach((s: any) => {
        const member = tripData.members.find((m: any) => m.id === s.member_id)
        if (member) {
          member.involved = true
          member.customAmount = Number(s.amount)
        }
      })
    } else {
        // Fallback for old data or default
        tripData.members.forEach((m:any) => m.involved = true)
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
  if (!form.value.items.some(item => item.name && item.unit_price > 0)) {
    alert('請至少填寫一個項目')
    return
  }
  if (!form.value.payer) {
    alert('請選擇付款人')
    return
  }

  // Validate total split if custom
  if (splitMode.value === 'custom') {
    const customTotal = trip.value.members
      .filter((m: any) => m.involved)
      .reduce((sum: number, m: any) => sum + (m.customAmount || 0), 0)
    
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
      amount: totalAmount.value,
      is_payer: true
    })
    
    // 2. Consumer Splits (Expenses)
    const involvedMembers = trip.value.members.filter((m: any) => m.involved)
    const count = involvedMembers.length
    
    let allocated = 0
    involvedMembers.forEach((m: any, index: number) => {
      let amount = 0
      if (splitMode.value === 'even') {
        if (index === count - 1) {
           amount = totalAmount.value - allocated
        } else {
           amount = Math.floor(totalAmount.value / count)
        }
        allocated += amount
      } else {
        amount = m.customAmount || 0
      }
      
      splits.push({
        member_id: m.id,
        amount: amount,
        is_payer: false
      })
    })

    const payload = {
      payer: payerMember?.name || 'Unknown',
      date: form.value.date.toISOString(),
      category: form.value.category,
      note: form.value.note,
      currency: trip.value.base_currency,
      items: form.value.items.filter(item => item.name && item.unit_price > 0),
      splits: splits
    }

    await api.put(`/api/ledgers/${trip.value.ledger_id}/transactions/${transaction.value.id}`, payload)
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
