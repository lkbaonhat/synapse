// §8 — Pure SRS composable: no imports from Vue, no side effects, fully unit-testable
// Implements SM-2 algorithm as used by Anki
import type { DifficultyRating } from '@/types'

export interface SRSCard {
  interval: number
  repetition: number
  easeFactor: number
}

export interface SRSResult extends SRSCard {
  nextReviewDate: string
}

/**
 * Quality mapping from human-friendly ratings to SM-2 q values (0–5)
 */
const QUALITY_MAP: Record<DifficultyRating, number> = {
  again: 0, // complete blackout
  hard: 3,  // correct but with serious difficulty
  good: 4,  // correct with effort
  easy: 5   // perfect recall
}

/**
 * Compute the next SM-2 schedule for a card given a difficulty rating.
 * Pure function — no DB access, no Vue reactivity.
 *
 * @param card - Current card schedule values
 * @param rating - User's difficulty rating
 * @returns Updated SRS fields and next review ISO date
 */
export function computeNextSRS(card: SRSCard, rating: DifficultyRating): SRSResult {
  const q = QUALITY_MAP[rating]

  // Update E-factor (ease factor), clamped to a minimum of 1.3
  let newEaseFactor = card.easeFactor + (0.1 - (5 - q) * (0.08 + (5 - q) * 0.02))
  if (newEaseFactor < 1.3) newEaseFactor = 1.3

  let newInterval: number
  let newRepetition: number

  if (q < 3) {
    // Incorrect response — reset
    newInterval = 0
    newRepetition = 0
  } else {
    newRepetition = card.repetition + 1
    if (card.repetition === 0) {
      newInterval = 1
    } else if (card.repetition === 1) {
      newInterval = 6
    } else {
      newInterval = Math.round(card.interval * newEaseFactor)
    }
  }

  const nextReview = new Date()
  if (newInterval === 0) {
    // "Again" — re-enter queue in 10 minutes
    nextReview.setMinutes(nextReview.getMinutes() + 10)
  } else {
    nextReview.setDate(nextReview.getDate() + newInterval)
  }

  return {
    interval: newInterval,
    repetition: newRepetition,
    easeFactor: newEaseFactor,
    nextReviewDate: nextReview.toISOString()
  }
}

/**
 * Preview the resulting interval string for each rating (for display in study UI).
 * Does NOT modify the card — purely visual.
 *
 * @returns Map from rating to human-readable interval string e.g. "6d", "10m"
 */
export function previewIntervals(card: SRSCard): Record<DifficultyRating, string> {
  const fmt = (result: SRSResult): string => {
    if (result.interval === 0) return '10m'
    if (result.interval === 1) return '1d'
    return `${result.interval}d`
  }

  return {
    again: fmt(computeNextSRS(card, 'again')),
    hard: fmt(computeNextSRS(card, 'hard')),
    good: fmt(computeNextSRS(card, 'good')),
    easy: fmt(computeNextSRS(card, 'easy'))
  }
}
