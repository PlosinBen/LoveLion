<template>
  <form @submit.prevent="$emit('submit')" class="flex flex-col gap-6">
    <!-- Basic Info -->
    <div class="flex flex-col gap-2">
      <label class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">基本資訊</label>
      <div class="bg-neutral-900 rounded-xl p-4 border border-neutral-800 flex flex-col gap-4">
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

        <!-- Title -->
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold text-neutral-500 uppercase px-1">標題</label>
          <BaseInput v-model="form.title" placeholder="例如: 小明補款給小美" />
        </div>

        <!-- Payer & Payee -->
        <div class="grid grid-cols-2 gap-4">
          <div class="flex flex-col gap-2">
            <label class="text-xs font-bold text-neutral-500 uppercase px-1">付款人</label>
            <BaseSelect
              v-model="form.payer_name"
              :options="memberOptions"
              placeholder="誰付"
            />
          </div>
          <div class="flex flex-col gap-2">
            <label class="text-xs font-bold text-neutral-500 uppercase px-1">收款人</label>
            <BaseSelect
              v-model="form.payee_name"
              :options="memberOptions"
              placeholder="收誰"
            />
          </div>
        </div>

        <!-- Amount -->
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold text-neutral-500 uppercase px-1">金額 ({{ baseCurrency }})</label>
          <BaseInput
            v-model.number="form.total_amount"
            type="number"
            placeholder="0"
            input-class="text-2xl font-bold text-center"
          />
        </div>

        <!-- Note -->
        <div class="flex flex-col gap-2">
          <label class="text-xs font-bold text-neutral-500 uppercase px-1">備註</label>
          <BaseTextarea
            v-model="form.note"
            placeholder="選填"
            rows="3"
          />
        </div>
      </div>
    </div>

    <slot name="actions">
      <BaseButton type="submit" variant="primary" class="w-full mt-4">
        儲存補款
      </BaseButton>
    </slot>
  </form>
</template>

<script setup lang="ts">
import { VueDatePicker } from '@vuepic/vue-datepicker'
import '@vuepic/vue-datepicker/dist/main.css'
import BaseInput from '~/components/BaseInput.vue'
import BaseSelect from '~/components/BaseSelect.vue'
import BaseTextarea from '~/components/BaseTextarea.vue'

export interface PaymentFormData {
  date: Date
  title: string
  payer_name: string
  payee_name: string
  total_amount: number
  note: string
}

defineProps<{
  baseCurrency: string
  memberOptions: string[]
}>()

defineEmits<{
  submit: []
}>()

const form = defineModel<PaymentFormData>('form', { required: true })
</script>
