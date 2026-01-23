<template>
  <div v-if="label" class="flex flex-col gap-2">
    <label class="block text-sm text-neutral-400">{{ label }}</label>
    <textarea
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      :rows="rows"
      class="w-full px-4 py-3 rounded-md border border-neutral-800 bg-neutral-800 text-white text-base focus:outline-none focus:border-indigo-500 placeholder-neutral-400 disabled:opacity-50 disabled:cursor-not-allowed resize-none"
      :class="textareaClass"
      v-bind="$attrs"
    />
  </div>
  <textarea
    v-else
    :value="modelValue"
    @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
    :placeholder="placeholder"
    :required="required"
    :disabled="disabled"
    :rows="rows"
    class="w-full px-4 py-3 rounded-md border border-neutral-800 bg-neutral-800 text-white text-base focus:outline-none focus:border-indigo-500 placeholder-neutral-400 disabled:opacity-50 disabled:cursor-not-allowed resize-none"
    :class="textareaClass"
    v-bind="$attrs"
  />
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
