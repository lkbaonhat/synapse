<template>
  <div class="stats-view">
    <header class="page-header">
      <h1 class="page-title">Statistics</h1>
      <p class="page-subtitle">Your learning progress at a glance</p>
    </header>

    <!-- Summary Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📚</div>
        <div class="stat-body">
          <div class="stat-value">{{ totalDecks }}</div>
          <div class="stat-label">Total Decks</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">🃏</div>
        <div class="stat-body">
          <div class="stat-value">{{ totalCards }}</div>
          <div class="stat-label">Total Cards</div>
        </div>
      </div>
      <div class="stat-card highlight">
        <div class="stat-icon">📅</div>
        <div class="stat-body">
          <div class="stat-value">{{ dueToday }}</div>
          <div class="stat-label">Due Today</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">✅</div>
        <div class="stat-body">
          <div class="stat-value">{{ masteredCards }}</div>
          <div class="stat-label">Mastered</div>
        </div>
      </div>
    </div>

    <!-- Deck Breakdown -->
    <div class="section">
      <h2 class="section-title">Deck Progress</h2>
    <div v-if="decks && decks.length > 0" class="deck-breakdown">
        <div v-for="deck in decks" :key="deck.id" class="deck-progress-row">
          <div class="deck-info">
            <span class="deck-progress-name">{{ deck.name }}</span>
            <span class="deck-progress-count">{{ getCardsByDeck(deck.id).length }} cards</span>
          </div>
          <div class="deck-progress-bar-wrapper">
            <div class="deck-progress-bar">
              <div
                class="deck-progress-fill"
                :style="{ width: getMastery(deck.id) + '%' }"
              ></div>
            </div>
            <span class="deck-progress-pct">{{ getMastery(deck.id) }}%</span>
          </div>
          <RouterLink :to="{ name: 'Study', params: { deckId: deck.id } }" class="btn btn-primary btn-sm" v-if="getDueByDeck(deck.id) > 0">
            Study ({{ getDueByDeck(deck.id) }})
          </RouterLink>
          <span class="done-badge" v-else>✓ Done</span>
        </div>
      </div>
      <div v-else class="empty-state">
        <p>No decks yet. <RouterLink to="/app/library">Create one →</RouterLink></p>
      </div>
    </div>

    <!-- Upcoming Reviews (7-day forecast) -->
    <div class="section">
      <h2 class="section-title">7-Day Review Forecast</h2>
      <div class="forecast-chart">
        <div v-for="day in forecast" :key="day.label" class="forecast-bar-wrapper">
          <div class="forecast-bar-container">
            <div
              class="forecast-bar"
              :style="{ height: maxForecast > 0 ? (day.count / maxForecast * 100) + '%' : '0%' }"
              :title="`${day.count} cards`"
            ></div>
          </div>
          <span class="forecast-count">{{ day.count }}</span>
          <span class="forecast-label">{{ day.label }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// §5 — Import order: vue → router → composables → components → types
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

import { useDecks } from '@/composables/useDecks'
import { useAllCards } from '@/composables/useCards'

// §7 — No stores here; all data from Vue Query composables
const { data: decks } = useDecks()
const { data: allCards } = useAllCards()

const totalDecks = computed(() => decks.value?.length ?? 0)
const totalCards = computed(() => allCards.value?.length ?? 0)

const dueToday = computed(() => {
  const now = new Date()
  return (allCards.value ?? []).filter(c => new Date(c.nextReviewDate) <= now).length
})

const masteredCards = computed(() => {
  return (allCards.value ?? []).filter(c => c.interval >= 21).length
})

function getCardsByDeck(deckId: string) {
  return (allCards.value ?? []).filter(c => c.deckId === deckId)
}

function getDueByDeck(deckId: string) {
  const now = new Date()
  return (allCards.value ?? []).filter(c => c.deckId === deckId && new Date(c.nextReviewDate) <= now).length
}

function getMastery(deckId: string) {
  const deckCards = getCardsByDeck(deckId)
  if (!deckCards.length) return 0
  const mastered = deckCards.filter(c => c.interval >= 21).length
  return Math.round((mastered / deckCards.length) * 100)
}

// 7-day forecast
const forecast = computed(() => {
  const days: { label: string; count: number }[] = []
  const dayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']

  for (let i = 0; i < 7; i++) {
    const d = new Date()
    d.setDate(d.getDate() + i)
    const dStart = new Date(d.getFullYear(), d.getMonth(), d.getDate())
    const dEnd = new Date(dStart.getTime() + 86400000)
    const count = (allCards.value ?? []).filter(c => {
      const rv = new Date(c.nextReviewDate)
      return rv >= dStart && rv < dEnd
    }).length
    days.push({ label: i === 0 ? 'Today' : (dayNames[d.getDay()] ?? ''), count })
  }
  return days
})

const maxForecast = computed(() => Math.max(...forecast.value.map(d => d.count), 1))
</script>


<style scoped>
.stats-view { max-width: 900px; margin: 0 auto; }

.page-header { margin-bottom: 2rem; }
.page-title { font-size: 1.8rem; font-weight: 700; color: var(--color-text); margin: 0 0 0.25rem; }
.page-subtitle { color: var(--color-text-muted); margin: 0; font-size: 0.9rem; }

.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 1rem; margin-bottom: 2.5rem; }

