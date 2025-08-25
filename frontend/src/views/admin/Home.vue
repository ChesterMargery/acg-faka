<template>
  <div class="admin-home">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon color="#409eff"><User /></el-icon>
            </div>
            <div class="stat-text">
              <h3>0</h3>
              <p>用户总数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon color="#67c23a"><ShoppingBag /></el-icon>
            </div>
            <div class="stat-text">
              <h3>{{ stats.commodities }}</h3>
              <p>商品总数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon color="#e6a23c"><CreditCard /></el-icon>
            </div>
            <div class="stat-text">
              <h3>{{ stats.cards }}</h3>
              <p>卡密总数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <el-icon color="#f56c6c"><Document /></el-icon>
            </div>
            <div class="stat-text">
              <h3>0</h3>
              <p>订单总数</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>最新商品</span>
          </template>
          <el-table :data="recentCommodities" style="width: 100%">
            <el-table-column prop="name" label="商品名称" />
            <el-table-column prop="price" label="价格">
              <template #default="scope">
                ¥{{ scope.row.price }}
              </template>
            </el-table-column>
            <el-table-column prop="create_time" label="创建时间">
              <template #default="scope">
                {{ formatDate(scope.row.create_time) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>系统信息</span>
          </template>
          <div class="system-info">
            <p><strong>系统版本:</strong> ACG-FAKA v2.0</p>
            <p><strong>技术栈:</strong> Go + Vue3 + SQLite</p>
            <p><strong>数据库:</strong> SQLite</p>
            <p><strong>服务状态:</strong> <el-tag type="success">运行中</el-tag></p>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { adminAPI } from '@/api'
import { User, ShoppingBag, CreditCard, Document } from '@element-plus/icons-vue'

const stats = ref({
  users: 0,
  commodities: 0,
  cards: 0,
  orders: 0,
})

const recentCommodities = ref<any[]>([])

const loadStats = async () => {
  try {
    // Load commodities count
    const commoditiesRes: any = await adminAPI.commodities.getAll({ limit: 1 })
    stats.value.commodities = commoditiesRes.total || 0
    
    // Load cards count
    const cardsRes: any = await adminAPI.cards.getAll({ limit: 1 })
    stats.value.cards = cardsRes.total || 0
    
    // Load recent commodities
    const recentRes: any = await adminAPI.commodities.getAll({ limit: 5 })
    recentCommodities.value = recentRes.data || []
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.stat-card {
  margin-bottom: 20px;
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  font-size: 40px;
  margin-right: 20px;
}

.stat-text h3 {
  margin: 0 0 8px 0;
  font-size: 24px;
  color: #333;
}

.stat-text p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.system-info p {
  margin: 10px 0;
  color: #666;
}

.system-info strong {
  color: #333;
}
</style>