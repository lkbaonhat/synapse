import { defineStore } from 'pinia'
import { ref } from 'vue'
import { nanoid } from 'nanoid'
import type { Folder } from '@/types'

export const useFolderStore = defineStore('folder', () => {
  const folders = ref<Folder[]>([])

  function addFolder(name: string, parentId?: string): Folder {
    const now = new Date().toISOString()
    const newFolder: Folder = {
      id: nanoid(),
      name,
      parentId,
      createdAt: now,
      updatedAt: now
    }
    folders.value.push(newFolder)
    return newFolder
  }

  function updateFolder(id: string, updates: Partial<Folder>): void {
    const folder = folders.value.find(f => f.id === id)
    if (folder) {
      Object.assign(folder, { ...updates, updatedAt: new Date().toISOString() })
    }
  }

  function deleteFolder(id: string): void {
    folders.value = folders.value.filter(f => f.id !== id)
  }

  return { folders, addFolder, updateFolder, deleteFolder }
})
