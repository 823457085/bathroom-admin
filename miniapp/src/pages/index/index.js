import { Component } from 'react'
import { View, Text, ScrollView, Image } from '@tarojs/components'
import { navigateTo } from '@tarojs/taro'
import { getProducts } from '../../api/product'
import './index.css'

export default class Index extends Component {
  state = { products: [], categories: [], activeCategory: 0 }

  componentDidMount() {
    this.loadProducts()
  }

  async loadProducts(categoryId) {
    const res = await getProducts({ page: 1, page_size: 20, category_id: categoryId })
    this.setState({ products: res.products || [] })
  }

  render() {
    const { products } = this.state
    return (
      <View className="container">
        <View className="search-bar">
          <Text className="search-tip" onClick={() => navigateTo({ url: '/pages/product/index' })}>搜索商品</Text>
        </View>
        <ScrollView scroll-y className="content">
          <View className="banner">
            <View className="banner-placeholder">卫浴精选</View>
          </View>
          <View className="section-title">热门商品</View>
          <View className="product-grid">
            {products.map(p => (
              <View className="product-card" key={p.id} onClick={() => navigateTo({ url: `/pages/detail/index?id=${p.id}` })}>
                <Image className="product-img" src={p.main_image || 'https://placehold.co/200x200?text=暂无图片'} mode="aspectFill" />
                <View className="product-info">
                  <Text className="product-name">{p.name}</Text>
                  <Text className="product-price">¥{p.price}</Text>
                </View>
              </View>
            ))}
          </View>
        </ScrollView>
      </View>
    )
  }
}
