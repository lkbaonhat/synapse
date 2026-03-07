<template>
  <div class="deck-detail-view">
    <!-- Header / breadcrumb -->
    <div class="page-header">
      <div class="breadcrumb">
        <RouterLink to="/app/library" class="back-link">← Library</RouterLink>
        <span class="breadcrumb-sep">/</span>
        <span class="breadcrumb-current">{{ deck?.name ?? '…' }}</span>
      </div>
      <div class="header-top">
        <h1 class="page-title">{{ deck?.name }}</h1>
        <div class="header-actions">
          <BaseButton variant="secondary" @click="uiStore.openAddCard()">+ Add Card</BaseButton>
          <BaseButton
            variant="primary"
            :disabled="dueCount === 0"
            @click="router.push({ name: 'Study', params: { deckId } })"
          >
            Study ({{ dueCount }})
          </BaseButton>
        </div>
      </div>
      <p v-if="deck?.description" class="deck-desc">{{ deck.description }}</p>
    </div>

    <!-- Format filter tabs -->
    <!-- §9 — filteredCards computed, no logic in template -->
    <div class="format-tabs">
      <button
        v-for="tab in formatTabs"
        :key="tab.value"
        class="format-tab"
        :class="{ 'format-tab--active': uiStore.selectedFormat === tab.value }"
        @click="uiStore.selectedFormat = tab.value"
      >
        {{ tab.label }}
        <span class="tab-count">{{ tab.count }}</span>
      </button>
    </div>

    <!-- Card list -->
    <div v-if="cardsLoading" class="card-list">
      <div v-for="n in 3" :key="n" class="card-item-skeleton"></div>
    </div>

    <div v-else-if="cardsError" class="error-state">
      <p>Failed to load cards.</p>
      <BaseButton variant="secondary" @click="refetchCards()">Retry</BaseButton>
    </div>

    <div v-else class="card-list">
      <template v-for="card in filteredCards" :key="card.id">
        <div class="card-item">
          <div class="card-item-meta">
            <BaseBadge :format="card.format" />
            <span class="card-interval">next in {{ nextReviewLabel(card) }}</span>
          </div>
          <p class="card-front">{{ card.front }}</p>
          <div class="card-item-actions">
            <BaseButton variant="ghost" size="sm" @click="startEdit(card)">Edit</BaseButton>
            <BaseButton variant="danger" size="sm" @click="handleDeleteCard(card.id)">Delete</BaseButton>
          </div>
        </div>
      </template>

      <div v-if="filteredCards.length === 0 && !cardsLoading" class="empty-state">
        <p class="empty-icon">🃏</p>
        <p class="empty-text">No cards yet. Add your first card!</p>
      </div>
    </div>

    <!-- Add / Edit Card Modal -->
    <BaseModal
      :open="uiStore.showCardModal"
      :title="uiStore.editingCardId ? 'Edit Card' : 'Add Card'"
      @close="uiStore.closeCardModal()"
    >
      <form @submit.prevent="handleSaveCard" class="card-form">
        <!-- Format selector -->
        <div class="form-group">
          <label class="form-label">Format</label>
          <select v-model="cardForm.format" class="input" :disabled="!!uiStore.editingCardId">
            <option value="flashcard">🃏 Flashcard</option>
            <option value="cloze">✏️ Cloze</option>
            <option value="free_response">💬 Free Response</option>
          </select>
        </div>

        <!-- Front / Question -->
        <BaseInput
          v-model="cardForm.front"
          :label="frontLabel"
          :placeholder="frontPlaceholder"
          :error="formErrors.front"
          required
        />

        <!-- Back / Answer (not shown for pure cloze) -->
        <BaseInput
          v-if="cardForm.format !== 'cloze'"
          v-model="cardForm.back"
          :label="backLabel"
          :placeholder="backPlaceholder"
          :error="formErrors.back"
          required
        />

        <!-- Tag selector -->
        <div class="form-group" v-if="tags && tags.length > 0">
          <label class="form-label">Tags</label>
          <div class="tag-checkboxes">
            <template v-for="tag in tags" :key="tag.id">
              <label class="tag-checkbox-row">
                <input
                  type="checkbox"
                  :value="tag.id"
                  v-model="cardForm.tagIds"
                />
                <span class="tag-dot" :style="{ background: tag.color }"></span>
                {{ tag.name }}
              </label>
            </template>
          </div>
        </div>
      </form>

      <template #footer>
        <BaseButton variant="secondary" @click="uiStore.closeCardModal()">Cancel</BaseButton>
        <BaseButton
          variant="primary"
          :loading="cardMutations.create.isPending.value || cardMutations.update.isPending.value"
          :disabled="!cardForm.front.trim()"
          @click="handleSaveCard"
        >
          {{ uiStore.editingCardId ? 'Save Changes' : 'Add Card' }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
// §5 — Import order: vue → router → composables → store → components → types
import { computed, reactive, watch } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import dayjs from 'dayjs'

import { useCards, useCardMutations } from '@/composables/useCards'
import { useDecks } from '@/composables/useDecks'
import { useTags } from '@/composables/useTags'

import { useDeckDetailUiStore } from '@/store/useDeckDetailUiStore'

import BaseButton from '@/components/base/BaseButton.vue'
import BaseModal from '@/components/base/BaseModal.vue'
import BaseInput from '@/components/base/BaseInput.vue'
import BaseBadge from '@/components/base/BaseBadge.vue'

import type { Card, QuestionFormat } from '@/types'

const route = useRoute()
const router = useRouter()
const deckId = route.params.deckId as string

// ── Server state (Vue Query) ───────────────────────────────────────────
const { data: cards, isLoading: cardsLoading, isError: cardsError, refetch: refetchCards } = useCards(deckId)
const { data: decks } = useDecks()
const { data: tags } = useTags()

const cardMutations = useCardMutations(deckId)

// ── UI state (Pinia — §6) ─────────────────────────────────────────────
const uiStore = useDeckDetailUiStore()

// ── Computed (§9) ─────────────────────────────────────────────────────
const deck = computed(() => decks.value?.find(d => d.id === deckId))

const now = new Date()
const dueCount = computed(() =>
  (cards.value ?? []).filter(c => new Date(c.nextReviewDate) <= now).length
)

// Format filter tabs with counts — §9: complex logic in computed
const formatTabs = computed(() => {
  const all = cards.value ?? []
  return [
    { label: 'All', value: 'all' as const, count: all.length },
    { label: '🃏 Flashcard', value: 'flashcard' as const, count: all.filter(c => c.format === 'flashcard').length },
    { label: '✏️ Cloze', value: 'cloze' as const, count: all.filter(c => c.format === 'cloze').length },
    { label: '💬 Free Response', value: 'free_response' as const, count: all.filter(c => c.format === 'free_response').length }
  ]
})

const filteredCards = computed(() => {
  const all = cards.value ?? []
  if (uiStore.selectedFormat === 'all') return all
  return all.filter(c => c.format === uiStore.selectedFormat)
})

// Dynamic labels based on format — §9: computed not template ternaries
const frontLabel = computed(() => {
  const map: Record<string, string> = {
    flashcard: 'Question (Front)',
    cloze: 'Cloze sentence — wrap blanks with {{word}}',
    free_response: 'Question'
  }
  return map[cardForm.format] ?? 'Front'
})

const frontPlaceholder = computed(() => {
  const map: Record<string, string> = {
    flashcard: 'e.g. What is a closure?',
    cloze: 'e.g. A {{closure}} has access to its outer scope.',
    free_response: 'e.g. Explain the difference between == and ===.'
  }
  return map[cardForm.format] ?? ''
})

const backLabel = computed(() => (cardForm.format === 'flashcard' ? 'Answer (Back)' : 'Model Answer'))
const backPlaceholder = computed(() => 'Enter the answer…')

// ── Form state ────────────────────────────────────────────────────────
const cardForm = reactive({
  format: 'flashcard' as QuestionFormat,
  front: '',
  back: '',
  tagIds: [] as string[]
})

const formErrors = reactive({ front: '', back: '' })

function resetForm() {
  cardForm.format = 'flashcard'
  cardForm.front = ''
  cardForm.back = ''
  cardForm.tagIds = []
  formErrors.front = ''
  formErrors.back = ''
}

function startEdit(card: Card) {
  cardForm.format = card.format
  cardForm.front = card.front
  cardForm.back = card.back
  cardForm.tagIds = [...card.tagIds]
  uiStore.openEditCard(card.id)
}

// ── Actions ───────────────────────────────────────────────────────────
function handleSaveCard() {
  // Client-side validation (server validates too)
  formErrors.front = cardForm.front.trim() ? '' : 'This field is required.'
  if (cardForm.format !== 'cloze') {
    formErrors.back = cardForm.back.trim() ? '' : 'This field is required.'
  }
  if (formErrors.front || formErrors.back) return

  if (uiStore.editingCardId) {
    cardMutations.update.mutate(
      { id: uiStore.editingCardId, payload: { front: cardForm.front, back: cardForm.back, tagIds: cardForm.tagIds } },
      { onSuccess: () => { resetForm(); uiStore.closeCardModal() } }
    )
  } else {
    cardMutations.create.mutate(
      { deckId, format: cardForm.format, front: cardForm.front, back: cardForm.back, tagIds: cardForm.tagIds },
      { onSuccess: () => { resetForm(); uiStore.closeCardModal() } }
    )
  }
}

function handleDeleteCard(id: string) {
  cardMutations.remove.mutate(id)
}

// ── Helpers ───────────────────────────────────────────────────────────
function nextReviewLabel(card: Card): string {
  const diff = dayjs(card.nextReviewDate).diff(dayjs(), 'day')
  if (diff <= 0) return 'now'
  if (diff === 1) return '1 day'
  return `${diff} days`
}

// Reset form when modal closes
watch(() => uiStore.showCardModal, open => {
  if (!open) resetForm()
})
</script>

<style scoped>
/* §10 — CSS tokens only */
.deck-detail-view { max-width: 900px; margin: 0 auto; }

.page-header { margin-bottom: 1.5rem; }

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
  color: var(--color-text-muted);
  margin-bottom: 0.75rem;
}

