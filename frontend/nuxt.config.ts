// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },

  // Enable SSR
  ssr: true,

  // Runtime config for API base URL
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080'
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
