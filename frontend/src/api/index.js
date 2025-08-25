import axios from 'axios'
import { ElMessage } from 'element-plus'

// Create axios instance
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add auth token from localStorage if available
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      
      // Handle specific error codes
      switch (status) {
        case 401:
          // Unauthorized - redirect to login
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          localStorage.removeItem('role')
          
          if (window.location.pathname.startsWith('/admin')) {
            window.location.href = '/admin/login'
          } else {
            window.location.href = '/login'
          }
          break
        case 403:
          ElMessage.error('Access denied')
          break
        case 404:
          ElMessage.error('Resource not found')
          break
        case 500:
          ElMessage.error('Server error')
          break
        default:
          ElMessage.error(data?.error || 'An error occurred')
      }
    } else if (error.request) {
      ElMessage.error('Network error')
    } else {
      ElMessage.error('Request failed')
    }
    
    return Promise.reject(error)
  }
)

export default api