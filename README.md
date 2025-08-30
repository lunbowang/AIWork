# AI智能办公助手

一个基于Go语言和LangChain框架开发的智能办公助手系统，集成了多种AI能力，为企业提供智能化的办公解决方案。

## 🚀 功能特性

### 核心功能
- **智能问答系统** - 基于大语言模型的智能对话
- **待办事项管理** - 智能创建、查询、修改待办事项
- **审批流程管理** - 支持请假、外出、补卡等审批类型
- **知识库管理** - 企业知识库的智能检索和更新
- **聊天记录分析** - 自动总结和分析聊天内容
- **文件分析处理** - 支持多种文件格式的智能分析
- **数据可视化** - 自动生成图表和报表

### AI能力
- **多模型支持** - 支持OpenAI、阿里通义千问等主流模型
- **智能路由** - 自动选择最适合的AI处理模块
- **上下文记忆** - 支持多轮对话和上下文理解
- **多语言支持** - 支持中文等多语言交互

## 🏗️ 技术架构

### 后端技术栈
- **语言**: Go 1.24.4
- **Web框架**: Gin
- **AI框架**: LangChain
- **数据库**: MongoDB
- **缓存**: Redis
- **认证**: JWT
- **配置管理**: Viper

### 核心模块
```
├── internal/           # 内部业务逻辑
│   ├── logic/         # 业务逻辑层
│   ├── handler/       # 请求处理器
│   ├── svc/           # 服务上下文
│   └── config/        # 配置管理
├── pkg/               # 公共包
│   ├── langchain/     # LangChain集成
│   ├── httpx/         # HTTP工具
│   └── xerr/          # 错误处理
└── web/               # 前端资源
```

## 📦 安装部署

### 环境要求
- Go 1.24.4+
- MongoDB 4.0+
- Redis 6.0+

### 快速开始

1. **克隆项目**
```bash
git clone <repository-url>
cd ai
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置环境**
```bash
cp etc/api.yaml.example etc/api.yaml
# 编辑配置文件，设置数据库连接、API密钥等
```

4. **运行服务**
```bash
# 开发模式
go run main.go

# 生产模式
make build
./ai.exe
```

### 配置说明

主要配置项包括：
- **MongoDB**: 数据库连接配置
- **LangChain**: AI模型配置
- **阿里GPT**: 通义千问配置
- **JWT**: 认证密钥配置

## 🔧 使用指南

### API接口

#### 聊天接口
```http
POST /v1/chat
Content-Type: application/json

{
  "prompts": "你好，请帮我分析一下这个月的销售数据",
  "relationId": "12345"
}
```

#### 待办事项
```http
POST /v1/todo
GET /v1/todo/list
PUT /v1/todo
DELETE /v1/todo/:id
```

#### 审批流程
```http
POST /v1/approval
GET /v1/approval/list
```

### 智能助手使用

系统支持多种智能场景：

1. **自然语言交互** - 直接描述需求，系统自动理解
2. **智能路由** - 自动选择最适合的处理模块
3. **上下文记忆** - 支持多轮对话
4. **多模态输入** - 支持文本、语音等多种输入方式

## 🧪 开发指南

### 项目结构
```
ai/
├── cmd/               # 命令行工具
├── internal/          # 内部包
│   ├── config/       # 配置定义
│   ├── domain/       # 领域模型
│   ├── handler/      # HTTP处理器
│   ├── logic/        # 业务逻辑
│   ├── model/        # 数据模型
│   └── svc/          # 服务上下文
├── pkg/              # 公共包
├── web/              # 前端资源
├── etc/              # 配置文件
├── deploy/           # 部署配置
└── Makefile          # 构建脚本
```

### 开发规范
- 遵循Go语言官方代码规范
- 使用统一的错误处理方式
- 编写完整的单元测试
- 保持代码注释的完整性

### 添加新功能
1. 在`internal/domain`中定义数据结构
2. 在`internal/model`中实现数据访问
3. 在`internal/logic`中实现业务逻辑
4. 在`internal/handler`中暴露HTTP接口

## 📊 性能监控

### 关键指标
- API响应时间
- 数据库查询性能
- AI模型调用成功率
- 系统资源使用率

### 日志系统
- 使用结构化日志记录
- 支持不同级别的日志输出
- 集成链路追踪

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

### 贡献流程
1. Fork项目
2. 创建功能分支
3. 提交代码
4. 创建Pull Request

### 代码审查
- 所有代码变更需要经过审查
- 确保测试覆盖率
- 遵循项目编码规范

## 📄 许可证

本项目采用 [MIT License](LICENSE) 许可证。

## 📞 联系我们

- 项目主页: [GitHub Repository]
- 问题反馈: [Issues]
- 邮箱: [your-email@example.com]

## 🙏 致谢

感谢以下开源项目的支持：
- [Gin](https://github.com/gin-gonic/gin) - Go Web框架
- [LangChain](https://github.com/tmc/langchaingo) - AI应用框架
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver) - MongoDB驱动

---

⭐ 如果这个项目对你有帮助，请给我们一个Star！
