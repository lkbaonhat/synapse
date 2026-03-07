// §7 — All HTTP calls through services, never direct axios in components
import api from '@/configs/axios'

import type { Card, QuestionFormat } from '@/types'

export interface CreateCardPayload {
  deckId: string
  format: QuestionFormat
  front: string
  back: string
  tagIds?: string[]
}

export type UpdateCardPayload = Partial<Omit<Card, 'id' | 'deckId' | 'createdAt'>>

export const cardService = {
  getAll: (): Promise<Card[]> =>
    api.get<Card[]>('/cards').then(r => r.data),

  getByDeck: (deckId: string): Promise<Card[]> =>
    api.get<Card[]>(`/decks/${deckId}/cards`).then(r => r.data),

  getById: (id: string): Promise<Card> =>
    api.get<Card>(`/cards/${id}`).then(r => r.data),

  create: (payload: CreateCardPayload): Promise<Card> =>
    api.post<Card>('/cards', payload).then(r => r.data),

  update: (id: string, payload: UpdateCardPayload): Promise<Card> =>
    api.patch<Card>(`/cards/${id}`, payload).then(r => r.data),

  remove: (id: string): Promise<void> =>
    api.delete(`/cards/${id}`).then(r => r.data),

  bulkCreate: (cards: CreateCardPayload[]): Promise<Card[]> =>
    api.post<Card[]>('/cards/bulk', { cards }).then(r => r.data)
}
