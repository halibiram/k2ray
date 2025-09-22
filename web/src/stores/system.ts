import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '../services/api'

interface V2rayStatus {
  status: 'running' | 'stopped'
  pid: number
}

export const useSystemStore = defineStore('system', () => {
  // State
  const apiStatus = ref<'online' | 'offline' | 'loading'>('loading')
  const v2rayStatus = ref<V2rayStatus | null>(null)
  const error = ref<string | null>(null)

  // Actions
  async function fetchApiStatus() {
    apiStatus.value = 'loading'
    try {
      // This is a public endpoint, but the interceptor will still work fine.
      await apiClient.get('/system/status')
      apiStatus.value = 'online'
    } catch (e) {
      apiStatus.value = 'offline'
    }
  }

  async function fetchV2rayStatus() {
    error.value = null
    try {
      // This is a protected endpoint. The interceptor will add the auth header.
      const response = await apiClient.get('/v2ray/status')
      v2rayStatus.value = response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to fetch V2Ray status.'
      v2rayStatus.value = null
    }
  }

  return {
    apiStatus,
    v2rayStatus,
    error,
    fetchApiStatus,
    fetchV2rayStatus,
  }
})
