<template>
  <div class="categories-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>分类管理</span>
          <el-button type="primary" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            新增分类
          </el-button>
        </div>
      </template>
      
      <el-table :data="categories" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="分类名称" />
        <el-table-column prop="icon" label="图标" width="100">
          <template #default="scope">
            <el-icon v-if="scope.row.icon"><component :is="scope.row.icon" /></el-icon>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.create_time) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="editCategory(scope.row)">编辑</el-button>
            <el-popconfirm
              title="确定删除这个分类吗？"
              @confirm="deleteCategory(scope.row.id)"
            >
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingCategory ? '编辑分类' : '新增分类'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
        </el-form-item>
        
        <el-form-item label="图标" prop="icon">
          <el-input v-model="form.icon" placeholder="图标名称（可选）" />
        </el-form-item>
        
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { adminAPI } from '@/api'
import { Plus } from '@element-plus/icons-vue'

const loading = ref(false)
const dialogVisible = ref(false)
const submitting = ref(false)
const editingCategory = ref<any>(null)
const categories = ref<any[]>([])

const formRef = ref<FormInstance>()
const form = reactive({
  name: '',
  icon: '',
  sort: 0,
  status: 1,
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
  ],
}

const loadCategories = async () => {
  loading.value = true
  try {
    const response = await adminAPI.categories.getAll()
    categories.value = response.data || []
  } catch (error: any) {
    ElMessage.error(error.error || '加载分类失败')
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  editingCategory.value = null
  resetForm()
  dialogVisible.value = true
}

const editCategory = (category: any) => {
  editingCategory.value = category
  form.name = category.name
  form.icon = category.icon || ''
  form.sort = category.sort
  form.status = category.status
  dialogVisible.value = true
}

const resetForm = () => {
  form.name = ''
  form.icon = ''
  form.sort = 0
  form.status = 1
}

const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (editingCategory.value) {
          await adminAPI.categories.update(editingCategory.value.id, form)
          ElMessage.success('分类更新成功')
        } else {
          await adminAPI.categories.create(form)
          ElMessage.success('分类创建成功')
        }
        dialogVisible.value = false
        loadCategories()
      } catch (error: any) {
        ElMessage.error(error.error || '操作失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

const deleteCategory = async (id: number) => {
  try {
    await adminAPI.categories.delete(id)
    ElMessage.success('分类删除成功')
    loadCategories()
  } catch (error: any) {
    ElMessage.error(error.error || '删除失败')
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>