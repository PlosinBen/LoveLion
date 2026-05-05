import { ref } from 'vue'
import { useApi } from './useApi'
import type {
  InvMember,
  InvSettlement,
  InvMemberTransaction,
  InvSettlementAllocation,
  InvSettlementDetail,
  InvStockTrade,
} from '~/types'

export function useInvestment() {
  const { get, post, put, del } = useApi()
  const members = ref<InvMember[]>([])
  const settlements = ref<InvSettlement[]>([])
  const allocations = ref<InvSettlementAllocation[]>([])
  const memberTransactions = ref<InvMemberTransaction[]>([])
  const stockTrades = ref<InvStockTrade[]>([])

  const fetchMembers = async () => {
    members.value = await get<InvMember[]>('/api/investments/members')
  }

  const createMember = async (data: { name: string; user_id?: string }) => {
    return await post<InvMember>('/api/investments/members', data)
  }

  const updateMember = async (id: string, data: Partial<InvMember>) => {
    return await put<InvMember>(`/api/investments/members/${id}`, data)
  }

  const fetchSettlements = async () => {
    settlements.value = await get<InvSettlement[]>('/api/investments/settlements')
  }

  const createSettlement = async (yearMonth: string) => {
    return await post<InvSettlement>('/api/investments/settlements', { year_month: yearMonth })
  }

  const getSettlement = async (ym: string) => {
    return await get<InvSettlementDetail>(`/api/investments/settlements/${ym}`)
  }

  const completeSettlement = async (ym: string) => {
    return await put<any>(`/api/investments/settlements/${ym}/complete`, {})
  }

  const reopenSettlement = async (ym: string) => {
    return await put<any>(`/api/investments/settlements/${ym}/reopen`, {})
  }

  const deleteSettlement = async (ym: string) => {
    return await del<any>(`/api/investments/settlements/${ym}`)
  }

  const upsertFutures = async (ym: string, data: any) => {
    return await put<any>(`/api/investments/settlements/${ym}/futures`, data)
  }

  const upsertStocks = async (ym: string, data: any) => {
    return await put<any>(`/api/investments/settlements/${ym}/stocks`, data)
  }

  const fetchAllocations = async (params?: { from?: string; to?: string }) => {
    let url = '/api/investments/allocations'
    if (params) {
      const qs = new URLSearchParams()
      if (params.from) qs.set('from', params.from)
      if (params.to) qs.set('to', params.to)
      const s = qs.toString()
      if (s) url += '?' + s
    }
    allocations.value = await get<InvSettlementAllocation[]>(url)
  }

  const fetchMemberTransactions = async (params?: { from?: string; to?: string; member_id?: string }) => {
    let url = '/api/investments/members/transactions'
    if (params) {
      const qs = new URLSearchParams()
      if (params.from) qs.set('from', params.from)
      if (params.to) qs.set('to', params.to)
      if (params.member_id) qs.set('member_id', params.member_id)
      const s = qs.toString()
      if (s) url += '?' + s
    }
    memberTransactions.value = await get<InvMemberTransaction[]>(url)
  }

  const createMemberTransaction = async (data: { member_id: string; date: string; type: string; amount: number; note?: string }) => {
    return await post<InvMemberTransaction>('/api/investments/members/transactions', data)
  }

  const updateMemberTransaction = async (id: string, data: any) => {
    return await put<InvMemberTransaction>(`/api/investments/members/transactions/${id}`, data)
  }

  const deleteMemberTransaction = async (id: string) => {
    return await del<any>(`/api/investments/members/transactions/${id}`)
  }

  const fetchStockTrades = async (params?: { from?: string; to?: string }) => {
    let url = '/api/investments/stocks/trades'
    if (params) {
      const qs = new URLSearchParams()
      if (params.from) qs.set('from', params.from)
      if (params.to) qs.set('to', params.to)
      const s = qs.toString()
      if (s) url += '?' + s
    }
    stockTrades.value = await get<InvStockTrade[]>(url)
  }

  const createStockTrade = async (data: any) => {
    return await post<InvStockTrade>('/api/investments/stocks/trades', data)
  }

  const updateStockTrade = async (id: string, data: any) => {
    return await put<InvStockTrade>(`/api/investments/stocks/trades/${id}`, data)
  }

  const deleteStockTrade = async (id: string) => {
    return await del<any>(`/api/investments/stocks/trades/${id}`)
  }

  return {
    members,
    settlements,
    allocations,
    memberTransactions,
    stockTrades,
    fetchMembers,
    createMember,
    updateMember,
    fetchSettlements,
    createSettlement,
    getSettlement,
    completeSettlement,
    reopenSettlement,
    deleteSettlement,
    upsertFutures,
    upsertStocks,
    fetchAllocations,
    fetchMemberTransactions,
    createMemberTransaction,
    updateMemberTransaction,
    deleteMemberTransaction,
    fetchStockTrades,
    createStockTrade,
    updateStockTrade,
    deleteStockTrade,
  }
}
