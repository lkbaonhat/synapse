# Synapse Client — Coding Standards & Best Practices

This document is the **single source of truth** for how code must be written in this repository. All PRs must adhere to these rules before being merged.

---

## 1. Vue 3 — Composition API Only

We exclusively use **Vue 3** with `<script setup lang="ts">`. The Options API is **forbidden**.

```vue
<!-- ✅ GOOD -->
<script setup lang="ts">
import { ref, computed } from 'vue'

const count = ref(0)
const double = computed(() => count.value * 2)
</script>

<!-- ❌ BAD — Options API -->
<script>
export default {
  data() { return { count: 0 } }
}
</script>
```

---

## 2. TypeScript — Strict Mode, No `any`

- Enable strict mode in `tsconfig.json`.
- **Never** use `any`. Use `unknown` for truly dynamic data and narrow it explicitly.
- Use `type` for unions/primitives, `interface` for object shapes.
- Use `import type` for type-only imports to keep runtime bundles clean.

```ts
// ✅ GOOD
import type { Card } from '@/types'

function getCard(id: string): Card | undefined { ... }

// ❌ BAD
function getCard(id: any): any { ... }
```

---

## 3. Component Naming

| Category | Convention | Example |
|---|---|---|
| Regular components | PascalCase file + template usage | `DeckCard.vue`, `<DeckCard />` |
| Base/reusable | `Base` prefix | `BaseButton.vue`, `BaseModal.vue` |
| Single-instance (layout) | `The` prefix | `TheSidebar.vue`, `TheNavbar.vue` |
| Layout wrappers | `Layout` suffix | `DefaultLayout.vue`, `AuthLayout.vue` |
| View components (routed) | `View` suffix | `LibraryView.vue`, `LoginView.vue` |

```vue
<!-- ✅ GOOD — PascalCase in template -->
<DeckCard :deck="deck" @study="handleStudy" />

<!-- ❌ BAD — kebab-case hides component identity -->
<deck-card :deck="deck" @study="handleStudy" />
```

---

## 4. Props & Emits — TypeScript Generic Syntax

Always use the generic `defineProps<T>()` and `defineEmits<T>()` syntax. Never use the runtime object form.

```vue
<script setup lang="ts">
// ✅ GOOD
const props = defineProps<{
  deck: Deck
  isActive?: boolean
}>()

const emit = defineEmits<{
  'study': [deckId: string]
  'delete': [deckId: string]
}>()
</script>
```

---

## 5. Import Order

Imports must follow this order, separated by blank lines:

```ts
// 1. Vue core
import { ref, computed, onMounted } from 'vue'

// 2. Third-party libraries
import { useQuery } from '@tanstack/vue-query'
import dayjs from 'dayjs'

// 3. Vue Router / Pinia
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/store/useAuthStore'

// 4. Services / composables
import { deckService } from '@/services/deckService'
import { useSRS } from '@/composables/useSRS'

// 5. Local types
import type { Deck, Card } from '@/types'

// 6. Assets / styles (only in main.ts or App.vue)
import '@/style.css'
```

---

## 6. State Management (Pinia) — Setup Store Only

- Use **Setup Store** syntax (function body, return statement).
- Store file naming: `use<Name>Store.ts` inside `src/store/`.
- **Never** mutate store state directly from a component. Use store actions.
- Use `pinia-plugin-persistedstate` with explicit `pick` to only persist what is needed.

```ts
// ✅ GOOD — src/store/useDeckStore.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Deck } from '@/types'

export const useDeckStore = defineStore('deck', () => {
  const decks = ref<Deck[]>([])

  const totalCount = computed(() => decks.value.length)

  function addDeck(deck: Deck) {
    decks.value.push(deck)
  }

  return { decks, totalCount, addDeck }
})
```

---

## 7. API Layer — Services, Never Direct Axios

All HTTP calls live in `src/services/`. **Components and composables must never call `axios` directly.**

```
src/
  configs/
    axios.ts          ← Singleton Axios instance with interceptors
  services/
    authService.ts    ← Auth endpoints
    deckService.ts    ← Deck CRUD endpoints
    cardService.ts    ← Card CRUD endpoints
```

```ts
// ✅ GOOD — service function
// src/services/deckService.ts
import api from '@/configs/axios'
import type { Deck } from '@/types'

export const deckService = {
  getAll: (): Promise<Deck[]> => api.get('/decks').then(r => r.data),
  getById: (id: string): Promise<Deck> => api.get(`/decks/${id}`).then(r => r.data),
}

// ✅ GOOD — consumed via Vue Query in a composable
// src/composables/useDecks.ts
import { useQuery } from '@tanstack/vue-query'
import { deckService } from '@/services/deckService'

export function useDecks() {
  return useQuery({ queryKey: ['decks'], queryFn: deckService.getAll })
}

// ❌ BAD — axios in a component
import axios from 'axios'
const decks = await axios.get('/decks')
```

---

## 8. Composables

