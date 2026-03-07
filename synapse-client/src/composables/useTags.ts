// §7 + §8 — Vue Query composable wrapping the tag service layer
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'

import { tagService } from '@/services/tagService'
import type { Tag } from '@/types'

// ── Query ─────────────────────────────────────────────────────────────
export function useTags() {
  return useQuery({
    queryKey: ['tags'],
    queryFn: tagService.getAll
  })
}

// ── Mutations ─────────────────────────────────────────────────────────
export function useTagMutations() {
  const client = useQueryClient()

  const invalidate = () => client.invalidateQueries({ queryKey: ['tags'] })

  const create = useMutation({
    mutationFn: (payload: Omit<Tag, 'id'>) => tagService.create(payload),
    onSuccess: invalidate
  })

  const update = useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: Partial<Omit<Tag, 'id'>> }) =>
      tagService.update(id, payload),
    onSuccess: invalidate
  })

  const remove = useMutation({
    mutationFn: (id: string) => tagService.remove(id),
    onSuccess: invalidate
  })

  return { create, update, remove }
}
