<template>
  <form @submit.prevent="$emit('submit')" class="flex flex-col gap-6">
    <!-- Basic Info -->
    <div class="flex flex-col gap-2">
      <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">基本資訊</label>
      <div class="bg-neutral-900 rounded-xl p-4 border border-neutral-800 flex flex-col gap-4">
        <!-- Title — kept enabled even with AI extract on, so the user can pre-label the receipt. -->
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold text-neutral-500 uppercase px-1">標題</label>
          <BaseInput v-model="form.title" placeholder="例如: 午餐、計程車" />
        </div>

        <fieldset :disabled="disabled" class="flex flex-col gap-4 border-0 p-0 m-0 min-w-0 disabled:opacity-60">
          <!-- Date & Time -->
          <div class="flex flex-col gap-2">
            <label class="text-xs font-bold text-neutral-500 uppercase px-1">日期與時間</label>
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

          <!-- Currency, Category & Payment Method -->
          <div class="grid grid-cols-3 gap-3">
            <div class="flex flex-col gap-2">
              <label class="text-xs font-bold text-neutral-500 uppercase px-1">幣別</label>
              <BaseSelect
                v-model="form.currency"
                :options="availableCurrencies"
              />
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-xs font-bold text-neutral-500 uppercase px-1">類別</label>
              <BaseSelect
                v-model="form.category"
                :options="categories"
              />
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-xs font-bold text-neutral-500 uppercase px-1">付款方式</label>
              <BaseSelect
                v-model="form.payment_method"
                :options="paymentMethods"
              />
            </div>
          </div>

          <!-- Total Amount -->
          <div class="flex flex-col gap-2">
            <label class="text-xs font-bold text-neutral-500 uppercase px-1">總額 ({{ form.currency }})</label>
            <BaseInput
              v-model.number="form.total_amount"
              type="number"
              placeholder="0"
              input-class="text-2xl font-bold text-center"
            />
          </div>
        </fieldset>
      </div>
    </div>

    <fieldset :disabled="disabled" class="flex flex-col gap-6 border-0 p-0 m-0 min-w-0 disabled:opacity-60">
      <!-- Items Section (collapsible) -->
      <div class="flex flex-col gap-2">
        <button
          type="button"
          @click="showItems = !showItems"
          class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1 flex justify-between items-center bg-transparent border-0 cursor-pointer w-full"
        >
          <span>項目明細（選填）</span>
          <Icon :icon="showItems ? 'mdi:chevron-up' : 'mdi:chevron-down'" class="text-lg" />
        </button>

        <div v-if="showItems" class="flex flex-col gap-3">
          <div v-for="(item, index) in form.items" :key="index" class="bg-neutral-900 rounded-xl p-4 border border-neutral-800 flex flex-col gap-4 relative">
            <BaseInput
              v-model="item.name"
              placeholder="項目名稱"
              input-class="font-bold"
            />

            <div class="grid grid-cols-2 gap-3">
              <div class="flex flex-col gap-2">
                <label class="text-xs font-bold text-neutral-500 uppercase px-1">單價</label>
                <BaseInput v-model.number="item.unit_price" type="number" placeholder="單價" input-class="text-right" />
              </div>
              <div class="flex flex-col gap-2">
                <label class="text-xs font-bold text-neutral-500 uppercase px-1">折扣</label>
                <BaseInput v-model.number="item.discount" type="number" placeholder="折扣" input-class="text-right" />
              </div>
            </div>

            <div class="flex items-end gap-3">
              <div class="w-24 flex flex-col gap-2">
                <label class="text-xs font-bold text-neutral-500 uppercase px-1">數量</label>
                <BaseInput v-model.number="item.quantity" type="number" placeholder="數量" min="1" input-class="text-right" />
              </div>

              <div class="flex-1 text-right pb-3">
                <div class="text-xs font-bold text-neutral-500 uppercase px-1 mb-1">項目小計</div>
                <div class="text-xl font-bold text-indigo-400">
                  {{ ((Number(item.unit_price) - Number(item.discount || 0)) * Number(item.quantity)).toLocaleString() }}
                </div>
              </div>

              <button
                v-if="form.items.length > 1"
                type="button"
                @click="removeItem(index)"
                class="flex justify-center items-center w-10 h-10 rounded bg-red-500/10 text-red-500 hover:bg-red-500 hover:text-white border border-red-500/20 cursor-pointer transition-colors active:scale-95 mb-1"
              >
                <Icon icon="mdi:close" class="text-xl" />
              </button>
            </div>
          </div>

          <button type="button" @click="addItem" class="flex justify-center items-center w-full py-2.5 text-sm rounded font-bold bg-transparent text-neutral-500 hover:text-indigo-400 border border-dashed border-neutral-700 cursor-pointer transition-all active:scale-95 gap-2">
            <Icon icon="mdi:plus" class="text-xl" /> 新增項目
          </button>

          <!-- Items total vs manual total comparison -->
          <div v-if="itemsTotal > 0" class="flex justify-between items-center px-1 text-xs font-bold">
            <span class="text-neutral-500">明細合計</span>
            <span :class="itemsTotal === form.total_amount ? 'text-green-400' : 'text-amber-400'">
              {{ itemsTotal.toLocaleString() }}
            </span>
          </div>
        </div>
      </div>

      <!-- Foreign Currency Settlement -->
      <div v-if="form.currency !== baseCurrency" class="flex flex-col gap-2">
        <button
          type="button"
          @click="showForeignCurrency = !showForeignCurrency"
          class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1 flex justify-between items-center bg-transparent border-0 cursor-pointer w-full"
        >
          <span>外幣結算</span>
          <Icon :icon="showForeignCurrency ? 'mdi:chevron-up' : 'mdi:chevron-down'" class="text-lg" />
        </button>
        <div v-if="showForeignCurrency" class="bg-neutral-900 rounded-xl p-4 border border-neutral-800 flex flex-col gap-4">
          <label class="flex items-center gap-2 cursor-pointer select-none">
            <input type="checkbox" v-model="form.manual_rate" class="w-4 h-4 rounded text-indigo-500 bg-neutral-800 border-neutral-700">
            <span class="text-xs font-bold text-neutral-400">手動輸入匯率</span>
          </label>

          <div v-if="form.manual_rate" class="flex flex-col gap-4">
            <div class="flex flex-col gap-2">
              <label class="text-xs font-bold text-neutral-500 uppercase px-1">匯率 (1 {{ baseCurrency }} = ? {{ form.currency }})</label>
              <BaseInput v-model.number="form.exchange_rate" type="number" step="0.0001" input-class="text-right" />
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-xs font-bold text-neutral-500 uppercase px-1">折合 {{ baseCurrency }}</label>
              <div class="w-full bg-neutral-800 border border-neutral-700 border-dashed text-white py-2 px-3 rounded-xl text-md font-bold text-right">
                {{ baseCurrency }} {{ calculatedBillingAmount.toLocaleString() }}
              </div>
            </div>
          </div>

          <div v-else class="flex flex-col gap-4">
            <div class="grid grid-cols-2 gap-3">
              <div class="flex flex-col gap-2">
                <label class="text-xs font-bold text-neutral-500 uppercase px-1">入帳金額 ({{ baseCurrency }})</label>
                <BaseInput v-model.number="form.billing_amount" type="number" input-class="text-right" />
              </div>
              <div class="flex flex-col gap-2">
                <label class="text-xs font-bold text-neutral-500 uppercase px-1">手續費 ({{ baseCurrency }})</label>
                <BaseInput v-model.number="form.handling_fee" type="number" input-class="text-right" />
              </div>
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-xs font-bold text-neutral-500 uppercase px-1">換算匯率</label>
              <div class="w-full bg-neutral-800 border border-neutral-700 border-dashed text-indigo-400 py-2 px-3 rounded-xl text-md font-bold text-right">
                {{ calculatedExchangeRate }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Debts -->
      <div class="flex flex-col gap-2">
        <button
          type="button"
          @click="showDebts = !showDebts"
          class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1 flex justify-between items-center bg-transparent border-0 cursor-pointer w-full"
        >
          <span>分帳</span>
          <Icon :icon="showDebts ? 'mdi:chevron-up' : 'mdi:chevron-down'" class="text-lg" />
        </button>
        <DebtEditor
          v-if="showDebts"
          :debts="debts"
          :total-amount="form.total_amount"
          :member-options="memberOptions"
          @update:debts="$emit('update:debts', $event)"
        />
      </div>
    </fieldset>

    <!-- Images — always enabled, the whole point of AI extract is to upload photos. -->
    <div class="flex flex-col gap-2">
      <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">收據 / 照片</label>
      <slot name="images" />
    </div>

    <fieldset :disabled="disabled" class="flex flex-col gap-6 border-0 p-0 m-0 min-w-0 disabled:opacity-60">
      <!-- Location URL -->
      <div class="flex flex-col gap-2">
        <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">地點連結</label>
        <BaseInput
          v-model="form.location_url"
          placeholder="貼上 Google Maps 連結"
          type="url"
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
    </fieldset>

    <slot name="actions">
      <BaseButton type="submit" variant="primary" class="w-full mt-4">
        儲存交易
      </BaseButton>
    </slot>
  </form>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Icon } from '@iconify/vue'
