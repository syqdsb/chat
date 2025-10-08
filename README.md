# 聊天后端系统

一个使用 Golang 构建的完整聊天后端系统，支持用户注册、登录、添加好友和实时私聊功能。

## 项目结构

```
chat-backend/
├── main.go              # 主程序入口
├── go.mod              # 依赖管理
├── config/
│   └── config.go       # 配置管理
├── database/
│   └── database.go     # 数据库初始化
├── models/
│   └── user.go         # 用户和消息模型
├── handlers/
│   ├── auth.go         # 认证处理器
│   └── websocket.go    # WebSocket 处理器
├── utils/
│   └── response.go     # 响应工具
└── routes/
    └── routes.go       # 路由设置
```

## 安装依赖

```bash
go mod download
```

## 运行项目

```bash
go run main.go
```

服务器默认运行在 `http://localhost:8080`

## API 文档

### 1. 用户注册

**接口**: `POST /api/signup`

**请求参数**:
```json
{
  "username": "用户名",
  "email": "邮箱地址",
  "password": "密码"
}
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "message": "注册成功",
    "token": "生成的token",
    "username": "用户名",
    "email": "邮箱地址"
  }
}
```

### 2. 用户登录

**接口**: `POST /api/login`

**请求参数**:
```json
{
  "email": "邮箱地址",
  "password": "密码"
}
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "message": "登录成功",
    "token": "用户token",
    "username": "用户名",
    "email": "邮箱地址"
  }
}
```

### 3. 添加好友

**接口**: `POST /api/adduser`

**请求参数**:
```json
{
  "username": "要添加的用户名",
  "token": "当前用户的token"
}
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "message": "添加好友成功",
    "username": "好友用户名"
  }
}
```

## WebSocket 使用说明

**连接地址**: `ws://localhost:8080/ws/chat`

### 连接流程

1. **建立 WebSocket 连接**
2. **发送鉴权消息**（仅需一次）
3. **发送和接收消息**

### 消息格式

#### 1. 鉴权（首次连接必须）

```json
{
  "type": "auth",
  "token": "用户token"
}
```

**响应**:
```json
{
  "success": true,
  "message": "鉴权成功",
  "data": {
    "user_id": 1,
    "username": "用户名"
  }
}
```

#### 2. 发送消息

```json
{
  "type": "send",
  "to_user": "接收者用户名",
  "content": "消息内容"
}
```

**响应**:
```json
{
  "success": true,
  "message": "消息发送成功"
}
```

**接收者收到的实时推送**:
```json
{
  "success": true,
  "message": "新消息",
  "data": {
    "from_username": "发送者用户名",
    "from_user_id": 1,
    "content": "消息内容"
  }
}
```

#### 3. 获取未读消息

```json
{
  "type": "get"
}
```

**响应**:
```json
{
  "success": true,
  "message": "获取消息成功",
  "data": {
    "count": 2,
    "messages": [
      {
        "id": 1,
        "from_user_id": 2,
        "to_user_id": 1,
        "content": "消息内容",
        "is_read": 0,
        "created_at": "2025-10-08 12:00:00"