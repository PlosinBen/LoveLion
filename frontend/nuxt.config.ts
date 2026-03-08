// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: false },
  modules: ['@pinia/nuxt', '@nuxtjs/tailwindcss'],
  
  // css: ['~/assets/css/main.css'],

  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },

  // Vite 效能優化
  vite: {
    server: {
      watch: {
        // 在 Docker/Windows 環境下必須開啟輪詢才能即時偵測檔案變動
        usePolling: true,
        interval: 100,
      },
    },
    // 預編譯常用套件，減少啟動時的二次編譯
    optimizeDeps: {
      include: ['@iconify/vue', 'vue', 'vue-router']
    },
    // 增加編譯快取
    build: {
      cacheDir: './node_modules/.vite'
    }
  },

  // Nuxt 效能優化
  nitro: {
    // 開啟快取以加速後續啟動
    storage: {
      cache: {
        driver: 'fs',
        base: './.nuxt/cache'
      }
    }
  },

  routeRules: {
    '/api/**': {
      proxy: 'http://backend:8080/api/**'
    }
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || ''
    }
  }
})