- Place in `src/composables/`, named `use<Feature>.ts`.
- Must return a reactive object or refs.
- Keep side effects (API calls, event listeners) contained inside the composable using `onMounted`/`onUnmounted`.

```ts
// ✅ GOOD — src/composables/useSRS.ts
export function useSRS() {
  function computeNextInterval(card: Card, rating: DifficultyRating): number { ... }
  return { computeNextInterval }
}
```

---

## 9. Template Rules

- **No logic in templates.** Extract conditions to `computed` properties.
- Keep `v-if` / `v-for` on separate elements via `<template>` wrapper.
- Never use `v-if` and `v-for` on the same element.
- Always provide `:key` on `v-for` — never use index as key unless items are static.

```vue
<!-- ✅ GOOD -->
<template v-for="deck in filteredDecks" :key="deck.id">
  <DeckCard v-if="deck.isVisible" :deck="deck" />
</template>

<!-- ❌ BAD — v-if + v-for on same element, index key -->
<li v-for="(item, i) in items" :key="i" v-if="item.active">{{ item.name }}</li>
```

---

## 10. CSS — Scoped + Design Tokens

- Always use `<style scoped>` in components.
- **Never hardcode colors** — use CSS custom properties from `style.css`.
- Use `:deep()` only for overriding styles in third-party components.
- Group CSS properties: layout → box model → typography → visual → animation.

```vue
<!-- ✅ GOOD -->
<style scoped>
.card {
  /* Layout */
  display: flex;
  flex-direction: column;
  /* Box model */
  padding: 1.5rem;
  border-radius: var(--radius-md);
  /* Visual */
  background: var(--color-surface);
  border: 1px solid var(--color-border);
}
</style>

<!-- ❌ BAD — hardcoded colors -->
<style scoped>
.card { background: #1a1a24; border: 1px solid rgba(255,255,255,0.08); }
</style>
```

---

## 11. Error Handling

- Catch errors at the **service or store level**, not in templates.
- Display user-friendly messages, never raw error objects.
- Use `try/catch` with explicit typed error handling.

```ts
// ✅ GOOD
async function login(credentials: LoginInput) {
  try {
    const data = await authService.login(credentials)
    token.value = data.token
  } catch (err: unknown) {
    if (err instanceof AxiosError) {
      error.value = err.response?.data?.message ?? 'Login failed.'
    }
  }
}
```

---

## 12. File & Directory Structure

```
src/
├── assets/                 # Static files (images, SVG icons)
├── components/
│   ├── base/               # BaseButton, BaseInput, BaseModal…
│   ├── card/               # CardViewer, ClozeCard…
│   └── deck/               # DeckCard, DeckGrid…
├── composables/            # useDecks, useSRS, useStreak…
├── configs/
│   └── axios.ts            # Singleton Axios instance
├── layouts/                # DefaultLayout, AuthLayout
├── router/
│   ├── index.ts            # createRouter
│   └── guards.ts           # Navigation guard functions
├── services/               # deckService, cardService, authService…
├── store/                  # useAuthStore, useDeckStore…
├── types/
│   └── index.ts            # All shared TS interfaces & Zod schemas
├── utils/                  # srs.ts, formatters.ts (pure functions)
├── views/
│   ├── auth/               # LoginView, RegisterView
│   ├── LibraryView.vue
│   ├── DeckDetailView.vue
│   ├── StudyView.vue
│   ├── StatsView.vue
│   └── NotFoundView.vue
├── App.vue
├── main.ts
└── style.css
```

---

*Following these standards ensures the codebase remains scalable, testable, and approachable for any developer joining the project.*

---

## 13. Git Commit Standards — Conventional Commits

