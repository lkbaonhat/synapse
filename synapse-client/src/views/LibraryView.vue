<template>
  <div class="library-view">
    <!-- Page header -->
    <header class="page-header">
      <div>
        <h1 class="page-title">Your Library</h1>
        <p class="page-subtitle">{{ deckCount }} decks · {{ cardCount }} cards</p>
      </div>
      <div class="header-actions">
        <BaseButton variant="secondary" @click="uiStore.openNewFolder()">+ New Folder</BaseButton>
        <BaseButton variant="primary" @click="uiStore.openNewDeck()">+ New Deck</BaseButton>
      </div>
    </header>

    <!-- Tag filter bar -->
    <!-- §9 — v-if + v-for on separate elements via <template> wrapper -->
    <div v-if="tags && tags.length > 0" class="tag-bar">
      <button
        class="tag-chip"
        :class="{ 'tag-chip--active': uiStore.selectedTagId === null }"
        @click="uiStore.selectedTagId = null"
      >All</button>
      <template v-for="tag in tags" :key="tag.id">
        <button
          class="tag-chip"
          :class="{ 'tag-chip--active': uiStore.selectedTagId === tag.id }"
          :style="{ '--tag-color': tag.color }"
          @click="uiStore.selectedTagId = tag.id"
        >{{ tag.name }}</button>
      </template>
    </div>

    <!-- Folders section -->
    <div v-if="folders && folders.length > 0" class="section">
      <h2 class="section-title">Folders</h2>
      <div class="folder-grid">
        <template v-for="folder in folders" :key="folder.id">
          <div class="folder-card">
            <span class="folder-icon">📁</span>
            <span class="folder-name">{{ folder.name }}</span>
            <span class="folder-count">{{ decksByFolder(folder.id).length }}</span>
            <button class="icon-btn danger" @click="handleDeleteFolder(folder.id)" title="Delete folder">✕</button>
          </div>
        </template>
      </div>
    </div>

    <!-- Decks section -->
    <div class="section">
      <h2 class="section-title">Decks</h2>

      <!-- Loading skeleton -->
      <div v-if="decksLoading" class="deck-grid">
        <div v-for="n in 4" :key="n" class="deck-skeleton"></div>
      </div>

      <!-- Error state — §11 err:unknown pattern -->
      <div v-else-if="decksError" class="error-state">
        <p class="error-icon">⚠️</p>
        <p class="error-text">Failed to load decks. Please try again.</p>
        <BaseButton variant="secondary" @click="refetchDecks()">Retry</BaseButton>
      </div>

      <!-- Deck grid -->
      <div v-else-if="filteredDecks.length > 0" class="deck-grid">
        <template v-for="deck in filteredDecks" :key="deck.id">
          <DeckCard
            :deck="deck"
            :card-count="cardCountByDeck(deck.id)"
            :due-count="dueCountByDeck(deck.id)"
            @study="router.push({ name: 'Study', params: { deckId: deck.id } })"
            @view="router.push({ name: 'DeckDetail', params: { deckId: deck.id } })"
            @delete="handleDeleteDeck(deck.id)"
          />
        </template>
      </div>

      <!-- Empty state -->
      <div v-else class="empty-state">
        <p class="empty-icon">📦</p>
        <p class="empty-text">No decks yet. Create your first deck to get started!</p>
      </div>
    </div>

    <!-- New Deck Modal — §3 BaseModal, §6 modal state in Pinia UI store -->
    <BaseModal :open="uiStore.showNewDeckModal" title="Create New Deck" @close="uiStore.closeModals()">
      <form @submit.prevent="handleCreateDeck" class="modal-form">
        <BaseInput
          v-model="newDeck.name"
          label="Deck Name"
          placeholder="e.g. JavaScript Fundamentals"
          required
        />
        <BaseInput
          v-model="newDeck.description"
          label="Description"
          placeholder="Optional description…"
        />
        <div class="form-group">
          <label class="form-label">Folder</label>
          <select v-model="newDeck.folderId" class="input">
            <option value="">No folder</option>
            <template v-for="f in folders" :key="f.id">
              <option :value="f.id">{{ f.name }}</option>
            </template>
          </select>
        </div>
      </form>
      <template #footer>
        <BaseButton variant="secondary" @click="uiStore.closeModals()">Cancel</BaseButton>
        <BaseButton
          variant="primary"
          :loading="deckMutations.create.isPending.value"
          :disabled="!newDeck.name.trim()"
          @click="handleCreateDeck"
        >Create</BaseButton>
      </template>
    </BaseModal>

    <!-- New Folder Modal -->
    <BaseModal :open="uiStore.showNewFolderModal" title="Create New Folder" @close="uiStore.closeModals()">
      <BaseInput
        v-model="newFolderName"
        label="Folder Name"
        placeholder="e.g. Computer Science"
        required
      />
      <template #footer>
        <BaseButton variant="secondary" @click="uiStore.closeModals()">Cancel</BaseButton>
        <BaseButton
          variant="primary"
          :loading="folderMutations.create.isPending.value"
          :disabled="!newFolderName.trim()"
          @click="handleCreateFolder"
        >Create</BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
