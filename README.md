# Zustchat

一个基于 Go 和 WebSocket 的实时聊天应用，支持单聊、群聊、AI 对话等功能。用于找实习，名字随便取的不定期更新

## 功能特性

- **用户管理**
  - 用户注册与登录
  - 用户信息查询与更新
  - 用户注销

- **实时通信**
  - WebSocket 实时消息推送
  - 单聊功能
  - 群聊功能
  - 消息历史记录（支持 Redis 缓存）

- **好友管理**
  - 添加好友
  - 查找好友
  - 好友列表加载

- **群组管理**
  - 创建群组
  - 加入群组
  - 群组列表加载

- **文件上传**
  - 支持图片、文件等附件上传

- **AI 对话**
  - 集成 Google Gemini AI，支持 AI 对话功能

- **在线状态管理**
  - 心跳检测机制
  - 在线用户状态缓存（Redis）

## 技术栈

### 后端

- **Go 1.24+** - 编程语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **MySQL** - 关系型数据库
- **Redis** - 缓存和在线状态管理
- **WebSocket (Gorilla)** - 实时通信
- **JWT** - 身份认证
- **Zap** - 结构化日志
- **Swagger** - API 文档

### 前端

- **HTML5** - 页面结构
- **jQuery** - JavaScript 库
- **MUI** - 移动端 UI 框架
- **WebSocket** - 实时通信客户端

## 项目结构

```
Zustchat/
├── api/v1/          # API 路由处理层
├── asset/           # 静态资源（CSS、JS、图片等）
├── common/          # 公共组件（JWT、追踪等）
├── conf/            # 配置文件
├── docs/            # Swagger 文档
├── global/          # 全局变量
├── model/           # 数据模型
├── router/          # 路由配置
├── service/         # 业务逻辑层
├── utils/           # 工具函数
└── views/           # HTML 模板文件
```

## 快速开始

### 环境要求

- Go 1.24+
- MySQL 8.0+
- Redis 6.0+

### 安装步骤

1. **克隆项目**

```bash
git clone https://github.com/your-org/Zustchat.git
cd Zustchat
```

2. **安装依赖**

```bash
go mod download
```

3. **配置数据库**

创建 MySQL 数据库：

```sql
CREATE DATABASE Zustchat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

4. **修改配置文件**

编辑 `conf/app.ini` 文件，配置数据库和 Redis 连接信息：

```ini
[app]
RUN_MODE = debug
GEMINI_KEY = your_gemini_api_key
JWT_SECRET = your_jwt_secret

[server]
HTTP_PORT = 8081
READ_TIMEOUT = 60
WRITE_TIMEOUT = 60

[mysql]
USER = root
PASSWORD = your_password
HOST = 127.0.0.1
PORT = 3306
DBNAME = gin_chat
CHARSET = utf8mb4

[redis]
ADDR = localhost:6379
PASSWORD = 
DB = 0

[timer]
DelayHeartbeat = 3
HeartbeatHz = 30
HeartbeatMaxTime = 30
RedisOnlineTime = 4
```

5. **运行项目**

```bash
go run main.go
```

6. **访问应用**

- Web 界面：http://localhost:8081/index
<!-- - API 文档：http://localhost:8081/swagger/index.html -->

## API 接口

### 用户相关

- `POST /login` - 用户登录
- `POST /register` - 用户注册
- `GET /user/getUserList` - 获取用户列表
- `POST /user/findUser` - 查找用户
- `POST /user/updateUserInfo` - 更新用户信息
- `DELETE /user/deleteUser` - 删除用户

### 好友相关

- `POST /findFriends` - 加载好友列表
- `POST /addFriend` - 添加好友
- `GET /findFriend` - 查找好友

### 群组相关

- `POST /group/createGroup` - 创建群组
- `POST /group/joinGroup` - 加入群组
- `POST /group/loadGroups` - 加载群组列表

### 消息相关

- `GET /chat` - WebSocket 连接端点
- `POST /user/getSingleMessagesFromRedis` - 获取单聊消息历史
- `POST /message/getGroupMessagesFromRedis` - 获取群聊消息历史
- `POST /message/getAiMessagesFromRedis` - 获取 AI 对话消息历史

### 文件相关

- `POST /attach/upload` - 上传文件

### TODO
- 1.日志完善，并追踪调用过程
- 2.冷热数据存储，近期消息存储在redis，全量数据存储在mysql(目前想法是消息队列异步实现)
- 3.docker部署（dockerfile编写）
- 4.优化结构和mysql语句


<!-- 详细的 API 文档请访问 Swagger 页面：http://localhost:8081/swagger/index.html -->

<!-- ## 配置说明

### 心跳配置

- `DelayHeartbeat`: 延迟心跳时间（秒）
- `HeartbeatHz`: 心跳频率（秒）
- `HeartbeatMaxTime`: 最大心跳超时时间（秒），超过此时间用户将被标记为离线
- `RedisOnlineTime`: 在线用户缓存时长（小时）

## 开发说明

### 数据库迁移

项目使用 GORM 进行数据库操作，首次运行时会自动创建表结构。 -->

<!-- ### WebSocket 连接

客户端连接 WebSocket 时需要发送认证消息：

```json
{
  "userId": "user_uuid"
}
``` -->

<!-- ### 日志

项目使用 Zap 进行日志记录，日志配置在 `utils/zap.go` 中。 -->


