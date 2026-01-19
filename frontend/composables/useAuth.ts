import { ref, computed } from 'vue'
import { useApi } from './useApi'

interface User {
    id: string
    username: string
    display_name: string
}

interface AuthResponse {
    token: string
    user: User
}

const user = ref<User | null>(null)
const token = ref<string | null>(null)

export function useAuth() {
    const api = useApi()
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
                } catch (e) {
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
            return response
        } catch (e) {
            logout()
            throw e
        }
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
    }
}
