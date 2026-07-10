import { http } from './http'

export interface AIConfig {
  id?: number
  name: string
  provider: string
  base_url: string
  api_key: string
  model: string
}

export interface AITestResponse {
  success: boolean
  model: string
  error?: string
}

export function getAIConfig() {
  return http.get<AIConfig>('/ai/config')
}

export function saveAIConfig(payload: AIConfig) {
  return http.post<AIConfig>('/ai/config', payload)
}

export function testAIConfig() {
  return http.post<AITestResponse>('/ai/test')
}
