import { test, expect } from './fixtures/auth'

/**
 * Helper: extract spaceId from current URL like /spaces/xxx/...
 */
function extractSpaceId(url: string): string {
  return url.match(/\/spaces\/([^/]+)/)?.[1] ?? ''
}

/**
 * Helper: enter a space by clicking its name on the home page.
 * Returns the spaceId.
 */
async function enterSpace(page: import('@playwright/test').Page, name: string): Promise<string> {
  await page.goto('/')
  await page.getByText(name).first().click()
  await page.waitForURL(/\/spaces\/.*\/stats/)
  return extractSpaceId(page.url())
}

test.describe('Space & Expense', () => {
  let spaceId: string

  test('create space', async ({ authedPage: page }) => {
    await page.goto('/spaces/add-new')
    await page.getByPlaceholder('例如：日本旅行、個人記帳').fill('E2E 測試空間')
    await page.getByRole('button', { name: '建立空間' }).click()
    await page.waitForURL('/')

    await page.getByText('E2E 測試空間').first().click()
    await page.waitForURL(/\/spaces\/.*\/stats/)
    spaceId = extractSpaceId(page.url())
    expect(spaceId).toBeTruthy()
  })

  test('add expense via UI', async ({ authedPage: page }) => {
    if (!spaceId) {
      spaceId = await enterSpace(page, 'E2E 測試空間')
    }

    await page.goto(`/spaces/${spaceId}/ledger/transaction/add`)
    await page.getByRole('button', { name: '消費' }).click()

    // Fill form
    await page.getByPlaceholder('例如: 午餐、計程車').fill('E2E 午餐')
    const totalInput = page.locator('label:has-text("總額")').locator('..').getByRole('spinbutton')
    await totalInput.click()
    await totalInput.fill('250')

    await page.getByRole('button', { name: '儲存交易' }).click()
    await page.waitForURL(/\/spaces\/.*\/ledger/)

    await expect(page.getByText('E2E 午餐')).toBeVisible()
  })

  test('view transaction detail', async ({ authedPage: page }) => {
    if (!spaceId) test.skip()

    await page.goto(`/spaces/${spaceId}/ledger`)
    await page.getByText('E2E 午餐').first().click()
    await page.waitForURL(/\/transaction\//)

    await expect(page.getByText('E2E 午餐')).toBeVisible()
  })

  test('edit transaction', async ({ authedPage: page }) => {
    if (!spaceId) test.skip()

    // Navigate directly to edit page via URL
    await page.goto(`/spaces/${spaceId}/ledger`)
    await page.getByText('E2E 午餐').first().click()
    await page.waitForURL(/\/transaction\//)

    // Get txnId from URL and navigate to edit
    const txnId = page.url().match(/\/transaction\/([^/]+)/)?.[1]
    await page.goto(`/spaces/${spaceId}/ledger/transaction/${txnId}/edit`)

    // Update title
    const titleInput = page.getByPlaceholder('例如: 午餐、計程車')
    await titleInput.clear()
    await titleInput.fill('E2E 晚餐')

    // Ensure total amount is valid (required for save)
    const totalInput = page.locator('label:has-text("總額")').locator('..').getByRole('spinbutton')
    await totalInput.clear()
    await totalInput.fill('250')

    await page.getByRole('button', { name: '儲存交易' }).click()

    // Edit page redirects to detail page after save
    await page.waitForURL(/\/transaction\/[^/]+$/)
    await expect(page.getByText('E2E 晚餐')).toBeVisible()
  })

  test('delete transaction', async ({ authedPage: page }) => {
    if (!spaceId) test.skip()

    await page.goto(`/spaces/${spaceId}/ledger`)
    await page.getByText('E2E 晚餐').first().click()
    await page.waitForURL(/\/transaction\//)

    await page.getByRole('button', { name: '刪除' }).click()
    await page.getByRole('button', { name: '確定' }).click()
    await page.waitForURL(/\/spaces\/.*\/ledger/)

    await expect(page.getByText('E2E 晚餐')).not.toBeVisible()
  })
})

test.describe('Payment', () => {
  test('add payment via UI', async ({ authedPage: page }) => {
    // Use trip space that has split_members (Antigravity, 小明, 小美)
    const spaceId = await enterSpace(page, '2024 東京春櫻季')

    await page.goto(`/spaces/${spaceId}/ledger/transaction/add`)
    await page.getByRole('button', { name: '付款' }).click()

    await page.getByPlaceholder('例如: 小明付款給小美').fill('E2E 測試付款')

    // Select payer and payee by member name
    const payerSelect = page.locator('select').first()
    const payeeSelect = page.locator('select').nth(1)
    await payerSelect.selectOption('Antigravity')
    await payeeSelect.selectOption('小明')

    // Fill amount
    const amountInput = page.locator('label:has-text("金額")').locator('..').getByRole('spinbutton')
    await amountInput.click()
    await amountInput.fill('500')

    await page.getByRole('button', { name: '儲存付款' }).click()
    await page.waitForURL(/\/spaces\/.*\/ledger/)

    await expect(page.getByText('E2E 測試付款')).toBeVisible()
  })
})

test.describe('Comparison (Stores & Products)', () => {
  let spaceId: string

  test('create store', async ({ authedPage: page }) => {
    spaceId = await enterSpace(page, '日常開銷')

    await page.goto(`/spaces/${spaceId}/stores/add`)
    await page.getByPlaceholder('例如：唐吉訶德、大國藥妝').fill('E2E 測試商店')
    await page.getByRole('button', { name: '建立店家' }).click()
    await page.waitForURL(/\/spaces\/.*\/stores/)

    await expect(page.getByText('E2E 測試商店')).toBeVisible()
  })

  test('add product to store', async ({ authedPage: page }) => {
    if (!spaceId) test.skip()

    await page.goto(`/spaces/${spaceId}/stores`)
    await page.getByText('E2E 測試商店').first().click()
    await page.waitForURL(/\/stores\//)

    const storeUrl = page.url()
    await page.goto(`${storeUrl}/products/add`)

    await page.getByPlaceholder('例如：pocky、午後紅茶').fill('E2E 測試商品')
    await page.getByRole('spinbutton').first().fill('99')
    await page.getByRole('button', { name: '新增商品' }).click()

    await expect(page.getByText('E2E 測試商品')).toBeVisible()
  })
})

test.describe('User Settings', () => {
  test('update display name', async ({ authedPage: page }) => {
    await page.goto('/settings')

    // Click edit pencil button in the profile section (帳戶資訊)
    const profileSection = page.locator('section').filter({ hasText: '帳戶資訊' })
    await profileSection.getByRole('button').click()

    // Wait for modal and fill new name
    await expect(page.getByText('編輯基本資料')).toBeVisible()
    const nameInput = page.getByPlaceholder('您希望如何被稱呼？')
    await nameInput.clear()
    await nameInput.fill('E2E Display Name')
    await page.getByRole('button', { name: '儲存修改' }).click()

    // Verify updated
    await expect(page.getByText('E2E Display Name')).toBeVisible()

    // Restore original
    await profileSection.getByRole('button').click()
    await expect(page.getByText('編輯基本資料')).toBeVisible()
    const restoreInput = page.getByPlaceholder('您希望如何被稱呼？')
    await restoreInput.clear()
    await restoreInput.fill('Antigravity')
    await page.getByRole('button', { name: '儲存修改' }).click()
  })

  test('change password shows error for mismatched confirmation', async ({ authedPage: page }) => {
    await page.goto('/settings')

    await page.getByPlaceholder('請輸入目前密碼').fill('dev123')
    await page.getByPlaceholder('請輸入新密碼 (至少 6 個字)').fill('newpass123')
    await page.getByPlaceholder('請再次輸入新密碼').fill('differentpass')
    await page.getByRole('button', { name: '確認修改' }).click()

    // Client-side validation: mismatched passwords show error
    await expect(page.getByText('新密碼與確認密碼不符')).toBeVisible()
  })
})

test.describe('Space Settings', () => {
  test('update space name and restore', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '日常開銷')

    await page.goto(`/spaces/${spaceId}/settings`)

    const nameInput = page.getByPlaceholder('請輸入空間名稱')
    await nameInput.clear()
    await nameInput.fill('E2E 更新名稱')
    await page.getByRole('button', { name: '儲存基本設定' }).click()

    // Verify by reloading
    await page.goto(`/spaces/${spaceId}/settings`)
    await expect(page.getByPlaceholder('請輸入空間名稱')).toHaveValue('E2E 更新名稱')

    // Restore
    const restoreInput = page.getByPlaceholder('請輸入空間名稱')
    await restoreInput.clear()
    await restoreInput.fill('日常開銷')
    await page.getByRole('button', { name: '儲存基本設定' }).click()
  })
})
