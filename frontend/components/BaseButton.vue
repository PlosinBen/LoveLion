<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    @click.stop="$emit('click', $event)"
    class="flex justify-center items-center transition-all active:scale-95 disabled:opacity-50 border-0 cursor-pointer font-bold shrink-0"
    :class="[
      variantClasses,
      sizeClasses,
      fullWidth ? 'w-full' : ''
    ]"
  >
    <Icon v-if="loading" icon="mdi:loading" class="animate-spin" :class="size === 'sm' ? 'text-sm' : 'text-xl'" />
    <slot v-else />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'

interface Props {
  type?: 'button' | 'submit' | 'reset'
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost' | 'white'
  size?: 'sm' | 'md' | 'lg' | 'icon'
  fullWidth?: boolean
  disabled?: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'button',
  variant: 'primary',
  size: 'md',
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
      return 'bg-neutral-900 text-indigo-400 hover:bg-neutral-800' // Custom default for our "Icon button" style
  }
})

const sizeClasses = computed(() => {
  if (props.size === 'icon') return 'w-10 h-10 rounded-xl'
  
  switch (props.size) {
    case 'sm':
      return 'px-4 py-2 text-xs rounded-lg'
    case 'lg':
      return 'px-8 py-5 text-lg rounded-2xl'
    default:
      return 'px-6 py-4 text-sm rounded-xl'
  }
})
</script>
