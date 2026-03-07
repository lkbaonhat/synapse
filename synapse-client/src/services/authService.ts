import api from '@/configs/axios'

import type { AuthResponse, LoginInput, RegisterInput } from '@/types'

export const authService = {
  login(payload: LoginInput): Promise<AuthResponse> {
    return api.post<AuthResponse>('/auth/login', payload).then(r => r.data)
  },

  register(payload: Omit<RegisterInput, 'confirmPassword'>): Promise<AuthResponse> {
    return api.post<AuthResponse>('/auth/register', payload).then(r => r.data)
  },

  me(): Promise<AuthResponse['user']> {
    return api.get<AuthResponse['user']>('/auth/me').then(r => r.data)
  },

  logout(): void {
    api.post('/auth/logout').catch(() => {})
  }
}
