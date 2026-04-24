import { test, expect } from './fixtures/auth'
import { login } from './fixtures/auth'

function extractSpaceId(url: string): string {
  return url.match(/\/spaces\/([^/]+)/)?.[1] ?? ''
}

async function enterSpace(page: import('@playwright/test').Page, name: string): Promise<string> {
  await page.goto('/')
  await page.getByText(name).first().click()
  await page.waitForURL(/\/spaces\/.*\/stats/)
  return extractSpaceId(page.url())
}

// The E2E stack runs with RECEIPT_EXTRACT_ENABLED=false, so the AI worker
// never picks up pending rows. That's what makes these tests deterministic —
// quick-text and image AI submissions stay in `pending` indefinitely and we
// can exercise the cancel flow without racing the worker.
test.describe('Quick text entry', () => {
  test('submitting quick text creates a pending row visible in the ledger', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '日常開銷')
    await page.goto(`/spaces/${spaceId}/ledger/transaction/add`)

    await page.getByPlaceholder('輸入一筆消費...').fill('E2E 停車費 120')
    await page.getByRole('button', { name: '送出' }).click()

    await page.waitForURL(/\/spaces\/[^/]+\/ledger$/)
    await expect(page.getByText('E2E 停車費 120')).toBeVisible()
  })
})

test.describe('AI cancel flow', () => {
  test('cancel on a pending transaction clears the banner and re-enables the form', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '日常開銷')

    // Seed a fresh pending row via quick text so this test is self-contained
    // (not reliant on the quick-text describe running first).
    await page.goto(`/spaces/${spaceId}/ledger/transaction/add`)
    await page.getByPlaceholder('輸入一筆消費...').fill('E2E 待取消的 AI 交易')
    await page.getByRole('button', { name: '送出' }).click()
    await page.waitForURL(/\/spaces\/[^/]+\/ledger$/)

    // Drill into the pending row → detail → edit (the cancel button lives on edit)
    await page.getByText('E2E 待取消的 AI 交易').first().click()
    await page.waitForURL(/\/transaction\/[^/]+$/)
    await page.getByRole('link', { name: '編輯' }).click()
    await page.waitForURL(/\/edit$/)

    // Pending banner (copy may be either 等待辨識中 or 辨識中 depending on timing)
    const cancelBtn = page.getByRole('button', { name: '取消辨識' })
    await expect(cancelBtn).toBeVisible()
    await cancelBtn.click()

    // After cancel, the banner disappears and the title field becomes editable.
    await expect(cancelBtn).not.toBeVisible()
    await expect(page.getByPlaceholder('例如: 午餐、計程車')).toBeEnabled()
  })
})

test.describe('Invite flow', () => {
  test('owner creates an invite and a second user joins the space', async ({ page }) => {
    // Phase 1 — dev (owner of 日常開銷) creates an invite.
    await login(page, 'dev', 'dev123')
    const spaceId = await enterSpace(page, '日常開銷')
    await page.goto(`/spaces/${spaceId}/settings`)

    await page.getByRole('button', { name: '+ 建立邀請' }).click()
    // The token isn't shown in full on screen (UI truncates), so capture it
    // by intercepting the invite-create response instead of scraping the DOM.
    const [response] = await Promise.all([
      page.waitForResponse(r => r.url().includes('/invites') && r.request().method() === 'POST'),
      page.getByRole('button', { name: '建立連結' }).click(),
    ])
    const invite = await response.json()
    const token = invite.token as string
    expect(token).toBeTruthy()

    // Phase 2 — switch to mei, accept the invite, verify the space appears.
    await page.evaluate(() => {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    })
    await login(page, 'mei', 'mei123')

    await page.goto(`/join/${token}`)
    await expect(page.getByText('受邀加入空間')).toBeVisible()
    await expect(page.getByText('日常開銷')).toBeVisible()
    await page.getByRole('button', { name: '接受邀請並加入' }).click()

    await page.waitForURL('/')
    // mei should now see 日常開銷 on her home page.
    await expect(page.getByText('日常開銷')).toBeVisible()
  })
})
