// §7 + §8 — Vue Query composable wrapping the folder service layer
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'

import { folderService } from '@/services/folderService'
import type { Folder } from '@/types'

// ── Query ─────────────────────────────────────────────────────────────
export function useFolders() {
  return useQuery({
    queryKey: ['folders'],
    queryFn: folderService.getAll
  })
}

// ── Mutations ─────────────────────────────────────────────────────────
export function useFolderMutations() {
  const client = useQueryClient()

  const invalidate = () => client.invalidateQueries({ queryKey: ['folders'] })

  const create = useMutation({
    mutationFn: (payload: Omit<Folder, 'id' | 'createdAt' | 'updatedAt'>) =>
      folderService.create(payload),
    onSuccess: invalidate
  })

  const update = useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: Partial<Folder> }) =>
      folderService.update(id, payload),
    onSuccess: invalidate
  })

  const remove = useMutation({
    mutationFn: (id: string) => folderService.remove(id),
    onSuccess: invalidate
  })

  return { create, update, remove }
}
