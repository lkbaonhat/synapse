/**
 * TDD: Tests for the useSRS composable (SM-2 algorithm)
 *
 * RED phase: these tests are written BEFORE the implementation.
 * They define the contract that computeNextSRS() must satisfy.
 * Run: pnpm test src/composables/__tests__/useSRS.test.ts
 */
import { describe, it, expect } from 'vitest'
import { computeNextSRS, previewIntervals } from '@/composables/useSRS'
import type { SRSCard } from '@/composables/useSRS'

// A freshly created card (first-time: never reviewed)
const freshCard: SRSCard = {
  interval: 0,
  repetition: 0,
  easeFactor: 2.5
}

describe('computeNextSRS — difficulty: again', () => {
  it('resets repetition to 0', () => {
    const result = computeNextSRS({ ...freshCard, repetition: 3, interval: 10 }, 'again')
    expect(result.repetition).toBe(0)
  })

  it('resets interval to 0 (10-minute requeue)', () => {
    const result = computeNextSRS({ ...freshCard, repetition: 3, interval: 10 }, 'again')
    expect(result.interval).toBe(0)
  })

  it('sets nextReviewDate roughly 10 minutes in the future', () => {
    const before = Date.now()
    const result = computeNextSRS(freshCard, 'again')
    const nextMs = new Date(result.nextReviewDate).getTime()
    const eightMin = 8 * 60 * 1000
    const twelveMin = 12 * 60 * 1000
    expect(nextMs).toBeGreaterThan(before + eightMin)
    expect(nextMs).toBeLessThan(before + twelveMin)
  })
})

describe('computeNextSRS — difficulty: hard', () => {
  it('increments repetition', () => {
    const result = computeNextSRS({ ...freshCard, repetition: 2, interval: 6 }, 'hard')
    expect(result.repetition).toBe(3)
  })

  it('decreases easeFactor by ~0.15', () => {
    const result = computeNextSRS(freshCard, 'hard')
    expect(result.easeFactor).toBeLessThan(freshCard.easeFactor)
  })

  it('does not let easeFactor drop below 1.3', () => {
    const lowEF: SRSCard = { interval: 1, repetition: 5, easeFactor: 1.31 }
    const result = computeNextSRS(lowEF, 'hard')
    expect(result.easeFactor).toBeGreaterThanOrEqual(1.3)
  })
})

describe('computeNextSRS — difficulty: good', () => {
  it('sets interval to 1 on first repetition (rep=0)', () => {
    const result = computeNextSRS(freshCard, 'good')
    expect(result.interval).toBe(1)
    expect(result.repetition).toBe(1)
  })

  it('sets interval to 6 on second repetition (rep=1)', () => {
    const result = computeNextSRS({ ...freshCard, repetition: 1, interval: 1 }, 'good')
    expect(result.interval).toBe(6)
    expect(result.repetition).toBe(2)
  })

  it('applies SM-2 formula on subsequent repetitions (rep≥2)', () => {
    const card: SRSCard = { interval: 6, repetition: 2, easeFactor: 2.5 }
    const result = computeNextSRS(card, 'good')
    // Good maps to quality=4; new EF = 2.5 + (0.1 - 1*0.08 + 1*0.02) = 2.5 + 0.1 - 0.1 = 2.5 (approx)
    expect(result.interval).toBe(Math.round(6 * result.easeFactor))
    expect(result.repetition).toBe(3)
  })
})

describe('computeNextSRS — difficulty: easy', () => {
  it('increments repetition', () => {
    const result = computeNextSRS({ ...freshCard, repetition: 2, interval: 6 }, 'easy')
    expect(result.repetition).toBe(3)
  })

  it('increases easeFactor', () => {
    const result = computeNextSRS(freshCard, 'easy')
    expect(result.easeFactor).toBeGreaterThan(freshCard.easeFactor)
  })

  it('sets nextReviewDate 1+ day in the future', () => {
    const result = computeNextSRS({ ...freshCard, repetition: 0 }, 'easy')
    const now = Date.now()
    const nextMs = new Date(result.nextReviewDate).getTime()
    expect(nextMs).toBeGreaterThan(now + 23 * 60 * 60 * 1000) // at least 23h from now
  })
})

describe('previewIntervals', () => {
  it('returns a string for each rating', () => {
    const previews = previewIntervals(freshCard)
    expect(previews).toHaveProperty('again')
    expect(previews).toHaveProperty('hard')
    expect(previews).toHaveProperty('good')
    expect(previews).toHaveProperty('easy')
  })

  it('returns "10m" for again (re-queue in 10 min)', () => {
    const previews = previewIntervals(freshCard)
    expect(previews.again).toBe('10m')
  })

  it('returns "1d" for good on first review', () => {
    const previews = previewIntervals(freshCard) // rep=0
    expect(previews.good).toBe('1d')
  })
})
