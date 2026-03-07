import { z } from 'zod'

// ═══════════════════════════════════════════════
// AUTH TYPES & SCHEMAS
// ═══════════════════════════════════════════════

export const loginSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters')
})

export const registerSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
  confirmPassword: z.string()
}).refine(data => data.password === data.confirmPassword, {
  message: 'Passwords do not match',
  path: ['confirmPassword']
})

export type LoginInput = z.infer<typeof loginSchema>
export type RegisterInput = z.infer<typeof registerSchema>

export interface User {
  id: string
  name: string
  email: string
  avatarUrl?: string
  createdAt: string
}

export interface AuthResponse {
  token: string
  user: User
}

// ═══════════════════════════════════════════════
// DOMAIN TYPES
// ═══════════════════════════════════════════════

export type DifficultyRating = 'again' | 'hard' | 'good' | 'easy'

export interface Tag {
  id: string
  name: string
  color?: string
}

export type QuestionFormat = 'flashcard' | 'cloze' | 'free_response'

export interface Card {
  id: string
  deckId: string
  format: QuestionFormat

  /** Front of card, question text, or cloze sentence */
  front: string
  /** Back of card, answer, or the hidden cloze value */
  back: string

  /** Tag IDs assigned */
  tagIds: string[]

  // SRS (SM-2) fields
  interval: number       // Current interval in days
  repetition: number     // Number of successful repetitions
  easeFactor: number     // SM-2 E-factor (starts at 2.5)
  nextReviewDate: string // ISO date string

  createdAt: string
  updatedAt: string
}

export interface Deck {
  id: string
  folderId?: string
  name: string
  description?: string
  createdAt: string
  updatedAt: string
}

export interface Folder {
  id: string
  name: string
  parentId?: string
  createdAt: string
  updatedAt: string
}
