import { useQuery, useMutation } from '@tanstack/vue-query'
import { studyService } from '@/services/studyService'

export function useQuizResults(sessionId: string) {
  return useQuery({
    queryKey: ['quizResults', sessionId],
    queryFn: () => studyService.getQuizResults(sessionId),
    enabled: !!sessionId
  })
}

export function useStudyMutations() {
  const startSession = useMutation({
    mutationFn: studyService.startSession
  })

  const answerCard = useMutation({
    mutationFn: ({ sessionId, payload }: { sessionId: string; payload: any }) =>
      studyService.answerCard(sessionId, payload)
  })

  const endSession = useMutation({
    mutationFn: studyService.endSession
  })

  return { startSession, answerCard, endSession }
}
