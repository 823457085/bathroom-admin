import request from './index'

export const getOrders = (params) => request.get('/orders', { params })
export const getOrder = (id) => request.get(`/orders/${id}`)
