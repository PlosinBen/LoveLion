import { usePromptStore } from '~/stores/prompt'

export function usePrompt() {
  const store = usePromptStore()
  return store.show
}
