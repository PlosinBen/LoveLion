import { test as base, type Page } from '@playwright/test'

/**
 * Login helper — fills the login form and waits for redirect to home page.
 */
export async function login(page: Page, username: string, password: string) {
  await page.goto('/login')
  await page.getByPlaceholder('請輸入您的帳號').fill(username)
  await page.getByPlaceholder('請輸入您的密碼').fill(password)
  await page.getByRole('button', { name: '登入' }).click()
  await page.waitForURL('/')
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
