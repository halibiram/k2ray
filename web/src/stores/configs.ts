import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '../services/api'

// Define the structure of a V2Ray config object from the frontend's perspective
export interface V2rayConfig {
  id: number
  user_id: number
  name: string
  protocol: string
  config_data: string // This is a JSON string
  created_at: string
  updated_at: string
}

// Define the payload for creating/updating a config
export interface ConfigPayload {
  name: string
  protocol: string
  config_data: any // The raw object, will be stringified
}

export const useConfigStore = defineStore('configs', () => {
  // State
  const configs = ref<V2rayConfig[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Actions
  async function fetchConfigs() {
    isLoading.value = true
    error.value = null
    try {
      const response = await apiClient.get('/configs')
      configs.value = response.data
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to fetch configurations.'
    } finally {
      isLoading.value = false
    }
  }

  async function createConfig(payload: ConfigPayload) {
    isLoading.value = true
    error.value = null
    try {
      const response = await apiClient.post('/configs', {
        ...payload,
        config_data: JSON.stringify(payload.config_data)
      })
      configs.value.push(response.data)
      return true
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to create configuration.'
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function updateConfig(id: number, payload: Partial<ConfigPayload>) {
    isLoading.value = true
    error.value = null
    try {
      const updateData: any = { ...payload }
      if (payload.config_data) {
        updateData.config_data = JSON.stringify(payload.config_data)
      }
      const response = await apiClient.put(`/configs/${id}`, updateData)
      const index = configs.value.findIndex(c => c.id === id)
      if (index !== -1) {
        configs.value[index] = response.data
      }
      return true
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to update configuration.'
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function deleteConfig(id: number) {
    isLoading.value = true
    error.value = null
    try {
      await apiClient.delete(`/configs/${id}`)
      configs.value = configs.value.filter(c => c.id !== id)
      return true
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to delete configuration.'
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function setActiveConfig(id: number) {
    // This action doesn't modify the local state directly, but could if we stored the active ID
    isLoading.value = true
    error.value = null
    try {
      await apiClient.post('/system/active-config', { config_id: id })
      return true
    } catch (e: any) {
      error.value = e.response?.data?.error || 'Failed to set active configuration.'
      return false
    } finally {
      isLoading.value = false
    }
  }

  return {
    configs,
    isLoading,
    error,
    fetchConfigs,
    createConfig,
    updateConfig,
    deleteConfig,
    setActiveConfig,
  }
})
