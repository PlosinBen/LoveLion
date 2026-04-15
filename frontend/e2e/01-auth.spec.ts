import { test, expect } from '@playwright/test'
import { login } from './fixtures/auth'

test.describe('Authentication', () => {
  test('login with valid credentials redirects to home', async ({ page }) => {
    await login(page, 'dev', 'dev123')
    await expect(page).toHaveURL('/')
    await expect(page.getByText('我的空間')).toBeVisible()
  })

  test('login with wrong password stays on login page', async ({ page }) => {
    await page.goto('/login')
    await page.getByPlaceholder('請輸入您的帳號').fill('dev')
    await page.getByPlaceholder('請輸入您的密碼').fill('wrongpassword')
    await page.getByRole('button', { name: '登入' }).click()

    // 401 triggers full page reload to /login (useApi clears auth + hard redirect)
    await page.waitForURL('/login')
    await expect(page.getByText('歡迎回來')).toBeVisible()
  })

  test('register new account and redirect to home', async ({ page }) => {
    const username = `e2e_${Date.now()}`

    await page.goto('/login')
    await page.getByRole('button', { name: '立即註冊' }).click()
    await expect(page.getByText('加入 LoveLion')).toBeVisible()

    await page.getByPlaceholder('設定登入帳號').fill(username)
    await page.getByPlaceholder('大家如何稱呼您').fill('E2E Tester')
    await page.getByPlaceholder('設定登入密碼').fill('test123')
    await page.getByRole('button', { name: '註冊帳號' }).click()

    await page.waitForURL('/')
    await expect(page.getByText('我的空間')).toBeVisible()
  })

  test('logout redirects to login page', async ({ page }) => {
    await login(page, 'dev', 'dev123')

    await page.goto('/settings')
    await page.getByRole('button', { name: '登出帳戶' }).click()

    await expect(page).toHaveURL('/login')
  })
})