.stat-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 16px;
  padding: 1.25rem;
  display: flex;
  gap: 1rem;
  align-items: center;
  transition: all 0.2s;
}

.stat-card:hover { border-color: var(--color-primary); }
.stat-card.highlight { border-color: var(--color-accent); background: rgba(167, 139, 250, 0.05); }

.stat-icon { font-size: 1.75rem; }
.stat-value { font-size: 1.6rem; font-weight: 700; color: var(--color-text); line-height: 1; }
.stat-label { font-size: 0.75rem; color: var(--color-text-muted); margin-top: 0.25rem; }

.section { margin-bottom: 2.5rem; }
.section-title { font-size: 1rem; font-weight: 600; color: var(--color-text-muted); text-transform: uppercase; letter-spacing: 0.05em; margin: 0 0 1rem; }

.deck-breakdown { display: flex; flex-direction: column; gap: 0.75rem; }

.deck-progress-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.25rem;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 12px;
}

.deck-info { width: 180px; flex-shrink: 0; }
.deck-progress-name { font-size: 0.9rem; font-weight: 500; color: var(--color-text); display: block; }
.deck-progress-count { font-size: 0.75rem; color: var(--color-text-muted); }

.deck-progress-bar-wrapper { flex: 1; display: flex; align-items: center; gap: 0.75rem; }
.deck-progress-bar { flex: 1; height: 8px; background: var(--color-surface-hover); border-radius: 999px; overflow: hidden; }
.deck-progress-fill { height: 100%; background: linear-gradient(90deg, var(--color-primary), var(--color-accent)); border-radius: 999px; transition: width 0.5s ease; }
.deck-progress-pct { font-size: 0.8rem; color: var(--color-text-muted); width: 3rem; text-align: right; flex-shrink: 0; }

.done-badge { font-size: 0.8rem; color: #22c55e; font-weight: 500; flex-shrink: 0; }

.empty-state { color: var(--color-text-muted); }
.empty-state a { color: var(--color-primary); }

/* Forecast chart */
.forecast-chart {
  display: flex;
  gap: 0.75rem;
  align-items: flex-end;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 16px;
  padding: 1.5rem 2rem 1rem;
  height: 180px;
}

.forecast-bar-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.35rem;
  height: 100%;
}

.forecast-bar-container {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: flex-end;
}

.forecast-bar {
  width: 100%;
  background: linear-gradient(180deg, var(--color-primary), var(--color-accent));
  border-radius: 6px 6px 0 0;
  min-height: 4px;
  transition: height 0.4s ease;
}

.forecast-count { font-size: 0.75rem; color: var(--color-text-muted); font-weight: 500; }
.forecast-label { font-size: 0.7rem; color: var(--color-text-muted); }
</style>
