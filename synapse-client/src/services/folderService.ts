// §7 — All HTTP calls through services, never direct axios in components
import api from '@/configs/axios'

import type { Folder } from '@/types'

type CreateFolderPayload = Omit<Folder, 'id' | 'createdAt' | 'updatedAt'>
type UpdateFolderPayload = Partial<CreateFolderPayload>

export const folderService = {
  getAll: (): Promise<Folder[]> =>
    api.get<Folder[]>('/folders').then(r => r.data),

  create: (payload: CreateFolderPayload): Promise<Folder> =>
    api.post<Folder>('/folders', payload).then(r => r.data),

  update: (id: string, payload: UpdateFolderPayload): Promise<Folder> =>
    api.patch<Folder>(`/folders/${id}`, payload).then(r => r.data),

  remove: (id: string): Promise<void> =>
    api.delete(`/folders/${id}`).then(r => r.data)
}
