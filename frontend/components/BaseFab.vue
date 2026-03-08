<template>
  <button
    @click="$emit('click')"
    class="fixed right-6 w-14 h-14 bg-indigo-500 hover:bg-indigo-600 shadow-lg rounded-full flex items-center justify-center text-white transition-transform active:scale-95 z-20 cursor-pointer border-0"
    :class="positionClass"
  >
    <Icon :icon="icon" class="text-3xl" />
  </button>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

interface Props {
  icon?: string
  bottomClass?: string
}

const props = withDefaults(defineProps<Props>(), {
  icon: 'mdi:plus',
  bottomClass: 'bottom-28'
})

defineEmits(['click'])

const positionClass = computed(() => props.bottomClass)
</script>
