import { ref } from 'vue'
import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock useState as a global (Nuxt auto-import)
vi.stubGlobal('useState', <T>(key: string, init: () => T) => ref(init()))

// Mock useApi
const mockGet = vi.fn()
const mockPatch = vi.fn()
const mockPost = vi.fn()
vi.mock('~/composables/useApi', () => ({
  useApi: () => ({
    get: mockGet,
    patch: mockPatch,
    post: mockPost,
    loading: { value: false },
    error: { value: null },
  }),
}))

import { useSpace } from '~/composables/useSpace'

const makeSpace = (overrides: Record<string, unknown> = {}) => ({
  id: 'sp-1',
  name: 'Test Space',
  type: 'personal',
  is_pinned: false,
  base_currency: 'TWD',
  ...overrides,
})

describe('useSpace', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Reset state
    const { allSpaces, currentSpaceId } = useSpace()
    allSpaces.value = []
    currentSpaceId.value = null
  })

  it('starts with empty state', () => {
    const { allSpaces, currentSpace, loading } = useSpace()
    expect(allSpaces.value).toEqual([])
    expect(currentSpace.value).toBeNull()
    expect(loading.value).toBe(false)
  })

  it('fetchSpaces populates allSpaces', async () => {
    const spaces = [makeSpace(), makeSpace({ id: 'sp-2', name: 'Trip' })]
    mockGet.mockResolvedValueOnce(spaces)

    const { fetchSpaces, allSpaces, currentSpaceId } = useSpace()
    await fetchSpaces(true)

    expect(mockGet).toHaveBeenCalledWith('/api/spaces')
    expect(allSpaces.value).toHaveLength(2)
    expect(currentSpaceId.value).toBe('sp-1')
  })

  it('fetchSpaces skips when already loaded and not forced', async () => {
    const spaces = [makeSpace()]
    mockGet.mockResolvedValueOnce(spaces)

    const { fetchSpaces, allSpaces } = useSpace()
    await fetchSpaces(true)
    expect(allSpaces.value).toHaveLength(1)

    await fetchSpaces() // should skip
    expect(mockGet).toHaveBeenCalledTimes(1)
  })

  it('selectSpace updates currentSpaceId', () => {
    const { selectSpace, currentSpaceId } = useSpace()
    selectSpace('sp-99')
    expect(currentSpaceId.value).toBe('sp-99')
  })

  it('computed: pinnedSpaces and otherSpaces', async () => {
    const spaces = [
      makeSpace({ id: 'sp-1', is_pinned: true }),
      makeSpace({ id: 'sp-2', is_pinned: false }),
      makeSpace({ id: 'sp-3', is_pinned: true }),
    ]
    mockGet.mockResolvedValueOnce(spaces)

    const { fetchSpaces, pinnedSpaces, otherSpaces } = useSpace()
    await fetchSpaces(true)

    expect(pinnedSpaces.value).toHaveLength(2)
    expect(otherSpaces.value).toHaveLength(1)
  })

  it('computed: organizedSpaces and personalSpaces', async () => {
    const spaces = [
      makeSpace({ id: 'sp-1', type: 'personal' }),
      makeSpace({ id: 'sp-2', type: 'trip' }),
      makeSpace({ id: 'sp-3', type: 'trip' }),
    ]
    mockGet.mockResolvedValueOnce(spaces)

    const { fetchSpaces, organizedSpaces, personalSpaces } = useSpace()
    await fetchSpaces(true)

    expect(organizedSpaces.value).toHaveLength(2)
    expect(personalSpaces.value).toHaveLength(1)
  })

  it('togglePin sends patch and updates local state', async () => {
    const space = makeSpace({ is_pinned: false })
    mockGet.mockResolvedValueOnce([space])

    const { fetchSpaces, togglePin, allSpaces } = useSpace()
    await fetchSpaces(true)

    const updated = { ...space, is_pinned: true }
    mockPatch.mockResolvedValueOnce(updated)
    await togglePin('sp-1')

    expect(mockPatch).toHaveBeenCalledWith('/api/spaces/sp-1', { is_pinned: true })
    expect(allSpaces.value[0].is_pinned).toBe(true)
  })

  it('leaveSpace removes space from local state', async () => {
    const spaces = [makeSpace({ id: 'sp-1' }), makeSpace({ id: 'sp-2' })]
    mockGet.mockResolvedValueOnce(spaces)

    const { fetchSpaces, leaveSpace, allSpaces, currentSpaceId } = useSpace()
    await fetchSpaces(true)
    currentSpaceId.value = 'sp-1'

    mockPost.mockResolvedValueOnce(undefined)
    await leaveSpace('sp-1')

    expect(allSpaces.value).toHaveLength(1)
    expect(allSpaces.value[0].id).toBe('sp-2')
    expect(currentSpaceId.value).toBe('sp-2')
  })

  it('currentSpace falls back to first space', async () => {
    const spaces = [makeSpace({ id: 'sp-1' }), makeSpace({ id: 'sp-2' })]
    mockGet.mockResolvedValueOnce(spaces)

    const { fetchSpaces, currentSpace, currentSpaceId } = useSpace()
    currentSpaceId.value = 'nonexistent'
    await fetchSpaces(true)

    // currentSpaceId was set to sp-1 by fetchSpaces, but let's test fallback
    currentSpaceId.value = 'nonexistent'
    expect(currentSpace.value).toEqual(spaces[0])
  })
})
