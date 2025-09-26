import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '../services/api'
import { useRouter } from 'vue-router'

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

  async function logout() {
    try {
      // The interceptor will add the access token to the header.
      await apiClient.post('/auth/logout')
    } catch (e) {
      console.error('Server logout failed, proceeding with client-side cleanup:', e)
    } finally {
      accessToken.value = null
      refreshToken.value = null
      isAuthenticated.value = false
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
      // In a real app, you might want to force a redirect.
      // const router = useRouter() // This doesn't work outside setup.
      // A common pattern is to do the redirect in the component that calls logout.
    }
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
