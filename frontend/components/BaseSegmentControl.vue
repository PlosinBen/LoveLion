<template>
  <div class="flex p-1 bg-neutral-800 rounded-xl mb-6 shadow-inner">
    <template v-for="option in options" :key="option.to || option.value">
      <!-- Link Mode -->
      <NuxtLink
        v-if="option.to"
        :to="option.to"
        class="flex-1 py-2 text-sm font-bold rounded-lg transition-all text-center no-underline"
        :class="route.path === option.to ? 'bg-indigo-500 text-white shadow-sm' : 'text-neutral-500 hover:text-neutral-300'"
      >
        {{ option.label }}
      </NuxtLink>

      <!-- Action Mode (BaseButton) -->
      <BaseButton
        v-else
        @click="$emit('update:modelValue', option.value)"
        :variant="modelValue === option.value ? 'primary' : 'ghost'"
        size="sm"
        class="flex-1 !shadow-none"
      >
        {{ option.label }}
      </BaseButton>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import BaseButton from '~/components/BaseButton.vue'

interface Option {
  label: string
  to?: string
  value?: string
}

defineProps<{
  options: Option[]
  modelValue?: string
}>()

defineEmits(['update:modelValue'])

const route = useRoute()
</script>
