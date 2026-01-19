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

    const request = async <T>(
        endpoint: string,
        options: RequestInit = {}
    ): Promise<T> => {
        loading.value = true
        error.value = null

        const token = getToken()
        const headers: HeadersInit = {
            'Content-Type': 'application/json',
            ...options.headers,
        }

        if (token) {
            (headers as Record<string, string>)['Authorization'] = `Bearer ${token}`
        }

        try {
            const response = await fetch(endpoint, {
                ...options,
                headers,
            })

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
    }
}
