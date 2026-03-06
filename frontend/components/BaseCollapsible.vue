<template>
  <div class="flex flex-col bg-neutral-900 border border-neutral-800/60 rounded-3xl overflow-hidden transition-all duration-300 shadow-sm">
    <!-- Header -->
    <button 
      type="button"
      class="flex items-center justify-between p-5 w-full bg-transparent border-0 cursor-pointer hover:bg-neutral-800/50 transition-colors active:bg-neutral-800 group"
      @click="toggle"
    >
      <slot name="header" :open="isOpen" :toggle="toggle">
        <div class="flex items-center gap-3">
           <div class="w-10 h-10 rounded-xl bg-neutral-800 flex items-center justify-center text-neutral-500 group-hover:text-indigo-400 transition-colors">
             <Icon 
               icon="mdi:chevron-right" 
               class="text-2xl transition-transform duration-300" 
               :class="{ 'rotate-90 text-indigo-500': isOpen }" 
             />
           </div>
           <span class="font-black text-white tracking-tight text-sm uppercase tracking-widest">詳細資訊</span>
        </div>
      </slot>
    </button>

    <!-- Body -->
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="max-h-0 opacity-0"
      enter-to-class="max-h-[1000px] opacity-100"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="max-h-[1000px] opacity-100"
      leave-to-class="max-h-0 opacity-0"
    >
      <div v-if="isOpen" class="overflow-hidden">
        <div class="px-5 pb-6 border-t border-neutral-800/30 pt-4">
          <slot :open="isOpen" />
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
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
