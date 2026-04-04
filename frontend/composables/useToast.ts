import { storeToRefs } from 'pinia'
import { useToastStore } from '~/stores/toast'

export function useToast() {
  const store = useToastStore()
  const { toasts } = storeToRefs(store)
  return { toasts, show: store.show, success: store.success, error: store.error, dismiss: store.dismiss }
}
