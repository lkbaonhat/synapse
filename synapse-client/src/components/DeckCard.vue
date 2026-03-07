<template>
  <div class="deck-card">
    <div class="deck-card-header">
      <div>
        <h3 class="deck-name">{{ deck.name }}</h3>
        <p class="deck-desc" v-if="deck.description">{{ deck.description }}</p>
      </div>
      <button class="icon-btn danger" @click.stop="$emit('delete')" title="Delete deck">✕</button>
    </div>

    <div class="deck-stats">
      <div class="stat">
        <span class="stat-value">{{ cardCount }}</span>
        <span class="stat-label">Total</span>
      </div>
      <div class="stat" :class="{ 'stat--due': dueCount > 0 }">
        <span class="stat-value">{{ dueCount }}</span>
        <span class="stat-label">Due</span>
      </div>
    </div>

    <!-- Mastery progress bar -->
    <div class="mastery-bar-wrapper">
      <div class="mastery-label">
        <span>Mastery</span>
        <span>{{ masteryPercent }}%</span>
      </div>
      <div class="mastery-bar">
        <div class="mastery-bar-fill" :style="{ width: masteryPercent + '%' }"></div>
      </div>
    </div>

    <div class="deck-actions">
      <button class="btn btn-secondary btn-sm" @click="$emit('view')">Browse</button>
      <button class="btn btn-primary btn-sm" :disabled="dueCount === 0" @click="$emit('study')">
        Study {{ dueCount > 0 ? `(${dueCount})` : '' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import type { Deck } from '@/types'

const props = defineProps<{
  deck: Deck
  cardCount: number
  dueCount: number
}>()

defineEmits<{
  study: []
  view: []
  delete: []
}>()

const masteryPercent = computed(() => {
  if (!props.cardCount) return 0
  // Cards with interval >= 21 days are considered "mastered"
  return Math.min(100, Math.round(((props.cardCount - props.dueCount) / props.cardCount) * 100))
})
</script>

<style scoped>
.deck-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 16px;
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  transition: all 0.2s ease;
  cursor: default;
}

.deck-card:hover {
  border-color: var(--color-primary);
  box-shadow: 0 4px 24px rgba(99, 102, 241, 0.12);
  transform: translateY(-2px);
}

.deck-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.deck-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 0.2rem;
}

.deck-desc {
  font-size: 0.8rem;
  color: var(--color-text-muted);
  margin: 0;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  max-width: 200px;
}

.deck-stats {
  display: flex;
  gap: 1.5rem;
}

.stat {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 1.4rem;
  font-weight: 700;
  color: var(--color-text);
}

.stat--due .stat-value { color: var(--color-accent); }

.stat-label {
  font-size: 0.7rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.mastery-bar-wrapper { width: 100%; }

.mastery-label {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: var(--color-text-muted);
  margin-bottom: 0.35rem;
}

.mastery-bar {
  height: 6px;
  background: var(--color-surface-hover);
  border-radius: 999px;
  overflow: hidden;
}

.mastery-bar-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-primary), var(--color-accent));
  border-radius: 999px;
  transition: width 0.4s ease;
}

.deck-actions {
  display: flex;
  gap: 0.5rem;
}

.icon-btn { background: none; border: none; cursor: pointer; padding: 0.25rem 0.5rem; border-radius: 6px; opacity: 0.4; transition: opacity 0.2s; font-size: 0.9rem; }
.icon-btn:hover { opacity: 1; }
.icon-btn.danger:hover { color: #ef4444; }
</style>
