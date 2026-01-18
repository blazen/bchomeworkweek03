# Gin 项目实战示例

这是一个完整的用户管理系统 API 示例，展示了如何使用 Gin 框架构建 RESTful API。

## 技术栈

- **Web 框架**: Gin v1.10.0
- **ORM**: GORM v1.25.12
- **数据库**: SQLite
- **认证**: JWT (github.com/golang-jwt/jwt/v5)
- **密码加密**: bcrypt (golang.org/x/crypto)

## 功能特性

- ✅ 统一响应格式
- ✅ 错误处理
- ✅ 中间件（日志、CORS、认证）
- ✅ 用户注册、登录（JWT 认证）、查询、更新
- ✅ 用户文章数统计（废弃AfterCreate，改为Transaction）
- ✅ 文章CURD
- ✅ 文章评论数统计（废弃AfterCreate、AfterDelete，改为Transaction），评论数为0时，文章评论状态显示：无评论
- ✅ 评论CURD

## 问题
- 1 废弃AfterCreate、AfterDelete，改为Transaction 保持一致性。钩子使用 context.WithValue 
- 2 

## 项目结构

```
project/
├── main.go              # 程序入口
├── config/              # 配置管理
│   └── config.go
├── handlers/            # 处理器（Controller）
│   |── comment_handler.go
│   |── post_handler.go
│   └── user_handler.go
├── middleware/          # 中间件
│   ├── auth.go
│   ├── logger.go
│   └── cors.go
├── models/              # 数据模型
│   ├── comment.go
│   ├── post.go
│   └── user.go
├── services/            # 业务逻辑层
│   ├── comment_service.go
│   ├── post_service.go
│   └── user_service.go
└── utils/               # 工具函数
    ├── errors.go
    ├── generate.go
    ├── jwt.go
    ├── page.go
    └── response.go
```

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 运行项目

```bash
go run main.go
```

服务器将在 `http://0.0.0.0:8080` 启动。

### 3. API 端点

| 类别 | 方法 | 路径 | 说明 | 认证 | 参数 |
|------|------|------|------|------|------|
| 通用 | GET | `/health` | 健康检查 | 否 | 无 |
| 用户 | POST | `/api/v1/users/register` | 用户注册 | 否 | JSON |
| - | POST | `/api/v1/users/login` | 用户登录 | 否 | JSON |
| - | GET | `/api/v1/users/me` | 获取登录用户信息 | 是 | 无 |
| - | PUT | `/api/v1/users/me` | 更新登录用户信息 | 是 | JSON |
| 文章 | GET | `/api/v1/posts` | 查询所有用户的文章 | 否 | Query |
| - | GET | `/api/v1/posts/:id` | 主键查询文章 | 否 | URL |
| - | POST | `/api/v1/posts/condition` | 条件查询文章 | 否 | JSON |
| - | POST | `/api/v1/posts/me` | 创建文章 | 是 | JSON |
| - | GET | `/api/v1/posts/me` | 查询登录用户的全部文章 | 是 | 无 |
| - | PUT | `/api/v1/posts/me` | 更新文章 | 是 | JSON |
| - | DELETE | `/api/v1/posts/me/:id` | 删除文章 | 是 | URL |
| 评论 | GET | `/api/v1/comments/:postId` | 查询文章的评论 | 否 | URL |
| - | POST | `/api/v1/comments` | 创建文章的评论 | 否 | JSON |
| - | DELETE | `/api/v1/comments/me/:postId/:id` | 删除文章的评论 | 否 | URL |


### 4. API 示例

#### 健康检查 

```bash
curl http://localhost:8080/health
```

#### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

#### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

#### 获取登录用户信息

```bash
curl http://localhost:8080/api/v1/users/me\
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 更新登录用户信息

```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "email": "newadmin@example.com"
  }'
```


## 注意事项

1. 数据库文件 `users.db` 会在首次运行时自动创建在项目根目录
2. JWT Secret 在生产环境中应该使用环境变量配置
3. 密码使用 bcrypt 加密存储
4. 所有 API 返回统一的 JSON 格式
5. 项目已配置独立的 `go.mod` 文件，可直接在 `project` 目录下运行

## 扩展建议

- [ ] 添加 Swagger API 文档
- [ ] 添加单元测试和集成测试
- [ ] 实现分页查询
- [ ] 添加日志系统（如 logrus）
- [ ] 使用 Viper 进行配置管理
- [ ] 添加 Redis 缓存
- [ ] 实现文件上传功能

