# 个人博客系统API

这是一个基于Go语言、Gin框架和GORM库开发的个人博客系统后端API。

## 功能特性

- 用户注册和登录（JWT认证）
- 文章的创建、读取、更新和删除（CRUD）操作
- 评论功能
- 权限控制（只有作者可以修改/删除自己的文章）

## 技术栈

- Go 1.16+
- Gin Web框架
- GORM ORM库
- MySQL数据库
- JWT认证

## 项目结构

```
.
├── config/         # 配置文件
├── controllers/    # 控制器
├── middleware/     # 中间件
├── models/         # 数据模型
├── routes/         # 路由
├── utils/          # 工具函数
├── main.go         # 入口文件
└── README.md       # 项目说明
```

## 安装与运行

### 前提条件

- Go 1.16+
- MySQL

### 安装步骤

1. 克隆项目

```bash
git clone https://github.com/yourusername/blog-api.git
cd blog-api
```

2. 安装依赖

```bash
go mod tidy
```

3. 配置数据库

修改 `config/config.go` 文件中的数据库配置：

```go
DatabaseConfig{
    Driver:   "mysql",
    Host:     "localhost",
    Port:     "3306",
    Username: "your_username",
    Password: "your_password",
    DBName:   "blog_db2",
    Charset:  "utf8mb4",
}
```

4. 运行项目

```bash
go run main.go
```

服务器将在 http://localhost:8090 上运行。

### 初始用户和数据

系统启动时会自动创建以下示例数据：

1. 用户账号：
   - 管理员：用户名 `admin`，密码 `admin123`
   - 普通用户：用户名 `user`，密码 `user123`

2. 示例文章：
   - 《欢迎使用个人博客系统》- 由管理员发布
   - 《Go语言学习笔记》- 由普通用户发布

3. 示例评论：
   - 每篇文章下有一条示例评论

## API接口

### 用户认证

- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录

### 文章管理

- `GET /api/posts` - 获取所有文章
- `GET /api/posts/:id` - 获取单个文章
- `POST /api/posts` - 创建文章（需要认证）
- `PUT /api/posts/:id` - 更新文章（需要认证和授权）
- `DELETE /api/posts/:id` - 删除文章（需要认证和授权）

### 评论管理

- `GET /api/posts/:id/comments` - 获取文章的所有评论
- `POST /api/posts/:id/comments` - 创建评论（需要认证）

## 测试

使用Postman或其他API测试工具测试接口。

### 示例测试流程

1. 使用初始用户登录获取JWT令牌：
   - 请求：`POST /api/login`
   - 请求体：`{"username": "admin", "password": "admin123"}`
   - 响应中获取token

2. 使用JWT令牌创建文章：
   - 请求：`POST /api/posts`
   - 请求头：`Authorization: Bearer <你的token>`
   - 请求体：`{"title": "测试文章", "content": "这是一篇测试文章"}`

3. 获取文章列表：
   - 请求：`GET /api/posts`

4. 获取单个文章详情：
   - 请求：`GET /api/posts/1`

5. 更新文章：
   - 请求：`PUT /api/posts/1`
   - 请求头：`Authorization: Bearer <你的token>`
   - 请求体：`{"title": "更新后的标题", "content": "更新后的内容"}`

6. 对文章发表评论：
   - 请求：`POST /api/posts/1/comments`
   - 请求头：`Authorization: Bearer <你的token>`
   - 请求体：`{"content": "这是一条评论"}`

7. 获取文章的评论列表：
   - 请求：`GET /api/posts/1/comments`

8. 删除文章：
   - 请求：`DELETE /api/posts/1`
   - 请求头：`Authorization: Bearer <你的token>`

## 许可证

MIT 