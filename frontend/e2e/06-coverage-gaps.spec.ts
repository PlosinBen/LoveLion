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

// ---------------------------------------------------------------------------
// Stats Page
// ---------------------------------------------------------------------------
test.describe('Stats Page', () => {
  test('personal space shows spending total and category breakdown', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '日常開銷')
    // enterSpace already lands on /stats
    await expect(page.getByText('消費統計')).toBeVisible()
    await expect(page.getByText('空間累積支出')).toBeVisible()
    await expect(page.getByText('類別分布')).toBeVisible()
    await expect(page.getByText('餐飲')).toBeVisible()
    await expect(page.getByText('交通')).toBeVisible()
  })

  test('trip space shows settlement and suggested transfers', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '2024 東京春櫻季')
    await expect(page.getByText('消費統計')).toBeVisible()
    await expect(page.getByText('空間累積支出')).toBeVisible()
    await expect(page.getByText('應收應付')).toBeVisible()
    await expect(page.getByText('建議付款')).toBeVisible()
  })
})

// ---------------------------------------------------------------------------
// Products Comparison Page
// ---------------------------------------------------------------------------
test.describe('Products Comparison', () => {
  let spaceId: string

  test('products page groups items by name', async ({ authedPage: page }) => {
    spaceId = await enterSpace(page, '2024 東京春櫻季')
    await page.goto(`/spaces/${spaceId}/products`)

    await expect(page.getByRole('heading', { name: '商品' })).toBeVisible()
    await expect(page.getByText('一蘭拉麵泡麵')).toBeVisible()
    await expect(page.getByText('Dyson 吹風機')).toBeVisible()
    await expect(page.getByText('2 個店面')).toBeVisible()
    await expect(page.getByText('1 個店面')).toBeVisible()
  })

  test('product detail shows price comparison with cheapest badge', async ({ authedPage: page }) => {
    if (!spaceId) spaceId = await enterSpace(page, '2024 東京春櫻季')
    await page.goto(`/spaces/${spaceId}/products`)

    await page.getByText('一蘭拉麵泡麵').first().click()
    await expect(page.getByText('價格對比')).toBeVisible()
    await expect(page.getByText('唐吉軻德 澀谷店')).toBeVisible()
    await expect(page.getByText('Bic Camera 新宿')).toBeVisible()
    await expect(page.getByText('最便宜')).toBeVisible()
    // Cheapest is 1,850 JPY at Donki
    await expect(page.getByText('1,850').first()).toBeVisible()
  })
})

// ---------------------------------------------------------------------------
// Expense Templates
// ---------------------------------------------------------------------------
test.describe('Expense Templates', () => {
  test('save transaction as template, apply it, then delete', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '日常開銷')
    await page.goto(`/spaces/${spaceId}/ledger`)

    // Navigate to 星巴克 transaction detail — click the card, not the inner text
    await page.locator('.cursor-pointer', { hasText: '星巴克' }).first().click()
    await page.waitForURL(/\/transaction\//)
    await expect(page.getByText('交易詳情')).toBeVisible()

    // Save as template via the amber button
    await page.locator('button[title="儲存為模板"]').click()
    await expect(page.getByText('儲存為模板')).toBeVisible()
    // The prompt pre-fills with the transaction title
    const promptInput = page.locator('input[placeholder="輸入模板名稱"]')
    await expect(promptInput).toHaveValue('星巴克')
    // Rename and confirm
    await promptInput.clear()
    await promptInput.fill('E2E 模板')
    await page.getByRole('button', { name: '確定' }).click()

    // Navigate to add transaction page
    await page.goto(`/spaces/${spaceId}/ledger/transaction/add`)
    const templateBtn = page.getByRole('button', { name: '套用模板' })
    await expect(templateBtn).toBeVisible()
    await templateBtn.click()

    // Template picker shows our template
    await expect(page.getByText('E2E 模板')).toBeVisible()

    // Apply the template
    await page.getByText('E2E 模板').click()
    // Verify the title field is populated from template
    await expect(page.getByPlaceholder('例如: 午餐、計程車')).toHaveValue('星巴克')

    // Clean up: reopen picker and delete the template
    await templateBtn.click()
    await expect(page.getByText('E2E 模板')).toBeVisible()
    // Click the delete button on the template row
    const templateRow = page.locator('div').filter({ hasText: /^E2E 模板/ })
    await templateRow.getByRole('button').click()
    // Confirm the deletion dialog
    await expect(page.getByText('確定要刪除此模板嗎？')).toBeVisible()
    await page.getByRole('button', { name: '確定' }).click()
    // After deletion, template should disappear
    await expect(page.getByText('E2E 模板')).not.toBeVisible()
    await expect(page.getByText('尚無模板')).toBeVisible()
  })
})

// ---------------------------------------------------------------------------
// Pagination
// ---------------------------------------------------------------------------
test.describe('Pagination', () => {
  test('transaction count displays and pagination appears after bulk insert', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '日常開銷')
    await page.goto(`/spaces/${spaceId}/ledger`)

    // Verify count label exists (seed has 2 expenses; earlier tests may have added more)
    await expect(page.getByText(/交易紀錄 \(\d+\)/)).toBeVisible()

    // Pagination should be hidden when under 50 items
    await expect(page.getByText('上一頁')).not.toBeVisible()

    // Bulk-create 52 expenses via API to trigger pagination
    const token = await page.evaluate(() => localStorage.getItem('token'))
    for (let i = 0; i < 52; i++) {
      await page.request.post(`/api/spaces/${spaceId}/expenses`, {
        headers: { Authorization: `Bearer ${token}` },
        data: {
          title: `Pagination #${i}`,
          date: new Date().toISOString(),
          currency: 'TWD',
          expense: {
            category: '餐飲',
            items: [{ name: 'item', unit_price: 10, quantity: 1 }],
          },
        },
      })
    }

    // Reload and verify pagination controls appear
    await page.goto(`/spaces/${spaceId}/ledger`)
    await expect(page.getByText('上一頁')).toBeVisible()
    await expect(page.getByText('下一頁')).toBeVisible()
    await expect(page.getByText('1 / 2')).toBeVisible()

    // Navigate to page 2
    await Promise.all([
      page.waitForResponse(r => r.url().includes('/transactions') && r.url().includes('offset=')),
      page.getByText('下一頁').click(),
    ])
    await expect(page.getByText('2 / 2')).toBeVisible()

    // Navigate back to page 1
    await Promise.all([
      page.waitForResponse(r => r.url().includes('/transactions')),
      page.getByText('上一頁').click(),
    ])
    await expect(page.getByText('1 / 2')).toBeVisible()
  })
})
