<template>
  <div>
    <h2 style="margin-bottom: 16px">订单管理</h2>
    <el-card>
      <el-table :data="orders" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="order_no" label="订单号" width="180" />
        <el-table-column prop="total_amount" label="金额" width="120">
          <template #default="{ row }">¥{{ row.total_amount }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="statusMap[row.status]">{{ statusText[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" />
        <el-table-column prop="created_at" label="下单时间" width="180" />
      </el-table>
      <el-pagination
        style="margin-top: 16px"
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadData"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getOrders } from '../../api/order'

const orders = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const statusMap = { 1: 'warning', 2: 'info', 3: 'success' }
const statusText = { 1: '待付款', 2: '已取消', 3: '已完成' }

const loadData = async () => {
  loading.value = true
  try {
    const res = await getOrders({ page: page.value, page_size: pageSize.value })
    orders.value = res.orders || []
    total.value = res.total || 0
  } catch (e) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>
