import { defineStore } from 'pinia'
import { ref } from 'vue'
import { nanoid } from 'nanoid'
import type { Tag } from '@/types'

export const useTagStore = defineStore('tag', () => {
  const tags = ref<Tag[]>([])

  function addTag(name: string, color?: string): Tag {
    const newTag: Tag = {
      id: nanoid(),
      name,
      color: color ?? '#808080'  // Use ?? (nullish coalescing) not ||
    }
    tags.value.push(newTag)
    return newTag
  }

  function updateTag(id: string, updates: Partial<Tag>): void {
    const tag = tags.value.find(t => t.id === id)
    if (tag) {
      Object.assign(tag, updates)
    }
  }

  function deleteTag(id: string): void {
    tags.value = tags.value.filter(t => t.id !== id)
  }

  return { tags, addTag, updateTag, deleteTag }
})
