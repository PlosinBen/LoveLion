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

test.describe('Ledger layout regressions', () => {
  test.use({ viewport: { width: 375, height: 812 } })

  test('a long title does not widen the viewport or displace fixed controls', async ({ authedPage: page }) => {
    const spaceId = await enterSpace(page, '日常開銷')
    const token = await page.evaluate(() => localStorage.getItem('token'))
    const title = `E2E-${'LONGTITLE'.repeat(9)}`

    const response = await page.request.post(`/api/spaces/${spaceId}/expenses`, {
      headers: { Authorization: `Bearer ${token}` },
      data: {
        title,
        total_amount: 1234,
        date: new Date().toISOString(),
        currency: 'TWD',
        expense: { category: '餐飲' },
      },
    })
    const transaction = await response.json()
    expect(transaction.id).toBeTruthy()

    await page.goto(`/spaces/${spaceId}/ledger`)
    const titleElement = page.getByRole('heading', { name: title })
    await expect(titleElement).toBeVisible()

    const layout = await titleElement.evaluate((element) => ({
      titleIsTruncated: element.scrollWidth > element.clientWidth,
      viewportHasNoHorizontalOverflow: document.documentElement.scrollWidth <= window.innerWidth,
    }))

    expect(layout.titleIsTruncated).toBe(true)
    expect(layout.viewportHasNoHorizontalOverflow).toBe(true)

    const footer = page.locator('footer')
    const addButton = page.getByRole('button', { name: '新增交易' })
    await expect(footer).toBeVisible()
    await expect(addButton).toBeVisible()

    const fixedControls = await page.evaluate(() => {
      const footerElement = document.querySelector('footer')
      const addButtonElement = document.querySelector('[aria-label="新增交易"]')
      if (!footerElement || !addButtonElement) return null

      const footerRect = footerElement.getBoundingClientRect()
      const addButtonRect = addButtonElement.getBoundingClientRect()
      return {
        footerPosition: getComputedStyle(footerElement).position,
        addButtonPosition: getComputedStyle(addButtonElement).position,
        footerIsInViewport: footerRect.top >= 0 && footerRect.bottom <= window.innerHeight,
        addButtonIsInViewport: addButtonRect.top >= 0 && addButtonRect.bottom <= window.innerHeight,
      }
    })

    expect(fixedControls).toEqual({
      footerPosition: 'fixed',
      addButtonPosition: 'fixed',
      footerIsInViewport: true,
      addButtonIsInViewport: true,
    })

    // Keep the shared E2E seed data clean for subsequent local runs.
    await page.goto(`/spaces/${spaceId}/ledger/transaction/${transaction.id}`)
    await page.locator('button[title="刪除"]').click()
    await page.getByRole('button', { name: '確定' }).click()
    await page.waitForURL(/\/ledger$/)
  })
})
