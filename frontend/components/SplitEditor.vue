<template>
  <div class="flex flex-col gap-3">
    <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1 flex justify-between items-center">
      <span>分帳</span>
      <button
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
          @click="split.is_payer = !split.is_payer"
          class="w-8 h-8 rounded-lg flex items-center justify-center border-0 cursor-pointer transition-colors shrink-0"
          :class="split.is_payer ? 'bg-indigo-500 text-white' : 'bg-neutral-800 text-neutral-600'"
          :title="split.is_payer ? '付款人' : '點擊設為付款人'"
        >
          <Icon icon="mdi:cash" class="text-base" />
        </button>

        <!-- Name -->
        <span class="text-sm font-bold text-white flex-1 truncate">{{ split.name }}</span>

        <!-- Amount -->
        <div class="w-28">
          <BaseInput
            :model-value="split.amount"
            @update:model-value="split.amount = Number($event)"
            type="number"
            placeholder="0"
            input-class="text-right text-sm"
          />
        </div>
      </div>
    </div>

    <!-- Remaining -->
    <div class="flex justify-between px-1 text-xs font-bold">
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
  member_id: string | null
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
