<template>
  <div class="flex flex-col gap-1.5">
    <label v-if="label" class="text-xs text-neutral-400 px-1">{{ label }}</label>
    <input
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      :type="type"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      class="w-full bg-neutral-800 border border-neutral-700 text-white py-3 px-4 rounded-xl outline-none focus:border-indigo-500 transition-colors placeholder-neutral-500 text-base"
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
  inputClass?: string
}>()

defineEmits(['update:modelValue'])
</script>
