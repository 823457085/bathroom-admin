<template>
  <div>
    <el-card>
      <template #header>
        <span>{{ isEdit ? '编辑商品' : '新增商品' }}</span>
      </template>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="请选择分类">
            <el-option v-for="c in categories" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number v-model="form.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="库存" prop="stock">
          <el-input-number v-model="form.stock" :min="0" />
        </el-form-item>
        <el-form-item label="主图URL">
          <el-input v-model="form.main_image" placeholder="请输入主图URL" />
        </el-form-item>
        <el-form-item label="副标题">
          <el-input v-model="form.subtitle" placeholder="请输入副标题" />
        </el-form-item>
        <el-form-item label="商品描述">
          <el-input v-model="form.description" type="textarea" :rows="4" placeholder="请输入商品描述" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving" @click="handleSave">{{ isEdit ? '保存' : '创建' }}</el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCategories, getProduct, createProduct, updateProduct } from '../../api/product'

const route = useRoute()
const router = useRouter()
const isEdit = !!route.params.id

const formRef = ref()
const saving = ref(false)
const categories = ref([])
const form = ref({
  name: '',
  category_id: null,
  price: 0,
  stock: 0,
  main_image: '',
  subtitle: '',
  description: ''
})
const rules = {
  name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }],
  price: [{ required: true, message: '请输入价格', trigger: 'blur' }]
}

onMounted(async () => {
  const cats = await getCategories()
  categories.value = cats.categories || []
  if (isEdit) {
    const p = await getProduct(route.params.id)
    form.value = {
      name: p.name,
      category_id: p.category_id,
      price: p.price,
      stock: p.stock,
      main_image: p.main_image,
      subtitle: p.subtitle,
      description: p.description
    }
  }
})

const handleSave = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    if (isEdit) {
      await updateProduct(route.params.id, form.value)
      ElMessage.success('保存成功')
    } else {
      await createProduct(form.value)
      ElMessage.success('创建成功')
    }
    router.push('/products')
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}
</script>
