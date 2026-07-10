import { http } from './http'

export interface LoginPayload {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
}

export function login(payload: LoginPayload) {
  return http.post<LoginResponse>('/auth/login', payload)
}
