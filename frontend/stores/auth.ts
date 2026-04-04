import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useApi } from '~/composables/useApi'
import type { User } from '~/types'

interface AuthResponse {
  token: string
  user: User
}

export const useAuthStore = defineStore('auth', () => {
  const api = useApi()

  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const isAuthenticated = computed(() => !!token.value)

  const initAuth = () => {
    if (typeof window !== 'undefined') {
      const storedToken = localStorage.getItem('token')
      const storedUser = localStorage.getItem('user')
      if (storedToken) {
        token.value = storedToken
      }
      if (storedUser) {
        try {
          user.value = JSON.parse(storedUser)
        } catch {
          // Invalid stored user
        }
      }
    }
  }

  const login = async (username: string, password: string) => {
    const response = await api.post<AuthResponse>('/api/users/login', {
      username,
      password,
    })
    token.value = response.token
    user.value = response.user
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
    }
    return response
  }

  const register = async (username: string, password: string, displayName: string) => {
    const response = await api.post<AuthResponse>('/api/users/register', {
      username,
      password,
      display_name: displayName,
    })
    token.value = response.token
    user.value = response.user
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
    }
    return response
  }

  const logout = () => {
    token.value = null
    user.value = null
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  }

  const fetchUser = async () => {
    try {
      const response = await api.get<User>('/api/users/me')
      user.value = response
      if (typeof window !== 'undefined') {
        localStorage.setItem('user', JSON.stringify(response))
      }
      return response
    } catch (e) {
      logout()
      throw e
    }
  }

  const updateProfile = async (data: { display_name?: string; current_password?: string; new_password?: string }) => {
    const response = await api.put<User>('/api/users/me', data)
    user.value = response
    if (typeof window !== 'undefined') {
      localStorage.setItem('user', JSON.stringify(response))
    }
    return response
  }

  return {
    user,
    token,
    isAuthenticated,
    initAuth,
    login,
    register,
    logout,
    fetchUser,
    updateProfile,
  }
})
