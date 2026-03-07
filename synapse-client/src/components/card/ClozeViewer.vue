<!-- Cloze card viewer — blanks shown before reveal, answers after -->
<template>
  <div class="cloze-viewer">
    <p class="card-label">FILL IN THE BLANK</p>
    <!-- §9 — Computed handles the rendering, not inline v-html expression -->
    <div class="cloze-content" v-html="renderedCloze"></div>
    <button v-if="!revealed" class="reveal-btn" @click="$emit('reveal')">
      Reveal answer
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import { renderCloze } from '@/utils/markdown'

// §4 — Generic syntax
const props = defineProps<{
  front: string
  back: string
  revealed: boolean
}>()

defineEmits<{
  reveal: []
}>()

// §9 — In computed, not template — switch between blank and revealed state
const renderedCloze = computed(() => renderCloze(props.front, props.revealed))
</script>

<style scoped>
/* §10 — CSS tokens only */
.cloze-viewer {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 2rem;
}

.card-label {
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 0.1em;
  color: var(--color-text-muted);
  margin: 0;
}

.cloze-content {
  font-size: 1.2rem;
  line-height: 1.8;
  color: var(--color-text);
}

/* Cloze blank style — used by renderCloze() */
:deep(.cloze-blank) {
  display: inline-block;
  min-width: 5rem;
  border-bottom: 2px solid var(--color-primary);
  margin: 0 0.25rem;
  color: transparent;
  user-select: none;
}

:deep(.cloze-answer) {
  color: var(--color-accent);
  font-weight: 600;
  border-bottom: 2px solid var(--color-accent);
  padding: 0 0.25rem;
}

.reveal-btn {
  align-self: flex-start;
  background: none;
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  font-size: 0.85rem;
  padding: 0.5rem 1.25rem;
  cursor: pointer;
  transition: all 0.2s;
}

.reveal-btn:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}
</style>
