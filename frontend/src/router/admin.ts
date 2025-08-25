import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const router = createRouter({
  history: createWebHistory('/admin'),
  routes: [
    {
      path: '/',
      name: 'AdminDashboard',
      component: () => import('@/views/admin/Dashboard.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          name: 'AdminHome',
          component: () => import('@/views/admin/Home.vue'),
        },
        {
          path: 'categories',
          name: 'AdminCategories',
          component: () => import('@/views/admin/Categories.vue'),
        },
        {
          path: 'commodities',
          name: 'AdminCommodities',
          component: () => import('@/views/admin/Commodities.vue'),
        },
        {
          path: 'cards',
          name: 'AdminCards',
          component: () => import('@/views/admin/Cards.vue'),
        },
        {
          path: 'users',
          name: 'AdminUsers',
          component: () => import('@/views/admin/Users.vue'),
        },
        {
          path: 'orders',
          name: 'AdminOrders',
          component: () => import('@/views/admin/Orders.vue'),
        },
      ],
    },
    {
      path: '/login',
      name: 'AdminLogin',
      component: () => import('@/views/admin/Login.vue'),
    },
  ],
})

// Navigation guard
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
  } else if (to.meta.requiresAdmin && authStore.user?.user_type !== 'admin') {
    next('/login')
  } else {
    next()
  }
})

export default router