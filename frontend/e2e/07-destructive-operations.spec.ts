import { test, expect } from './fixtures/auth'
import { login } from './fixtures/auth'
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
// Payment Edit
// ---------------------------------------------------------------------------
test.describe('Payment Edit', () => {
  test('edit payment title and restore in trip space', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '2024 東京春櫻季')
    await page.goto(`/spaces/${spaceId}/ledger`)

    // Click the payment row
    await page.locator('.cursor-pointer', { hasText: '小明付款給 Antigravity' }).first().click()
    await page.waitForURL(/\/transaction\//)
    await expect(page.getByText('交易詳情')).toBeVisible()
    await expect(page.getByText('小明付款給 Antigravity')).toBeVisible()

    // Navigate to edit page
    await page.locator('a[title="編輯"]').click()
    await page.waitForURL(/\/edit/)
    await expect(page.getByText('編輯交易')).toBeVisible()

    // Verify form is populated
    const titleInput = page.getByPlaceholder('例如: 小明付款給小美')
    await expect(titleInput).toHaveValue('小明付款給 Antigravity')

    // Change title
    await titleInput.clear()
    await titleInput.fill('E2E 改名付款')
    await page.getByRole('button', { name: '儲存付款' }).click()

    // Should navigate back to detail with new title
    await page.waitForURL(/\/transaction\/[^/]+$/)
    await expect(page.locator('h2', { hasText: 'E2E 改名付款' })).toBeVisible()

    // Restore: edit again
    await page.locator('a[title="編輯"]').click()
    await page.waitForURL(/\/edit/)
    const titleInput2 = page.getByPlaceholder('例如: 小明付款給小美')
    await titleInput2.clear()
    await titleInput2.fill('小明付款給 Antigravity')
    await page.getByRole('button', { name: '儲存付款' }).click()
    await page.waitForURL(/\/transaction\/[^/]+$/)
    await expect(page.locator('h2', { hasText: '小明付款給 Antigravity' })).toBeVisible()
  })
})

// ---------------------------------------------------------------------------
// Payment Delete (create via API first, then delete via UI)
// ---------------------------------------------------------------------------
test.describe('Payment Delete', () => {
  test('create payment via API and delete it via UI', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '2024 東京春櫻季')

    // Create a temporary payment via API
    const token = await page.evaluate(() => localStorage.getItem('token'))
    const res = await page.request.post(`/api/spaces/${spaceId}/payments`, {
      headers: { Authorization: `Bearer ${token}` },
      data: {
        title: 'E2E 待刪付款',
        date: new Date().toISOString(),
        total_amount: 100,
        payer_name: '小明',
        payee_name: 'Antigravity',
      },
    })
    const created = await res.json()

    // Navigate to the transaction detail
    await page.goto(`/spaces/${spaceId}/ledger/transaction/${created.id}`)
    await expect(page.locator('h2', { hasText: 'E2E 待刪付款' })).toBeVisible()

    // Delete it
    await page.locator('button[title="刪除"]').click()
    await expect(page.getByText('確定要刪除此交易嗎？')).toBeVisible()
    await page.getByRole('button', { name: '確定' }).click()

    // Should redirect back to ledger without the deleted payment
    await page.waitForURL(/\/ledger$/)
    await expect(page.locator('h4', { hasText: 'E2E 待刪付款' })).not.toBeVisible()
  })
})

