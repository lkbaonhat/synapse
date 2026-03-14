// MSW request handlers shared across all tests
// Add handlers here that should be available by default.
// Individual tests can add or override handlers via server.use(…)

import { http, HttpResponse } from 'msw'

const BASE = 'http://localhost:8080/api/v1'

export const handlers = [
  // ── Auth ────────────────────────────────────────────────────────────
  http.post(`${BASE}/auth/login`, () => {
    return HttpResponse.json({
      token: 'test-jwt-token',
      user: { id: 'user-1', name: 'Test User', email: 'test@example.com', createdAt: '2024-01-01T00:00:00Z' }
    })
  }),

  http.post(`${BASE}/auth/register`, () => {
    return HttpResponse.json({
      token: 'test-jwt-token',
      user: { id: 'user-1', name: 'Test User', email: 'test@example.com', createdAt: '2024-01-01T00:00:00Z' }
    })
  }),

  // ── Decks ───────────────────────────────────────────────────────────
  http.get(`${BASE}/decks`, () => {
    return HttpResponse.json([
      { id: 'deck-1', name: 'JavaScript', description: 'JS basics', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
      { id: 'deck-2', name: 'TypeScript', description: 'TS types', createdAt: '2024-01-02T00:00:00Z', updatedAt: '2024-01-02T00:00:00Z' }
    ])
  }),

  http.post(`${BASE}/decks`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>
    return HttpResponse.json({
      id: 'deck-new',
      ...body,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }, { status: 201 })
  }),

  http.delete(`${BASE}/decks/:id`, () => {
    return new HttpResponse(null, { status: 204 })
  }),

  // ── Cards ───────────────────────────────────────────────────────────
  http.get(`${BASE}/decks/:deckId/cards`, ({ params }) => {
    const deckId = params.deckId as string
    return HttpResponse.json([
      {
        id: 'card-1', deckId, format: 'flashcard',
        front: 'What is a closure?', back: 'A function with access to its outer scope.',
        tagIds: [], interval: 0, repetition: 0, easeFactor: 2.5,
        nextReviewDate: '2024-01-01T00:00:00Z',
        createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z'
      }
    ])
  }),

  // ── Tags ─────────────────────────────────────────────────────────────
  http.get(`${BASE}/tags`, () => {
    return HttpResponse.json([
      { id: 'tag-1', name: 'JavaScript', color: '#f7df1e' }
    ])
  }),

  // ── Folders ──────────────────────────────────────────────────────────
  http.get(`${BASE}/folders`, () => {
    return HttpResponse.json([])
  })
]
