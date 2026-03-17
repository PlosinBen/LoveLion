<template>
  <div class="flex flex-col gap-3">
    <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1 flex justify-between items-center">
      <span>分帳</span>
      <button
        v-if="splits.length > 0"
        type="button"
        @click="splitEqually"
        class="text-xs font-bold text-indigo-400 bg-transparent border-0 cursor-pointer hover:text-indigo-300 transition-colors"
      >
        均分
      </button>
    </label>

    <div class="flex flex-col gap-2">
      <div
        v-for="(split, index) in splits"
        :key="index"
        class="bg-neutral-900 rounded-xl p-3 border border-neutral-800 flex items-center gap-3"
      >
        <!-- Payer toggle -->
        <button
          type="button"
          @click="togglePayer(index)"
          class="w-8 h-8 rounded-lg flex items-center justify-center border-0 cursor-pointer transition-colors shrink-0"
          :class="split.is_payer ? 'bg-indigo-500 text-white' : 'bg-neutral-800 text-neutral-600'"
          :title="split.is_payer ? '付款人' : '點擊設為付款人'"
        >
          <Icon icon="mdi:cash" class="text-base" />
        </button>

        <!-- Name -->
        <div class="flex-1 min-w-0">
          <BaseInput
            :model-value="split.name"
            @update:model-value="updateSplit(index, 'name', $event)"
            placeholder="姓名"
            input-class="text-sm font-bold"
          />
        </div>

        <!-- Amount -->
        <div class="w-24">
          <BaseInput
            :model-value="split.amount"
            @update:model-value="updateSplit(index, 'amount', Number($event))"
            type="number"
            placeholder="0"
            input-class="text-right text-sm"
          />
        </div>

        <!-- Remove -->
        <button
          type="button"
          @click="removeSplit(index)"
          class="w-8 h-8 rounded-lg flex items-center justify-center bg-transparent text-neutral-600 hover:text-red-400 border-0 cursor-pointer transition-colors shrink-0"
        >
          <Icon icon="mdi:close" class="text-base" />
        </button>
      </div>
    </div>

    <!-- Add person -->
    <button
      type="button"
      @click="addSplit"
      class="flex justify-center items-center w-full py-2.5 text-sm rounded font-bold bg-transparent text-neutral-500 hover:text-indigo-400 border border-dashed border-neutral-700 cursor-pointer transition-all active:scale-95 gap-2"
    >
      <Icon icon="mdi:account-plus-outline" class="text-lg" /> 新增對象
    </button>

    <!-- Remaining -->
    <div v-if="splits.length > 0" class="flex justify-between px-1 text-xs font-bold">
      <span class="text-neutral-500">剩餘未分配</span>
      <span :class="remaining === 0 ? 'text-green-400' : 'text-amber-400'">
        {{ remaining.toLocaleString() }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'
import BaseInput from '~/components/BaseInput.vue'

export interface SplitItem {
  name: string
  amount: number
  is_payer: boolean
}

const props = defineProps<{
  splits: SplitItem[]
  totalAmount: number
}>()

const emit = defineEmits<{
  'update:splits': [splits: SplitItem[]]
}>()

const remaining = computed(() => {
  const sum = props.splits.reduce((s, item) => s + Number(item.amount || 0), 0)
  return Math.round((props.totalAmount - sum) * 100) / 100
})

const addSplit = () => {
  emit('update:splits', [...props.splits, { name: '', amount: 0, is_payer: false }])
}

const removeSplit = (index: number) => {
  const updated = props.splits.filter((_, i) => i !== index)
  emit('update:splits', updated)
}

const updateSplit = (index: number, field: string, value: any) => {
  const updated = props.splits.map((s, i) =>
    i === index ? { ...s, [field]: value } : s
  )
  emit('update:splits', updated)
}

const togglePayer = (index: number) => {
  const updated = props.splits.map((s, i) =>
    i === index ? { ...s, is_payer: !s.is_payer } : s
  )
  emit('update:splits', updated)
}

const splitEqually = () => {
  const count = props.splits.length
  if (count === 0) return
  const each = Math.floor(props.totalAmount / count)
  const remainder = props.totalAmount - each * count

  const updated = props.splits.map((s, i) => ({
    ...s,
    amount: i === 0 ? each + remainder : each,
  }))
  emit('update:splits', updated)
}
</script>
