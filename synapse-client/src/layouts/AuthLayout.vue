<template>
  <div class="auth-shell">
    <!-- Decorative background blobs -->
    <div class="blob blob-1"></div>
    <div class="blob blob-2"></div>

    <div class="auth-container">
      <!-- Left branding panel -->
      <div class="auth-brand">
        <div class="brand-logo">
          <span class="logo-icon">⚡</span>
          <span class="logo-text">Synapse</span>
        </div>
        <h1 class="brand-headline">Learn Smarter,<br />Not Harder.</h1>
        <p class="brand-sub">
          Built on the science of spaced repetition and active recall — the only study tool you need.
        </p>
        <div class="brand-stats">
          <div class="brand-stat">
            <span class="brand-stat-value">SM-2</span>
            <span class="brand-stat-label">Algorithm</span>
          </div>
          <div class="brand-stat">
            <span class="brand-stat-value">3×</span>
            <span class="brand-stat-label">Retention</span>
          </div>
          <div class="brand-stat">
            <span class="brand-stat-value">∞</span>
            <span class="brand-stat-label">Subjects</span>
          </div>
        </div>
      </div>

      <!-- Right form panel -->
      <div class="auth-form-panel">
        <RouterView />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<style scoped>
.auth-shell {
  min-height: 100vh;
  background: var(--color-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  position: relative;
  overflow: hidden;
}

/* Decorative blobs */
.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.15;
  pointer-events: none;
}

.blob-1 {
  width: 500px;
  height: 500px;
  background: var(--color-primary);
  top: -100px;
  left: -100px;
}

.blob-2 {
  width: 400px;
  height: 400px;
  background: var(--color-accent);
  bottom: -80px;
  right: -80px;
}

.auth-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  max-width: 960px;
  width: 100%;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 24px;
  overflow: hidden;
  position: relative;
  z-index: 1;
  box-shadow: 0 24px 80px rgba(0, 0, 0, 0.4);
}

/* Left brand panel */
.auth-brand {
  padding: 3rem;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(167, 139, 250, 0.08));
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 1.5rem;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.logo-icon { font-size: 2rem; }

.logo-text {
  font-size: 1.6rem;
  font-weight: 700;
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.brand-headline {
  font-size: 2rem;
  font-weight: 700;
  color: var(--color-text);
  line-height: 1.25;
  margin: 0;
}

.brand-sub {
  font-size: 0.9rem;
  color: var(--color-text-muted);
  line-height: 1.7;
  margin: 0;
}

.brand-stats {
  display: flex;
  gap: 2rem;
  margin-top: 0.5rem;
}

.brand-stat { display: flex; flex-direction: column; }
.brand-stat-value { font-size: 1.4rem; font-weight: 700; color: var(--color-primary); }
.brand-stat-label { font-size: 0.75rem; color: var(--color-text-muted); }

/* Right form panel */
.auth-form-panel {
  padding: 3rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

/* Responsive: stack on small screens */
@media (max-width: 700px) {
  .auth-container { grid-template-columns: 1fr; }
  .auth-brand { display: none; }
  .auth-form-panel { padding: 2rem; }
}
</style>
