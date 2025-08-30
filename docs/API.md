# API接口文档

## 概述

AI智能办公助手系统提供RESTful API接口，支持智能问答、待办事项管理、审批流程等功能。

## 基础信息

- **基础URL**: `http://localhost:8080`
- **认证方式**: JWT Token
- **数据格式**: JSON
- **字符编码**: UTF-8

## 认证

### 获取Token

```http
POST /v1/auth/login
Content-Type: application/json

{
  "username": "your_username",
  "password": "your_password"
}
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "123",
      "username": "your_username",
      "email": "user@example.com"
    }
  }
}
```

### 使用Token

在请求头中添加：
```
Authorization: Bearer <your_token>
```

## 智能聊天接口

### 发送消息

```http
POST /v1/chat
Authorization: Bearer <token>
Content-Type: application/json

{
  "prompts": "你好，请帮我分析一下这个月的销售数据",
  "relationId": "12345",
  "startTime": 0,
  "endTime": 0
}
```

**请求参数**:
- `prompts` (string, 必填): 用户输入的消息内容
- `relationId` (string, 可选): 关联ID，用于多轮对话
- `startTime` (int64, 可选): 开始时间戳
- `endTime` (int64, 可选): 结束时间戳

**响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "chatType": 0,
    "data": "根据您提供的信息，我来分析一下这个月的销售数据..."
  }
}
```

**聊天类型说明**:
- `0`: 默认处理器
- `1`: 待办事项
- `2`: 审批流程
- `3`: 知识库
- `4`: 聊天记录分析
- `5`: 图片生成

## 待办事项管理

### 创建待办事项

```http
POST /v1/todo
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "完成项目报告",
  "desc": "编写本季度项目总结报告",
  "deadlineAt": 1735689600,
  "executeIds": ["user1", "user2"]
}
```

**请求参数**:
- `title` (string, 必填): 待办事项标题
- `desc` (string, 可选): 待办事项描述
- `deadlineAt` (int64, 可选): 截止时间戳
- `executeIds` ([]string, 可选): 执行人ID列表

**响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": "64f8a1b2c3d4e5f6a7b8c9d0"
  }
}
```

### 查询待办事项

```http
GET /v1/todo/list?page=1&size=10&status=pending
Authorization: Bearer <token>
```

**查询参数**:
- `page` (int, 可选): 页码，默认1
- `size` (int, 可选): 每页大小，默认10
- `status` (string, 可选): 状态筛选 (pending/completed/cancelled)
- `keyword` (string, 可选): 关键词搜索

### 更新待办事项

```http
PUT /v1/todo
Authorization: Bearer <token>
Content-Type: application/json

{
  "id": "64f8a1b2c3d4e5f6a7b8c9d0",
  "title": "完成项目报告（已更新）",
  "status": "completed"
}
```

### 删除待办事项

```http
DELETE /v1/todo/64f8a1b2c3d4e5f6a7b8c9d0
Authorization: Bearer <token>
```

### 完成待办事项

```http
POST /v1/todo/finish
Authorization: Bearer <token>
Content-Type: application/json

{
  "id": "64f8a1b2c3d4e5f6a7b8c9d0",
  "remark": "任务已完成"
}
```

## 审批流程管理

### 创建审批

```http
POST /v1/approval
Authorization: Bearer <token>
Content-Type: application/json

{
  "type": "leave",
  "title": "请假申请",
  "content": "因家中有事，申请请假3天",
  "startTime": 1735689600,
  "endTime": 1735776000,
  "approverIds": ["manager1", "manager2"]
}
```

**审批类型**:
- `leave`: 请假审批
- `goout`: 外出审批
- `makecard`: 补卡审批

### 查询审批列表

```http
GET /v1/approval/list?page=1&size=10&type=leave
Authorization: Bearer <token>
```

### 审批操作

```http
POST /v1/approval/process
Authorization: Bearer <token>
Content-Type: application/json

{
  "id": "64f8a1b2c3d4e5f6a7b8c9d0",
  "action": "approve",
  "remark": "同意请假申请"
}
```

**操作类型**:
- `approve`: 同意
- `reject`: 拒绝
- `return`: 退回

## 文件管理

### 上传文件

```http
POST /v1/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: [文件内容]
```

**响应示例**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "filename": "document.pdf",
    "size": 1024000,
    "url": "/uploads/document.pdf"
  }
}
```

### 文件分析

```http
POST /v1/chat
Authorization: Bearer <token>
Content-Type: application/json

{
  "prompts": "请分析这个文件的内容",
  "files": ["file_id_1", "file_id_2"]
}
```

## 用户管理

### 获取用户信息

```http
GET /v1/user/info
Authorization: Bearer <token>
```

### 更新用户信息

```http
PUT /v1/user
Authorization: Bearer <token>
Content-Type: application/json

{
  "nickname": "新昵称",
  "email": "newemail@example.com"
}
```

## 部门管理

### 获取部门列表

```http
GET /v1/department/list
Authorization: Bearer <token>
```

### 创建部门

```http
POST /v1/department
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "技术部",
  "parentId": "parent_dept_id",
  "description": "负责技术研发工作"
}
```

## 错误处理

### 错误响应格式

```json
{
  "code": 400,
  "msg": "参数错误",
  "data": null
}
```

### 常见错误码

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权，Token无效 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

### 错误处理示例

```go
// 客户端错误处理
if response.code != 200 {
    fmt.Printf("API调用失败: %s\n", response.msg)
    return
}

// 处理成功响应
data := response.data
```

## 限流策略

- **普通用户**: 100次/分钟
- **VIP用户**: 1000次/分钟
- **管理员**: 无限制

## 最佳实践

### 1. 错误处理
- 始终检查响应状态码
- 实现重试机制
- 记录详细的错误日志

### 2. 性能优化
- 使用连接池
- 实现请求缓存
- 批量处理数据

### 3. 安全性
- 定期更新Token
- 验证所有输入参数
- 使用HTTPS传输

### 4. 监控告警
- 监控API响应时间
- 设置错误率告警
- 记录调用统计

## 更新日志

### v1.0.0 (2025-08-30)
- 初始版本发布
- 支持基础聊天功能
- 实现待办事项管理
- 添加审批流程功能

---

如有问题，请联系技术支持团队。