// ---------------------------------------------------------------------------
// Member Removal (owner removes a member)
// ---------------------------------------------------------------------------
test.describe('Member Removal', () => {
  test('owner removes member and re-invites them back', async ({ authedPage: page, browser }) => {
    const spaceId = await enterSpace(page, '2024 東京春櫻季')
    await page.goto(`/spaces/${spaceId}/settings`)

    // Verify 小美 is listed (use @mei which is unique)
    await expect(page.getByText('@mei')).toBeVisible()

    // Click the remove button on 小美's row (each member row has p-5 + border-b)
    const meiRow = page.locator('.p-5.border-b', { hasText: '@mei' })
    await meiRow.locator('button').last().click()

    // Confirm removal
    await expect(page.getByText('確定要移除成員 小美 嗎？')).toBeVisible()
    await page.getByRole('button', { name: '確定' }).click()

    // 小美 should be gone
    await expect(page.getByText('@mei')).not.toBeVisible()

    // Re-add 小美: create invite, then login as mei and join
    await page.getByText('+ 建立邀請').click()
    await expect(page.getByText('建立邀請連結')).toBeVisible()
    await page.getByRole('button', { name: '建立連結' }).click()

    // Wait for invite to be created and get the token
    await page.waitForTimeout(500)
    const inviteToken = await page.evaluate(async () => {
      const token = localStorage.getItem('token')
      const res = await fetch(window.location.origin + '/api/spaces/' +
        window.location.pathname.split('/')[2] + '/invites', {
        headers: { Authorization: `Bearer ${token}` }
      })
      const invites = await res.json()
      return invites[invites.length - 1]?.token
    })

    // Open new context, login as mei, and join via invite
    const meiPage = await browser.newPage()
    await login(meiPage, 'mei', 'mei123')
    await meiPage.goto(`/join/${inviteToken}`)
    await expect(meiPage.getByText('受邀加入空間')).toBeVisible()
    await meiPage.getByRole('button', { name: '接受邀請並加入' }).click()
    await meiPage.waitForURL('/')
    await meiPage.close()

    // Verify 小美 is back in the members list
    await page.goto(`/spaces/${spaceId}/settings`)
    await expect(page.getByText('@mei')).toBeVisible()
  })
})

// ---------------------------------------------------------------------------
// Product Deletion
// ---------------------------------------------------------------------------
test.describe('Product Deletion', () => {
  test('delete product from store and verify it disappears', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '2024 東京春櫻季')

    // First, add a temporary product via API
    const token = await page.evaluate(() => localStorage.getItem('token'))

    // Get the first store (唐吉軻德 澀谷本店)
    const storesRes = await page.request.get(`/api/spaces/${spaceId}/stores`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    const stores = await storesRes.json()
    const storeId = stores[0].id

    // Add a temporary product
    await page.request.post(`/api/spaces/${spaceId}/stores/${storeId}/products`, {
      headers: { Authorization: `Bearer ${token}` },
      data: { name: 'E2E 待刪商品', price: 999, currency: 'JPY' },
    })

    // Navigate to store detail
    await page.goto(`/spaces/${spaceId}/stores/${storeId}`)
    await expect(page.getByText('E2E 待刪商品')).toBeVisible()

    // Click the trash icon on the E2E product row
    const productCard = page.locator('div.p-4', { hasText: 'E2E 待刪商品' })
    await productCard.locator('button').last().click()

    // Confirm deletion
    await expect(page.getByText('確定要刪除此商品紀錄嗎？')).toBeVisible()
    await page.getByRole('button', { name: '確定' }).click()

    // Product should disappear
    await expect(page.getByText('E2E 待刪商品')).not.toBeVisible()
  })
})

// ---------------------------------------------------------------------------
// Space Deletion (create a temporary space, then delete it)
// ---------------------------------------------------------------------------
test.describe('Space Deletion', () => {
  test('create space and delete it from settings', async ({ authedPage: page }) => {
    // Create a temporary space via API
    const token = await page.evaluate(() => localStorage.getItem('token'))
    const res = await page.request.post('/api/spaces', {
      headers: { Authorization: `Bearer ${token}` },
      data: { name: 'E2E 待刪空間', type: 'personal', base_currency: 'TWD' },
    })
    const created = await res.json()
    const spaceId = created.id

    // Navigate to settings
    await page.goto(`/spaces/${spaceId}/settings`)
    await expect(page.getByText('E2E 待刪空間')).toBeVisible()

    // Click delete
    await page.getByRole('button', { name: '刪除此空間' }).click()

    // Confirm deletion
    await expect(page.getByText('警告：確定要刪除此空間嗎？')).toBeVisible()
    await page.getByRole('button', { name: '確定' }).click()

    // Should redirect to home
    await page.waitForURL('/')
    await expect(page.getByText('我的空間')).toBeVisible()
    // Deleted space should not appear
    await expect(page.getByText('E2E 待刪空間')).not.toBeVisible()
  })
})
