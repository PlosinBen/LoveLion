<template>
  <div class="flex flex-col gap-2">
    <label v-if="label" class="text-[10px] font-black text-neutral-500 uppercase tracking-widest px-1">{{ label }}</label>
    <div class="relative">
      <select
        :value="modelValue"
        @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
        :required="required"
        :disabled="disabled"
        class="w-full bg-neutral-900 border border-neutral-800 text-white py-4 px-5 pr-12 rounded-2xl outline-none focus:border-indigo-500/50 focus:ring-4 focus:ring-indigo-500/5 transition-all appearance-none disabled:opacity-50 disabled:cursor-not-allowed font-black text-sm tracking-tight cursor-pointer"
        :class="selectClass"
        v-bind="$attrs"
      >
        <option v-if="placeholder" value="" disabled selected>{{ placeholder }}</option>
        <option
          v-for="option in options"
          :key="typeof option === 'string' ? option : option.value"
          :value="typeof option === 'string' ? option : option.value"
          class="bg-neutral-900 text-white py-2"
        >
          {{ typeof option === 'string' ? option : option.label }}
        </option>
      </select>
      <!-- Custom Chevron Icon -->
      <div class="absolute inset-y-0 right-0 flex items-center px-4 pointer-events-none text-neutral-600">
        <Icon icon="mdi:chevron-down" class="text-xl" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

export interface SelectOption {
  label: string
  value: string | number
}

defineProps<{
  modelValue: string | number | null | undefined
  label?: string
  placeholder?: string
  required?: boolean
  disabled?: boolean
  options: (string | SelectOption)[]
  selectClass?: string
}>()

defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()
</script>
