<template>
  <div class="placeholder-view">
    <el-result
      icon="warning"
      title="功能开发中"
      sub-title="该功能正在开发中，敬请期待"
    >
      <template #extra>
        <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup lang="ts">
// Placeholder component for features under development
</script>

<style scoped>
.placeholder-view {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}
</style>
