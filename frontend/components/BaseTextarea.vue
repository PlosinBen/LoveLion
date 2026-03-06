<template>
  <div class="flex flex-col gap-1.5">
    <label v-if="label" class="text-xs font-bold text-neutral-500 uppercase tracking-wider px-1">{{ label }}</label>
    <textarea
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      :rows="rows"
      class="w-full bg-neutral-800 border border-neutral-700 text-white py-3 px-4 rounded-xl outline-none focus:border-indigo-500 transition-colors placeholder-neutral-500 resize-none font-medium text-base"
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

defineEmits(['update:modelValue'])
</script>
