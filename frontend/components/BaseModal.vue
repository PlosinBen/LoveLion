<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-500 ease-[cubic-bezier(0.16,1,0.3,1)]"
      enter-from-class="opacity-0 translate-y-full"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-300 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-full"
    >
      <div v-if="modelValue" class="fixed inset-0 z-[200] flex items-end sm:items-center justify-center p-0 sm:p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/70 backdrop-blur-md" @click="$emit('update:modelValue', false)"></div>

        <!-- Content Area -->
        <div class="relative w-full max-w-lg bg-neutral-900 border-t sm:border border-neutral-800 rounded-t-[2.5rem] sm:rounded-3xl shadow-[0_-20px_40px_rgba(0,0,0,0.5)] overflow-hidden">
          <!-- Header (Optional) -->
          <div v-if="title" class="p-6 border-b border-neutral-800/50 flex items-center justify-between bg-neutral-900/80 backdrop-blur-xl sticky top-0 z-10">
            <h3 class="font-black text-white text-xl tracking-tight px-2">{{ title }}</h3>
            <button @click="$emit('update:modelValue', false)" class="w-12 h-12 rounded-2xl bg-neutral-800/50 text-neutral-400 flex items-center justify-center hover:text-white transition-all border-0 cursor-pointer active:scale-90">
              <Icon icon="mdi:close" class="text-2xl" />
            </button>
          </div>

          <!-- Body -->
          <div class="max-h-[85vh] overflow-y-auto">
            <slot />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'

defineProps({
  modelValue: {
    type: Boolean,
    required: true
  },
  title: {
    type: String,
    default: ''
  }
})

defineEmits(['update:modelValue'])
</script>
