// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },

  // Modules
  modules: ['@nuxtjs/tailwindcss'],

  // Enable SSR
  ssr: true,

  // Runtime config for API base URL
  runtimeConfig: {
    public: {
      apiBase: ''  // Empty for same-origin requests
    }
  },

  // Nitro server configuration
  nitro: {
    routeRules: {
      '/api/**': {
        proxy: 'http://backend:8080/api/**'
      },
      '/health': {
        proxy: 'http://backend:8080/health'
      }
    }
  },

  // Development server settings
  devServer: {
    host: '0.0.0.0',
    port: 3000
  },

  // App configuration
  app: {
    head: {
      title: 'LoveLion - Personal Bookkeeping & Travel Assistant',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Personal bookkeeping and travel expense tracking application' }
      ]
    }
  }
})
