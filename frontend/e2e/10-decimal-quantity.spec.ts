import { test, expect } from './fixtures/auth'
import type { Page } from '@playwright/test'

function extractSpaceId(url: string): string {
  return url.match(/\/spaces\/([^/]+)/)?.[1] ?? ''
}

async function enterSpace(page: Page, name: string): Promise<string> {
  await page.goto('/')
  await page.getByText(name).first().click()
  await page.waitForURL(/\/spaces\/.*\/stats/)
  return extractSpaceId(page.url())
}

test('an item quantity accepts and persists two decimal places', async ({ authedPage: page }) => {
  const spaceId = await enterSpace(page, '日常開銷')
  const token = await page.evaluate(() => localStorage.getItem('token'))

  const response = await page.request.post(`/api/spaces/${spaceId}/expenses`, {
    headers: { Authorization: `Bearer ${token}` },
    data: {
      title: 'E2E 小數數量',
      total_amount: 107.42,
      date: new Date().toISOString(),
      currency: 'TWD',
      expense: {
        category: '交通',
        items: [{ name: '汽油', unit_price: 31.5, quantity: 1 }],
      },
    },
  })
  const transaction = await response.json()
  expect(transaction.id).toBeTruthy()

  await page.goto(`/spaces/${spaceId}/ledger/transaction/${transaction.id}/edit`)
  const quantityInput = page.locator('label', { hasText: '數量' }).locator('..').getByRole('spinbutton')
  await expect(quantityInput).toHaveAttribute('step', '0.01')
  await expect(quantityInput).toHaveAttribute('min', '0.01')
  await quantityInput.fill('3.41')
  await page.getByRole('button', { name: '儲存交易' }).click()
  await page.waitForURL(/\/transaction\/[^/]+$/)

  const savedResponse = await page.request.get(`/api/spaces/${spaceId}/transactions/${transaction.id}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  const savedTransaction = await savedResponse.json()
  expect(savedTransaction.expense.items[0].quantity).toBe('3.41')

  await page.locator('button[title="刪除"]').click()
  await page.getByRole('button', { name: '確定' }).click()
  await page.waitForURL(/\/ledger$/)
})
