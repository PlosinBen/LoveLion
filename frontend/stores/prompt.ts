import { ref } from 'vue'
import { defineStore } from 'pinia'

export const usePromptStore = defineStore('prompt', () => {
  const visible = ref(false)
  const title = ref('')
  const placeholder = ref('')
  const inputValue = ref('')
  const confirmLabel = ref('確定')
  const cancelLabel = ref('取消')

  let resolvePromise: ((value: string | null) => void) | null = null

  const show = (options: {
    title: string
    placeholder?: string
    defaultValue?: string
    confirmLabel?: string
    cancelLabel?: string
  }): Promise<string | null> => {
    title.value = options.title
    placeholder.value = options.placeholder ?? ''
    inputValue.value = options.defaultValue ?? ''
    confirmLabel.value = options.confirmLabel ?? '確定'
    cancelLabel.value = options.cancelLabel ?? '取消'
    visible.value = true

    return new Promise<string | null>((resolve) => {
      resolvePromise = resolve
    })
  }

  const handleConfirm = () => {
    visible.value = false
    resolvePromise?.(inputValue.value.trim() || null)
    resolvePromise = null
  }

  const handleCancel = () => {
    visible.value = false
    resolvePromise?.(null)
    resolvePromise = null
  }

  return { visible, title, placeholder, inputValue, confirmLabel, cancelLabel, show, handleConfirm, handleCancel }
})
