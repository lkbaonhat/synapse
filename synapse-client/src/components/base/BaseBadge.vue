<template>
  <!-- §9 — format resolved to label/class in computed, not inline in template -->
  <span :class="['badge', variantClass]">{{ label }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import type { QuestionFormat } from '@/types'

// §3 — Base prefix; §4 — Generic props syntax
const props = defineProps<{
  format: QuestionFormat
}>()

// §9 — Computed instead of ternary in template
const label = computed(() => {
  const map: Record<QuestionFormat, string> = {
    flashcard: '🃏 Flashcard',
    cloze: '✏️ Cloze',
    free_response: '💬 Free Response'
  }
  return map[props.format]
})

const variantClass = computed(() => {
  const map: Record<QuestionFormat, string> = {
    flashcard: 'badge--flashcard',
    cloze: 'badge--cloze',
    free_response: 'badge--free-response'
  }
  return map[props.format]
})
</script>

<style scoped>
/* §10 — CSS tokens */
.badge {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  padding: 0.2rem 0.65rem;
  border-radius: 999px;
  font-size: 0.72rem;
  font-weight: 600;
  letter-spacing: 0.02em;
}

.badge--flashcard    { background: rgba(99,102,241,0.1); color: var(--color-primary); }
.badge--cloze        { background: rgba(245,158,11,0.1); color: #f59e0b; }
.badge--free-response { background: rgba(34,197,94,0.1); color: #22c55e; }
</style>
