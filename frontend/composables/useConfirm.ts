import { useConfirmStore } from '~/stores/confirm'

export function useConfirm() {
  const store = useConfirmStore()
  return store.show
}
