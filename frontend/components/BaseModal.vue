<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0 translate-y-4"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 translate-y-4"
    >
      <div v-if="modelValue" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4">
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/60" @click="$emit('update:modelValue', false)"></div>

        <!-- Content Area -->
        <div class="relative w-full max-w-lg bg-neutral-900 border-t sm:border border-neutral-800 rounded-t-2xl sm:rounded-2xl shadow-2xl overflow-hidden">
          <!-- Header -->
          <div v-if="title" class="p-4 border-b border-neutral-800 flex items-center justify-between bg-neutral-900 sticky top-0 z-10">
            <h3 class="font-bold text-white text-lg px-2">{{ title }}</h3>
            <button @click="$emit('update:modelValue', false)" class="flex justify-center items-center w-10 h-10 rounded-full bg-neutral-800 text-neutral-400 hover:text-white hover:bg-neutral-700 border-0 cursor-pointer transition-colors active:scale-95">
              <Icon icon="mdi:close" class="text-xl" />
            </button>
          </div>

          <!-- Body -->
          <div class="max-h-full overflow-y-auto">
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
