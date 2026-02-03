<template>
  <div class="flex flex-col transition-all duration-300">
    <!-- Header -->
    <div 
      class="flex items-center justify-between p-4 cursor-pointer hover:bg-neutral-800/50 transition-colors"
      @click="toggle"
    >
      <slot name="header" :open="isOpen" :toggle="toggle">
        <div class="flex items-center gap-2">
           <Icon :icon="isOpen ? 'mdi:chevron-down' : 'mdi:chevron-right'" class="text-xl text-neutral-400" />
           <h3 class="font-bold text-white">詳細資訊</h3>
        </div>
      </slot>
    </div>

    <!-- Body -->
    <div v-show="isOpen" class="animate-fade-in">
        <slot :open="isOpen" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: undefined
  },
  defaultOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const internalOpen = ref(props.defaultOpen)

// If controlled (v-model), use prop; else use internal state
const isOpen = computed(() => {
  return props.modelValue !== undefined ? props.modelValue : internalOpen.value
})

const toggle = () => {
  if (props.modelValue !== undefined) {
    emit('update:modelValue', !props.modelValue)
  } else {
    internalOpen.value = !internalOpen.value
  }
}

watch(() => props.defaultOpen, (newVal) => {
    if (props.modelValue === undefined) {
        internalOpen.value = newVal
    }
})
</script>
