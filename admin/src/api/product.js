import request from './index'

export const getProducts = (params) => request.get('/products', { params })
export const getProduct = (id) => request.get(`/products/${id}`)
export const createProduct = (data) => request.post('/products', data)
export const updateProduct = (id, data) => request.put(`/products/${id}`, data)
export const deleteProduct = (id) => request.delete(`/products/${id}`)

export const getCategories = () => request.get('/categories')
