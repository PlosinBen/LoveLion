import { describe, it, expect } from 'vitest'
import { useLoading } from '~/composables/useLoading'

describe('useLoading', () => {
  it('starts with loading false', () => {
    const { isLoading } = useLoading()
    expect(isLoading.value).toBe(false)
  })

  it('showLoading sets isLoading to true', () => {
    const { isLoading, showLoading } = useLoading()
    showLoading()
    expect(isLoading.value).toBe(true)
  })

  it('hideLoading sets isLoading to false', () => {
    const { isLoading, showLoading, hideLoading } = useLoading()
    showLoading()
    expect(isLoading.value).toBe(true)
    hideLoading()
    expect(isLoading.value).toBe(false)
  })

  it('shares global state across instances', () => {
    const a = useLoading()
    const b = useLoading()
    a.showLoading()
    expect(b.isLoading.value).toBe(true)
    b.hideLoading()
    expect(a.isLoading.value).toBe(false)
  })
})
