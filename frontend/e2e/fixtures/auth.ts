import { test as base, expect, type Page } from '@playwright/test'

export async function login(page: Page, username: string, password: string) {
  await page.goto('/login')
  await page.getByPlaceholder('請輸入您的帳號').fill(username)
  await page.getByPlaceholder('請輸入您的密碼').fill(password)
  await page.getByRole('button', { name: '登入' }).click()
  await expect(page.getByText('我的空間')).toBeVisible({ timeout: 15000 })
}

/**
 * Fixture that provides a page already logged in as dev/dev123.
 */
export const test = base.extend<{ authedPage: Page }>({
  authedPage: async ({ page }, use) => {
    await login(page, 'dev', 'dev123')
    await use(page)
  },
})

export { expect } from '@playwright/test'
