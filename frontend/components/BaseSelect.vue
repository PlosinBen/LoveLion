<template>
  <div class="flex flex-col gap-1.5">
    <label v-if="label" class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">{{ label }}</label>
    <div class="relative">
      <select
        :value="modelValue"
        @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
        :required="required"
        :disabled="disabled"
        class="w-full bg-neutral-800 border border-neutral-700 text-white py-3 px-4 pr-10 rounded-xl outline-none focus:border-indigo-500 transition-colors appearance-none disabled:opacity-50 disabled:cursor-not-allowed font-medium text-base cursor-pointer"
        :class="selectClass"
        v-bind="$attrs"
      >
        <option v-if="placeholder" value="" disabled selected>{{ placeholder }}</option>
        <option
          v-for="option in options"
          :key="typeof option === 'string' ? option : option.value"
          :value="typeof option === 'string' ? option : option.value"
          class="bg-neutral-800 text-white py-2"
        >
          {{ typeof option === 'string' ? option : option.label }}
        </option>
      </select>
      <div class="absolute inset-y-0 right-0 flex items-center px-3 pointer-events-none text-neutral-500">
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

defineEmits(['update:modelValue'])
</script>
