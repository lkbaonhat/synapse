import { createRouter, createWebHistory } from 'vue-router'

import { useAuthStore } from '@/store/useAuthStore'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // ── Auth routes (unauthenticated only) ────────────────────────────────
    {
      path: '/',
      component: () => import('@/layouts/AuthLayout.vue'),
      meta: { requiresGuest: true },
      children: [
        {
          path: '',
          redirect: '/login'
        },
        {
          path: 'login',
          name: 'Login',
          component: () => import('@/views/auth/LoginView.vue'),
          meta: { title: 'Sign In — Synapse' }
        },
        {
          path: 'register',
          name: 'Register',
          component: () => import('@/views/auth/RegisterView.vue'),
          meta: { title: 'Create Account — Synapse' }
        }
      ]
    },

    // ── App routes (authenticated only) ──────────────────────────────────
    {
      path: '/app',
      component: () => import('@/layouts/DefaultLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/app/library'
        },
        {
          path: 'library',
          name: 'Library',
          component: () => import('@/views/LibraryView.vue'),
          meta: { title: 'Library — Synapse' }
        },
        {
          path: 'deck/:deckId',
          name: 'DeckDetail',
          component: () => import('@/views/DeckDetailView.vue'),
          meta: { title: 'Deck — Synapse' }
        },
        {
          path: 'study/:deckId',
          name: 'Study',
          component: () => import('@/views/StudyView.vue'),
          meta: { title: 'Study — Synapse' }
        },
        {
          path: 'decks/:deckId/quiz-results/:sessionId',
          name: 'QuizResult',
          component: () => import('@/views/QuizResultView.vue'),
          meta: { title: 'Quiz Result — Synapse' }
        },
        {
          path: 'stats',
          name: 'Stats',
          component: () => import('@/views/StatsView.vue'),
          meta: { title: 'Statistics — Synapse' }
        }
      ]
    },

    // Redirect old /library paths to /app/library
    {
      path: '/library',
      redirect: '/app/library'
    },
    {
      path: '/stats',
      redirect: '/app/stats'
    },

    // ── 404 ───────────────────────────────────────────────────────────────
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFoundView.vue')
    }
  ]
})

// ── Navigation Guards ────────────────────────────────────────────────────
router.beforeEach((to) => {
  // Update page title
  if (to.meta.title) {
    document.title = to.meta.title as string
  }

  const authStore = useAuthStore()

  // Route requires authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { name: 'Login', query: { redirect: to.fullPath } }
  }

  // Route requires guest (e.g. login page) — redirect logged-in users away
  if (to.meta.requiresGuest && authStore.isAuthenticated) {
    return { name: 'Library', path: '/app/library' }
  }
})

export default router
