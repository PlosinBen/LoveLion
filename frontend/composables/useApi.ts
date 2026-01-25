import { ref } from 'vue'

export function useApi() {
    const loading = ref(false)
    const error = ref<string | null>(null)

    const getToken = () => {
        if (typeof window !== 'undefined') {
            return localStorage.getItem('token')
        }
        return null
    }

    const clearAuthAndRedirect = () => {
        if (typeof window !== 'undefined') {
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            // Redirect to login page
            window.location.href = '/login'
        }
    }

    const request = async <T>(
        endpoint: string,
        options: RequestInit = {}
    ): Promise<T> => {
        loading.value = true
        error.value = null

        const token = getToken()
        const headers: HeadersInit = {
            ...options.headers,
        }

        if (!(options.body instanceof FormData)) {
            (headers as Record<string, string>)['Content-Type'] = 'application/json'
        }

        if (token) {
            (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`
        }

        try {
            const response = await fetch(endpoint, {
                ...options,
                headers,
            })

            // Handle 401 Unauthorized - token is invalid or user doesn't exist
            if (response.status === 401) {
                clearAuthAndRedirect()
                throw new Error('Session expired. Please login again.')
            }

            const data = await response.json()

            if (!response.ok) {
                throw new Error(data.error || 'Request failed')
            }

            return data as T
        } catch (err: any) {
            error.value = err.message
            throw err
        } finally {
            loading.value = false
        }
    }

    const get = <T>(endpoint: string) => request<T>(endpoint, { method: 'GET' })

    const post = <T>(endpoint: string, body: any) =>
        request<T>(endpoint, { method: 'POST', body: JSON.stringify(body) })

    const upload = <T>(endpoint: string, body: FormData) =>
        request<T>(endpoint, { method: 'POST', body })

    const put = <T>(endpoint: string, body: any) =>
        request<T>(endpoint, { method: 'PUT', body: JSON.stringify(body) })

    const del = <T>(endpoint: string) => request<T>(endpoint, { method: 'DELETE' })

    return {
        loading,
        error,
        get,
        post,
        put,
        del,
        upload,
    }
}
