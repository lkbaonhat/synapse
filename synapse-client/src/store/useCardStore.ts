import { defineStore } from 'pinia'
import { ref } from 'vue'
import { nanoid } from 'nanoid'
import type { Card, QuestionFormat } from '@/types'

export const useCardStore = defineStore('card', () => {
  const cards = ref<Card[]>([])

  function addCard(
    deckId: string,
    format: QuestionFormat,
    front: string,
    back: string,
    tagIds: string[] = []
  ): Card {
    const now = new Date().toISOString()
    const newCard: Card = {
      id: nanoid(),
      deckId,
      format,
      front,
      back,
      tagIds,
      interval: 0,
      repetition: 0,
      easeFactor: 2.5,
      nextReviewDate: now,
      createdAt: now,
      updatedAt: now
    }
    cards.value.push(newCard)
    return newCard
  }

  function updateCard(id: string, updates: Partial<Card>): void {
    const card = cards.value.find(c => c.id === id)
    if (card) {
      Object.assign(card, { ...updates, updatedAt: new Date().toISOString() })
    }
  }

  function deleteCard(id: string): void {
    cards.value = cards.value.filter(c => c.id !== id)
  }

  function getByDeck(deckId: string): Card[] {
    return cards.value.filter(c => c.deckId === deckId)
  }

  return { cards, addCard, updateCard, deleteCard, getByDeck }
})
