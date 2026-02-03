<template>
  <div :class="{'flex flex-col gap-2': label}">
    <label v-if="label" class="block text-xs text-neutral-400">{{ label }}</label>
    <div class="relative">
      <select
        :value="modelValue"
        @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
        :required="required"
        :disabled="disabled"
        class="w-full cursor-pointer px-2 py-1.5 rounded border border-neutral-800 bg-neutral-800 text-white text-base focus:outline-none focus:border-neutral-400 placeholder-neutral-400 disabled:opacity-50 disabled:cursor-not-allowed appearance-none"
        :class="selectClass"
        v-bind="$attrs"
      >
        <option v-if="placeholder" value="" disabled selected>{{ placeholder }}</option>
        <option 
          v-for="option in options" 
          :key="typeof option === 'string' ? option : option.value" 
          :value="typeof option === 'string' ? option : option.value"
        >
          {{ typeof option === 'string' ? option : option.label }}
        </option>
      </select>
      <!-- Custom arrow icon -->
      <div class="absolute inset-y-0 right-0 flex items-center px-4 pointer-events-none text-neutral-400">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface SelectOption {
  label: string
  value: string | number
}

defineProps<{
  modelValue: string | number
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
