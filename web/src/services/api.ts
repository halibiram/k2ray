import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api/v1', // The Vite dev server will proxy this to the backend
  headers: {
    'Content-Type': 'application/json',
  },
})

// We can add interceptors here later for attaching auth tokens

export default apiClient
