# Bathroom Admin 卫浴后台管理系统

## 项目简介

卫浴 B2C 电商一期 MVP，包含：
- **后端**：Go + Gin + MySQL
- **PC 后台**：Vue3 + Element Plus
- **小程序演示版**：Taro + React

## 技术栈

| 模块 | 技术 |
|------|------|
| 后端 | Go 1.21 + Gin + GORM + MySQL 8 |
| 后台前端 | Vue3 + Vite + Element Plus |
| 小程序 | Taro + React + NutUI |
| 数据库 | MySQL 8.0 |
| 认证 | JWT |

## 快速启动

### 后端

```bash
cd backend
cp config.yaml.example config.yaml
# 修改 config.yaml 中的数据库配置

# 初始化数据库（如还未建表）
mysql -u bathroom -p bathroom_admin < migrations/001_init.sql

# 启动服务
go run cmd/server/main.go
```

服务默认运行在 http://localhost:8080

### 一期里程碑

- [ ] Issue #1 - 项目初始化（仓库结构 + Go 骨架 + MySQL 建表）**进行中**
- [ ] Issue #2 - 用户模块：注册 + 登录 + JWT 鉴权
- [ ] Issue #3 - 商品模块：分类 + 商品列表 + 商品详情
- [ ] Issue #4 - 订单模块：购物车 + 下单 + 订单列表
- [ ] Issue #5 - Vue3 后台管理前端
- [ ] Issue #6 - Taro 小程序演示版

## 开发规范

- 提交信息格式：`[Issue #N] 简短描述`
- API 路由：`/api/v1/{module}/{action}`
- 分支命名：`feature/{issue-number}-{description}`

## License

MIT
