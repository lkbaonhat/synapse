// ── Vue core ────────────────────────────────────────────────────────────
import { createApp } from 'vue'
import { createPinia } from 'pinia'

// ── Third-party plugins ────────────────────────────────────────────────
import { VueQueryPlugin } from '@tanstack/vue-query'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

// ── Internal — router ──────────────────────────────────────────────────
import router from '@/router'

// ── Fonts & global styles (side-effect imports last) ───────────────────
import '@fontsource/inter/400.css'
import '@fontsource/inter/500.css'
import '@fontsource/inter/600.css'
import '@fontsource/inter/700.css'
import './style.css'

// ── Root component ─────────────────────────────────────────────────────
import App from '@/App.vue'

const app = createApp(App)

async function prepareApp() {
  if (import.meta.env.DEV) {
    const { worker } = await import('@/test/msw/browser')
    await worker.start({ onUnhandledRequest: 'bypass' })
  }
}

prepareApp().then(() => {
  const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

  app.use(pinia)
  app.use(router)
  app.use(VueQueryPlugin)

  app.mount('#app')
})
