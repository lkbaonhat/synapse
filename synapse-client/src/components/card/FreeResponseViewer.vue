<!-- Free response card — user types their answer, then compares -->
<template>
  <div class="free-response-viewer">
    <div class="question-area">
      <p class="card-label">QUESTION</p>
      <div class="card-content" v-html="renderedQuestion"></div>
    </div>

    <div v-if="!revealed" class="response-area">
      <label class="response-label" for="user-response">Your answer:</label>
      <textarea
        id="user-response"
        v-model="userAnswer"
        class="input response-textarea"
        placeholder="Type your answer here…"
        rows="3"
      ></textarea>
      <button class="reveal-btn" @click="$emit('reveal')">
        Show model answer
      </button>
    </div>

    <Transition name="fade">
      <div v-if="revealed" class="model-answer">
        <p class="card-label answer-label">MODEL ANSWER</p>
        <div class="card-content" v-html="renderedAnswer"></div>
        <div v-if="userAnswer" class="your-answer-preview">
          <p class="your-label">Your answer:</p>
          <p class="your-text">{{ userAnswer }}</p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

import { renderMarkdown } from '@/utils/markdown'

// §4 — Generic defineProps/defineEmits syntax
const props = defineProps<{
  front: string
  back: string
  revealed: boolean
}>()

defineEmits<{
  reveal: []
}>()

const userAnswer = ref('')

// §9 — Computed for all markdown rendering
const renderedQuestion = computed(() => renderMarkdown(props.front))
const renderedAnswer = computed(() => renderMarkdown(props.back))
</script>

<style scoped>
/* §10 — CSS tokens only */
.free-response-viewer {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.question-area {
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
  margin: 0 0 1rem;
}

.answer-label { color: var(--color-accent); }
.your-label { font-size: 0.78rem; color: var(--color-text-muted); margin: 0 0 0.3rem; }
.your-text { font-size: 0.9rem; color: var(--color-text); margin: 0; }

.card-content {
  font-size: 1.1rem;
  line-height: 1.7;
  color: var(--color-text);
}

.response-area {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.response-label {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--color-text-muted);
}

.response-textarea { resize: vertical; min-height: 90px; }

.model-answer {
  background: rgba(167, 139, 250, 0.05);
  border: 1px solid var(--color-accent);
  border-radius: var(--radius-lg);
  padding: 2rem;
}

.your-answer-preview {
  margin-top: 1.25rem;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
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

.reveal-btn:hover { border-color: var(--color-primary); color: var(--color-primary); }

.fade-enter-active, .fade-leave-active { transition: opacity 0.25s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
