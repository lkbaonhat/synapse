// Global test setup for vitest
// - Polyfills browser APIs missing in jsdom
// - Configures MSW (Mock Service Worker) to intercept HTTP requests

import { afterAll, afterEach, beforeAll } from 'vitest'
import { server } from './msw/server'

// Start MSW server before all tests
beforeAll(() => server.listen({ onUnhandledRequest: 'warn' }))

// Reset handlers after each test to avoid test pollution
afterEach(() => server.resetHandlers())

// Clean up after all tests
afterAll(() => server.close())