import { VueDatePicker } from '@vuepic/vue-datepicker'
import '@vuepic/vue-datepicker/dist/main.css'
import BaseInput from '~/components/BaseInput.vue'
import BaseSelect from '~/components/BaseSelect.vue'
import BaseTextarea from '~/components/BaseTextarea.vue'
import DebtEditor from '~/components/DebtEditor.vue'
import type { DebtItem } from '~/components/DebtEditor.vue'

export interface ExpenseFormData {
  date: Date
  title: string
  total_amount: number
  category: string
  payment_method: string
  note: string
  items: { name: string; unit_price: number; quantity: number; discount: number }[]
  currency: string
  manual_rate: boolean
  exchange_rate: number
  billing_amount: number
  handling_fee: number
  location_url: string
}

const props = withDefaults(
  defineProps<{
    baseCurrency: string
    categories: { label: string; value: string }[]
    availableCurrencies: { label: string; value: string }[]
    paymentMethods: { label: string; value: string }[]
    memberOptions: string[]
    debts: DebtItem[]
    /**
     * When true, all controls except Title and the Images slot become disabled.
     * Used while the AI extract toggle is on — title stays editable as a manual
     * override, and the user still needs to be able to add/remove receipts.
     */
    disabled?: boolean
  }>(),
  { disabled: false },
)

