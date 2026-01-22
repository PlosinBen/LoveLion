<template>
  <div v-if="label" class="flex flex-col gap-2">
    <label class="block text-sm text-neutral-400">{{ label }}</label>
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
      class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white text-base focus:outline-none focus:border-indigo-500 placeholder-neutral-400 disabled:opacity-50 disabled:cursor-not-allowed"
      :class="inputClass"
      v-bind="$attrs"
    />
  </div>
  <input
    v-else
    :value="modelValue"
    @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    :type="type"
    :placeholder="placeholder"
    :required="required"
    :disabled="disabled"
    :min="min"
    :max="max"
    :step="step"
    class="w-full px-4 py-3 rounded-xl border border-neutral-800 bg-neutral-800 text-white text-base focus:outline-none focus:border-indigo-500 placeholder-neutral-400 disabled:opacity-50 disabled:cursor-not-allowed"
    :class="inputClass"
    v-bind="$attrs"
  />
</template>

<script setup lang="ts">
defineOptions({
  inheritAttrs: false
})

defineProps<{
  modelValue: string | number | null | undefined
  label?: string
  type?: 'text' | 'password' | 'email' | 'number' | 'date' | 'tel'
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

