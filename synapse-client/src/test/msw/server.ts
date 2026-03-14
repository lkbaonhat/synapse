// MSW Node.js server — used in vitest (Node-compatible via jsdom)
import { setupServer } from 'msw/node'
import { handlers } from './handlers'

export const server = setupServer(...handlers)
