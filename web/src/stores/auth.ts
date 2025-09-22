import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '../services/api'

// Define the structure of the login credentials
interface LoginCredentials {
  username: string
  password: string
}

export const useAuthStore = defineStore('auth', () => {
  // State
  const accessToken = ref<string | null>(localStorage.getItem('accessToken'))
  const refreshToken = ref<string | null>(localStorage.getItem('refreshToken'))
  const error = ref<string | null>(null)
  const loading = ref(false)

  const isAuthenticated = ref(!!accessToken.value)

  // Actions
  async function login(credentials: LoginCredentials) {
    loading.value = true
    error.value = null
    try {
      const response = await apiClient.post('/auth/login', credentials)
      const data = response.data

      accessToken.value = data.access_token
      refreshToken.value = data.refresh_token
      isAuthenticated.value = true

      localStorage.setItem('accessToken', data.access_token)
      localStorage.setItem('refreshToken', data.refresh_token)

      return true
    } catch (e: any) {
      error.value = e.response?.data?.error || 'An unknown error occurred.'
      isAuthenticated.value = false
      return false
    } finally {
      loading.value = false
    }
  }

  function logout() {
    // In a real app, we would also call a /logout endpoint on the server
    // to invalidate the token on the backend via the blocklist.
    accessToken.value = null
    refreshToken.value = null
    isAuthenticated.value = false
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    // Here we would redirect to the login page
    // router.push('/login')
  }

  return {
    accessToken,
    refreshToken,
    isAuthenticated,
    error,
    loading,
    login,
    logout,
  }
})
