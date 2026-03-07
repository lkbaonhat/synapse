// §7 — All HTTP calls through services, never direct axios in components
import api from '@/configs/axios'

import type { Tag } from '@/types'

type CreateTagPayload = Omit<Tag, 'id'>
type UpdateTagPayload = Partial<CreateTagPayload>

export const tagService = {
  getAll: (): Promise<Tag[]> =>
    api.get<Tag[]>('/tags').then(r => r.data),

  create: (payload: CreateTagPayload): Promise<Tag> =>
    api.post<Tag>('/tags', payload).then(r => r.data),

  update: (id: string, payload: UpdateTagPayload): Promise<Tag> =>
    api.patch<Tag>(`/tags/${id}`, payload).then(r => r.data),

  remove: (id: string): Promise<void> =>
    api.delete(`/tags/${id}`).then(r => r.data)
}
