import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../store/auth'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../views/login/Login.vue') },
  {
    path: '/',
    component: () => import('../views/layout/Layout.vue'),
    children: [
      { path: '', redirect: '/dashboard' },
      { path: 'dashboard', name: 'Dashboard', component: () => import('../views/dashboard/Dashboard.vue') },
      { path: 'products', name: 'Products', component: () => import('../views/product/ProductList.vue') },
      { path: 'products/create', name: 'ProductCreate', component: () => import('../views/product/ProductForm.vue') },
      { path: 'products/:id/edit', name: 'ProductEdit', component: () => import('../views/product/ProductForm.vue') },
      { path: 'orders', name: 'Orders', component: () => import('../views/order/OrderList.vue') },
      { path: 'users', name: 'Users', component: () => import('../views/user/UserList.vue') },
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()
  if (to.path !== '/login' && !auth.token) {
    next('/login')
  } else {
    next()
  }
})

export default router
