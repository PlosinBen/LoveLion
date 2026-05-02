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

// Run in Asia/Taipei so the Z-suffix bug causes a visible +8h shift.
test.use({ timezoneId: 'Asia/Taipei' })

test.describe('Date round-trip', () => {
  test('transaction date survives save → display → edit without timezone shift', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '日常開銷')
    const token = await page.evaluate(() => localStorage.getItem('token'))

    // Use today's date at 20:30 so it appears at the top of the ledger.
    // The Z suffix tells Go it's RFC3339; the backend stores naive 20:30:00.
    const today = new Date()
    const yyyy = today.getFullYear()
    const mm = String(today.getMonth() + 1).padStart(2, '0')
    const dd = String(today.getDate()).padStart(2, '0')
    const dateStr = `${yyyy}-${mm}-${dd}T20:30:00Z`

    const res = await page.request.post(`/api/spaces/${spaceId}/expenses`, {
      headers: { Authorization: `Bearer ${token}` },
      data: {
        title: 'E2E 日期測試',
        total_amount: 100,
        date: dateStr,
        currency: 'TWD',
        expense: { category: '餐飲' },
      },
    })
    const txn = await res.json()
    expect(txn.id).toBeTruthy()

    // 1. Ledger list: date should show 20:30, not 04:30
    await page.goto(`/spaces/${spaceId}/ledger`)
    const row = page.locator('.cursor-pointer', { hasText: 'E2E 日期測試' })
    await expect(row).toBeVisible()
    await expect(row).toContainText('20:30')

    // 2. Detail page: date should show 20:30
    await row.click()
    await page.waitForURL(/\/transaction\//)
    await expect(page.getByText('20:30')).toBeVisible()

    // 3. Edit page: date picker should pre-populate with 20:30
    await page.locator('a[title="編輯"]').click()
    await page.waitForURL(/\/edit/)
    const dateInput = page.locator('.dp__input')
    await expect(dateInput).toHaveValue(/20:30/)

    // Cleanup: delete the transaction
    await page.goto(`/spaces/${spaceId}/ledger/transaction/${txn.id}`)
    await page.locator('button[title="刪除"]').click()
    await page.getByRole('button', { name: '確定' }).click()
    await page.waitForURL(/\/ledger$/)
  })
})
