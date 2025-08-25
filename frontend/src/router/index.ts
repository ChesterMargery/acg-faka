import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('@/views/Home.vue'),
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/Register.vue'),
    },
    {
      path: '/category/:id',
      name: 'Category',
      component: () => import('@/views/Category.vue'),
    },
    {
      path: '/commodity/:id',
      name: 'Commodity',
      component: () => import('@/views/Commodity.vue'),
    },
    {
      path: '/user',
      name: 'UserDashboard',
      component: () => import('@/views/user/Dashboard.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'UserHome',
          component: () => import('@/views/user/Home.vue'),
        },
        {
          path: 'profile',
          name: 'UserProfile',
          component: () => import('@/views/user/Profile.vue'),
        },
        {
          path: 'categories',
          name: 'UserCategories',
          component: () => import('@/views/user/Categories.vue'),
        },
        {
          path: 'commodities',
          name: 'UserCommodities',
          component: () => import('@/views/user/Commodities.vue'),
        },
        {
          path: 'cards',
          name: 'UserCards',
          component: () => import('@/views/user/Cards.vue'),
        },
      ],
    },
  ],
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next('/login')
  } else {
    next()
  }
})

export default router