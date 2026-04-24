import type { RouterConfig } from '@nuxt/schema'

export default <RouterConfig>{
  scrollBehavior(to, from, savedPosition) {
    if (to.matched.length > from.matched.length) {
      return false
    }
    if (to.matched.length < from.matched.length) {
      return false
    }
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0 }
  },
}
