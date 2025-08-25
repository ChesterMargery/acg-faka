<template>
  <div id="admin-app">
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'

const authStore = useAuthStore()

onMounted(() => {
  // Try to get admin profile if logged in
  if (authStore.isLoggedIn && !authStore.user) {
    authStore.getProfile().catch(() => {
      // If profile fetch fails, logout
      authStore.logout()
    })
  }
})
</script>

<style>
#admin-app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  min-height: 100vh;
  background: #f0f2f5;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
</style>