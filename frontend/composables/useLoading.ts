import { storeToRefs } from 'pinia'
import { useLoadingStore } from '~/stores/loading'

export function useLoading() {
  const store = useLoadingStore()
  const { isLoading } = storeToRefs(store)
  return { isLoading, showLoading: store.showLoading, hideLoading: store.hideLoading }
}
