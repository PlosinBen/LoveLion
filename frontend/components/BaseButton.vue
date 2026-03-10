<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    @click.stop="$emit('click', $event)"
    class="flex justify-center items-center transition-all active:scale-95 disabled:opacity-50 border-0 cursor-pointer font-bold shrink-0 px-6 py-4 text-sm rounded-xl"
    :class="[
      variantClasses,
      fullWidth ? 'w-full' : ''
    ]"
  >
    <Icon v-if="loading" icon="mdi:loading" class="animate-spin text-xl" />
    <slot v-else />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

interface Props {
  type?: 'button' | 'submit' | 'reset'
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost' | 'white'
  fullWidth?: boolean
  disabled?: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'button',
  variant: 'primary',
  fullWidth: false,
  disabled: false,
  loading: false
})

defineEmits(['click'])

const variantClasses = computed(() => {
  switch (props.variant) {
    case 'primary':
      return 'bg-indigo-500 text-white hover:bg-indigo-600 shadow-lg'
    case 'secondary':
      return 'bg-neutral-800 text-neutral-400 hover:text-white hover:bg-neutral-700'
    case 'danger':
      return 'bg-red-500/10 text-red-500 hover:bg-red-500 hover:text-white border border-red-500/20'
    case 'ghost':
      return 'bg-transparent text-neutral-500 hover:text-neutral-300'
    case 'white':
      return 'bg-white text-black hover:bg-neutral-200'
    default:
      return 'bg-neutral-900 text-indigo-400 hover:bg-neutral-800'
  }
})
</script>
