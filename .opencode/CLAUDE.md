# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

ZeroBank 是一个基于 go-zero 微服务框架的银行系统，采用微服务架构，包含账户管理、交易处理、风控和报表等核心功能。

## 常用命令

### 开发命令
```bash
# 构建所有服务
make build

# 运行所有服务
make run

# 停止所有服务
make stop

# 快速启动（停止 -> 构建 -> 运行）
make quickstart

# 构建所有 Docker 镜像
make docker
```

### Go 命令
```bash
# 在特定服务目录下构建单个服务
cd service/account/api && go build

# 运行单个服务
cd service/account/api && go run account.go

# 安装依赖
go mod download
```

### Docker 命令
```bash
# 启动所有基础设施服务（MySQL, Redis, Consul, Nginx）
docker-compose up -d

# 停止所有服务
docker-compose down

# 查看服务日志
docker-compose logs -f [service_name]
```

## 架构说明

### 微服务结构

项目采用 API + RPC 双层架构：

- **API 层**：对外提供 HTTP REST 接口，端口 8xxx
- **RPC 层**：服务间通过 gRPC 通信，端口 9xxx

### 服务列表

| 服务 | 类型 | 端口 | 说明 |
|------|------|------|------|
| account-api | HTTP | 8001 | 账户管理 API（注册、登录、用户信息） |
| account-rpc | gRPC | 9001 | 账户 RPC 服务 |
| transaction-api | HTTP | 8002 | 交易 API |
| transaction-rpc | gRPC | 9002 | 交易 RPC 服务 |
| riskcontrol-rpc | gRPC | 9003 | 风控 RPC 服务 |
| report-api | HTTP | 8003 | 报表 API |

### 基础设施服务

- **MySQL** (33306): 主数据库
- **Redis** (36379): 缓存和会话存储
- **Consul** (8500): 服务注册与发现
- **Nginx** (8888): API 网关

### 目录结构

```
service/
├── account/          # 账户服务
│   ├── api/         # HTTP API 服务
│   ├── rpc/         # gRPC 服务
│   └── model/       # 数据模型
├── transaction/     # 交易服务
├── riskcontrol/     # 风控服务
└── report/          # 报表服务

pkg/                 # 共享包
├── auth/           # JWT 认证
├── errorx/         # 统一错误码定义
├── format/         # 格式化工具（手机号等）
├── idgen/          # 分布式 ID 生成器
├── password/       # 密码加密
├── response/       # 统一响应封装
└── vars/           # 全局变量
```

### go-zero 服务定义

- **API 服务**：使用 `.api` 文件定义 HTTP 接口，通过 `goctl` 生成代码
- **RPC 服务**：使用 `.proto` 文件定义 gRPC 接口，通过 `protoc` 和 `goctl` 生成代码

### 代码生成

项目使用 goctl (go-zero 代码生成工具) 生成脚手架代码：

```bash
# 从 .api 文件生成 API 服务代码
goctl api go -api account.api -dir .

# 从 .proto 文件生成 RPC 服务代码
goctl rpc protoc account.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```

**注意**：带有 `// Code scaffolded by goctl` 注释的文件为生成文件，部分标记为 `Safe to edit` 的可以修改。

### 错误处理

项目使用统一的错误码系统（`pkg/errorx`）：

- **通用错误** (1000-1999): 参数错误、未登录、权限不足等
- **业务错误** (2000+): 账户不存在、余额不足、风控拒绝等

错误码定义在 `pkg/errorx/code.go` 中，错误消息在 `pkg/errorx/msg.go` 中。

### 服务间通信

- API 服务通过 Consul 发现并调用 RPC 服务
- RPC 服务注册到 Consul 实现服务发现
- 使用 gRPC 进行服务间通信

### 配置文件

每个服务的配置文件位于 `service/{service_name}/{api|rpc}/etc/` 目录，包含：

- 服务监听端口
- 数据库连接字符串
- JWT 密钥配置
- Consul 连接信息
- Redis 连接信息

### 身份认证

- 使用 JWT 进行身份认证
- JWT 工具封装在 `pkg/auth/jwt.go`
- 密码使用 bcrypt 加密，工具在 `pkg/password/password.go`
