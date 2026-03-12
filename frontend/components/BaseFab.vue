<template>
  <button
    @click="$emit('click')"
    class="fixed right-6 bottom-24 w-12 h-12 bg-indigo-500 shadow-lg rounded-full flex items-center justify-center text-white cursor-pointer border-0"
  >
    <Icon :icon="icon" class="text-3xl" />
  </button>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

interface Props {
  icon?: string
}

withDefaults(defineProps<Props>(), {
  icon: 'mdi:plus',
})

defineEmits(['click'])
</script>
