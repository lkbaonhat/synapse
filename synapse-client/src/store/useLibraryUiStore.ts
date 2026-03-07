// §6 — Pinia store holds ONLY local UI state (modal open/close, filter selection)
// No API calls, no axios, no server data in this store
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLibraryUiStore = defineStore('library-ui', () => {
  // ── Modal state ───────────────────────────────────────────────────────
  const showNewDeckModal = ref(false)
  const showNewFolderModal = ref(false)

  // ── Filter state ──────────────────────────────────────────────────────
  const selectedTagId = ref<string | null>(null)

  // ── Actions ───────────────────────────────────────────────────────────
  function openNewDeck(): void {
    showNewDeckModal.value = true
    showNewFolderModal.value = false
  }

  function openNewFolder(): void {
    showNewFolderModal.value = true
    showNewDeckModal.value = false
  }

  function closeModals(): void {
    showNewDeckModal.value = false
    showNewFolderModal.value = false
  }

  return {
    showNewDeckModal,
    showNewFolderModal,
    selectedTagId,
    openNewDeck,
    openNewFolder,
    closeModals
  }
})
