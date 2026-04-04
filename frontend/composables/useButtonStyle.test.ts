import { describe, it, expect } from 'vitest'
import { useButtonStyle } from './useButtonStyle'

describe('useButtonStyle', () => {
  it('returns primary variant by default', () => {
    const classes = useButtonStyle()
    expect(classes).toContain('bg-indigo-500')
    expect(classes).toContain('text-white')
  })

  it('returns secondary variant classes', () => {
    const classes = useButtonStyle('secondary')
    expect(classes).toContain('bg-neutral-800')
    expect(classes).toContain('text-neutral-400')
  })

  it('returns danger variant classes', () => {
    const classes = useButtonStyle('danger')
    expect(classes).toContain('text-red-500')
  })

  it('returns ghost variant classes', () => {
    const classes = useButtonStyle('ghost')
    expect(classes).toContain('bg-transparent')
  })

  it('falls back to primary for unknown variant', () => {
    const classes = useButtonStyle('nonexistent')
    expect(classes).toContain('bg-indigo-500')
  })

  it('always includes base classes', () => {
    const variants = ['primary', 'secondary', 'danger', 'ghost']
    for (const variant of variants) {
      const classes = useButtonStyle(variant)
      expect(classes).toContain('inline-flex')
      expect(classes).toContain('rounded')
      expect(classes).toContain('transition-all')
      expect(classes).toContain('active:scale-95')
      expect(classes).toContain('disabled:opacity-50')
    }
  })
})
