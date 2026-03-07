// §7 + §8 — Vue Query composable wrapping the card service layer
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import type { Ref } from 'vue'

import { cardService } from '@/services/cardService'
import type { CreateCardPayload, UpdateCardPayload } from '@/services/cardService'

// ── Query ─────────────────────────────────────────────────────────────
export function useAllCards() {
  return useQuery({
    queryKey: ['cards'],
    queryFn: cardService.getAll
  })
}

export function useCards(deckId: string | Ref<string>) {
  return useQuery({
    queryKey: ['cards', deckId],
    queryFn: () => {
      const id = typeof deckId === 'string' ? deckId : deckId.value
      return cardService.getByDeck(id)
    },
    enabled: !!deckId
  })
}

// ── Mutations ─────────────────────────────────────────────────────────
export function useCardMutations(deckId: string) {
  const client = useQueryClient()

  const invalidate = () =>
    client.invalidateQueries({ queryKey: ['cards', deckId] })

  const create = useMutation({
    mutationFn: (payload: CreateCardPayload) => cardService.create(payload),
    onSuccess: invalidate
  })

  const update = useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateCardPayload }) =>
      cardService.update(id, payload),
    onSuccess: invalidate
  })

  const remove = useMutation({
    mutationFn: (id: string) => cardService.remove(id),
    onSuccess: invalidate
  })

  const bulkCreate = useMutation({
    mutationFn: (cards: CreateCardPayload[]) => cardService.bulkCreate(cards),
    onSuccess: invalidate
  })

  return { create, update, remove, bulkCreate }
}
