<template>
  <div class="add-transaction-page">
    <PageTitle
      title="新增交易"
      :back-to="`/spaces/${route.params.id}/ledger`"
      :breadcrumbs="[{ label: detailStore.space?.name || '空間', to: `/spaces/${route.params.id}/ledger` }]"
    />

    <div>
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
              
              <!-- Row 1: Price Factors (Unit Price - Discount) -->
              <div class="grid grid-cols-2 gap-3">
                <div class="flex flex-col gap-1">
                  <label class="text-[10px] font-bold text-neutral-500 uppercase px-1">單價</label>
                  <BaseInput
                    v-model.number="item.unit_price"
                    type="number"
                    placeholder="單價"
                  />
                </div>
                <div class="flex flex-col gap-1">
                  <label class="text-[10px] font-bold text-neutral-500 uppercase px-1">折扣</label>
                  <BaseInput
                    v-model.number="item.discount"
                    type="number"
                    placeholder="折扣"
                  />
                </div>
              </div>
              
              <!-- Row 2: Result Factors (Quantity & Subtotal) -->
              <div class="flex items-end gap-3">
                <div class="w-24 flex flex-col gap-1">
                  <label class="text-[10px] font-bold text-neutral-500 uppercase px-1">數量</label>
                  <BaseInput
                    v-model.number="item.quantity"
                    type="number"
                    placeholder="數量"
                    min="1"
                  />
                </div>
                
                <div class="flex-1 text-right pb-3">
                  <div class="text-[10px] font-bold text-neutral-500 uppercase px-1 mb-1">項目小計</div>
                  <div class="text-xl font-bold text-indigo-400">
                    {{ ((Number(item.unit_price) - Number(item.discount || 0)) * Number(item.quantity)).toLocaleString() }}
                  </div>
                </div>

                <BaseButton 
                  v-if="form.items.length > 1"
                  type="button" 
                  @click="removeItem(index)" 
                  variant="danger"
                  class="!p-0 w-10 h-10 mb-1"
                >
                  <Icon icon="mdi:close" class="text-xl" />
                </BaseButton>
              </div>
            </div>
          </div>
          
          <BaseButton type="button" @click="addItem" variant="ghost" class="border border-dashed border-neutral-700 !text-neutral-500 hover:!text-indigo-400">
            <Icon icon="mdi:plus" class="text-xl" /> 新增項目
          </BaseButton>
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

          <!-- Manual Rate Mode -->
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

          <!-- Billing Amount Mode -->
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

        <!-- Images -->
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">收據 / 照片</label>
          <ImageManager
            entity-id="pending"
            entity-type="transaction"
            :instant-upload="false"
            :instant-delete="false"
            ref="imageManagerRef"
          />
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

        <!-- Submit -->
        <BaseButton 
          type="submit" 
          variant="primary"
          fullWidth
          class="mt-4"
          :loading="submitting"
        >
          儲存交易
        </BaseButton>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Icon } from '@iconify/vue'
import { VueDatePicker } from '@vuepic/vue-datepicker'
import '@vuepic/vue-datepicker/dist/main.css'
import BaseButton from '~/components/BaseButton.vue'
import BaseInput from '~/components/BaseInput.vue'
import BaseSelect from '~/components/BaseSelect.vue'
import BaseTextarea from '~/components/BaseTextarea.vue'
import PageTitle from '~/components/PageTitle.vue'
import ImageManager from '~/components/ImageManager.vue'
import { useSpaceDetailStore } from '~/stores/spaceDetail'

definePageMeta({
  path: '/spaces/:id/ledger/transaction/add',
  layout: 'default'
})

const router = useRouter()
const route = useRoute()
const api = useApi()
const { isAuthenticated, initAuth } = useAuth()
const detailStore = useSpaceDetailStore()

const categories = ref<{label: string, value: string}[]>([])
const availableCurrencies = ref<{label: string, value: string}[]>([])

const baseCurrency = ref('TWD')
const submitting = ref(false)
const imageManagerRef = ref<InstanceType<typeof ImageManager> | null>(null)

const form = ref({
  date: new Date(),
  category: '其他',
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

    const created = await api.post<any>(`/api/spaces/${route.params.id}/transactions`, payload)
    if (imageManagerRef.value) {
      await imageManagerRef.value.commit(created.id)
    }
    detailStore.invalidate('transactions')
    router.push(`/spaces/${route.params.id}/ledger`)
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
  
  detailStore.setSpaceId(route.params.id as string)
  detailStore.fetchSpace()

  try {
    const space = await api.get<any>(`/api/spaces/${route.params.id}`)
    baseCurrency.value = space.base_currency || 'TWD'
    form.value.currency = baseCurrency.value

    // Load dynamic options from space
    if (space.categories) {
      categories.value = space.categories.map((c: string) => ({ label: c, value: c }))
    } else {
      categories.value = [
        { label: '餐飲', value: '餐飲' },
        { label: '交通', value: '交通' },
        { label: '購物', value: '購物' },
        { label: '娛樂', value: '娛樂' },
        { label: '生活', value: '生活' },
        { label: '其他', value: '其他' }
      ]
    }

    if (space.currencies) {
      availableCurrencies.value = space.currencies.map((c: string) => ({ label: c, value: c }))
    } else {
      availableCurrencies.value = [
        { label: 'TWD', value: 'TWD' },
        { label: 'JPY', value: 'JPY' },
        { label: 'USD', value: 'USD' }
      ]
    }
    
    // Set default category
    if (categories.value && categories.value.length > 0 && !form.value.category) {
      const firstCat = categories.value[0]
      if (firstCat) {
        form.value.category = firstCat.value
      }
    }
  } catch (e) {
    console.error('Failed to fetch space data', e)
  }
})
</script>
