<template>
  <div class="home">
    <el-container>
      <!-- Header -->
      <el-header height="80px" class="header">
        <div class="header-content">
          <div class="logo">
            <h2>{{ config.shop_name || 'ACG-FAKA' }}</h2>
          </div>
          <div class="nav-menu">
            <el-menu mode="horizontal" :default-active="activeIndex">
              <el-menu-item index="1">
                <router-link to="/">首页</router-link>
              </el-menu-item>
              <el-menu-item index="2">
                <router-link to="/shop">商店</router-link>
              </el-menu-item>
              <el-menu-item index="3" v-if="!authStore.isAuthenticated">
                <router-link to="/login">登录</router-link>
              </el-menu-item>
              <el-menu-item index="4" v-if="!authStore.isAuthenticated">
                <router-link to="/register">注册</router-link>
              </el-menu-item>
              <el-submenu index="5" v-if="authStore.isAuthenticated">
                <template #title>{{ authStore.user?.username }}</template>
                <el-menu-item index="5-1">
                  <router-link to="/profile">个人资料</router-link>
                </el-menu-item>
                <el-menu-item index="5-2">
                  <router-link to="/orders">我的订单</router-link>
                </el-menu-item>
                <el-menu-item index="5-3" @click="handleLogout">退出登录</el-menu-item>
              </el-submenu>
            </el-menu>
          </div>
        </div>
      </el-header>

      <!-- Main Content -->
      <el-main>
        <!-- Hero Section -->
        <div class="hero-section">
          <div class="hero-content">
            <h1>{{ config.title || 'ACG-FAKA - 异次元店铺系统' }}</h1>
            <p>{{ config.description || '现代化的虚拟商品销售平台' }}</p>
            <el-button type="primary" size="large" @click="$router.push('/shop')">
              开始购物
            </el-button>
          </div>
        </div>

        <!-- Features Section -->
        <div class="features-section">
          <el-row :gutter="24">
            <el-col :span="8">
              <el-card class="feature-card">
                <template #header>
                  <div class="card-header">
                    <el-icon><ShoppingBag /></el-icon>
                    <span>商品管理</span>
                  </div>
                </template>
                <p>支持多种商品类型，自动发货，库存管理</p>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card class="feature-card">
                <template #header>
                  <div class="card-header">
                    <el-icon><CreditCard /></el-icon>
                    <span>支付系统</span>
                  </div>
                </template>
                <p>支持多种支付方式，安全快捷的支付体验</p>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card class="feature-card">
                <template #header>
                  <div class="card-header">
                    <el-icon><User /></el-icon>
                    <span>用户体系</span>
                  </div>
                </template>
                <p>完整的用户等级体系，会员折扣，推广返利</p>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <!-- Popular Products -->
        <div class="products-section" v-if="recommendedProducts.length > 0">
          <h2>推荐商品</h2>
          <el-row :gutter="24">
            <el-col :span="6" v-for="product in recommendedProducts" :key="product.id">
              <el-card class="product-card" @click="$router.push(`/product/${product.id}`)">
                <img :src="product.cover || '/default-product.jpg'" class="product-image" />
                <div class="product-info">
                  <h3>{{ product.name }}</h3>
                  <p class="product-price">¥{{ product.price }}</p>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </el-main>

      <!-- Footer -->
      <el-footer height="60px" class="footer">
        <div class="footer-content">
          <p>&copy; 2024 ACG-FAKA. All rights reserved. Powered by Go + Vue.js</p>
        </div>
      </el-footer>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import { ShoppingBag, CreditCard, User } from '@element-plus/icons-vue'
import api from '@/api'

const authStore = useAuthStore()
const activeIndex = ref('1')
const config = ref({})
const recommendedProducts = ref([])

const loadConfig = async () => {
  try {
    const response = await api.get('/public/configs')
    config.value = response.data
  } catch (error) {
    console.error('Failed to load config:', error)
  }
}

const loadRecommendedProducts = async () => {
  try {
    const response = await api.get('/public/commodities?limit=4')
    recommendedProducts.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load products:', error)
  }
}

const handleLogout = () => {
  authStore.logout()
  ElMessage.success('退出登录成功')
  window.location.reload()
}

onMounted(() => {
  loadConfig()
  loadRecommendedProducts()
})
</script>

<style scoped>
.header {
  background: #fff;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.logo h2 {
  margin: 0;
  color: #409eff;
}

.nav-menu {
  flex: 1;
  display: flex;
  justify-content: flex-end;
}

.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 100px 0;
  text-align: center;
}

.hero-content h1 {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.hero-content p {
  font-size: 1.2rem;
  margin-bottom: 2rem;
}

.features-section {
  padding: 80px 0;
  background: #f8f9fa;
}

.feature-card {
  text-align: center;
  height: 200px;
}

.card-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.card-header .el-icon {
  font-size: 2rem;
  color: #409eff;
}

.products-section {
  padding: 80px 0;
  max-width: 1200px;
  margin: 0 auto;
}

.products-section h2 {
  text-align: center;
  margin-bottom: 3rem;
}

.product-card {
  cursor: pointer;
  transition: transform 0.3s;
}

.product-card:hover {
  transform: translateY(-5px);
}

.product-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.product-info {
  padding: 1rem 0;
}

.product-info h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
}

.product-price {
  color: #e74c3c;
  font-weight: bold;
  font-size: 1.2rem;
  margin: 0;
}

.footer {
  background: #333;
  color: white;
}

.footer-content {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.footer-content p {
  margin: 0;
}
</style>