import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '../services/api'

interface LoginCredentials {
  username: string
  password: string
}

export const useAuthStore = defineStore('auth', () => {
  // State
  const accessToken = ref<string | null>(localStorage.getItem('accessToken'))
  const refreshToken = ref<string | null>(localStorage.getItem('refreshToken'))
  const twoFactorToken = ref<string | null>(sessionStorage.getItem('twoFactorToken'))
  const error = ref<string | null>(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!accessToken.value)

  // Actions
  async function login(credentials: LoginCredentials): Promise<{ success: boolean; twoFactor: boolean }> {
    loading.value = true
    error.value = null
    try {
      const response = await apiClient.post('/auth/login', credentials)
      const data = response.data

      if (data.two_factor_token) {
        twoFactorToken.value = data.two_factor_token
        sessionStorage.setItem('twoFactorToken', data.two_factor_token)
        return { success: true, twoFactor: true }
      }

      accessToken.value = data.access_token
      refreshToken.value = data.refresh_token
      localStorage.setItem('accessToken', data.access_token)
      localStorage.setItem('refreshToken', data.refresh_token)
      sessionStorage.removeItem('twoFactorToken')

      return { success: true, twoFactor: false }
    } catch (e: any) {
      error.value = e.response?.data?.error || 'An unknown error occurred.'
      return { success: false, twoFactor: false }
    } finally {
      loading.value = false
    }
  }

  async function verify2FA(code: string) {
    if (!twoFactorToken.value) {
      throw new Error('2FA token not found.')
    }
    loading.value = true
    error.value = null
    try {
      const response = await apiClient.post('/auth/login/2fa', {
        two_factor_token: twoFactorToken.value,
        code,
      })
      const data = response.data
      accessToken.value = data.access_token
      refreshToken.value = data.refresh_token
      localStorage.setItem('accessToken', data.access_token)
      localStorage.setItem('refreshToken', data.refresh_token)
      sessionStorage.removeItem('twoFactorToken')
    } catch (e: any) {
      const errorMessage = e.response?.data?.error || 'Failed to verify 2FA code.'
      error.value = errorMessage
      throw new Error(errorMessage)
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await apiClient.post('/auth/logout')
    } catch (e) {
      console.error('Server logout failed, proceeding with client-side cleanup:', e)
    } finally {
      accessToken.value = null
      refreshToken.value = null
      twoFactorToken.value = null
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
      sessionStorage.removeItem('twoFactorToken')
    }
  }

  async function enable2FA() {
    const response = await apiClient.post('/2fa/enable')
    return response.data // { qrCode: "data:image/png;base64,...", secret: "..." }
  }

  async function verifyAndEnable2FA(code: string) {
    const response = await apiClient.post('/2fa/verify', { code })
    return response.data // { recoveryCodes: [...] }
  }

  async function disable2FA(password: string) {
    await apiClient.post('/2fa/disable', { password })
  }


  return {
    accessToken,
    refreshToken,
    twoFactorToken,
    isAuthenticated,
    error,
    loading,
    login,
    logout,
    verify2FA,
    enable2FA,
    verifyAndEnable2FA,
    disable2FA,
  }
})
