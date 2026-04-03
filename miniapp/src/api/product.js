const BASE_URL = 'http://localhost:8080/api/v1'

const request = (path, data, method = 'GET') => {
  const token = wx.getStorageSync('token')
  const header = { 'Content-Type': 'application/json' }
  if (token) header['Authorization'] = `Bearer ${token}`
  return wx.request({ url: BASE_URL + path, data, method, header }).then(r => r.data)
}

export const getProducts = (params) => request('/products', params)
export const getProduct = (id) => request(`/products/${id}`)
export const getCategories = () => request('/categories')
export const addCart = (data) => request('/cart', data, 'POST')
export const getCart = () => request('/cart')
export const updateCart = (id, data) => request(`/cart/${id}`, data, 'PUT')
export const removeCart = (id) => request(`/cart/${id}`, {}, 'DELETE')
export const createOrder = (data) => request('/orders', data, 'POST')
export const getOrders = (params) => request('/orders', params)
export const getAddresses = () => request('/addresses')
export const createAddress = (data) => request('/addresses', data, 'POST')
export const login = (data) => request('/auth/login', data, 'POST')
export const register = (data) => request('/auth/register', data, 'POST')
