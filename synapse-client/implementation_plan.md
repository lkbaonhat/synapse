# Synapse — Front-End Implementation Plan (Redefined)

> **Scope**: This plan covers the **front-end client only** (`synapse-client`). It assumes a separate back-end REST/GraphQL API will handle persistence, authentication, and SRS scheduling. The front-end will consume that API via a typed HTTP layer.

---

## 1. Technology Stack & Library List

This is the definitive list of all libraries, separated by category.

### Core Framework

| Library | Version | Purpose |
|---|---|---|
| [vue](file:///d:/project-client/synapse-client/src/App.vue) | `^3.5` | Core framework, Composition API, `<script setup>` |
| `vite` | `^7.x` | Build tool & dev server |
| `typescript` | `~5.9` | Full static typing across the codebase |
| `vue-tsc` | `^3.x` | TypeScript type checking for [.vue](file:///d:/project-client/synapse-client/src/App.vue) files |

### Routing

| Library | Version | Purpose |
|---|---|---|
| `vue-router` | `^5.x` | SPA routing, lazy-loaded views, route guards |

### State Management

| Library | Version | Purpose |
|---|---|---|
| `pinia` | `^3.x` | Global stores — auth, decks, cards, study session |
| `pinia-plugin-persistedstate` | `^4.x` | Persist state (auth token, user preferences) to `localStorage` |

### HTTP Client & API Layer

| Library | Version | Purpose |
|---|---|---|
| `axios` | `^1.x` | HTTP client with interceptors for auth tokens, error handling |
| `@tanstack/vue-query` | `^5.x` | Server state management: caching, refetching, loading/error states for all API calls |

> **Rationale**: Pinia handles *local UI state*. Vue Query handles *server state* (API responses, loading, caching). This split avoids mixing concerns.

### Form Handling & Validation

| Library | Version | Purpose |
|---|---|---|
| `vee-validate` | `^4.x` | Composable-first form state management |
| `zod` | `^3.x` | Schema-based validation — shared with back-end |

### UI Component Library

| Library | Version | Purpose |
|---|---|---|
| `@headlessui/vue` | `^1.x` | Unstyled, accessible UI primitives (Modal, Menu, Disclosure) |

> We write our own styled components on top of Headless UI — consistent with the design system in [CODING_STANDARDS.md](file:///d:/project-client/synapse-client/CODING_STANDARDS.md).

### Styling & Theming

| Library | Version | Purpose |
|---|---|---|
| *(Vanilla CSS via [style.css](file:///d:/project-client/synapse-client/src/style.css))* | — | CSS custom properties design tokens (already established) |
| `@fontsource/inter` | `^5.x` | Self-hosted Inter font — no external CDN dependency |

### Markdown & Rich Content

| Library | Version | Purpose |
|---|---|---|
| `marked` | `^15.x` | Fast Markdown-to-HTML parser for card content |
| `highlight.js` | `^11.x` | Syntax highlighting inside code blocks |
| `dompurify` | `^3.x` | Sanitize HTML output before injecting into DOM (XSS prevention) |

### Charts & Data Visualization (Phase 3)

| Library | Version | Purpose |
|---|---|---|
| `chart.js` | `^4.x` | Statistics graphs — bar, line, doughnut |
| `vue-chartjs` | `^5.x` | Vue 3 wrapper for Chart.js |

### File Handling (Phase 4)

| Library | Version | Purpose |
|---|---|---|
| `xlsx` | `^0.18.x` | Parse `.xlsx` / `.csv` files for bulk card import |
| `file-saver` | `^2.x` | Trigger file download for export |

### Utilities

| Library | Version | Purpose |
|---|---|---|
| `dayjs` | `^1.x` | Lightweight date formatting and arithmetic (SRS date math) |
| `@vueuse/core` | `^13.x` | Collection of composables — `useLocalStorage`, `useDebounce`, `useIntersectionObserver`, etc. |
| `nanoid` | `^5.x` | Collision-resistant ID generation for optimistic UI |

### Testing

| Library | Version | Purpose |
|---|---|---|
| `vitest` | `^3.x` | Unit testing — fast, Vite-native |
| `@vue/test-utils` | `^2.x` | Vue component testing utilities |
| `msw` | `^2.x` | Mock Service Worker — intercept HTTP calls in tests without real API |
| `@testing-library/vue` | `^8.x` | User-centric component testing (replaces brittle selector tests) |

### Code Quality

| Library | Version | Purpose |
|---|---|---|
| `eslint` | `^9.x` | Static analysis |
| `@vue/eslint-config-typescript` | `^14.x` | Vue + TypeScript ESLint rules |
| `prettier` | `^3.x` | Code formatting |
| `husky` | `^9.x` | Git hooks |
| `lint-staged` | `^15.x` | Run linters only on staged files |

---

## 2. Application Directory Structure

```
src/
├── assets/               # Fonts, images, global SVG icons
├── components/           # Dumb/presentational reusable components
│   ├── base/             # BaseButton, BaseInput, BaseModal, BaseBadge…
│   ├── card/             # CardViewer, ClozeCard, FlashcardFace…
│   ├── deck/             # DeckCard, DeckGrid…
│   └── charts/           # ForecastChart, MasteryRing…
├── composables/          # useStudySession, useSRS, useImport, useStreak…
├── configs/              # axios instance, app constants
├── directive/            # v-markdown, v-focus…
├── hooks/                # (alias for composables — keep for legacy, merge eventually)
├── layouts/              # DefaultLayout, AuthLayout
├── router/               # index.ts, guards.ts, routes.ts
├── services/             # api.ts, deckService.ts, cardService.ts, authService.ts
├── store/                # useAuthStore, useDeckStore, useCardStore, useTagStore…
├── types/                # index.ts — all shared TypeScript interfaces & Zod schemas
├── utils/                # srs.ts, formatters.ts, parsers.ts
├── views/                # Route-level page components
│   ├── auth/             # LoginView, RegisterView
│   ├── LibraryView.vue
│   ├── DeckDetailView.vue
│   ├── StudyView.vue
│   └── StatsView.vue
├── App.vue
├── main.ts
└── style.css
```

---

## 3. Component Architecture Strategy (Core vs Project)

To ensure high reusability and prepare for potential future project overrides (without splitting repositories), we enforce a strict separation of concerns within `src/components/`:

#### 1. Core / Base Components (`src/components/base/`)
These act as the universal design system primitives. They must be 100% project-agnostic.
- **Rules:**
  - **No Business Logic:** Cannot import Pinia stores, Vue Router, or Axios/services.
  - **No Domain Types:** Cannot use Synapse-specific types like `Deck` or `Card`. Only primitive types (`string`, `boolean`, `number`) or generic structural types.
  - **Data Down, Events Up:** All data is passed via `defineProps()`. All interactions bubble up via `defineEmits()`.
  - **Styling:** Exclusively use CSS custom properties (`var(--color-*)`) for theming so that a different project only needs to swap `style.css` to completely re-theme the components.
  - **Flexibility:** Use Slots (named & scoped) heavily to allow the parent to customize inner content.

#### 2. Project-Specific Components (`src/components/deck/`, `src/components/card/`, etc.)
These are domain-specific wrappers that implement Synapse's business logic.
- **Rules:**
  - **Logic Allowed:** Can import stores, call API composables (`useDecks`), and use domain types.
  - **Composition:** They should compose multiple `Base*` components together to form a feature (e.g., `DeckCard.vue` uses `BaseCard`, `BaseBadge`, `BaseButton`).

#### 3. Unique Project Overrides
If a specific project requires a component that diverges significantly from the `Base*` standard:
1. **Light Override:** The project simply overrides the specific `Base*` component CSS using `:deep()` in the parent component's `<style scoped>`.
2. **Heavy Override/Unique Component:** The project creates its own specialized component (e.g., `src/components/overrides/SpecialButton.vue`) that wraps or replaces the `Base` component entirely, leaving `src/components/base/` untouched.

---

## 4. API Communication Strategy

All HTTP calls live in `src/services/`. Components and composables **never** call `axios` directly.

```ts
// src/configs/axios.ts — Axios instance with interceptors
const api = axios.create({ baseURL: import.meta.env.VITE_API_BASE_URL })

api.interceptors.request.use(config => {
  const { token } = useAuthStore()
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(null, error => {
  if (error.response?.status === 401) useAuthStore().logout()
  return Promise.reject(error)
})
```

Vue Query wraps service calls at the composable layer for caching:
```ts
// src/composables/useDecks.ts
export function useDecks() {
  return useQuery({ queryKey: ['decks'], queryFn: deckService.getAll })
}
```

---

## 5. Implementation Phases

### Test-Driven Development (TDD) Requirement

All feature development globally must follow a strict Test-Driven Development (TDD) cycle:
1. **Red**: Write a failing test for the next piece of functionality (e.g., a `vitest` unit test for a composable, or a `@testing-library/vue` component test). MSW should be used to mock API calls before the API wrapper is even written.
2. **Green**: Write the minimum amount of Vue/TypeScript code necessary to make the test pass.
3. **Refactor**: Clean up the code while ensuring the tests continue to pass, strictly adhering to the project's coding standards.

Tests are not an afterthought; they are the primary driver of the design and implementation. Every Phase must begin with writing its corresponding tests.

---

### Phase 1 — Project Foundation & Auth *(1 week)*
- [ ] Install and configure all libraries (axios, vue-query, pinia-persistedstate, vee-validate, zod, dayjs, vueuse)
- [ ] Set up ESLint, Prettier, Husky, lint-staged
- [ ] Build `AuthLayout` — Login & Register pages with form validation (vee-validate + zod)
- [ ] `useAuthStore` with persisted JWT token
- [ ] Route guards: redirect unauthenticated users to `/login`
- [ ] Axios instance with token interceptor and 401 auto-logout

---

### Phase 2 — Library (Decks, Folders, Tags) *(1 week)*

> Every file produced in this phase must pass the following checklist before being committed:

#### ✅ Phase 2 — CODING_STANDARDS Compliance Checklist
| Rule | Check |
|---|---|
| All imports use `@/` alias — no `../` relative paths | §5 Import Order |
| No `any` — use typed interfaces or `unknown` + narrowing | §2 TypeScript |
| `<script setup lang="ts">` only — no Options API | §1 Vue 3 |
| `defineProps<T>()` and `defineEmits<{ event: [] }>()` generics | §4 Props & Emits |
| All HTTP calls through `src/services/` — never direct axios | §7 API Layer |
| Vue Query (`useQuery`/`useMutation`) in composables — not stores | §7 API Layer |
| Pinia stores hold local UI state only (no API data) | §6 State Management |
| CSS uses `var(--color-*)` tokens — no hardcoded hex values | §10 CSS |
| `<style scoped>` on every component | §10 CSS |
| `err: unknown` + `AxiosError` cast in catch blocks | §11 Error Handling |
| `v-for` + `v-if` on separate elements via `<template>` wrapper | §9 Template Rules |
| Complex template conditions extracted to `computed` | §9 Template Rules |

---

#### 2.1 — Services Layer (`src/services/`)
*Standard: Rule §7 — API Layer*

- [ ] **`deckService.ts`** — typed CRUD functions (`getAll`, `getById`, `create`, `update`, `delete`). Each returns a typed Promise via `api.get().then(r => r.data)`. Import with `@/configs/axios` and `@/types`.
- [ ] **`cardService.ts`** — same pattern for Card endpoints (`getByDeck`, `create`, `update`, `delete`, `bulkCreate`).
- [ ] **`tagService.ts`** — CRUD for Tags.
- [ ] **`folderService.ts`** — CRUD for Folders.
- [ ] All service functions have **explicit return types** — no implicit `any`.

```ts
// ✅ Standard pattern for every service
// src/services/deckService.ts
import api from '@/configs/axios'
import type { Deck } from '@/types'

export const deckService = {
  getAll: (): Promise<Deck[]> => api.get('/decks').then(r => r.data),
  getById: (id: string): Promise<Deck> => api.get(`/decks/${id}`).then(r => r.data),
  create: (payload: Omit<Deck, 'id' | 'createdAt' | 'updatedAt'>): Promise<Deck> =>
    api.post('/decks', payload).then(r => r.data),
  update: (id: string, payload: Partial<Deck>): Promise<Deck> =>
    api.patch(`/decks/${id}`, payload).then(r => r.data),
  remove: (id: string): Promise<void> => api.delete(`/decks/${id}`).then(r => r.data)
}
```

#### 2.2 — Vue Query Composables (`src/composables/`)
*Standard: Rule §7 (API Layer) + Rule §8 (Composables)*

- [ ] **`useDecks.ts`** — wraps `deckService.getAll` with `useQuery`. Exposes `{ data, isLoading, isError, error }`.
- [ ] **`useDeckMutations.ts`** — wraps create/update/delete with `useMutation`. Each mutation **invalidates** the `['decks']` query key on success.
- [ ] **`useCards.ts`** — `useQuery` keyed by `['cards', deckId]`.
- [ ] **`useCardMutations.ts`** — create/update/delete card mutations with query invalidation.
- [ ] **`useTags.ts`** + **`useTagMutations.ts`** — same pattern.
- [ ] Composables live in `src/composables/`, named `use<Feature>.ts` — §8.
- [ ] Composables **never** use `useRouter`, `useRoute` directly — those belong in views.

```ts
// ✅ Pattern for every Vue Query composable
// src/composables/useDecks.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'

import { deckService } from '@/services/deckService'
import type { Deck } from '@/types'

export function useDecks() {
  return useQuery({ queryKey: ['decks'], queryFn: deckService.getAll })
}

export function useDeckMutations() {
  const client = useQueryClient()
  const create = useMutation({
    mutationFn: (payload: Omit<Deck, 'id' | 'createdAt' | 'updatedAt'>) =>
      deckService.create(payload),
    onSuccess: () => client.invalidateQueries({ queryKey: ['decks'] })
  })
  return { create }
}
```

#### 2.3 — Pinia Stores (UI State Only)
*Standard: Rule §6 — State Management*

- [ ] Stores hold **only local UI state**: selected tag filter, modal open/close, expanded folder ID.
- [ ] No API calls or `axios` imports in any store.
- [ ] All stores use Setup Store syntax (`defineStore('name', () => { ... })`).
- [ ] Persistence via `pinia-plugin-persistedstate` with explicit `pick` array.

#### 2.4 — LibraryView Refactor
*Standard: Rules §1, §9, §10*

- [ ] Replace direct `deckStore.decks` data source with `useDecks()` composable result.
- [ ] Show loading skeleton while `isLoading === true` (no raw spinners in template logic).
- [ ] Show error alert when `isError === true` — `err: unknown` pattern, display `error.message`.
- [ ] Empty state when `data.length === 0`.
- [ ] All modals extracted to `BaseModal.vue` component — §3 Component Naming (`Base` prefix).
- [ ] `v-for` + `v-if` separated using `<template>` wrapper — §9.
- [ ] CSS uses only `var(--color-*)` — §10.

#### 2.5 — DeckDetailView Refactor
*Standard: Rules §1, §4, §9, §10*

- [ ] Replace direct `cardStore.cards` with `useCards(deckId)` composable.
- [ ] Card format filter extracted to a `computed` property — §9.
- [ ] Add Card form uses `vee-validate` + `zod` schema for validation — consistent with Auth.
- [ ] Each card format (Flashcard, Cloze, Free Response) rendered by a dedicated child component:
  - `FlashcardViewer.vue`, `ClozeViewer.vue`, `FreeResponseViewer.vue` in `src/components/card/`.
  - Each uses `defineProps<T>()` generic syntax — §4.
- [ ] Markdown rendered via `marked + highlight.js + dompurify` in a `v-markdown` directive or `MarkdownRenderer.vue` component.

#### 2.6 — Base Components (`src/components/base/`)
*Standard: Rule §3 — Component Naming (`Base` prefix)*

- [ ] **`BaseModal.vue`** — accessible modal with `@headlessui/vue Dialog`. Accepts `title` prop, emits `close` event using tuple syntax `{ close: [] }`.
- [ ] **`BaseButton.vue`** — typed props `variant: 'primary' | 'secondary'`, `size: 'sm' | 'md'`, `loading: boolean`. Uses `<style scoped>` + CSS tokens.
- [ ] **`BaseInput.vue`** — label, error message slot, wraps `<input class="input">`.
- [ ] **`BaseBadge.vue`** — format badge (Flashcard / Cloze / Free Response).

#### 2.7 — Tag Management
*Standard: Rules §7, §8, §10*

- [ ] `useTags()` composable for query, `useTagMutations()` for mutations.
- [ ] Tag color picker uses `<input type="color">` — persists via `tagService.create`.
- [ ] Tags displayed as chips with `--tag-color` CSS custom property — §10.

#### 2.8 — Markdown Rendering
*Standard: Rules §2, §9, §10*

- [ ] `src/utils/markdown.ts` — pure function `renderMarkdown(text: string): string` using `marked + highlight.js + dompurify`.
- [ ] Consumed via `v-html` with sanitised output only — `vue/no-v-html` ESLint warning acknowledged.
- [ ] No logic in templates — call `renderMarkdown()` from a `computed` property.

---

### Phase 3 — Study Engine *(1 week)*
*All Phase 2 standards apply. Additional rules:*

- [ ] **`useSRS.ts`** composable in `src/composables/` — pure SM-2 function, no side effects, no Pinia imports, fully unit-testable.
- [ ] **`StudyView`** consumes `useCards(deckId)` and `useCardMutations()` — no direct store access.
- [ ] Rating action calls `cardMutations.update.mutate(...)` — never mutates store directly.
- [ ] Session state (current card index, flipped, cram mode) in local `ref`s — not Pinia (not shared globally).
- [ ] Mobile card flip: CSS-only `transform: rotateY()` — no JS animation library.

---

### Phase 4 — Statistics & Motivation *(4–5 days)*
*All Phase 2 standards apply. Additional rules:*

- [ ] Chart data derived from `useCards()` composable result — no duplicate API calls.
- [ ] `useStreak.ts` composable — reads last-study date from persisted Pinia store, computes streak with `dayjs`.
- [ ] `StatsView` uses `vue-chartjs` — chart config in a `computed` property, not inline in template.
- [ ] `chart.js` lazy-loaded: `const { Chart } = await import('chart.js')` — reduces initial bundle.

---

### Phase 5 — Extended Utilities & Production Readiness *(4–5 days)*
*All Phase 2 standards apply. Additional rules:*

- [ ] `useImport.ts` composable — wraps `xlsx` parsing, returns typed `Card[]`, validates with `zod` before calling `cardService.bulkCreate`.
- [ ] `useExport.ts` composable — builds export payload from `useDecks()` / `useCards()` data, triggers download via `file-saver`.
- [ ] `GlobalErrorBoundary.vue` — Vue `errorCaptured` lifecycle hook, renders fallback UI.
- [ ] Skeleton loaders: `DeckCardSkeleton.vue`, `CardListSkeleton.vue` — shown while `isLoading === true`.
- [ ] `xlsx` lazy-loaded for import/export: `const xlsx = await import('xlsx')`.
- [ ] Unit tests: `vitest` for `useSRS`, `renderMarkdown`, all service functions (with MSW mocks).

---

## 6. Environment Variables

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_APP_NAME=Synapse
```

---

## Verification Plan

### Automated Tests
- `vitest` unit tests for `useSRS()` — all SM-2 algorithm edge cases
- `@testing-library/vue` component tests for `StudyView` card flip flow
- MSW intercepts all service tests so no real API is required

### Manual Verification
- Full user journey: Register → Create Deck → Add Cards → Study → View Stats
- Test all 3 card formats (Flashcard, Cloze, Free Response)
- Validate CSV import with malformed data to verify error handling
- Verify JWT expiry triggers redirect to login
