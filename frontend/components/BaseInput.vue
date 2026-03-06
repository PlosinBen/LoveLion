<template>
  <div class="flex flex-col gap-2">
    <label v-if="label" class="text-[10px] font-black text-neutral-500 uppercase tracking-widest px-1">{{ label }}</label>
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
      class="w-full bg-neutral-900 border border-neutral-800 text-white py-4 px-5 rounded-2xl outline-none focus:border-indigo-500/50 focus:ring-4 focus:ring-indigo-500/5 transition-all placeholder-neutral-600 disabled:opacity-50 disabled:cursor-not-allowed font-medium text-sm"
      :class="inputClass"
      v-bind="$attrs"
    />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string | number | null | undefined
  label?: string
  type?: string
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