.back-link { color: var(--color-text-muted); text-decoration: none; transition: color 0.2s; }
.back-link:hover { color: var(--color-primary); }
.breadcrumb-sep { opacity: 0.4; }
.breadcrumb-current { color: var(--color-text); font-weight: 500; }

.header-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}

.page-title { font-size: 1.6rem; font-weight: 700; color: var(--color-text); margin: 0; }
.deck-desc { font-size: 0.875rem; color: var(--color-text-muted); margin: 0.4rem 0 0; }
.header-actions { display: flex; gap: 0.75rem; flex-shrink: 0; }

/* Format tabs */
.format-tabs {
  display: flex;
  gap: 0.25rem;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid var(--color-border);
  padding-bottom: 0.5rem;
}

.format-tab {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.45rem 1rem;
  border: none;
  background: none;
  color: var(--color-text-muted);
  font-size: 0.875rem;
  cursor: pointer;
  border-radius: var(--radius-sm);
  transition: all 0.2s;
}

.format-tab:hover { background: var(--color-surface); color: var(--color-text); }
.format-tab--active { background: var(--color-primary-subtle); color: var(--color-primary); font-weight: 600; }

.tab-count {
  background: var(--color-surface-hover);
  padding: 0.1rem 0.45rem;
  border-radius: 999px;
  font-size: 0.72rem;
  color: var(--color-text-muted);
}

