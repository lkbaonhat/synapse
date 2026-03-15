import api from '@/configs/axios'
import type { QuizResult, Card } from '@/types'

export interface StartSessionPayload {
  deckId: string
  mode: 'learn' | 'review' | 'cram' | 'quiz'
}

export interface StartSessionResponse {
  session: {
    id: string
    deckId: string
    mode: string
    startedAt: string
  }
  cards: Card[]
}

export interface AnswerCardPayload {
  cardId: string
  rating: number
  timeTaken: number
}

// §7 API Layer - All API calls funnel through this service
export const studyService = {
  startSession: (payload: StartSessionPayload): Promise<StartSessionResponse> =>
    api.post('/study/sessions', payload).then(r => r.data),

  getNextCards: (sessionId: string): Promise<Card[]> =>
    api.get(`/study/sessions/${sessionId}/next`).then(r => r.data),

  answerCard: (sessionId: string, payload: AnswerCardPayload): Promise<void> =>
    api.post(`/study/sessions/${sessionId}/answer`, payload).then(r => r.data),

  endSession: (sessionId: string): Promise<void> =>
    api.post(`/study/sessions/${sessionId}/end`).then(r => r.data),

  getQuizResults: (sessionId: string): Promise<QuizResult> =>
    api.get(`/study/sessions/${sessionId}/results`).then(r => r.data)
}
