import { describe, it, expect, vi, beforeEach } from 'vitest'

const mockGet = vi.fn()
const mockPost = vi.fn()
const mockPut = vi.fn()
const mockDel = vi.fn()

vi.mock('~/composables/useApi', () => ({
  useApi: () => ({
    get: mockGet,
    post: mockPost,
    put: mockPut,
    del: mockDel,
    loading: { value: false },
    error: { value: null },
  }),
}))

import { useInvestment } from '~/composables/useInvestment'

describe('useInvestment', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('members', () => {
    it('fetchMembers populates members ref', async () => {
      const fakeMembers = [
        { id: 'm1', name: 'Owner', is_owner: true, active: true },
        { id: 'm2', name: 'Alice', is_owner: false, active: true },
      ]
      mockGet.mockResolvedValueOnce(fakeMembers)

      const { members, fetchMembers } = useInvestment()
      await fetchMembers()

      expect(mockGet).toHaveBeenCalledWith('/api/investments/members')
      expect(members.value).toEqual(fakeMembers)
    })

    it('createMember posts correct data', async () => {
      const newMember = { id: 'm3', name: 'Bob' }
      mockPost.mockResolvedValueOnce(newMember)

      const { createMember } = useInvestment()
      const result = await createMember({ name: 'Bob' })

      expect(mockPost).toHaveBeenCalledWith('/api/investments/members', { name: 'Bob' })
      expect(result).toEqual(newMember)
    })

    it('updateMember sends PUT with correct path', async () => {
      mockPut.mockResolvedValueOnce({ id: 'm1', active: false })

      const { updateMember } = useInvestment()
      await updateMember('m1', { active: false })

      expect(mockPut).toHaveBeenCalledWith('/api/investments/members/m1', { active: false })
    })
  })

  describe('settlements', () => {
    it('fetchSettlements populates settlements ref', async () => {
      const fakeList = [
        { year_month: '2026-05', status: 'draft', total_profit_loss: 0 },
      ]
      mockGet.mockResolvedValueOnce(fakeList)

      const { settlements, fetchSettlements } = useInvestment()
      await fetchSettlements()

      expect(mockGet).toHaveBeenCalledWith('/api/investments/settlements')
      expect(settlements.value).toEqual(fakeList)
    })

    it('createSettlement posts year_month', async () => {
      mockPost.mockResolvedValueOnce({ year_month: '2026-05', status: 'draft' })

      const { createSettlement } = useInvestment()
      await createSettlement('2026-05')

      expect(mockPost).toHaveBeenCalledWith('/api/investments/settlements', { year_month: '2026-05' })
    })

    it('getSettlement fetches by year_month', async () => {
      const fakeDetail = { year_month: '2026-05', status: 'draft', futures_statement: null, stock_statement: null, allocations: [] }
      mockGet.mockResolvedValueOnce(fakeDetail)

      const { getSettlement } = useInvestment()
      const result = await getSettlement('2026-05')

      expect(mockGet).toHaveBeenCalledWith('/api/investments/settlements/2026-05')
      expect(result).toEqual(fakeDetail)
    })

    it('completeSettlement sends PUT', async () => {
      mockPut.mockResolvedValueOnce({ status: 'completed' })

      const { completeSettlement } = useInvestment()
      await completeSettlement('2026-05')

      expect(mockPut).toHaveBeenCalledWith('/api/investments/settlements/2026-05/complete', {})
    })

    it('reopenSettlement sends PUT', async () => {
      mockPut.mockResolvedValueOnce({ status: 'draft' })

      const { reopenSettlement } = useInvestment()
      await reopenSettlement('2026-05')

      expect(mockPut).toHaveBeenCalledWith('/api/investments/settlements/2026-05/reopen', {})
    })

    it('deleteSettlement sends DELETE', async () => {
      mockDel.mockResolvedValueOnce({})

      const { deleteSettlement } = useInvestment()
      await deleteSettlement('2026-05')

      expect(mockDel).toHaveBeenCalledWith('/api/investments/settlements/2026-05')
    })
  })

  describe('futures & stocks', () => {
    it('upsertFutures sends PUT with data', async () => {
      const data = { ending_equity: 100000, floating_profit_loss: 5000, realized_profit_loss: 15000, deposit: 0, withdrawal: 0 }
      mockPut.mockResolvedValueOnce({})

      const { upsertFutures } = useInvestment()
      await upsertFutures('2026-05', data)

      expect(mockPut).toHaveBeenCalledWith('/api/investments/settlements/2026-05/futures', data)
    })

    it('upsertStocks sends PUT with holdings', async () => {
      const data = { account_balance: 80000, deposit: 0, withdrawal: 0, holdings: [{ symbol: '2330', shares: 1000, closing_price: 600.5 }] }
      mockPut.mockResolvedValueOnce({})

      const { upsertStocks } = useInvestment()
      await upsertStocks('2026-05', data)

      expect(mockPut).toHaveBeenCalledWith('/api/investments/settlements/2026-05/stocks', data)
    })
  })

  describe('allocations', () => {
    it('fetchAllocations without params', async () => {
      const fakeData = [{ year_month: '2026-05', member_id: 'm1', weight: 10, amount: 8431, balance: 58431 }]
      mockGet.mockResolvedValueOnce(fakeData)

      const { allocations, fetchAllocations } = useInvestment()
      await fetchAllocations()

      expect(mockGet).toHaveBeenCalledWith('/api/investments/allocations')
      expect(allocations.value).toEqual(fakeData)
    })

    it('fetchAllocations with date range', async () => {
      mockGet.mockResolvedValueOnce([])

      const { fetchAllocations } = useInvestment()
      await fetchAllocations({ from: '2026-01', to: '2026-06' })

      expect(mockGet).toHaveBeenCalledWith('/api/investments/allocations?from=2026-01&to=2026-06')
    })
  })

  describe('member transactions', () => {
    it('fetchMemberTransactions populates ref', async () => {
      const fakeTxns = [{ id: 'uuid-1', member_id: 'm2', type: 'deposit', amount: 10000, date: '2026-05-01' }]
      mockGet.mockResolvedValueOnce(fakeTxns)

      const { memberTransactions, fetchMemberTransactions } = useInvestment()
      await fetchMemberTransactions()

      expect(mockGet).toHaveBeenCalledWith('/api/investments/members/transactions')
      expect(memberTransactions.value).toEqual(fakeTxns)
    })

    it('createMemberTransaction posts data', async () => {
      const data = { member_id: 'm2', date: '2026-05-01', type: 'deposit', amount: 10000 }
      mockPost.mockResolvedValueOnce({ id: 'uuid-new', ...data })

      const { createMemberTransaction } = useInvestment()
      await createMemberTransaction(data)

      expect(mockPost).toHaveBeenCalledWith('/api/investments/members/transactions', data)
    })

    it('deleteMemberTransaction sends DELETE', async () => {
      mockDel.mockResolvedValueOnce({})

      const { deleteMemberTransaction } = useInvestment()
      await deleteMemberTransaction('uuid-1')

      expect(mockDel).toHaveBeenCalledWith('/api/investments/members/transactions/uuid-1')
    })
  })

  describe('stock trades', () => {
    it('fetchStockTrades with date range', async () => {
      mockGet.mockResolvedValueOnce([])

      const { fetchStockTrades } = useInvestment()
      await fetchStockTrades({ from: '2026-01-01', to: '2026-05-31' })

      expect(mockGet).toHaveBeenCalledWith('/api/investments/stocks/trades?from=2026-01-01&to=2026-05-31')
    })

    it('createStockTrade posts data', async () => {
      const data = { trade_date: '2026-05-02', symbol: '2330', shares: 1000, price: 600, fee: 585, tax: 0 }
      mockPost.mockResolvedValueOnce({ id: 'uuid-t1', ...data })

      const { createStockTrade } = useInvestment()
      await createStockTrade(data)

      expect(mockPost).toHaveBeenCalledWith('/api/investments/stocks/trades', data)
    })

    it('updateStockTrade sends PUT', async () => {
      const data = { shares: -1000 }
      mockPut.mockResolvedValueOnce({ id: 'uuid-t1', shares: -1000 })

      const { updateStockTrade } = useInvestment()
      await updateStockTrade('uuid-t1', data)

      expect(mockPut).toHaveBeenCalledWith('/api/investments/stocks/trades/uuid-t1', data)
    })

    it('deleteStockTrade sends DELETE', async () => {
      mockDel.mockResolvedValueOnce({})

      const { deleteStockTrade } = useInvestment()
      await deleteStockTrade('uuid-t1')

      expect(mockDel).toHaveBeenCalledWith('/api/investments/stocks/trades/uuid-t1')
    })
  })
})
