<template>
  <div class="study-view">
    <!-- No cards due -->
    <div v-if="!currentCard" class="finished-state">
      <p class="finished-icon">🎉</p>
      <h2>Session Complete!</h2>
      <p class="finished-text">No cards are due right now. Come back later or use Cram mode to review all cards.</p>
      <div class="finished-actions">
        <RouterLink to="/app/library" class="btn btn-secondary">Back to Library</RouterLink>
        <button class="btn btn-primary" @click="startCram">Cram All Cards</button>
        <button class="btn btn-primary bg-[var(--color-primary-light)] border-[var(--color-primary)]" @click="startQuiz">Take as Quiz</button>
      </div>
    </div>

    <!-- Active Study Session -->
    <div v-else class="session-wrapper">
      <!-- Progress bar -->
      <div class="session-header">
        <RouterLink to="/app/library" class="back-link">← Exit</RouterLink>
        <div class="progress-wrapper">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
          </div>
          <span class="progress-text">{{ reviewedCount }}/{{ queue.length }}</span>
        </div>
        <div class="mode-badge">{{ quizSessionId ? '📝 Quiz' : cramMode ? '🔥 Cram' : '📅 Review' }}</div>
      </div>

      <!-- Multiple Choice -->
      <div v-if="currentCard?.format === 'multiple_choice'" class="card-area">
        <div class="study-card">
          <div class="card-format-label">📝 Quiz Test</div>
          <MultipleChoiceViewer :card="currentCard" @answer="rate" />
        </div>
      </div>

      <!-- Flashcard / Cloze / Free Response -->
      <div v-else class="card-area">
        <div class="study-card" :class="{ 'is-flipped': isFlipped }" @click="reveal">
          <div class="card-face card-front-face">
            <div class="card-format-label">{{ formatLabel(currentCard.format) }}</div>
            <div class="card-text" v-html="renderContent(currentCard.front, currentCard.format)"></div>
            <div class="click-hint" v-if="!isFlipped">Click to reveal answer</div>
          </div>
          <div class="card-face card-back-face">
            <div class="card-text">{{ currentCard.back }}</div>
          </div>
        </div>
      </div>

      <!-- Rating Buttons (shown after reveal) -->
      <Transition name="slide-up">
        <div v-if="isFlipped && currentCard?.format !== 'multiple_choice'" class="rating-panel">
          <p class="rating-prompt">How well did you remember?</p>
          <div class="rating-buttons">
            <button class="rating-btn rating-again" @click="rate('again')">
              <span class="rating-icon">🔁</span>
              <span class="rating-label">Again</span>
              <span class="rating-interval">&lt;1m</span>
            </button>
            <button class="rating-btn rating-hard" @click="rate('hard')">
              <span class="rating-icon">😓</span>
              <span class="rating-label">Hard</span>
              <span class="rating-interval">{{ nextIntervals.hard }}</span>
            </button>
            <button class="rating-btn rating-good" @click="rate('good')">
              <span class="rating-icon">😊</span>
              <span class="rating-label">Good</span>
              <span class="rating-interval">{{ nextIntervals.good }}</span>
            </button>
            <button class="rating-btn rating-easy" @click="rate('easy')">
              <span class="rating-icon">🚀</span>
              <span class="rating-label">Easy</span>
              <span class="rating-interval">{{ nextIntervals.easy }}</span>
            </button>
          </div>
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'

import { useCards, useCardMutations } from '@/composables/useCards'
import { computeNextSRS, previewIntervals } from '@/composables/useSRS'
import { useStudyMutations } from '@/composables/useStudySession'
import { renderCloze } from '@/utils/markdown'
import type { Card, DifficultyRating } from '@/types'
import MultipleChoiceViewer from '@/components/card/MultipleChoiceViewer.vue'

const route = useRoute()
const router = useRouter()
const deckId = route.params.deckId as string

// §6 — Session state in local refs only, not Pinia (not shared globally)
const { data: deckCards } = useCards(deckId)
const cardMutations = useCardMutations(deckId)
const studyMutations = useStudyMutations()

const isFlipped = ref(false)
const cramMode = ref(false)
const quizSessionId = ref<string | null>(null)
const reviewedCount = ref(0)
const queue = ref<Card[]>([])
const currentIndex = ref(0)

const currentCard = computed<Card | null>(() => queue.value[currentIndex.value] ?? null)
const progressPercent = computed(() => queue.value.length === 0 ? 100 : (reviewedCount.value / queue.value.length) * 100)

// Preview intervals from useSRS — no duplication of SM-2 logic
const nextIntervals = computed(() => {
  if (!currentCard.value) return { again: '', hard: '', good: '', easy: '' }
  return previewIntervals({
    interval: currentCard.value.interval,
    repetition: currentCard.value.repetition,
    easeFactor: currentCard.value.easeFactor
  })
})

function formatLabel(format: string) {
  return { flashcard: '🃏 Flashcard', cloze: '✏️ Fill in the Blank', free_response: '💬 Free Response' }[format] ?? format
}

function renderContent(text: string, format: string) {
  if (format === 'cloze') return renderCloze(text, false)
  return text
}

function reveal() {
  if (!isFlipped.value) isFlipped.value = true
}

