import { useState, useEffect } from 'react'
import { View, Text, Image, Button } from '@tarojs/components'
import Taro from '@tarojs/taro'
import { getCart, updateCart, removeCart, createOrder } from '../../api/product'
import './index.css'

export default function Cart() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadCart()
  }, [])

  const loadCart = async () => {
    setLoading(true)
    try {
      const res = await getCart()
      setItems(res.items || [])
    } catch (e) {
      Taro.showToast({ title: '请先登录', icon: 'none' })
    }
    setLoading(false)
  }

  const handleQuantityChange = async (id, qty) => {
    await updateCart(id, { quantity: qty })
    loadCart()
  }

  const handleRemove = async (id) => {
    await removeCart(id)
    loadCart()
  }

  const total = items.reduce((sum, item) => sum + item.product_price * item.quantity, 0)

  return (
    <View className="container">
      <View className="cart-list">
        {items.map(item => (
          <View className="cart-item" key={item.id}>
            <Image className="item-img" src={item.main_image || 'https://placehold.co/80x80?text=图'} mode="aspectFill" />
            <View className="item-info">
              <Text className="item-name">{item.product_name}</Text>
              <Text className="item-price">¥{item.product_price}</Text>
              <View className="qty-row">
                <Button size="mini" onClick={() => handleQuantityChange(item.id, item.quantity - 1)}>-</Button>
                <Text>{item.quantity}</Text>
                <Button size="mini" onClick={() => handleQuantityChange(item.id, item.quantity + 1)}>+</Button>
              </View>
            </View>
            <Button size="mini" type="warn" onClick={() => handleRemove(item.id)}>删除</Button>
          </View>
        ))}
      </View>
      {items.length > 0 && (
        <View className="footer">
          <Text className="total">合计: ¥{total.toFixed(2)}</Text>
          <Button type="primary" onClick={() => Taro.navigateTo({ url: '/pages/order/confirm' })}>去结算</Button>
        </View>
      )}
    </View>
  )
}
