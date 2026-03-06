import { ref, computed } from 'vue'
import { useApi } from './useApi'

export const useSpace = () => {
  const api = useApi()
  
  // Use useState for SSR safety and global state sharing
  const allSpaces = useState<any[]>('all-spaces', () => [])
  const currentSpaceId = useState<string | null>('current-space-id', () => null)
  const loading = useState<boolean>('spaces-loading', () => false)

  const currentSpace = computed(() => {
    return allSpaces.value.find(s => s.id === currentSpaceId.value) || allSpaces.value[0] || null
  })

  // Grouped spaces for Dashboard
  const pinnedSpaces = computed(() => allSpaces.value.filter(s => s.is_pinned))
  const otherSpaces = computed(() => allSpaces.value.filter(s => !s.is_pinned))
  
  const tripSpaces = computed(() => allSpaces.value.filter(s => s.type === 'trip'))
  const personalSpaces = computed(() => allSpaces.value.filter(s => s.type === 'personal'))

  const fetchSpaces = async (force = false) => {
    if (allSpaces.value.length > 0 && !force) return
    
    loading.value = true
    try {
      // Fetch from unified spaces API
      const spaces = await api.get<any[]>('/api/spaces')
      allSpaces.value = spaces
      
      if (!currentSpaceId.value && spaces.length > 0) {
        currentSpaceId.value = spaces[0].id
      }
    } catch (e) {
      console.error('Failed to fetch spaces:', e)
    } finally {
      loading.value = false
    }
  }

  const selectSpace = (id: string) => {
    currentSpaceId.value = id
  }

  const togglePin = async (id: string) => {
    const space = allSpaces.value.find(s => s.id === id)
    if (!space) return

    try {
      const updated = await api.patch(`/api/spaces/${id}`, {
        is_pinned: !space.is_pinned
      })
      
      // Update local state
      const index = allSpaces.value.findIndex(s => s.id === id)
      if (index !== -1) {
        allSpaces.value[index] = updated
      }
    } catch (e) {
      console.error('Failed to toggle pin:', e)
    }
  }

  return {
    allSpaces,
    currentSpace,
    currentSpaceId,
    pinnedSpaces,
    otherSpaces,
    tripSpaces,
    personalSpaces,
    loading,
    fetchSpaces,
    selectSpace,
    togglePin
  }
}

// Export useLedger as an alias for backward compatibility during transition
export const useLedger = () => {
  const space = useSpace()
  return {
    allLedgers: space.allSpaces,
    currentLedger: space.currentSpace,
    currentLedgerId: space.currentSpaceId,
    loading: space.loading,
    fetchLedgers: space.fetchSpaces,
    selectLedger: space.selectSpace
  }
}
