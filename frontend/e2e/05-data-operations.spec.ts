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

test.describe('Ledger Filters & Settings', () => {
  let spaceId: string

  test('search filters transactions by text', async ({ authedPage: page }) => {
    spaceId = await enterSpace(page, '日常開銷')
    await page.goto(`/spaces/${spaceId}/ledger`)
    await expect(page.getByText('星巴克')).toBeVisible()
    await expect(page.getByText('捷運定期票')).toBeVisible()

    const searchInput = page.getByPlaceholder('搜尋交易...')
    await searchInput.click()
    await Promise.all([
      page.waitForResponse(r => r.url().includes('/transactions') && r.url().includes('search=')),
      searchInput.fill('星巴克'),
    ])

    await expect(page.getByRole('heading', { name: '星巴克' })).toBeVisible()
    await expect(page.getByText('捷運定期票')).not.toBeVisible()

    // Clear search — restore all
    await Promise.all([
      page.waitForResponse(r => r.url().includes('/transactions') && !r.url().includes('search=')),
      searchInput.clear(),
    ])
    await expect(page.getByText('星巴克')).toBeVisible()
    await expect(page.getByText('捷運定期票')).toBeVisible()
  })

  test('type filter shows only matching type', async ({ authedPage: page }) => {
    if (!spaceId) spaceId = await enterSpace(page, '日常開銷')
    await page.goto(`/spaces/${spaceId}/ledger`)

    const paymentBtn = page.getByRole('button', { name: '付款' }).first()
    await Promise.all([
      page.waitForResponse(r => r.url().includes('type=payment')),
      paymentBtn.click(),
    ])
    await expect(page.getByText('沒有符合條件的交易')).toBeVisible()

    await Promise.all([
      page.waitForResponse(r => r.url().includes('/transactions') && !r.url().includes('type=')),
      paymentBtn.click(),
    ])
    await expect(page.getByText('星巴克')).toBeVisible()

    const expenseBtn = page.getByRole('button', { name: '消費' }).first()
    await Promise.all([
      page.waitForResponse(r => r.url().includes('type=expense')),
      expenseBtn.click(),
    ])
    await expect(page.getByText('星巴克')).toBeVisible()
    await expect(page.getByText('捷運定期票')).toBeVisible()

    await expenseBtn.click()
  })

  test('category filter narrows to single category', async ({ authedPage: page }) => {
    if (!spaceId) spaceId = await enterSpace(page, '日常開銷')
    await page.goto(`/spaces/${spaceId}/ledger`)

    const transportBtn = page.getByRole('button', { name: '交通' })
    await Promise.all([
      page.waitForResponse(r => r.url().includes('category=')),
      transportBtn.click(),
    ])
    await expect(page.getByText('捷運定期票')).toBeVisible()
    await expect(page.getByRole('heading', { name: '星巴克' })).not.toBeVisible()

    await Promise.all([
      page.waitForResponse(r => r.url().includes('/transactions') && !r.url().includes('category=')),
      transportBtn.click(),
    ])

    const foodBtn = page.getByRole('button', { name: '餐飲' })
    await Promise.all([
      page.waitForResponse(r => r.url().includes('category=')),
      foodBtn.click(),
    ])
    await expect(page.getByRole('heading', { name: '星巴克' })).toBeVisible()
    await expect(page.getByText('捷運定期票')).not.toBeVisible()

    await foodBtn.click()
  })

  test('ListEditor: add, duplicate check, and remove items', async ({ authedPage: page }) => {
    if (!spaceId) spaceId = await enterSpace(page, '日常開銷')
    await page.goto(`/spaces/${spaceId}/settings`)

    // --- Category ---
    const catInput = page.getByPlaceholder('新增分類...')
    await catInput.fill('E2E分類')
    await catInput.press('Enter')

    const catTag = page.locator('div').filter({ hasText: /^E2E分類$/ })
    await expect(catTag).toBeVisible()

    await catInput.fill('E2E分類')
    await catInput.press('Enter')
    await expect(page.getByText('E2E分類')).toHaveCount(1)

    await catTag.getByRole('button').click()
    await expect(page.getByText('E2E分類')).not.toBeVisible()

    // --- Payment method ---
    const payInput = page.getByPlaceholder('例如: 現金, 信用卡...')
    await payInput.fill('E2E支付')
    await payInput.press('Enter')

    const payTag = page.locator('div').filter({ hasText: /^E2E支付$/ })
    await expect(payTag).toBeVisible()

    await payInput.fill('E2E支付')
    await payInput.press('Enter')
    await expect(page.getByText('E2E支付')).toHaveCount(1)

    await payTag.getByRole('button').click()
    await expect(page.getByText('E2E支付')).not.toBeVisible()
  })

  test('ListEditor: save and persist changes', async ({ authedPage: page }) => {
    if (!spaceId) spaceId = await enterSpace(page, '日常開銷')
    await page.goto(`/spaces/${spaceId}/settings`)

    const catInput = page.getByPlaceholder('新增分類...')
    await catInput.fill('E2E持久分類')
    await catInput.press('Enter')
    await expect(page.getByText('E2E持久分類')).toBeVisible()

    await Promise.all([
      page.waitForResponse(r => r.request().method() === 'PATCH' && r.url().includes('/spaces/')),
      page.getByRole('button', { name: '儲存分類設定' }).click(),
    ])

    await page.goto(`/spaces/${spaceId}/settings`)
    await expect(page.getByText('E2E持久分類')).toBeVisible()

    const tag = page.locator('div').filter({ hasText: /^E2E持久分類$/ })
    await tag.getByRole('button').click()
    await Promise.all([
      page.waitForResponse(r => r.request().method() === 'PATCH' && r.url().includes('/spaces/')),
      page.getByRole('button', { name: '儲存分類設定' }).click(),
    ])

    await page.goto(`/spaces/${spaceId}/settings`)
    await expect(page.getByText('E2E持久分類')).not.toBeVisible()
  })
})

