import { useState, useEffect } from 'react'
import { View, Text, Button } from '@tarojs/components'
import Taro from '@tarojs/taro'
import { getCart, createOrder, getAddresses, createAddress } from '../../api/product'
import './index.css'

export default function OrderConfirm() {
  const [cart, setCart] = useState([])
  const [address, setAddress] = useState(null)
  const [addresses, setAddresses] = useState([])
  const [showAddressPicker, setShowAddressPicker] = useState(false)
  const [newAddr, setNewAddr] = useState({ receiver_name: '', phone: '', province: '', city: '', district: '', detail: '' })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([getCart(), getAddresses()]).then(([c, a]) => {
      setCart(c.items || [])
      setAddresses(a.addresses || [])
      if (a.addresses?.length > 0) setAddress(a.addresses[0])
      setLoading(false)
    })
  }, [])

  const total = cart.reduce((sum, item) => sum + item.product_price * item.quantity, 0)

  const handleCreateOrder = async () => {
    if (!address) {
      Taro.showToast({ title: '请选择收货地址', icon: 'none' })
      return
    }
    try {
      const res = await createOrder({ address_id: address.id })
      Taro.redirectTo({ url: `/pages/order/list?success=1` })
    } catch (e) {
      Taro.showToast({ title: '下单失败', icon: 'none' })
    }
  }

  const handleAddAddress = async () => {
    await createAddress(newAddr)
    getAddresses().then(([c, a]) => {
      setAddresses(a.addresses || [])
      if (a.addresses?.length > 0) setAddress(a.addresses[0])
    })
    setShowAddressPicker(false)
  }

  if (loading) return <View>加载中...</View>

  return (
    <View className="container">
      <View className="section">
        <Text className="section-title">收货地址</Text>
        {address ? (
          <View className="address-card">
            <Text>{address.receiver_name} {address.phone}</Text>
            <Text>{address.province}{address.city}{address.district}{address.detail}</Text>
          </View>
        ) : (
          <View className="address-tip">暂无收货地址，请添加</View>
        )}
        <Button size="mini" onClick={() => setShowAddressPicker(!showAddressPicker)}>+ 添加地址</Button>
        {showAddressPicker && (
          <View className="addr-form">
            {['receiver_name', 'phone', 'province', 'city', 'district', 'detail'].map(field => (
              <View key={field} className="addr-field">
                <Text>{field}:</Text>
                <input className="addr-input" value={newAddr[field]} onInput={e => setNewAddr({ ...newAddr, [field]: e.detail.value })} />
              </View>
            ))}
            <Button size="mini" type="primary" onClick={handleAddAddress}>保存</Button>
          </View>
        )}
      </View>

      <View className="section">
        <Text className="section-title">商品清单</Text>
        {cart.map(item => (
          <View className="order-item" key={item.id}>
            <Text>{item.product_name} x{item.quantity}</Text>
            <Text>¥{(item.product_price * item.quantity).toFixed(2)}</Text>
          </View>
        ))}
      </View>

      <View className="footer">
        <Text className="total">合计: ¥{total.toFixed(2)}</Text>
        <Button type="primary" onClick={handleCreateOrder}>提交订单</Button>
      </View>
    </View>
  )
}
