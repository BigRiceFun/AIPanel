import { http } from './http'

export interface AuditLog {
  id: number
  user_id: number
  username: string
  action: string
  target: string
  ip: string
  created_at: string
}

export function getAuditLogs(limit = 100) {
  return http.get<AuditLog[]>('/audit/logs', {
    params: { limit },
  })
}
