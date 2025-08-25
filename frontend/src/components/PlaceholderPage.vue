<template>
  <div class="placeholder-page">
    <h1>{{ title }}</h1>
    <p>This page is under construction.</p>
    <el-button @click="$router.go(-1)">Go Back</el-button>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: 'Page'
  }
})
</script>

<style scoped>
.placeholder-page {
  text-align: center;
  padding: 100px 20px;
}
</style>