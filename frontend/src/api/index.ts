import axios from 'axios'

// Create axios instance
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
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
    return response.data
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error.message)
  }
)

export default api

// API endpoints
export const authAPI = {
  login: (data: { username: string; password: string }) =>
    api.post('/auth/login', data),
  register: (data: { username: string; email: string; password: string; phone?: string; qq?: string }) =>
    api.post('/auth/register', data),
  adminLogin: (data: { username: string; password: string }) =>
    api.post('/auth/admin/login', data),
  getProfile: () => api.get('/user/profile'),
}

export const categoryAPI = {
  getCategories: (params?: any) => api.get('/public/categories', { params }),
  getCategory: (id: number) => api.get(`/public/categories/${id}`),
  createCategory: (data: any) => api.post('/user/categories', data),
  updateCategory: (id: number, data: any) => api.put(`/user/categories/${id}`, data),
  deleteCategory: (id: number) => api.delete(`/user/categories/${id}`),
}

export const commodityAPI = {
  getCommodities: (params?: any) => api.get('/public/commodities', { params }),
  getCommodity: (id: number) => api.get(`/public/commodities/${id}`),
  createCommodity: (data: any) => api.post('/user/commodities', data),
  updateCommodity: (id: number, data: any) => api.put(`/user/commodities/${id}`, data),
  deleteCommodity: (id: number) => api.delete(`/user/commodities/${id}`),
}

export const cardAPI = {
  getCards: (params?: any) => api.get('/user/cards', { params }),
  getCard: (id: number) => api.get(`/user/cards/${id}`),
  createCard: (data: any) => api.post('/user/cards', data),
  batchCreateCards: (data: any) => api.post('/user/cards/batch', data),
  updateCard: (id: number, data: any) => api.put(`/user/cards/${id}`, data),
  deleteCard: (id: number) => api.delete(`/user/cards/${id}`),
  batchDeleteCards: (data: { ids: number[] }) => api.delete('/user/cards/batch', { data }),
}

// Admin API endpoints
export const adminAPI = {
  categories: {
    getAll: (params?: any) => api.get('/admin/categories', { params }),
    create: (data: any) => api.post('/admin/categories', data),
    update: (id: number, data: any) => api.put(`/admin/categories/${id}`, data),
    delete: (id: number) => api.delete(`/admin/categories/${id}`),
  },
  commodities: {
    getAll: (params?: any) => api.get('/admin/commodities', { params }),
    create: (data: any) => api.post('/admin/commodities', data),
    update: (id: number, data: any) => api.put(`/admin/commodities/${id}`, data),
    delete: (id: number) => api.delete(`/admin/commodities/${id}`),
  },
  cards: {
    getAll: (params?: any) => api.get('/admin/cards', { params }),
    create: (data: any) => api.post('/admin/cards', data),
    batchCreate: (data: any) => api.post('/admin/cards/batch', data),
    update: (id: number, data: any) => api.put(`/admin/cards/${id}`, data),
    delete: (id: number) => api.delete(`/admin/cards/${id}`),
    batchDelete: (data: { ids: number[] }) => api.delete('/admin/cards/batch', { data }),
  },
}