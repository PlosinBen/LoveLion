<template>
  <div :class="{'flex flex-col gap-2': label}">
    <label v-if="label" class="block text-xs text-neutral-400">{{ label }}</label>
    <textarea
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      :rows="rows"
      class="w-full px-4 py-3 rounded border border-neutral-800 bg-neutral-800 text-white text-base focus:outline-none focus:border-neutral-400 placeholder-neutral-400 disabled:opacity-50 disabled:cursor-not-allowed resize-none"
      :class="textareaClass"
      v-bind="$attrs"
    />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: string
  label?: string
  placeholder?: string
  required?: boolean
  disabled?: boolean
  rows?: number | string
  textareaClass?: string
}>()

defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()
</script>
