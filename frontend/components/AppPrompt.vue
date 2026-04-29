<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60" @click="handleCancel"></div>
        <div class="relative w-full max-w-sm bg-neutral-900 border border-neutral-800 rounded-2xl shadow-2xl overflow-hidden">
          <div class="p-6 flex flex-col gap-3">
            <p class="text-white text-sm font-bold">{{ title }}</p>
            <input
              ref="inputRef"
              v-model="inputValue"
              :placeholder="placeholder"
              class="w-full bg-neutral-800 border border-neutral-700 text-white py-2 px-3 rounded-xl text-sm focus:outline-none focus:border-indigo-500"
              @keydown.enter="handleConfirm"
            >
          </div>
          <div class="flex border-t border-neutral-800">
            <button
              class="flex-1 py-3 text-sm text-neutral-400 hover:bg-neutral-800 transition-colors cursor-pointer border-0 bg-transparent"
              @click="handleCancel"
            >
              {{ cancelLabel }}
            </button>
            <button
              class="flex-1 py-3 text-sm font-medium text-blue-400 hover:bg-blue-950 transition-colors cursor-pointer border-0 border-l border-neutral-800 bg-transparent"
              @click="handleConfirm"
            >
              {{ confirmLabel }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { watch, ref, nextTick } from 'vue'
import { usePromptStore } from '~/stores/prompt'

const store = usePromptStore()
const { visible, title, placeholder, inputValue, confirmLabel, cancelLabel } = storeToRefs(store)
const { handleConfirm, handleCancel } = store

const inputRef = ref<HTMLInputElement | null>(null)

watch(visible, (val) => {
  if (val) {
    nextTick(() => inputRef.value?.focus())
  }
})
</script>