async function rate(rating: DifficultyRating) {
  if (!currentCard.value) return
  const card = currentCard.value

  const ratingMap: Record<DifficultyRating, number> = { again: 1, hard: 2, good: 3, easy: 4 }

  if (quizSessionId.value) {
    try {
      await studyMutations.answerCard.mutateAsync({
        sessionId: quizSessionId.value,
        payload: { cardId: card.id, rating: ratingMap[rating], timeTaken: 1000 }
      })
    } catch (err) {
      console.error('Failed to save answer', err)
    }
  } else if (!cramMode.value) {
    // §7 — compute via useSRS pure function, submit via mutation
    const next = computeNextSRS(
      { interval: card.interval, repetition: card.repetition, easeFactor: card.easeFactor },
      rating
    )
    cardMutations.update.mutate({ id: card.id, payload: next })
  }

  reviewedCount.value++
  currentIndex.value++
  isFlipped.value = false

  if (currentIndex.value >= queue.value.length) {
    if (quizSessionId.value) {
      await studyMutations.endSession.mutateAsync(quizSessionId.value)
      router.push(`/app/decks/${deckId}/quiz-results/${quizSessionId.value}`)
    }
  }
}

function loadDue() {
  const now = new Date()
  queue.value = (deckCards.value ?? [])
    .filter((c: Card) => new Date(c.nextReviewDate) <= now)
  currentIndex.value = 0
}

function startCram() {
  quizSessionId.value = null
  cramMode.value = true
  queue.value = [...(deckCards.value ?? [])]
  currentIndex.value = 0
  reviewedCount.value = 0
}

async function startQuiz() {
  try {
    const res = await studyMutations.startSession.mutateAsync({ deckId, mode: 'quiz' })
    quizSessionId.value = res.session.id
    cramMode.value = false
    queue.value = res.cards
    currentIndex.value = 0
    reviewedCount.value = 0
  } catch (err) {
    console.error('Failed to start quiz session', err)
  }
}

onMounted(() => {
  loadDue()
})
</script>


<style scoped>
.study-view {
  min-height: 80vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  max-width: 760px;
  margin: 0 auto;
}

.session-wrapper { width: 100%; display: flex; flex-direction: column; gap: 1.5rem; }

.session-header {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.back-link { color: var(--color-text-muted); text-decoration: none; font-size: 0.9rem; flex-shrink: 0; transition: color 0.2s; }
.back-link:hover { color: var(--color-primary); }

.progress-wrapper { flex: 1; display: flex; align-items: center; gap: 0.75rem; }

.progress-bar { flex: 1; height: 8px; background: var(--color-surface-hover); border-radius: 999px; overflow: hidden; }
.progress-fill { height: 100%; background: linear-gradient(90deg, var(--color-primary), var(--color-accent)); border-radius: 999px; transition: width 0.4s ease; }
.progress-text { font-size: 0.8rem; color: var(--color-text-muted); flex-shrink: 0; }

.mode-badge { font-size: 0.75rem; padding: 0.25rem 0.65rem; border-radius: 999px; background: var(--color-surface); border: 1px solid var(--color-border); color: var(--color-text-muted); flex-shrink: 0; }

.card-area { display: flex; justify-content: center; }

.study-card {
  width: 100%;
  min-height: 280px;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 24px;
  padding: 3rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s;
  user-select: none;
  position: relative;
}

.study-card:hover { border-color: var(--color-primary); box-shadow: 0 8px 32px rgba(99, 102, 241, 0.15); }
.study-card.is-flipped { border-color: var(--color-accent); cursor: default; }

.card-face { width: 100%; display: flex; flex-direction: column; align-items: center; gap: 1.5rem; }
.card-format-label { font-size: 0.8rem; color: var(--color-text-muted); text-transform: uppercase; letter-spacing: 0.08em; }

.card-text {
  font-size: 1.4rem;
  font-weight: 500;
  color: var(--color-text);
  text-align: center;
  line-height: 1.6;
}

:deep(.cloze-blank) {
  display: inline-block;
  width: 6rem;
  height: 0.15em;
  background: var(--color-primary);
  border-radius: 2px;
  vertical-align: middle;
  margin: 0 0.25rem;
}

.click-hint { font-size: 0.8rem; color: var(--color-text-muted); margin-top: 0.5rem; }

.card-back-face .card-text { color: var(--color-accent); }

/* Rating panel */
.rating-panel { text-align: center; }
.rating-prompt { font-size: 0.9rem; color: var(--color-text-muted); margin: 0 0 1rem; }

.rating-buttons { display: flex; gap: 0.75rem; justify-content: center; }

.rating-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
  padding: 0.9rem 1.25rem;
  border-radius: 14px;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 80px;
}

.rating-icon { font-size: 1.4rem; }
.rating-label { font-size: 0.8rem; font-weight: 600; }
.rating-interval { font-size: 0.7rem; opacity: 0.7; }

.rating-again  { background: rgba(239,68,68,0.1);  color: #ef4444;  } .rating-again:hover  { border-color: #ef4444;  background: rgba(239,68,68,0.2);  }
.rating-hard   { background: rgba(245,158,11,0.1); color: #f59e0b;  } .rating-hard:hover   { border-color: #f59e0b;  background: rgba(245,158,11,0.2);  }
.rating-good   { background: rgba(99,102,241,0.1); color: #6366f1;  } .rating-good:hover   { border-color: #6366f1;  background: rgba(99,102,241,0.2);  }
.rating-easy   { background: rgba(34,197,94,0.1);  color: #22c55e;  } .rating-easy:hover   { border-color: #22c55e;  background: rgba(34,197,94,0.2);  }

/* Finished */
.finished-state { text-align: center; }
.finished-icon { font-size: 4rem; margin: 0; }
.finished-state h2 { font-size: 1.5rem; margin: 0.5rem 0; }
.finished-text { color: var(--color-text-muted); max-width: 400px; margin: 0 auto 2rem; }
.finished-actions { display: flex; gap: 1rem; justify-content: center; }

/* Transitions */
.slide-up-enter-active { transition: all 0.3s ease; }
.slide-up-enter-from { opacity: 0; transform: translateY(12px); }
</style>
