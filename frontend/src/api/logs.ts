import { http } from './http'

export interface SystemLogsResponse {
  type: string
  logs: string[]
}

export function getSystemLogs(type = 'system', limit = 100) {
  return http.get<SystemLogsResponse>('/system/logs', {
    params: { type, limit },
  })
}
