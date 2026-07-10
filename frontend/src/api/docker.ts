import { http } from './http'

export interface DockerContainer {
  id: string
  name: string
  image: string
  status: string
}

export function getContainers() {
  return http.get<DockerContainer[]>('/docker/containers')
}

export function startContainer(id: string) {
  return http.post(`/docker/start/${id}`)
}

export function stopContainer(id: string) {
  return http.post(`/docker/stop/${id}`)
}

export function restartContainer(id: string) {
  return http.post(`/docker/restart/${id}`)
}

export function removeContainer(id: string) {
  return http.delete(`/docker/remove/${id}`)
}
