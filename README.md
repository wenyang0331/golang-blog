# Homework Blog API

一个使用Go语言、Gin框架和GORM库开发的个人博客系统后端API。

## 功能特性

- 用户注册和登录
- JWT认证和授权
- 文章的CRUD操作
- 评论功能
- 统一的错误处理
- 日志记录
- CORS跨域支持
- 数据库自动迁移
- 基于Viper的灵活配置管理

## 技术栈

- **Go 1.26**
- **Gin** - Web框架
- **GORM** - ORM库
- **SQLite** - 数据库
- **JWT (golang-jwt/jwt/v5)** - 身份认证
- **Viper** - 配置管理
- **bcrypt (golang.org/x/crypto)** - 密码加密

## 项目结构

```
homework_blog/
├── main.go                    # 程序入口
├── config/
│   ├── config.go              # 配置结构体和加载逻辑
│   └── config.yaml            # YAML配置文件
├── database/
│   ├── database.go            # 数据库初始化
│   ├── post_db.go             # 文章数据库操作
│   └── user_db.go             # 用户数据库操作
├── handlers/
│   ├── user_handler.go        # 用户处理器（注册、登录、个人中心）
│   ├── post_handler.go        # 文章处理器（CRUD）
│   └── comment_handler.go     # 评论处理器
├── middleware/
│   ├── auth.go                # JWT认证中间件
│   ├── cors.go                # CORS跨域中间件
│   └── logger.go              # 日志中间件
├── models/
│   ├── user.go                # 用户模型
│   ├── post.go                # 文章模型
│   └── comment.go             # 评论模型
├── routes/
│   └── routes.go              # 路由配置
├── utils/
│   ├── errors.go              # 统一错误处理
│   ├── jwt.go                 # JWT工具函数
│   └── response.go            # 统一响应格式
├── go.mod
├── go.sum
└── README.md
```

## 数据库设计

### Users 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint (主键) | 用户ID |
| username | string (唯一, 非空) | 用户名 |
| email | string (唯一, 非空) | 邮箱 |
| password | string (非空) | 加密密码 |
| created_at | time.Time | 创建时间 |
| updated_at | time.Time | 更新时间 |
| deleted_at | gorm.DeletedAt | 软删除时间 |

### Posts 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint (主键) | 文章ID |
| title | string (非空) | 标题 |
| content | text (非空) | 内容 |
| user_id | uint (非空, 外键) | 作者ID |
| created_at | time.Time | 创建时间 |
| updated_at | time.Time | 更新时间 |
| deleted_at | gorm.DeletedAt | 软删除时间 |

### Comments 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint (主键) | 评论ID |
| content | text | 评论内容 |
| user_id | uint (非空, 外键) | 评论者ID |
| post_id | uint (非空, 外键) | 文章ID |
| created_at | time.Time | 创建时间 |
| updated_at | time.Time | 更新时间 |
| deleted_at | gorm.DeletedAt | 软删除时间 |

## API 接口

### 公开接口（无需认证）

#### 认证相关
- `POST /api/v1/public/register` - 用户注册
- `POST /api/v1/public/login` - 用户登录

#### 文章浏览
- `GET /api/v1/public/posts` - 获取文章列表（支持分页）
- `GET /api/v1/public/posts/:id` - 获取文章详情

#### 评论浏览
- `GET /api/v1/comments/post/:post_id` - 获取文章评论列表

### 需要认证的接口（JWT Token）

#### 用户信息
- `GET /api/v1/profile` - 获取当前用户信息

#### 文章管理
- `POST /api/v1/posts` - 创建文章
- `PUT /api/v1/posts/:id` - 更新文章（仅作者）
- `DELETE /api/v1/posts/:id` - 删除文章（仅作者）

#### 评论管理
- `POST /api/v1/posts/:post_id/comments` - 创建评论

### 其他接口
- `GET /health` - 健康检查

## 运行项目

### 环境要求

- Go 1.21+

### 安装依赖

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite
go get github.com/spf13/viper
go get -u github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt/v5
```

或直接运行：

```bash
go mod tidy
```

### 配置说明

编辑 `config/config.yaml` 文件，修改数据库连接、服务器端口、JWT密钥等配置：

```yaml
server:
  port: "8080"
  host: "0.0.0.0"
  mode: "debug"  # debug, release, test

jwt:
  secret: "your-secret-key"
  expire: "24h"
```

> 也可通过环境变量覆盖配置，环境变量前缀为 `APP_`（如 `APP_SERVER_PORT=9090`）。

### 启动服务

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。

## 使用示例

### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/public/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/public/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 创建文章

```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "我的第一篇文章",
    "content": "这是文章内容..."
  }'
```

### 发表评论

```bash
curl -X POST http://localhost:8080/api/v1/posts/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "content": "写得很棒！"
  }'
```

## 注意事项

- JWT 密钥在生产环境中应该从环境变量读取或使用强密钥
- 所有密码使用 bcrypt 进行加密存储，通过 GORM 的 `BeforeCreate` 钩子自动完成
- API 返回统一的 JSON 格式响应
- 支持软删除数据（使用 GORM 的 `DeletedAt` 字段）
