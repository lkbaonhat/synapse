<!-- §3 Component Naming: card/ prefix for card-format-specific viewers -->
<template>
  <div class="card-viewer flashcard-viewer">
    <div class="card-side card-front">
      <p class="side-label">FRONT</p>
      <div class="card-content" v-html="renderedFront"></div>
    </div>
    <Transition name="flip">
      <div v-if="revealed" class="card-side card-back">
        <p class="side-label answer-label">ANSWER</p>
        <div class="card-content" v-html="renderedBack"></div>
      </div>
    </Transition>
    <button v-if="!revealed" class="reveal-btn" @click="$emit('reveal')">
      Click to reveal answer
    </button>
  </div>
</template>

<script setup lang="ts">
// §4 — Generic defineProps/defineEmits syntax
// §9 — Computed for markdown rendering, not inline in template
import { computed } from 'vue'

import { renderMarkdown } from '@/utils/markdown'

const props = defineProps<{
  front: string
  back: string
  revealed: boolean
}>()

defineEmits<{
  reveal: []
}>()

// §9 — Computed instead of method calls in template
const renderedFront = computed(() => renderMarkdown(props.front))
const renderedBack = computed(() => renderMarkdown(props.back))
</script>

<style scoped>
/* §10 — CSS tokens only */
.flashcard-viewer {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.card-side {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 2rem;
  min-height: 140px;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.card-back {
  border-color: var(--color-accent);
  background: rgba(167, 139, 250, 0.05);
}

.side-label {
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 0.1em;
  color: var(--color-text-muted);
  margin: 0;
}

.answer-label { color: var(--color-accent); }

.card-content {
  font-size: 1.1rem;
  line-height: 1.7;
  color: var(--color-text);
}

:deep(code) {
  background: rgba(99, 102, 241, 0.1);
  padding: 0.1em 0.4em;
  border-radius: 4px;
  font-size: 0.9em;
}

:deep(pre) {
  background: #0d0d16;
  border-radius: var(--radius-sm);
  padding: 1rem;
  overflow-x: auto;
}

.reveal-btn {
  align-self: center;
  background: none;
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  font-size: 0.85rem;
  padding: 0.65rem 1.5rem;
  cursor: pointer;
  transition: all 0.2s;
}

.reveal-btn:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.flip-enter-active { transition: all 0.25s ease; }
.flip-enter-from { opacity: 0; transform: translateY(8px); }
</style>
