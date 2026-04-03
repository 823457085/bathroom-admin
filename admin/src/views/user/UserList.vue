<template>
  <div>
    <h2 style="margin-bottom: 16px">用户管理</h2>
    <el-card>
      <el-table :data="users" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="phone" label="手机号" width="150" />
        <el-table-column prop="nickname" label="昵称" width="150" />
        <el-table-column prop="created_at" label="注册时间" width="180" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import request from '../../api/index'

const users = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    const res = await request.get('/users')
    users.value = res.users || []
  } catch (e) {
    // placeholder - backend user list API not implemented yet
    users.value = []
  } finally {
    loading.value = false
  }
})
</script>
