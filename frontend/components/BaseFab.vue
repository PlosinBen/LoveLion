<template>
  <button
    @click="$emit('click')"
    class="fixed right-6 bottom-24 w-14 h-14 bg-indigo-500 rounded-full flex items-center justify-center text-white z-20 cursor-pointer"
  >
    <Icon :icon="icon" class="text-3xl" />
  </button>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

interface Props {
  icon?: string
}

const props = withDefaults(defineProps<Props>(), {
  icon: 'mdi:plus',
})

defineEmits(['click'])
</script>