// §5 — Import order: vue → router → composables → store → components → types
import { ref, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'

import { useDecks, useDeckMutations } from '@/composables/useDecks'
import { useAllCards } from '@/composables/useCards'
import { useTags } from '@/composables/useTags'
import { useFolders, useFolderMutations } from '@/composables/useFolders'

import { useLibraryUiStore } from '@/store/useLibraryUiStore'

import BaseButton from '@/components/base/BaseButton.vue'
import BaseModal from '@/components/base/BaseModal.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import DeckCard from '@/components/DeckCard.vue'

const router = useRouter()

// ── Server state (Vue Query) ───────────────────────────────────────────
const { data: decks, isLoading: decksLoading, isError: decksError, refetch: refetchDecks } = useDecks()
const { data: allCards } = useAllCards()
const { data: tags } = useTags()
const { data: folders } = useFolders()

const deckMutations = useDeckMutations()
const folderMutations = useFolderMutations()

// ── UI state (Pinia store — §6) ───────────────────────────────────────
const uiStore = useLibraryUiStore()

// ── Local form state ──────────────────────────────────────────────────
const newDeck = reactive({ name: '', description: '', folderId: '' })
const newFolderName = ref('')

// ── Computed (§9) ─────────────────────────────────────────────────────
const deckCount = computed(() => decks.value?.length ?? 0)
const cardCount = computed(() => allCards.value?.length ?? 0)

const filteredDecks = computed(() => {
  if (!decks.value) return []
  if (!uiStore.selectedTagId) return decks.value
  // Filter decks that have at least one card tagged with the selected tag
  const taggedDeckIds = new Set(
    (allCards.value ?? [])
      .filter(c => c.tagIds.includes(uiStore.selectedTagId!))
      .map(c => c.deckId)
  )
  return decks.value.filter(d => taggedDeckIds.has(d.id))
})

function decksByFolder(folderId: string) {
  return (decks.value ?? []).filter(d => d.folderId === folderId)
}

function cardCountByDeck(deckId: string) {
  return (allCards.value ?? []).filter(c => c.deckId === deckId).length
}

function dueCountByDeck(deckId: string) {
  const now = new Date()
  return (allCards.value ?? []).filter(
    c => c.deckId === deckId && new Date(c.nextReviewDate) <= now
  ).length
}

// ── Actions ───────────────────────────────────────────────────────────
function handleCreateDeck() {
  if (!newDeck.name.trim()) return
  deckMutations.create.mutate(
    {
      name: newDeck.name.trim(),
      description: newDeck.description || undefined,
      folderId: newDeck.folderId || undefined
    },
    {
      onSuccess: () => {
        newDeck.name = ''
        newDeck.description = ''
        newDeck.folderId = ''
        uiStore.closeModals()
      }
    }
  )
}

function handleDeleteDeck(id: string) {
  deckMutations.remove.mutate(id)
}

function handleCreateFolder() {
  if (!newFolderName.value.trim()) return
  folderMutations.create.mutate(
    { name: newFolderName.value.trim() },
    {
      onSuccess: () => {
        newFolderName.value = ''
        uiStore.closeModals()
      }
    }
  )
}

function handleDeleteFolder(id: string) {
  folderMutations.remove.mutate(id)
}
</script>

<style scoped>
/* §10 — CSS tokens only */
.library-view { max-width: 1100px; margin: 0 auto; }

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
}

.page-title { font-size: 1.8rem; font-weight: 700; color: var(--color-text); margin: 0 0 0.25rem; }
.page-subtitle { font-size: 0.9rem; color: var(--color-text-muted); margin: 0; }

.header-actions { display: flex; gap: 0.75rem; }

.tag-bar { display: flex; gap: 0.5rem; flex-wrap: wrap; margin-bottom: 1.5rem; }

.tag-chip {
  padding: 0.3rem 0.9rem;
  border-radius: 999px;
  border: 1px solid var(--color-border);
  background: transparent;
  color: var(--color-text-muted);
  font-size: 0.8rem;
  cursor: pointer;
  transition: all 0.2s;
}

.tag-chip:hover { border-color: var(--color-primary); color: var(--color-primary); }
.tag-chip--active { background: var(--tag-color, var(--color-primary)); border-color: transparent; color: white; }

.section { margin-bottom: 2rem; }
.section-title {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0 0 1rem;
}

.folder-grid { display: flex; flex-wrap: wrap; gap: 0.75rem; margin-bottom: 1.5rem; }
.folder-card {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1rem;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.9rem;
}
.folder-card:hover { border-color: var(--color-primary); }
.folder-name { font-weight: 500; color: var(--color-text); }
.folder-count { color: var(--color-text-muted); font-size: 0.8rem; }

.deck-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}

/* Loading skeleton */
.deck-skeleton {
  height: 180px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  animation: shimmer 1.4s ease infinite;
}

@keyframes shimmer {
  0%, 100% { opacity: 0.7; }
  50% { opacity: 0.3; }
}

/* Error state */
.error-state {
  text-align: center;
  padding: 3rem 2rem;
  background: rgba(239,68,68,0.05);
  border: 1px solid rgba(239,68,68,0.2);
  border-radius: var(--radius-lg);
}
.error-icon { font-size: 2.5rem; margin: 0 0 0.5rem; }
.error-text { color: var(--color-text-muted); margin: 0 0 1rem; }

/* Empty state */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  background: var(--color-surface);
  border: 2px dashed var(--color-border);
  border-radius: var(--radius-lg);
}
.empty-icon { font-size: 3rem; margin: 0 0 0.5rem; }
.empty-text { color: var(--color-text-muted); font-size: 0.95rem; }

/* Modal form */
.modal-form { display: flex; flex-direction: column; gap: 1rem; }
.form-group { display: flex; flex-direction: column; gap: 0.35rem; }
.form-label { font-size: 0.85rem; font-weight: 500; color: var(--color-text-muted); }

/* Icon buttons */
.icon-btn { background: none; border: none; cursor: pointer; padding: 0.25rem 0.5rem; border-radius: var(--radius-sm); opacity: 0.4; transition: opacity 0.2s; }
.icon-btn:hover { opacity: 1; }
.icon-btn.danger:hover { color: #ef4444; }
</style>
