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
          <div class="p-6">
            <p class="text-white text-sm leading-relaxed whitespace-pre-line">{{ message }}</p>
          </div>
          <div class="flex border-t border-neutral-800">
            <button
              class="flex-1 py-3 text-sm text-neutral-400 hover:bg-neutral-800 transition-colors cursor-pointer border-0 bg-transparent"
              @click="handleCancel"
            >
              {{ cancelLabel }}
            </button>
            <button
              :class="[
                'flex-1 py-3 text-sm font-medium transition-colors cursor-pointer border-0 border-l border-neutral-800 bg-transparent',
                destructive ? 'text-red-400 hover:bg-red-950' : 'text-blue-400 hover:bg-blue-950'
              ]"
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
import { useConfirmStore } from '~/stores/confirm'

const store = useConfirmStore()
const { visible, message, confirmLabel, cancelLabel, destructive } = storeToRefs(store)
const { handleConfirm, handleCancel } = store
</script>
