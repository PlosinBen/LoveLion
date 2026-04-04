import { useAuth } from '~/composables/useAuth'

export default defineNuxtRouteMiddleware((to) => {
  // Pages accessible without authentication
  if (to.path === '/login' || to.path.startsWith('/join/')) return

  const { initAuth, isAuthenticated } = useAuth()
  initAuth()

  if (!isAuthenticated.value) {
    return navigateTo('/login')
  }
})
