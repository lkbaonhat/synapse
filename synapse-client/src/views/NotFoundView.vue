<template>
  <div class="not-found">
    <div class="not-found-glow"></div>
    <p class="not-found-code">404</p>
    <p class="not-found-msg">Page not found</p>
    <p class="not-found-sub">The page you're looking for doesn't exist or has been moved.</p>
    <RouterLink to="/app/library" class="btn btn-primary">← Back to Library</RouterLink>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<style scoped>
.not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  gap: 0.75rem;
  text-align: center;
  padding: 2rem;
  position: relative;
  overflow: hidden;
}

.not-found-glow {
  position: absolute;
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(99,102,241,0.2) 0%, transparent 70%);
  border-radius: 50%;
  pointer-events: none;
}

.not-found-code {
  font-size: 8rem;
  font-weight: 700;
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
  line-height: 1;
}

.not-found-msg {
  font-size: 1.4rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0;
}

.not-found-sub {
  font-size: 0.9rem;
  color: var(--color-text-muted);
  margin: 0 0 0.5rem;
}
</style>
