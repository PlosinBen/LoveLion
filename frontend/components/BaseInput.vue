<template>
  <div :class="{'flex flex-col gap-2': label}">
    <label v-if="label" class="block text-xs text-neutral-400">{{ label }}</label>
    <input
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      :type="type"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      :min="min"
      :max="max"
      :step="step"
      class="w-full cursor-pointer px-2 py-1.5 rounded border border-neutral-800 bg-neutral-800 text-white text-base focus:outline-none focus:border-neutral-400 placeholder-neutral-400 disabled:opacity-50 disabled:cursor-not-allowed"
      :class="inputClass"
      v-bind="$attrs"
    />
  </div>
</template>

<script setup lang="ts">


defineProps<{
  modelValue: string | number | null | undefined
  label?: string
  type?: 'text' | 'password' | 'email' | 'number' | 'date' | 'tel' | 'datetime-local' | 'time'
  placeholder?: string
  required?: boolean
  disabled?: boolean
  min?: string | number
  max?: string | number
  step?: string | number
  inputClass?: string
}>()

defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()
</script>

