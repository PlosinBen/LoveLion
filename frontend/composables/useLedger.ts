import { ref, computed } from 'vue'
import { useApi } from './useApi'

export const useLedger = () => {
  const api = useApi()
  
  // Global state stored in a singleton pattern or shared via Provide/Inject if needed, 
  // but for simple Nuxt composables, we can use refs defined outside the export or 
  // rely on useState if we want SSR safety.
  const allLedgers = useState<any[]>('all-ledgers', () => [])
  const currentLedgerId = useState<string | null>('current-ledger-id', () => null)
  const loading = useState<boolean>('ledgers-loading', () => false)

  const currentLedger = computed(() => {
    return allLedgers.value.find(l => l.id === currentLedgerId.value) || allLedgers.value[0] || null
  })

  const fetchLedgers = async (force = false) => {
    if (allLedgers.value.length > 0 && !force) return
    
    loading.value = true
    try {
      const ledgers = await api.get<any[]>('/api/ledgers')
      allLedgers.value = ledgers
      
      // If no ledger is selected, default to the first one or a preference
      if (!currentLedgerId.value && ledgers.length > 0) {
        currentLedgerId.value = ledgers[0].id
      }
    } catch (e) {
      console.error('Failed to fetch ledgers:', e)
    } finally {
      loading.value = false
    }
  }

  const selectLedger = (id: string) => {
    currentLedgerId.value = id
  }

  return {
    allLedgers,
    currentLedger,
    currentLedgerId,
    loading,
    fetchLedgers,
    selectLedger
  }
}
