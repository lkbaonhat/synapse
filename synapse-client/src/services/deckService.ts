// §7 — All HTTP calls through services, never direct axios in components
import api from '@/configs/axios'

import type { Deck } from '@/types'

type CreateDeckPayload = Omit<Deck, 'id' | 'createdAt' | 'updatedAt'>
type UpdateDeckPayload = Partial<CreateDeckPayload>

export const deckService = {
  getAll: (): Promise<Deck[]> =>
    api.get<Deck[]>('/decks').then(r => r.data),

  getById: (id: string): Promise<Deck> =>
    api.get<Deck>(`/decks/${id}`).then(r => r.data),

  create: (payload: CreateDeckPayload): Promise<Deck> =>
    api.post<Deck>('/decks', payload).then(r => r.data),

  update: (id: string, payload: UpdateDeckPayload): Promise<Deck> =>
    api.patch<Deck>(`/decks/${id}`, payload).then(r => r.data),

  remove: (id: string): Promise<void> =>
    api.delete(`/decks/${id}`).then(r => r.data)
}
