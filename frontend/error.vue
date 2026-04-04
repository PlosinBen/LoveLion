<template>
  <div class="flex flex-col min-h-screen bg-neutral-900 text-neutral-50 items-center justify-center px-6">
    <div class="flex flex-col items-center gap-6 text-center max-w-sm">
      <div class="text-7xl font-bold text-neutral-700">{{ error?.statusCode || 500 }}</div>
      <h1 class="text-xl font-bold text-white">
        {{ title }}
      </h1>
      <p class="text-sm text-neutral-400">
        {{ description }}
      </p>
      <button
        class="inline-flex justify-center items-center px-6 py-2.5 text-sm rounded font-bold bg-indigo-500 text-white hover:bg-indigo-600 shadow-lg transition-all active:scale-95 border-0 cursor-pointer"
        @click="handleError"
      >
        回到首頁
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { NuxtError } from '#app'

const props = defineProps<{
  error: NuxtError
}>()

const title = computed(() => {
  if (props.error?.statusCode === 404) return '找不到頁面'
  return '發生錯誤'
})

const description = computed(() => {
  if (props.error?.statusCode === 404) return '你要找的頁面不存在，可能已被移除或網址有誤。'
  return '伺服器發生了一些問題，請稍後再試。'
})

const handleError = () => clearError({ redirect: '/' })
</script>
