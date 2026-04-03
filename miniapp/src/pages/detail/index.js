import { useState, useEffect } from 'react'
import { View, Text, Image, Button } from '@tarojs/components'
import Taro from '@tarojs/taro'
import { getProduct, addCart } from '../../api/product'
import './index.css'

export default function Detail() {
  const [product, setProduct] = useState(null)
  const [quantity, setQuantity] = useState(1)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const { id } = Taro.getCurrentInstance().router.params
    getProduct(id).then(res => {
      setProduct(res)
      setLoading(false)
    })
  }, [])

  const handleAddCart = async () => {
    await addCart({ product_id: product.id, quantity })
    Taro.showToast({ title: '已加入购物车' })
  }

  if (loading) return <View>加载中...</View>
  if (!product) return <View>商品不存在</View>

  return (
    <View className="container">
      <Image className="main-img" src={product.main_image || 'https://placehold.co/375x400?text=暂无图片'} mode="widthFix" />
      <View className="info">
        <Text className="price">¥{product.price}</Text>
        <Text className="name">{product.name}</Text>
        <Text className="subtitle">{product.subtitle}</Text>
      </View>
      <View className="footer">
        <View className="qty">
          <Button size="mini" onClick={() => setQuantity(Math.max(1, quantity - 1))}>-</Button>
          <Text>{quantity}</Text>
          <Button size="mini" onClick={() => setQuantity(quantity + 1)}>+</Button>
        </View>
        <Button className="add-btn" type="primary" onClick={handleAddCart}>加入购物车</Button>
      </View>
    </View>
  )
}