/* Card list */
.card-list { display: flex; flex-direction: column; gap: 0.75rem; }

.card-item {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: 1rem 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  transition: border-color 0.2s;
}

.card-item:hover { border-color: rgba(255,255,255,0.15); }

.card-item-meta {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.card-interval { font-size: 0.75rem; color: var(--color-text-muted); }

.card-front {
  font-size: 0.95rem;
  color: var(--color-text);
  margin: 0;
  line-height: 1.5;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
}

.card-item-actions {
  display: flex;
  gap: 0.5rem;
  align-self: flex-end;
}

/* Skeleton */
.card-item-skeleton {
  height: 80px;
  background: var(--color-surface);
  border-radius: var(--radius-md);
  animation: shimmer 1.4s ease infinite;
}

@keyframes shimmer { 0%, 100% { opacity: 0.7; } 50% { opacity: 0.3; } }

/* Error / empty state */
.error-state, .empty-state {
  text-align: center;
  padding: 3rem 2rem;
  background: var(--color-surface);
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-lg);
}
.empty-icon { font-size: 3rem; margin: 0 0 0.5rem; }
.empty-text { color: var(--color-text-muted); }

/* Card form */
.card-form { display: flex; flex-direction: column; gap: 1rem; }
.form-group { display: flex; flex-direction: column; gap: 0.35rem; }
.form-label { font-size: 0.85rem; font-weight: 500; color: var(--color-text-muted); }

/* Tag checkboxes */
.tag-checkboxes { display: flex; flex-direction: column; gap: 0.5rem; }
.tag-checkbox-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  color: var(--color-text);
  cursor: pointer;
}
.tag-dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
</style>
