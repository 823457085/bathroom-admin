<template>
  <div>
    <h2 style="margin-bottom: 20px">仪表盘</h2>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card>
          <div style="text-align: center">
            <div style="font-size: 32px; font-weight: bold; color: #409eff">{{ stats.totalProducts }}</div>
            <div style="color: #888">商品总数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div style="text-align: center">
            <div style="font-size: 32px; font-weight: bold; color: #67c23a">{{ stats.totalOrders }}</div>
            <div style="color: #888">订单总数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div style="text-align: center">
            <div style="font-size: 32px; font-weight: bold; color: #f56c6c">{{ stats.pendingOrders }}</div>
            <div style="color: #888">待处理订单</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div style="text-align: center">
            <div style="font-size: 32px; font-weight: bold; color: #e6a23c">{{ stats.totalUsers }}</div>
            <div style="color: #888">用户总数</div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getProducts } from '../../api/product'
import { getOrders } from '../../api/order'

const stats = ref({ totalProducts: 0, totalOrders: 0, pendingOrders: 0, totalUsers: 0 })

onMounted(async () => {
  try {
    const [p, o] = await Promise.all([
      getProducts({ page: 1, page_size: 1 }),
      getOrders({ page: 1, page_size: 1 })
    ])
    stats.value.totalProducts = p.total || 0
    stats.value.totalOrders = o.total || 0
  } catch (e) {
    console.error(e)
  }
})
</script>
