import { test, expect } from './fixtures/auth'

const API_BASE = process.env.API_BASE || 'http://backend:8080'

/**
 * Seed a published announcement via the backend API directly.
 */
async function seedAnnouncement(authToken: string) {
  const createRes = await fetch(`${API_BASE}/api/admin/announcements`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${authToken}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      title: 'E2E 測試公告',
      content: '這是一則 **E2E 測試** 公告內容。',
      status: 'published',
    }),
  })

  return { status: createRes.status, isAdmin: createRes.status === 201 }
}

test.describe('Announcements', () => {
  let announcementSeeded = false

  test.beforeAll(async ({ browser }) => {
    // Login to get auth token
    const page = await browser.newPage()
    await page.goto('/login')
    await page.getByPlaceholder('請輸入您的帳號').fill('dev')
    await page.getByPlaceholder('請輸入您的密碼').fill('dev123')
    await page.getByRole('button', { name: '登入' }).click()
    await page.waitForURL('/')

    // Extract token from localStorage
    const token = await page.evaluate(() => localStorage.getItem('token'))
    if (token) {
      const result = await seedAnnouncement(token)
      announcementSeeded = result.isAdmin
    }

    await page.close()
  })

  test('announcement list page shows announcements', async ({ authedPage: page }) => {
    test.skip(!announcementSeeded, 'dev user is not admin — cannot seed announcement')

    await page.goto('/announcements')
    await expect(page.getByText('E2E 測試公告')).toBeVisible()
  })

  test('announcement detail page renders content', async ({ authedPage: page }) => {
    test.skip(!announcementSeeded, 'dev user is not admin — cannot seed announcement')

    await page.goto('/announcements')
    await page.getByText('E2E 測試公告').click()
    await page.waitForURL(/\/announcements\//)

    await expect(page.getByRole('heading', { name: 'E2E 測試公告' }).first()).toBeVisible()
    await expect(page.getByText('E2E 測試', { exact: true })).toBeVisible()
  })

  test('settings page links to announcements', async ({ authedPage: page }) => {
    await page.goto('/settings')
    await page.getByText('公告').first().click()
    await page.waitForURL('/announcements')
  })
})
