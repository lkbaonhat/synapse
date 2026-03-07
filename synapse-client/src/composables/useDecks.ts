// §7 + §8 — Vue Query composable wrapping the deck service layer
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'

import { deckService } from '@/services/deckService'
import type { Deck } from '@/types'

// ── Query ─────────────────────────────────────────────────────────────
export function useDecks() {
  return useQuery({
    queryKey: ['decks'],
    queryFn: deckService.getAll
  })
}

export function useDeck(id: string) {
  return useQuery({
    queryKey: ['decks', id],
    queryFn: () => deckService.getById(id),
    enabled: !!id
  })
}

// ── Mutations ─────────────────────────────────────────────────────────
export function useDeckMutations() {
  const client = useQueryClient()

  const invalidate = () => client.invalidateQueries({ queryKey: ['decks'] })

  const create = useMutation({
    mutationFn: (payload: Omit<Deck, 'id' | 'createdAt' | 'updatedAt'>) =>
      deckService.create(payload),
    onSuccess: invalidate
  })

  const update = useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: Partial<Deck> }) =>
      deckService.update(id, payload),
    onSuccess: invalidate
  })

  const remove = useMutation({
    mutationFn: (id: string) => deckService.remove(id),
    onSuccess: invalidate
  })

  return { create, update, remove }
}
