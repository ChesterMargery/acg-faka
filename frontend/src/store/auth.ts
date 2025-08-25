import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authAPI } from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<any>(null)
  const isLoggedIn = ref<boolean>(!!token.value)

  // Load user from localStorage
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      user.value = JSON.parse(userStr)
    } catch (e) {
      console.error('Error parsing user from localStorage:', e)
    }
  }

  const login = async (credentials: { username: string; password: string }) => {
    try {
      const response = await authAPI.login(credentials)
      token.value = response.token
      user.value = response.user
      isLoggedIn.value = true
      
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
      
      return response
    } catch (error) {
      throw error
    }
  }

  const adminLogin = async (credentials: { username: string; password: string }) => {
    try {
      const response = await authAPI.adminLogin(credentials)
      token.value = response.token
      user.value = response.user
      isLoggedIn.value = true
      
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
      
      return response
    } catch (error) {
      throw error
    }
  }

  const register = async (userData: {
    username: string
    email: string
    password: string
    phone?: string
    qq?: string
  }) => {
    try {
      const response = await authAPI.register(userData)
      token.value = response.token
      user.value = response.user
      isLoggedIn.value = true
      
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(response.user))
      
      return response
    } catch (error) {
      throw error
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    isLoggedIn.value = false
    
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  const getProfile = async () => {
    try {
      const response = await authAPI.getProfile()
      user.value = response
      localStorage.setItem('user', JSON.stringify(response))
      return response
    } catch (error) {
      throw error
    }
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    adminLogin,
    register,
    logout,
    getProfile,
  }
})