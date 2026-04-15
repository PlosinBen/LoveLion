import { test, expect } from './fixtures/auth'

const API_BASE = process.env.API_BASE || 'http://backend:8080'

test.describe('Core Flow', () => {
  test('create space → add expense via API → verify in ledger', async ({ authedPage: page }) => {
    // Step 1: Create a new space via UI
    await page.goto('/spaces/add-new')
    await page.getByPlaceholder('例如：日本旅行、個人記帳').fill('E2E 測試空間')
    await page.getByRole('button', { name: '建立空間' }).click()
    await page.waitForURL('/')

    // Step 2: Enter the space
    await page.getByText('E2E 測試空間').first().click()
    await page.waitForURL(/\/spaces\/.*\/stats/)
    const spaceId = page.url().match(/\/spaces\/([^/]+)/)?.[1]

    // Step 3: Create expense via API (avoids Playwright number input issue)
    const token = await page.evaluate(() => localStorage.getItem('token'))
    const res = await fetch(`${API_BASE}/api/spaces/${spaceId}/expenses`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        title: 'E2E 測試消費',
        total_amount: 150,
        date: new Date().toISOString(),
        currency: 'TWD',
        expense: {
          category: '餐飲',
          exchange_rate: 1,
          billing_amount: 150,
          handling_fee: 0,
          payment_method: '現金',
          items: [{ name: '測試項目', unit_price: 150, quantity: 1, discount: 0 }],
        },
      }),
    })
    expect(res.status).toBe(201)

    // Step 4: Navigate to ledger and verify
    await page.getByText('記帳').click()
    await page.waitForURL(/\/spaces\/.*\/ledger/)

    await expect(page.getByText('E2E 測試消費')).toBeVisible()
    await expect(page.getByText('150')).toBeVisible()
  })
})
