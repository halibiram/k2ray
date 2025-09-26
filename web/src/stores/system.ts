import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '../services/api'

interface V2rayStatus {
  status: 'running' | 'stopped'
  pid: number
}

// Corresponds to the Go SystemInfo struct
export interface SystemInfo {
	hostname: string
	os: string
	kernel: string
	cpu: string
	cpu_cores: number
	cpu_usage: number
	memory_total_mb: number
	memory_used_mb: number
	memory_usage: number
	uptime: string
	keenetic_model: string
	firmware_version: string
}

export const useSystemStore = defineStore('system', () => {
  // State
  const apiStatus = ref<'online' | 'offline' | 'loading'>('loading')
  const v2rayStatus = ref<V2rayStatus | null>(null)
  const systemInfo = ref<SystemInfo | null>(null)
  const systemLogs = ref<string>('')
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

  async function fetchSystemInfo() {
    error.value = null
    try {
      const response = await apiClient.get('/system/info')
      systemInfo.value = response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to fetch system info.'
    }
  }

  async function fetchSystemLogs() {
    error.value = null
    try {
      const response = await apiClient.get('/system/logs')
      systemLogs.value = response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to fetch system logs.'
    }
  }

  return {
    apiStatus,
    v2rayStatus,
    systemInfo,
    systemLogs,
    error,
    fetchApiStatus,
    fetchV2rayStatus,
    fetchSystemInfo,
    fetchSystemLogs,
  }
})
