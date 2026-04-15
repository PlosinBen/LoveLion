import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  fullyParallel: false,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  maxFailures: process.env.CI ? 0 : 2,
  workers: 1,
  reporter: [['list'], ['html', { open: 'never' }]],
  timeout: 30_000,

  use: {
    baseURL: process.env.BASE_URL || 'http://frontend:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { browserName: 'chromium' },
    },
  ],
})
