<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Card, DifficultyRating } from '@/types'

const props = defineProps<{
  card: Card
}>()

const emit = defineEmits<{
  answer: [rating: DifficultyRating]
}>()

// Safely parse JSON payload specific to Multiple Choice
const parsedContent = computed(() => {
  try {
    // We assume the front and back actually contain the encoded JSON,
    // or if the frontend provides a specific 'content' field we use it.
    // Based on Phase 4 spec backend provides Content as JSONB.
    // In our client interface, card has `front` and `back`.
    // Let's assume the JSON is passed raw in `front` or we construct it.
    // Actually, in standard synapse client patterns, a card's format determines its use.
    // We should parse `props.card.front` as JSON if it's the raw payload.
    // For safety, let's try parsing `card.front`
    let raw = props.card.front
    if (raw.startsWith('{')) {
      return JSON.parse(raw) as { prompt: string; options: string[]; correctIndex: number }
    }
  } catch (err) {
    console.error('Failed to parse multiple choice card content', err)
  }
  return { prompt: props.card.front, options: [], correctIndex: 0 }
})

const selectedIndex = ref<number | null>(null)
const hasAnswered = ref(false)

function onSelectOption(index: number) {
  if (hasAnswered.value) return

  selectedIndex.value = index
  hasAnswered.value = true

  const isCorrect = index === parsedContent.value.correctIndex
  // Emit 'easy' for correct (rating 4), 'again' for wrong (rating 1)
  emit('answer', isCorrect ? 'easy' : 'again')
}
</script>

<template>
  <div class="multiple-choice-viewer flex flex-col gap-6">
    <div class="prompt text-lg font-medium text-center">
      {{ parsedContent.prompt }}
    </div>

    <div class="options flex flex-col gap-3">
      <button
        v-for="(option, idx) in parsedContent.options"
        :key="idx"
        @click="onSelectOption(idx)"
        class="option-btn p-4 rounded-xl border-2 text-left transition-colors font-medium hover:bg-[var(--color-surface-hover)] hover:border-[var(--color-primary-light)]"
        :class="{
          'cursor-not-allowed opacity-80': hasAnswered,
          'border-[var(--color-success)] bg-[var(--color-success-light)]': hasAnswered && idx === parsedContent.correctIndex,
          'border-[var(--color-error)] bg-[var(--color-error-light)]': hasAnswered && selectedIndex === idx && idx !== parsedContent.correctIndex,
          'border-[var(--color-border)]': !hasAnswered || (idx !== parsedContent.correctIndex && idx !== selectedIndex)
        }"
        :disabled="hasAnswered"
      >
        <span class="mr-3 text-[var(--color-text-muted)]">{{ String.fromCharCode(65 + idx) }}.</span>
        {{ option }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.option-btn {
  background-color: var(--color-surface);
  color: var(--color-text);
}
</style>
