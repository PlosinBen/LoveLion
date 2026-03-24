<template>
  <div class="flex flex-col gap-2">
    <div class="bg-neutral-900 rounded-xl p-4 border border-neutral-800 flex flex-col gap-3">
      <!-- Header with split equally button -->
      <div class="flex justify-between items-center">
        <span class="text-xs font-bold text-neutral-500 uppercase">選擇付款人</span>
        <button
          v-if="debtRows.length > 0"
          type="button"
          @click="splitEqually"
          class="text-xs font-bold text-indigo-400 bg-transparent border-0 cursor-pointer hover:text-indigo-300 transition-colors"
        >
          均分
        </button>
      </div>

      <!-- Debt rows (all members) -->
      <div class="flex flex-col gap-1">
        <div
          v-for="row in debtRows"
          :key="row.payer_name"
          class="flex flex-col"
        >
          <div class="flex items-center gap-3 py-1.5">
            <label class="flex items-center gap-2 flex-1 min-w-0 cursor-pointer">
              <input
                type="radio"
                name="debt-payee"
                :value="row.payer_name"
                :checked="row.payer_name === sharedPayee"
                @change="updateSharedPayee(row.payer_name)"
                class="w-4 h-4 text-indigo-500 bg-neutral-800 border-neutral-700"
              >
              <span class="text-sm text-neutral-300 truncate">{{ row.payer_name }}</span>
            </label>
            <button
              type="button"
              @click="updateSpotPaid(row.payer_name, !row.is_spot_paid)"
              class="text-xs font-bold px-2 py-1 rounded-lg border cursor-pointer transition-colors shrink-0"
              :class="row.is_spot_paid
                ? 'bg-emerald-500/15 text-emerald-400 border-emerald-500/30'
                : 'bg-transparent text-neutral-600 border-neutral-700 hover:text-neutral-400'"
            >
              現付
            </button>
            <div class="w-28">
              <BaseInput
                :model-value="row.amount"
                @update:model-value="updateAmount(row.payer_name, Number($event))"
                type="number"
                placeholder="0"
                input-class="text-right"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Remaining -->
      <div v-if="debtRows.length > 0" class="flex justify-between pt-2 border-t border-neutral-800 text-xs font-bold">
        <span class="text-neutral-500">剩餘未分配</span>
        <span :class="remaining === 0 ? 'text-green-400' : 'text-amber-400'">
          {{ remaining.toLocaleString() }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import BaseInput from '~/components/BaseInput.vue'

export interface DebtItem {
  payer_name: string
  payee_name: string
  amount: number
  is_spot_paid: boolean
}

const props = defineProps<{
  debts: DebtItem[]
  totalAmount: number
  memberOptions: string[]
}>()

const emit = defineEmits<{
  'update:debts': [debts: DebtItem[]]
}>()

const sharedPayee = computed(() => {
  if (props.debts.length === 0) return ''
  return props.debts[0]?.payee_name || ''
})

// Build rows from memberOptions, merging with existing debts
const debtRows = computed(() => {
  const payee = sharedPayee.value
  const debtMap = new Map(props.debts.map(d => [d.payer_name, d]))
  return props.memberOptions
    .map(name => {
      const existing = debtMap.get(name)
      return {
        payer_name: name,
        payee_name: payee,
        amount: existing?.amount ?? 0,
        is_spot_paid: existing?.is_spot_paid ?? false,
      }
    })
})

const remaining = computed(() => {
  const sum = debtRows.value.reduce((s, item) => s + Number(item.amount || 0), 0)
  return Math.round((props.totalAmount - sum) * 100) / 100
})

const emitDebts = (rows: DebtItem[]) => {
  emit('update:debts', rows)
}

const updateSharedPayee = (value: string) => {
  // Rebuild debts: all members except new payee
  const debtMap = new Map(props.debts.map(d => [d.payer_name, d]))
  const updated = props.memberOptions
    .map(name => {
      const existing = debtMap.get(name)
      return {
        payer_name: name,
        payee_name: value,
        amount: existing?.amount ?? 0,
        is_spot_paid: existing?.is_spot_paid ?? false,
      }
    })
  emitDebts(updated)
}

const updateAmount = (payerName: string, amount: number) => {
  const updated = props.debts.map(d =>
    d.payer_name === payerName ? { ...d, amount } : d
  )
  emitDebts(updated)
}

const updateSpotPaid = (payerName: string, checked: boolean) => {
  const updated = props.debts.map(d =>
    d.payer_name === payerName ? { ...d, is_spot_paid: checked } : d
  )
  emitDebts(updated)
}

const splitEqually = () => {
  const rows = debtRows.value
  const count = rows.length
  if (count === 0) return
  const each = Math.floor(props.totalAmount / count)
  const remainder = props.totalAmount - each * count

  const updated = rows.map((d, i) => ({
    ...d,
    amount: i === 0 ? each + remainder : each,
  }))
  emitDebts(updated)
}

// Initialize debts when memberOptions become available and debts are empty
watch(() => props.memberOptions, (members) => {
  if (members.length > 0 && props.debts.length === 0) {
    const updated = members.map(name => ({
      payer_name: name,
      payee_name: '',
      amount: 0,
      is_spot_paid: false,
    }))
    emitDebts(updated)
  }
}, { immediate: true })
</script>
