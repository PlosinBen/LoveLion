import { defineStore } from 'pinia'
import { useApi } from '~/composables/useApi'

export const useSpaceDetailStore = defineStore('spaceDetail', () => {
  const api = useApi()

  // State
  const spaceId = ref<string | null>(null)
  const space = ref<any>(null)
  const transactions = ref<any[]>([])
  const stores = ref<any[]>([])
  const members = ref<any[]>([])
  const invites = ref<any[]>([])

  const loading = ref({
    space: false,
    transactions: false,
    stores: false,
    members: false,
    invites: false
  })

  // Track if data has been fetched at least once for this spaceId
  const fetched = ref({
    space: false,
    transactions: false,
    stores: false,
    members: false,
    invites: false
  })

  // Reset when switching to a different space
  const setSpaceId = (id: string) => {
    if (spaceId.value === id) return
    spaceId.value = id
    space.value = null
    transactions.value = []
    stores.value = []
    members.value = []
    invites.value = []
    fetched.value = { space: false, transactions: false, stores: false, members: false, invites: false }
  }

  // Fetch helpers — skip if already loaded unless force=true
  const fetchSpace = async (force = false) => {
    if (!spaceId.value) return
    if (fetched.value.space && !force) return
    loading.value.space = true
    try {
      space.value = await api.get<any>(`/api/spaces/${spaceId.value}`)
      members.value = space.value.members || []
      fetched.value.space = true
    } catch (e) {
      console.error('Failed to fetch space:', e)
      throw e
    } finally {
      loading.value.space = false
    }
  }

  const fetchTransactions = async (force = false) => {
    if (!spaceId.value) return
    if (fetched.value.transactions && !force) return
    loading.value.transactions = true
    try {
      transactions.value = await api.get<any[]>(`/api/spaces/${spaceId.value}/transactions`) || []
      fetched.value.transactions = true
    } catch (e) {
      console.error('Failed to fetch transactions:', e)
      throw e
    } finally {
      loading.value.transactions = false
    }
  }

  const fetchStores = async (force = false) => {
    if (!spaceId.value) return
    if (fetched.value.stores && !force) return
    loading.value.stores = true
    try {
      stores.value = await api.get<any[]>(`/api/spaces/${spaceId.value}/stores`) || []
      fetched.value.stores = true
    } catch (e) {
      console.error('Failed to fetch stores:', e)
      throw e
    } finally {
      loading.value.stores = false
    }
  }

  const fetchMembers = async (force = false) => {
    if (!spaceId.value) return
    if (fetched.value.members && !force) return
    loading.value.members = true
    try {
      members.value = await api.get<any[]>(`/api/spaces/${spaceId.value}/members`) || []
      fetched.value.members = true
    } catch (e) {
      console.error('Failed to fetch members:', e)
      throw e
    } finally {
      loading.value.members = false
    }
  }

  const fetchInvites = async (force = false) => {
    if (!spaceId.value) return
    if (fetched.value.invites && !force) return
    loading.value.invites = true
    try {
      invites.value = await api.get<any[]>(`/api/spaces/${spaceId.value}/invites`) || []
      fetched.value.invites = true
    } catch (e) {
      console.error('Failed to fetch invites:', e)
      throw e
    } finally {
      loading.value.invites = false
    }
  }

  // Mark specific data as stale so next fetch will reload
  const invalidate = (...keys: Array<'space' | 'transactions' | 'stores' | 'members' | 'invites'>) => {
    for (const key of keys) {
      fetched.value[key] = false
    }
  }

  return {
    spaceId,
    space,
    transactions,
    stores,
    members,
    invites,
    loading,
    fetched,
    setSpaceId,
    fetchSpace,
    fetchTransactions,
    fetchStores,
    fetchMembers,
    fetchInvites,
    invalidate
  }
})
