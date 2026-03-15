<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useQuizResults } from '@/composables/useStudySession'

const route = useRoute()
const deckId = route.params.deckId as string
const sessionId = route.params.sessionId as string

const { data: results, isLoading, isError } = useQuizResults(sessionId)

const scorePercentage = computed(() => {
  if (!results.value) return 0
  const total = results.value.totalCorrect + results.value.totalWrong
  if (total === 0) return 0
  return Math.round((results.value.totalCorrect / total) * 100)
})
</script>

<template>
  <div class="quiz-result-view flex flex-col items-center gap-8 py-8 w-full max-w-3xl mx-auto">
    <div v-if="isLoading" class="text-[var(--color-text-muted)] p-8">
      Loading your quiz results...
    </div>

    <div v-else-if="isError" class="text-[var(--color-error)] p-8">
      Failed to load quiz results.
    </div>

    <template v-else-if="results">
      <!-- Scorecard Header -->
      <div class="scorecard bg-[var(--color-surface)] border border-[var(--color-border)] rounded-2xl p-8 w-full shadow-sm text-center">
        <h1 class="text-3xl font-bold mb-6 text-[var(--color-text)]">Quiz Complete! 🎉</h1>
        
        <div class="stats flex justify-center gap-12 mb-6">
          <div class="stat flex flex-col items-center">
            <span class="text-4xl font-extrabold text-[var(--color-success)]">{{ results.totalCorrect }}</span>
            <span class="text-sm uppercase tracking-wide text-[var(--color-text-muted)] mt-1">Correct</span>
          </div>
          <div class="stat flex flex-col items-center">
            <span class="text-4xl font-extrabold text-[var(--color-error)]">{{ results.totalWrong }}</span>
            <span class="text-sm uppercase tracking-wide text-[var(--color-text-muted)] mt-1">Wrong</span>
          </div>
          <div class="stat flex flex-col items-center">
            <span class="text-4xl font-extrabold text-[var(--color-primary)]">{{ scorePercentage }}%</span>
            <span class="text-sm uppercase tracking-wide text-[var(--color-text-muted)] mt-1">Score</span>
          </div>
        </div>

        <div class="actions flex justify-center gap-4 mt-6">
          <RouterLink :to="`/app/deck/${deckId}`" class="px-6 py-3 rounded-lg font-medium bg-[var(--color-surface-hover)] hover:bg-[var(--color-border)] transition-colors">
            Back to Deck
          </RouterLink>
          <RouterLink :to="`/app/study/${deckId}`" class="px-6 py-3 rounded-lg font-medium bg-[var(--color-primary)] text-white hover:opacity-90 transition-opacity">
            Study Again
          </RouterLink>
        </div>
      </div>

      <!-- Review Section -->
      <div v-if="results.wrongAnswers.length > 0" class="review-section w-full border-t border-[var(--color-border)] pt-8">
        <h2 class="text-xl font-semibold mb-6">Let's review what you missed</h2>
        
        <div class="wrong-answers flex flex-col gap-4">
          <div 
            v-for="ans in results.wrongAnswers" 
            :key="ans.cardId" 
            class="missed-card bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-6"
          >
            <div class="prompt text-lg font-medium mb-4 pb-4 border-b border-[var(--color-border)]">
              {{ ans.front }}
            </div>
            <div>
              <span class="text-sm font-semibold uppercase tracking-wide text-[var(--color-success)] mr-2">Correct Answer:</span>
              <span class="text-md text-[var(--color-text)]">{{ ans.correctBack }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="text-center text-[var(--color-success)] font-medium text-xl mt-4">
        Perfect score! You didn't miss any questions. 💯
      </div>
    </template>
  </div>
</template>

<style scoped>
/* Scoped styles use tailwind utils mostly, fallback to vars if needed */
</style>
