// §6 — Pinia store holds ONLY local UI state for DeckDetailView
// No API calls, no server data
import { defineStore } from 'pinia'
import { ref } from 'vue'

import type { QuestionFormat } from '@/types'

export const useDeckDetailUiStore = defineStore('deck-detail-ui', () => {
  // ── Modal state ────────────────────────────────────────────────────
  const showCardModal = ref(false)
  const editingCardId = ref<string | null>(null)

  // ── Filter state ───────────────────────────────────────────────────
  const selectedFormat = ref<QuestionFormat | 'all'>('all')

  // ── Actions ────────────────────────────────────────────────────────
  function openAddCard(): void {
    editingCardId.value = null
    showCardModal.value = true
  }

  function openEditCard(cardId: string): void {
    editingCardId.value = cardId
    showCardModal.value = true
  }

  function closeCardModal(): void {
    showCardModal.value = false
    editingCardId.value = null
  }

  return {
    showCardModal,
    editingCardId,
    selectedFormat,
    openAddCard,
    openEditCard,
    closeCardModal
  }
})