defineEmits<{
  submit: []
  'update:debts': [debts: DebtItem[]]
}>()

const form = defineModel<ExpenseFormData>('form', { required: true })

const hasExistingItems = form.value.items.some(i => i.name && Number(i.unit_price) > 0)
const showItems = ref(hasExistingItems)
const showForeignCurrency = ref(form.value.currency !== props.baseCurrency)
const hasExistingDebts = props.debts.some(d => d.amount > 0)
const showDebts = ref(hasExistingDebts)

const itemsTotal = computed(() => {
  return form.value.items.reduce((sum, item) => {
    const subtotal = (Number(item.unit_price) - Number(item.discount || 0)) * Number(item.quantity)
    return sum + subtotal
  }, 0)
})

const calculatedBillingAmount = computed(() => {
  if (form.value.manual_rate && form.value.exchange_rate > 0) {
    return Math.round(form.value.total_amount * form.value.exchange_rate)
  }
  return 0
})

const calculatedExchangeRate = computed(() => {
  if (!form.value.manual_rate && form.value.total_amount > 0 && form.value.billing_amount > 0) {
    return (form.value.billing_amount / form.value.total_amount).toFixed(4)
  }
  return 0
})

watch(() => form.value.currency, (currency) => {
  if (currency !== props.baseCurrency) showForeignCurrency.value = true
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

defineExpose({ itemsTotal })
</script>
