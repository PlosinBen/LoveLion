// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: false },
  ssr: false,
  modules: ['@pinia/nuxt', '@nuxtjs/tailwindcss'],

  app: {
    head: {
      title: 'LoveLion',
      meta: [
        { name: 'description', content: '共享記帳與比價應用程式' },
        { name: 'theme-color', content: '#171717' },
        { name: 'mobile-web-app-capable', content: 'yes' },
        { name: 'apple-mobile-web-app-capable', content: 'yes' },
        { name: 'apple-mobile-web-app-status-bar-style', content: 'black-translucent' },
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
        { rel: 'manifest', href: '/manifest.json' },
        { rel: 'apple-touch-icon', href: '/favicon-155.png' },
      ],
    },
  },
  
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
      allowedHosts: ['frontend'],
      watch: {
        // 在 Docker/Windows 環境下必須開啟輪詢才能即時偵測檔案變動
        usePolling: process.env.NODE_ENV !== 'production',
        interval: 100,
      },
    },
    // 預編譯常用套件，減少啟動時的二次編譯
    optimizeDeps: {
      include: ['@iconify/vue', 'vue', 'vue-router']
    },
    // 增加編譯快取
    // Removed invalid cacheDir
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
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '',
      appVersion: process.env.NUXT_PUBLIC_APP_VERSION || ''
    }
  }
})
