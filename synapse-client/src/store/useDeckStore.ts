import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { nanoid } from 'nanoid'
import type { Deck } from '@/types'

export const useDeckStore = defineStore('deck', () => {
  const decks = ref<Deck[]>([])

  const totalCount = computed(() => decks.value.length)

  function addDeck(name: string, description?: string, folderId?: string): Deck {
    const now = new Date().toISOString()
    const newDeck: Deck = {
      id: nanoid(),
      name,
      description,
      folderId,
      createdAt: now,
      updatedAt: now
    }
    decks.value.push(newDeck)
    return newDeck
  }

  function updateDeck(id: string, updates: Partial<Deck>): void {
    const deck = decks.value.find(d => d.id === id)
    if (deck) {
      Object.assign(deck, { ...updates, updatedAt: new Date().toISOString() })
    }
  }

  function deleteDeck(id: string): void {
    decks.value = decks.value.filter(d => d.id !== id)
  }

  function getById(id: string): Deck | undefined {
    return decks.value.find(d => d.id === id)
  }

  return { decks, totalCount, addDeck, updateDeck, deleteDeck, getById }
})
