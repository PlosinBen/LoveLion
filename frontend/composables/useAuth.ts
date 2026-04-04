import { storeToRefs } from 'pinia'
import { useAuthStore } from '~/stores/auth'

export function useAuth() {
  const store = useAuthStore()
  const { user, token, isAuthenticated } = storeToRefs(store)
  return {
    user,
    token,
    isAuthenticated,
    initAuth: store.initAuth,
    login: store.login,
    register: store.register,
    logout: store.logout,
    fetchUser: store.fetchUser,
    updateProfile: store.updateProfile,
  }
}
