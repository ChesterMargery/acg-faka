<template>
  <div class="home-page">
    <el-container>
      <el-header height="60px">
        <nav class="navbar">
          <div class="nav-brand">
            <h2>ACG-FAKA</h2>
          </div>
          <div class="nav-links">
            <el-button @click="$router.push('/login')" type="primary">登录</el-button>
            <el-button @click="$router.push('/register')" type="default">注册</el-button>
          </div>
        </nav>
      </el-header>
      
      <el-main>
        <div class="hero-section">
          <h1>欢迎来到 ACG-FAKA</h1>
          <p>专业的虚拟发卡系统</p>
        </div>
        
        <div class="categories-section" v-if="categories.length">
          <h2>商品分类</h2>
          <el-row :gutter="20">
            <el-col :span="6" v-for="category in categories" :key="category.id">
              <el-card @click="$router.push(`/category/${category.id}`)" class="category-card">
                <div class="category-icon">
                  <el-icon v-if="category.icon" :name="category.icon" />
                  <el-icon v-else><Shop /></el-icon>
                </div>
                <h3>{{ category.name }}</h3>
              </el-card>
            </el-col>
          </el-row>
        </div>
        
        <div class="commodities-section" v-if="commodities.length">
          <h2>热门商品</h2>
          <el-row :gutter="20">
            <el-col :span="6" v-for="commodity in commodities" :key="commodity.id">
              <el-card @click="$router.push(`/commodity/${commodity.id}`)" class="commodity-card">
                <img v-if="commodity.cover" :src="commodity.cover" class="commodity-cover" />
                <div class="commodity-content">
                  <h3>{{ commodity.name }}</h3>
                  <p class="commodity-desc">{{ commodity.description }}</p>
                  <div class="commodity-price">
                    <span class="price">¥{{ commodity.price }}</span>
                    <span v-if="commodity.vip_price" class="vip-price">VIP价: ¥{{ commodity.vip_price }}</span>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { categoryAPI, commodityAPI } from '@/api'
import { Shop } from '@element-plus/icons-vue'

const categories = ref<any[]>([])
const commodities = ref<any[]>([])

const loadCategories = async () => {
  try {
    const response = await categoryAPI.getCategories({ status: 1 })
    categories.value = response.data || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

const loadCommodities = async () => {
  try {
    const response = await commodityAPI.getCommodities({ status: 1, limit: 8 })
    commodities.value = response.data || []
  } catch (error) {
    console.error('Failed to load commodities:', error)
  }
}

onMounted(() => {
  loadCategories()
  loadCommodities()
})
</script>

<style scoped>
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 100%;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.nav-brand h2 {
  color: #409eff;
}

.hero-section {
  text-align: center;
  padding: 80px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  margin-bottom: 40px;
}

.hero-section h1 {
  font-size: 3rem;
  margin-bottom: 16px;
}

.hero-section p {
  font-size: 1.2rem;
  opacity: 0.9;
}

.categories-section, .commodities-section {
  padding: 0 20px;
  margin-bottom: 40px;
}

.categories-section h2, .commodities-section h2 {
  margin-bottom: 20px;
  color: #333;
}

.category-card {
  cursor: pointer;
  text-align: center;
  padding: 20px;
  transition: transform 0.3s;
}

.category-card:hover {
  transform: translateY(-5px);
}

.category-icon {
  font-size: 2rem;
  margin-bottom: 10px;
  color: #409eff;
}

.commodity-card {
  cursor: pointer;
  transition: transform 0.3s;
}

.commodity-card:hover {
  transform: translateY(-5px);
}

.commodity-cover {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.commodity-content {
  padding: 16px;
}

.commodity-desc {
  color: #666;
  margin: 8px 0;
  font-size: 14px;
}

.commodity-price {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price {
  color: #e74c3c;
  font-weight: bold;
  font-size: 18px;
}

.vip-price {
  color: #f39c12;
  font-size: 14px;
}
</style>