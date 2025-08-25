import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || null)
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const role = ref(localStorage.getItem('role') || null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => role.value === 'admin')
  const isUser = computed(() => role.value === 'user')

  const login = async (credentials) => {
    try {
      const response = await api.post('/auth/login', credentials)
      const { token: newToken, user: userData } = response.data
      
      setAuth(newToken, userData, 'user')
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Login failed' }
    }
  }

  const adminLogin = async (credentials) => {
    try {
      const response = await api.post('/auth/admin/login', credentials)
      const { token: newToken, manager } = response.data
      
      setAuth(newToken, manager, 'admin')
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Admin login failed' }
    }
  }

  const register = async (userData) => {
    try {
      const response = await api.post('/auth/register', userData)
      const { token: newToken, user: newUser } = response.data
      
      setAuth(newToken, newUser, 'user')
      return { success: true, data: response.data }
    } catch (error) {
      return { success: false, error: error.response?.data?.error || 'Registration failed' }
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    role.value = null
    
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('role')
    
    // Reset axios authorization header
    delete api.defaults.headers.common['Authorization']
  }

  const setAuth = (newToken, userData, userRole) => {
    token.value = newToken
    user.value = userData
    role.value = userRole
    
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(userData))
    localStorage.setItem('role', userRole)
    
    // Set axios authorization header
    api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
  }

  const updateUser = (userData) => {
    user.value = { ...user.value, ...userData }
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  // Initialize axios header if token exists
  if (token.value) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  return {
    token,
    user,
    role,
    isAuthenticated,
    isAdmin,
    isUser,
    login,
    adminLogin,
    register,
    logout,
    updateUser
  }
})