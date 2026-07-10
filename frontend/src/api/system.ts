import { http } from './http'

export interface SystemStatus {
  cpu: number
  memory: number
  disk: number
  hostname: string
  uptime: string
}

export function getSystemStatus() {
  return http.get<SystemStatus>('/system/status')
}
