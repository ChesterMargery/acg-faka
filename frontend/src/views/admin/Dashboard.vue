<template>
  <el-container class="admin-dashboard">
    <el-aside width="260px">
      <div class="logo">
        <h2>ACG-FAKA 管理</h2>
      </div>
      <el-menu
        :default-active="$route.path"
        class="el-menu-vertical"
        router
        :collapse="false"
      >
        <el-menu-item index="/">
          <el-icon><House /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/categories">
          <el-icon><Folder /></el-icon>
          <span>分类管理</span>
        </el-menu-item>
        <el-menu-item index="/commodities">
          <el-icon><ShoppingBag /></el-icon>
          <span>商品管理</span>
        </el-menu-item>
        <el-menu-item index="/cards">
          <el-icon><CreditCard /></el-icon>
          <span>卡密管理</span>
        </el-menu-item>
        <el-menu-item index="/orders">
          <el-icon><Document /></el-icon>
          <span>订单管理</span>
        </el-menu-item>
        <el-menu-item index="/users">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header>
        <div class="header-content">
          <div class="header-left">
            <h3>{{ getPageTitle() }}</h3>
          </div>
          <div class="header-right">
            <el-dropdown @command="handleCommand">
              <span class="el-dropdown-link">
                {{ user?.username }}
                <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </el-header>
      
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { ElMessage } from 'element-plus'
import { 
  House, 
  Folder, 
  ShoppingBag, 
  CreditCard, 
  Document, 
  User,
  ArrowDown 
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const user = computed(() => authStore.user)

const getPageTitle = () => {
  const titles: Record<string, string> = {
    '/': '仪表盘',
    '/categories': '分类管理',
    '/commodities': '商品管理',
    '/cards': '卡密管理',
    '/orders': '订单管理',
    '/users': '用户管理',
  }
  return titles[route.path] || '管理后台'
}

const handleCommand = (command: string) => {
  if (command === 'logout') {
    authStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  }
}
</script>

<style scoped>
.admin-dashboard {
  height: 100vh;
}

.logo {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid #e4e7ed;
  background: #001529;
  color: white;
}

.logo h2 {
  margin: 0;
  font-size: 18px;
}

.el-aside {
  background: #001529;
}

.el-menu-vertical {
  border-right: none;
  background: #001529;
}

.el-menu-vertical .el-menu-item {
  color: rgba(255, 255, 255, 0.8);
}

.el-menu-vertical .el-menu-item:hover {
  background-color: #1890ff !important;
  color: white;
}

.el-menu-vertical .el-menu-item.is-active {
  background-color: #1890ff !important;
  color: white;
}

.el-header {
  background: white;
  border-bottom: 1px solid #e4e7ed;
  padding: 0 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.header-left h3 {
  margin: 0;
  color: #333;
}

.el-dropdown-link {
  cursor: pointer;
  color: #333;
  display: flex;
  align-items: center;
}

.el-main {
  background: #f0f2f5;
  padding: 20px;
}
</style>