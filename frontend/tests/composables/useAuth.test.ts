import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

// Mock useApi before importing useAuth
const mockPost = vi.fn()
const mockGet = vi.fn()
const mockPut = vi.fn()
vi.mock('~/composables/useApi', () => ({
  useApi: () => ({
    post: mockPost,
    get: mockGet,
    put: mockPut,
    loading: { value: false },
    error: { value: null },
  }),
}))

// Mock localStorage and window for Node environment
const store: Record<string, string> = {}
const localStorageMock = {
  getItem: vi.fn((key: string) => store[key] ?? null),
  setItem: vi.fn((key: string, value: string) => { store[key] = value }),
  removeItem: vi.fn((key: string) => { delete store[key] }),
}
Object.defineProperty(globalThis, 'localStorage', { value: localStorageMock })

// useAuth checks `typeof window !== 'undefined'` for SSR safety
if (typeof globalThis.window === 'undefined') {
  (globalThis as any).window = globalThis
}

import { useAuth } from '~/composables/useAuth'

describe('useAuth', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    Object.keys(store).forEach(k => delete store[k])
  })

  it('starts unauthenticated', () => {
    const { isAuthenticated, user, token } = useAuth()
    expect(isAuthenticated.value).toBe(false)
    expect(user.value).toBeNull()
    expect(token.value).toBeNull()
  })

  it('login stores token and user', async () => {
    const fakeResponse = {
      token: 'jwt-token-123',
      user: { id: '1', username: 'dev', display_name: 'Dev' },
    }
    mockPost.mockResolvedValueOnce(fakeResponse)

    const { login, isAuthenticated, user, token } = useAuth()
    await login('dev', 'dev123')

    expect(mockPost).toHaveBeenCalledWith('/api/users/login', {
      username: 'dev',
      password: 'dev123',
    })
    expect(token.value).toBe('jwt-token-123')
    expect(user.value).toEqual(fakeResponse.user)
    expect(isAuthenticated.value).toBe(true)
    expect(localStorageMock.setItem).toHaveBeenCalledWith('token', 'jwt-token-123')
  })

  it('register stores token and user', async () => {
    const fakeResponse = {
      token: 'jwt-token-456',
      user: { id: '2', username: 'newuser', display_name: 'New' },
    }
    mockPost.mockResolvedValueOnce(fakeResponse)

    const { register, token, user } = useAuth()
    await register('newuser', 'pass123', 'New')

    expect(mockPost).toHaveBeenCalledWith('/api/users/register', {
      username: 'newuser',
      password: 'pass123',
      display_name: 'New',
    })
    expect(token.value).toBe('jwt-token-456')
    expect(user.value).toEqual(fakeResponse.user)
  })

  it('logout clears state and localStorage', async () => {
    // First login
    mockPost.mockResolvedValueOnce({
      token: 'tok',
      user: { id: '1', username: 'dev', display_name: 'Dev' },
    })
    const { login, logout, isAuthenticated, token, user } = useAuth()
    await login('dev', 'dev123')
    expect(isAuthenticated.value).toBe(true)

    logout()
    expect(token.value).toBeNull()
    expect(user.value).toBeNull()
    expect(isAuthenticated.value).toBe(false)
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('token')
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('user')
  })

  it('fetchUser updates user and falls back to logout on error', async () => {
    // Setup authenticated state
    mockPost.mockResolvedValueOnce({
      token: 'tok',
      user: { id: '1', username: 'dev', display_name: 'Dev' },
    })
    const { login, fetchUser, user } = useAuth()
    await login('dev', 'dev123')

    // fetchUser success
    const updatedUser = { id: '1', username: 'dev', display_name: 'Updated' }
    mockGet.mockResolvedValueOnce(updatedUser)
    await fetchUser()
    expect(user.value).toEqual(updatedUser)

    // fetchUser failure triggers logout
    mockGet.mockRejectedValueOnce(new Error('401'))
    await expect(fetchUser()).rejects.toThrow('401')
    expect(user.value).toBeNull()
  })

  it('updateProfile updates stored user', async () => {
    mockPost.mockResolvedValueOnce({
      token: 'tok',
      user: { id: '1', username: 'dev', display_name: 'Dev' },
    })
    const { login, updateProfile, user } = useAuth()
    await login('dev', 'dev123')

    const updated = { id: '1', username: 'dev', display_name: 'New Name' }
    mockPut.mockResolvedValueOnce(updated)
    await updateProfile({ display_name: 'New Name' })
    expect(user.value).toEqual(updated)
    expect(mockPut).toHaveBeenCalledWith('/api/users/me', { display_name: 'New Name' })
  })

  it('initAuth restores state from localStorage', () => {
    // Set localStorage BEFORE calling initAuth
    localStorageMock.getItem.mockImplementation((key: string) => {
      if (key === 'token') return 'stored-token'
      if (key === 'user') return JSON.stringify({ id: '1', username: 'dev', display_name: 'Dev' })
      return null
    })

    const { initAuth, token, user, isAuthenticated } = useAuth()
    initAuth()

    expect(token.value).toBe('stored-token')
    expect(user.value).toEqual({ id: '1', username: 'dev', display_name: 'Dev' })
    expect(isAuthenticated.value).toBe(true)

    // Restore default mock
    localStorageMock.getItem.mockImplementation((key: string) => store[key] ?? null)
  })

  it('initAuth handles invalid stored user gracefully', () => {
    localStorageMock.getItem.mockImplementation((key: string) => {
      if (key === 'token') return 'stored-token'
      if (key === 'user') return 'not-json'
      return null
    })

    const { initAuth, token } = useAuth()
    initAuth()

    expect(token.value).toBe('stored-token')

    // Restore default mock
    localStorageMock.getItem.mockImplementation((key: string) => store[key] ?? null)
  })
})