test.describe('Member Alias', () => {
  test('edit and restore member alias', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '2024 東京春櫻季')
    await page.goto(`/spaces/${spaceId}/settings`)

    await expect(page.getByText('@ming')).toBeVisible({ timeout: 15000 })

    const mingRow = page.locator('.flex.items-center.justify-between').filter({ hasText: '小明' }).filter({ hasText: '@ming' })
    await mingRow.locator('button').first().click()

    await expect(page.getByText('修改成員別名')).toBeVisible()
    const aliasInput = page.getByPlaceholder('請輸入成員別名')
    await aliasInput.clear()
    await aliasInput.fill('MingE2E')
    await page.getByRole('button', { name: '更新別名' }).click()

    await expect(page.getByText('修改成員別名')).not.toBeVisible({ timeout: 10000 })
    await expect(page.getByText('MingE2E')).toBeVisible()

    const restoreRow = page.locator('.flex.items-center.justify-between').filter({ hasText: 'MingE2E' }).filter({ hasText: '@ming' })
    await restoreRow.locator('button').first().click()
    await expect(page.getByText('修改成員別名')).toBeVisible()
    const restoreInput = page.getByPlaceholder('請輸入成員別名')
    await restoreInput.clear()
    await restoreInput.fill('小明')
    await page.getByRole('button', { name: '更新別名' }).click()
    await expect(page.getByText('修改成員別名')).not.toBeVisible({ timeout: 10000 })
  })
})

test.describe('Store Edit', () => {
  test('edit store name and restore', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '日常開銷')
    await page.goto(`/spaces/${spaceId}/stores`)

    await page.getByText('E2E 測試商店').first().click()
    await page.waitForURL(/\/stores\//)

    const storeDetailUrl = page.url()
    await page.goto(`${storeDetailUrl}/edit`)

    const nameInput = page.getByPlaceholder('例如：全家便利商店、宜得利')
    await nameInput.clear()
    await nameInput.fill('E2E 商店改名')
    await page.getByRole('button', { name: '儲存變更' }).click()
    await expect(page.getByRole('heading', { name: 'E2E 商店改名' })).toBeVisible({ timeout: 10000 })

    await page.goto(`${storeDetailUrl}/edit`)
    const restoreInput = page.getByPlaceholder('例如：全家便利商店、宜得利')
    await restoreInput.clear()
    await restoreInput.fill('E2E 測試商店')
    await page.getByRole('button', { name: '儲存變更' }).click()
    await expect(page.getByRole('heading', { name: 'E2E 測試商店' })).toBeVisible({ timeout: 10000 })
  })
})
