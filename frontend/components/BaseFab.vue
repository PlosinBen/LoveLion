<template>
  <button
    @click="$emit('click')"
    class="fixed right-6 bottom-28 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg rounded-full flex items-center justify-center text-white transition-transform active:scale-95 z-20 cursor-pointer border-0"
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