We follow the **[Conventional Commits](https://www.conventionalcommits.org/)** specification, enforced automatically by **commitlint** + **Husky**.

> Every commit is linted by `.husky/commit-msg` before it is written to history. Non-conforming commits are **rejected** immediately.

---

### 13.1 Commit Message Format

```
<type>(<scope>): <subject>

[optional body — wrap at 120 chars]

[optional footer — BREAKING CHANGE: ... / Closes #123]
```

| Part | Rule |
|---|---|
| **type** | Required. Must be one of the allowed types below. |
| **scope** | Optional but strongly recommended. Lowercase, hyphenated. |
| **subject** | Required. Imperative mood, no capital letter, no period at end. Min 10 chars, max 100. |
| **body** | Optional. Explain *why*, not *what*. Wrap at 120 chars. |
| **footer** | Optional. `BREAKING CHANGE:` prefix or `Closes #<issue>`. |

---

### 13.2 Allowed Types

| Type | Emoji | When to use |
|---|---|---|
| `feat` | ✨ | New feature added (visible to users) |
| `fix` | 🐛 | Bug fix |
| `docs` | 📝 | Documentation only — README, CODING_STANDARDS, comments |
| `style` | 💅 | Formatting, whitespace — **no logic change** |
| `refactor` | ♻️ | Code change that is neither a fix nor a feature |
| `perf` | ⚡ | Performance improvement |
| `test` | ✅ | Adding or correcting tests |
| `build` | 🏗 | Build system or dependency changes (`package.json`, `vite.config`) |
| `ci` | 🤖 | CI/CD pipeline config |
| `chore` | 🔧 | Tooling, config, scripts (`eslint.config.js`, `.prettierrc`, Husky) |
| `revert` | ⏪ | Reverting a previous commit |

> ❌ **`fix` is NOT for refactors.** ❌ **`chore` is NOT for features.**

---

### 13.3 Allowed Scopes

Use a scope to identify which part of the system changed:

| Scope | Files affected |
|---|---|
| `auth` | Login, Register, useAuthStore, authService |
| `deck` | DeckDetailView, deckService, useDecks |
| `card` | card viewers, cardService, useCards |
| `library` | LibraryView, useLibraryUiStore |
| `study` | StudyView, useSRS |
| `stats` | StatsView, useStreak |
| `router` | src/router/index.ts, guards |
| `store` | Any Pinia store file |
| `types` | src/types/index.ts |
| `ui` | Base components (BaseModal, BaseButton, etc.) |
| `style` | style.css, CSS tokens |
| `config` | vite.config.ts, tsconfig, eslint, prettier, commitlint |
| `deps` | package.json dependency changes |

---

### 13.4 Atomic Commit Rules

> **One commit = one logical change.** Never bundle unrelated changes into a single commit.

| ✅ DO | ❌ DON'T |
|---|---|
| One commit per new service file | Mix service + view + store in one commit |
| One commit per bug fix | Fix a bug while also adding a feature |
| One commit for all related style tweaks | Mix style changes with logic changes |
| Separate `refactor` from `feat` | Add a feature and refactor the same file together |

---

### 13.5 How to Split Commits — Practical Guide

Use `git add -p` (interactive patch staging) to stage only part of a file:

```bash
# Stage specific hunks from a file interactively
git add -p src/views/LibraryView.vue

# Then commit only what's staged
git commit
```

**Typical split order when implementing a feature:**

```
1. build(deps): add dayjs and nanoid packages
2. chore(config): add @/ path alias to vite.config and tsconfig
3. refactor(types): merge auth + domain types into src/types/index.ts
4. feat(auth): add authService with login, register, me, logout
5. feat(auth): add useAuthStore with pinia-persistedstate
6. feat(auth): add LoginView with vee-validate + zod validation
7. feat(auth): add RegisterView with password confirmation schema
8. feat(router): add requiresAuth + requiresGuest navigation guards
```

---

### 13.6 Commit Examples

```bash
# ✅ GOOD — clear type, scope, imperative subject
feat(deck): add deckService with getAll, getById, create, update, remove
fix(auth): redirect to login when 401 response is received
refactor(library): replace cardStore.cards with useCards composable
docs(standards): add §13 git commit standards section
chore(config): configure commitlint with conventional commits ruleset
test(srs): add vitest unit tests for SM-2 algorithm edge cases
build(deps): add marked, highlight.js, dompurify for markdown rendering
style(ui): apply CSS token var(--color-primary) to BaseButton hover state

# ❌ BAD
git commit -m "fix stuff"                    # no type, no scope, vague
git commit -m "WIP"                          # never commit WIP to main
git commit -m "feat: Add Deck Component."    # capital + period
git commit -m "feat(deck): added deck"       # past tense, too short subject
git commit -m "update"                       # meaningless
```

---

### 13.7 Interactive Commit with Commitizen

Instead of typing the format manually, use the interactive prompt:

```bash
# Use the interactive commit wizard (reads .czrc)
pnpm commit

# Or via git directly after setup
git cz
```

This guides you through: **type → scope → subject → body → breaking change → issue reference**.

---

### 13.8 Tooling Setup (one-time, per machine)

```bash
# Install commitlint + commitizen
pnpm add -D @commitlint/cli @commitlint/config-conventional commitizen cz-conventional-changelog

# Add commit script to package.json
# "scripts": { "commit": "cz" }

# Initialise Husky (if not already done)
pnpm exec husky init

# The commit-msg hook is already in .husky/commit-msg
```

---

### 13.9 Branch & PR Strategy

| Branch | Purpose |
|---|---|
| `main` | Production-ready code. Protected — no direct push. |
| `develop` | Integration branch. All feature branches merge here first. |
| `feat/<scope>-<description>` | Feature work. e.g. `feat/deck-service` |
| `fix/<scope>-<description>` | Bug fix. e.g. `fix/auth-401-redirect` |
| `refactor/<scope>` | Refactor. e.g. `refactor/library-view-composables` |
| `chore/<description>` | Config/tooling. e.g. `chore/commitlint-setup` |

> Branch names mirror the commit type prefix — no ambiguity about what a branch contains.

---

*Following these standards ensures a clean, readable git history that the whole team can navigate and tools like `git log`, changelogs, and CI pipelines can parse automatically.*
