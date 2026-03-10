import { ref } from 'vue'

const isLoading = ref(false)

export function useLoading() {
  const showLoading = () => {
    isLoading.value = true
  }

  const hideLoading = () => {
    isLoading.value = false
  }

  return { isLoading, showLoading, hideLoading }
}
