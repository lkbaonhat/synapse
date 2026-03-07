<template>
  <div class="app-shell">
    <aside class="sidebar">
      <div class="sidebar-logo">
        <span class="logo-icon">⚡</span>
        <span class="logo-text">Synapse</span>
      </div>

      <nav class="sidebar-nav">
        <RouterLink to="/app/library" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">📚</span>
          <span>Library</span>
        </RouterLink>
        <RouterLink to="/app/stats" class="nav-item" active-class="nav-item--active">
          <span class="nav-icon">📊</span>
          <span>Statistics</span>
        </RouterLink>
      </nav>

      <div class="sidebar-footer">
        <div class="user-avatar">{{ userInitial }}</div>
        <div class="user-info">
          <span class="user-name">{{ authStore.user?.name ?? 'Student' }}</span>
          <span class="user-streak">🔥 0 day streak</span>
        </div>
        <button class="logout-btn" @click="handleLogout" title="Sign out">⎋</button>
      </div>
    </aside>

    <main class="main-content">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'

import { useAuthStore } from '@/store/useAuthStore'

const router = useRouter()
const authStore = useAuthStore()

const userInitial = computed(() => {
  const name = authStore.user?.name ?? 'S'
  return name.charAt(0).toUpperCase()
})

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.app-shell {
  display: flex;
  height: 100vh;
  background: var(--color-bg);
  overflow: hidden;
}

.sidebar {
  width: 240px;
  background: var(--color-surface);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  padding: 1.5rem 1rem;
  flex-shrink: 0;
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 2.5rem;
  padding: 0 0.5rem;
}

.logo-icon {
  font-size: 1.75rem;
}

.logo-text {
  font-size: 1.4rem;
  font-weight: 700;
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.65rem 0.75rem;
  border-radius: 10px;
  color: var(--color-text-muted);
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s ease;
}

.nav-item:hover {
  background: var(--color-surface-hover);
  color: var(--color-text);
}

.nav-item--active {
  background: var(--color-primary-subtle);
  color: var(--color-primary);
}

.nav-icon {
  font-size: 1.1rem;
}

.sidebar-footer {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  background: var(--color-surface-hover);
  border-radius: 12px;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--color-primary), var(--color-accent));
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  color: white;
  font-size: 0.9rem;
  flex-shrink: 0;
}

.user-info {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.user-name {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-text);
}

.user-streak {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 2rem 2.5rem;
}

.logout-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--color-text-muted);
  font-size: 1.1rem;
  padding: 0.25rem;
  border-radius: 6px;
  transition: color 0.2s;
  flex-shrink: 0;
  margin-left: auto;
}

.logout-btn:hover { color: #ef4444; }
</style>
