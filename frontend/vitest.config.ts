import { defineConfig } from 'vitest/config'
import { resolve } from 'path'

export default defineConfig({
  resolve: {
    alias: {
      '~': resolve(__dirname, '.'),
      '#imports': resolve(__dirname, '.vitest/mock-imports.ts'),
    },
  },
  test: {
    environment: 'node',
    include: ['**/*.test.ts'],
    globals: true,
  },
})
