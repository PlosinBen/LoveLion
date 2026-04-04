import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useConfirmStore = defineStore('confirm', () => {
  const visible = ref(false)
  const message = ref('')
  const confirmLabel = ref('確定')
  const cancelLabel = ref('取消')
  const destructive = ref(false)

  let resolvePromise: ((value: boolean) => void) | null = null

  const show = (options: {
    message: string
    confirmLabel?: string
    cancelLabel?: string
    destructive?: boolean
  }): Promise<boolean> => {
    message.value = options.message
    confirmLabel.value = options.confirmLabel ?? '確定'
    cancelLabel.value = options.cancelLabel ?? '取消'
    destructive.value = options.destructive ?? false
    visible.value = true

    return new Promise<boolean>((resolve) => {
      resolvePromise = resolve
    })
  }

  const handleConfirm = () => {
    visible.value = false
    resolvePromise?.(true)
    resolvePromise = null
  }

  const handleCancel = () => {
    visible.value = false
    resolvePromise?.(false)
    resolvePromise = null
  }

  return { visible, message, confirmLabel, cancelLabel, destructive, show, handleConfirm, handleCancel }
})
