import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const apiClient = axios.create({
  baseURL: '/api/v1', // The Vite dev server will proxy this to the backend
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add the auth token to headers
apiClient.interceptors.request.use(
  (config) => {
    // We must use the store inside the callback to avoid issues with initialization order.
    const authStore = useAuthStore()
    const token = authStore.accessToken

    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

export default apiClient
