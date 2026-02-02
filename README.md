# ZeroBank - 微服务银行系统

基于 go-zero 框架构建的分布式银行系统，采用微服务架构，提供账户管理、交易处理、风控和报表等核心功能。

## 🚀 快速开始

```bash
# 启动基础设施（MySQL, Redis, Consul, DTM）
docker-compose up -d

# 构建所有服务
make build

# 运行所有服务
make run

# 快速启动（停止 -> 构建 -> 运行）
make quickstart
```

访问 Nginx 网关：`http://localhost:8888`

## 📚 文档

详细文档位于 [`.opencode/`](./.opencode/) 目录：

- **[AGENTS.md](./.opencode/AGENTS.md)** - AI 代理开发规范（构建、测试、代码风格）
- **[CLAUDE.md](./.opencode/CLAUDE.md)** - 项目架构和概述
- **[DTM_INTEGRATION.md](./.opencode/DTM_INTEGRATION.md)** - 分布式事务集成说明

## 🏗️ 架构

### 微服务列表

| 服务 | 类型 | 端口 | 说明 |
|------|------|------|------|
| account-api | HTTP | 8001 | 账户管理 API |
| account-rpc | gRPC | 9001 | 账户 RPC 服务 |
| transaction-api | HTTP | 8002 | 交易 API |
| transaction-rpc | gRPC | 9002 | 交易 RPC 服务 |
| riskcontrol-rpc | gRPC | 9003 | 风控 RPC 服务 |
| report-api | HTTP | 8003 | 报表 API |

### 基础设施

- **MySQL** (33306) - 主数据库
- **Redis** (36379) - 缓存和会话
- **Consul** (8500) - 服务注册与发现
- **DTM** (36789/36790) - 分布式事务管理器
- **Nginx** (8888) - API 网关

## 🛠️ 技术栈

- **框架**: [go-zero](https://go-zero.dev/) v1.9.3
- **语言**: Go 1.25.4
- **数据库**: MySQL 8.0
- **缓存**: Redis 8.4
- **服务发现**: Consul
- **分布式事务**: [DTM](https://dtm.pub/)
- **通信协议**: gRPC, HTTP/REST

## 📁 项目结构

```
ZeroBank/
├── .opencode/              # 项目文档
├── service/                # 微服务目录
│   ├── account/           # 账户服务
│   │   ├── api/          # HTTP API
│   │   ├── rpc/          # gRPC 服务
│   │   └── model/        # 数据模型
│   ├── transaction/       # 交易服务
│   ├── riskcontrol/       # 风控服务
│   └── report/            # 报表服务
├── pkg/                   # 共享包
│   ├── auth/             # JWT 认证
│   ├── errorx/           # 错误码
│   ├── password/         # 密码加密
│   └── vars/             # 全局变量
├── deploy/                # 部署配置
├── data/                  # 数据持久化目录
├── docker-compose.yml     # Docker Compose 配置
└── Makefile              # 构建脚本
```

## 🔑 核心功能

### ✅ 已实现

- 账户注册（个人/企业）
- 账户登录（JWT 认证）
- 存款/取款操作
- 转账功能（DTM 分布式事务保护）
- 风控检查（交易限额、频率控制）
- 交易记录查询
- 余额查询
- 报表生成（Excel 导出）

### 🔐 安全特性

- **密码加密**: bcrypt 哈希
- **JWT 认证**: 无状态身份验证
- **悲观锁**: `SELECT ... FOR UPDATE` 防止并发问题
- **分布式事务**: DTM SAGA 模式确保数据一致性
- **风控系统**: 实时交易风险检测

## 🧪 开发

### 运行单个服务

```bash
cd service/account/api
go build
./account-api -f etc/account-api.yaml
```

### 运行测试

```bash
# 所有测试
go test ./...

# 单个测试
go test -run TestFunctionName ./path/to/package

# 测试覆盖率
go test -cover ./...
```

### 代码生成

```bash
# 从 .api 文件生成代码
cd service/account/api
goctl api go -api account.api -dir .

# 从 .proto 文件生成代码
cd service/account/rpc
goctl rpc protoc account.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```

## 📊 监控

- **DTM 管理界面**: http://localhost:36789
- **Consul UI**: http://localhost:8500

## 🤝 贡献

请参考 [AGENTS.md](./.opencode/AGENTS.md) 了解代码规范和最佳实践。

## 📄 许可证

[MIT License](LICENSE)

---

**作者**: Drengin  
**最后更新**: 2026-01-31
