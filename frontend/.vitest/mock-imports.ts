// Mock Nuxt auto-imports for vitest
export { ref, computed, reactive, watch, watchEffect, toRef, toRefs } from 'vue'

// Mock useState to behave like ref for testing
import { ref } from 'vue'
export function useState<T>(key: string, init: () => T) {
  return ref(init()) as ReturnType<typeof ref<T>>
}

// Mock useRuntimeConfig
export function useRuntimeConfig() {
  return {
    public: {
      apiBase: '',
    },
  }
}
