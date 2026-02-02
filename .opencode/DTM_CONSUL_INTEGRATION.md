# DTM + Consul 服务发现集成指南

## 📌 当前实现状态

**ZeroBank 项目当前使用的方案**：DTM 通过 Docker 网络直接访问服务（不使用 Consul）

**原因**：DTM 官方 Docker 镜像 (`yedf/dtm:latest`) **不包含** Consul 服务发现驱动

## 🎯 架构说明

### 当前架构（Docker 网络直连）

```
┌─────────────────────────────────────────┐
│         Transaction API                  │
│  发起 SAGA 事务                          │
│  accountRpcTarget = "account-rpc:9001"  │
└──────────────┬──────────────────────────┘
               │ gRPC
               ↓
┌─────────────────────────────────────────┐
│              DTM Server                  │
│  直接通过 Docker DNS 解析服务名          │
│  account-rpc:9001 → 172.21.0.8:9001     │
└──────────────┬──────────────────────────┘
               │ gRPC
               ↓
┌─────────────────────────────────────────┐
│          Account RPC (容器)              │
│  服务名: account-rpc                     │
│  IP: 172.21.0.8                         │
│  Port: 9001                             │
└─────────────────────────────────────────┘
```

**优点**：
- ✅ 配置简单，开箱即用
- ✅ 性能高（无需额外的服务发现查询）
- ✅ 适合单机 Docker Compose 部署

**缺点**：
- ❌ 不支持动态服务发现
- ❌ 服务扩容后无法自动负载均衡
- ❌ 服务迁移需要修改配置

---

## 🚀 方案一：使用自定义 DTM 镜像（推荐）

要让 DTM 支持 Consul，需要编译包含 Consul 驱动的自定义镜像。

### 1. 创建自定义 DTM Dockerfile

在项目根目录创建 `deploy/dtm/Dockerfile.consul`：

```dockerfile
# 基于 Golang 构建镜像
FROM golang:1.21-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /build

# 克隆 DTM 源码
RUN git clone --depth 1 https://github.com/dtm-labs/dtm.git .

# 下载 DTM Consul 驱动
RUN go get github.com/dtm-labs/dtmdriver-clients

# 创建包含 Consul 驱动的 main.go
RUN cat > /build/main_consul.go <<EOF
package main

import (
	_ "github.com/dtm-labs/dtmdriver-clients/driver_gozero"
	"github.com/dtm-labs/dtm/dtm"
)

func main() {
	dtm.Main()
}
EOF

# 编译 DTM with Consul
RUN CGO_ENABLED=0 GOOS=linux go build -o dtm-consul main_consul.go

# 运行时镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 复制二进制文件
COPY --from=builder /build/dtm-consul /app/dtm
COPY --from=builder /build/admin /app/admin

# 暴露端口
EXPOSE 36789 36790

# 启动 DTM
ENTRYPOINT ["/app/dtm"]
```

### 2. 构建自定义镜像

```bash
cd deploy/dtm
docker build -t zerobank-dtm:consul -f Dockerfile.consul .
```

### 3. 修改 docker-compose.yml

```yaml
dtm:
  image: zerobank-dtm:consul   # 使用自定义镜像
  container_name: zerobank_dtm
  environment:
    STORE_DRIVER: mysql
    STORE_HOST: mysql
    STORE_PORT: 3306
    STORE_USER: root
    STORE_PASSWORD: 123456
    # Consul 服务发现配置
    MICRO_SERVICE_DRIVER: dtm-driver-gozero
    MICRO_SERVICE_TARGET: consul:8500
    MICRO_SERVICE_END_POINT: dtm:36790
  ports:
    - "36789:36789"
    - "36790:36790"
  depends_on:
    - mysql
    - consul
  networks:
    - zerobank_net
  restart: unless-stopped
```

### 4. 修改 transferlogic.go 使用 Consul 服务名

```go
// 使用 Consul 服务名
accountRpcTarget := "account.rpc"

saga.Add(
    accountRpcTarget+"/account.Account/DeductBalance",
    accountRpcTarget+"/account.Account/CompensateDeduct",
    &account.DeductBalanceRequest{...},
)
```

### 5. 验证集成

```bash
# 重启 DTM
docker-compose up -d dtm

# 查看日志，应该看到 Consul 连接成功
docker-compose logs dtm | grep consul

# 测试转账
curl -X POST http://localhost:8888/api/transaction/transfer \
  -H "Authorization: Bearer TOKEN" \
  -d '{"account_to": "xxx", "amount": 100}'
```

---

## 🔧 方案二：使用 etcd 服务发现（官方支持）

DTM 官方镜像支持 etcd，如果你愿意引入 etcd，这是最简单的方案。

### 1. 添加 etcd 到 docker-compose.yml

```yaml
services:
  etcd:
    image: quay.io/coreos/etcd:v3.5.0
    container_name: zerobank_etcd
    command:
      - etcd
      - --name=etcd0
      - --advertise-client-urls=http://etcd:2379
      - --listen-client-urls=http://0.0.0.0:2379
      - --initial-cluster=etcd0=http://etcd:2380
    ports:
      - "2379:2379"
    networks:
      - zerobank_net
```

### 2. 配置 DTM 使用 etcd

```yaml
dtm:
  environment:
    MICRO_SERVICE_DRIVER: etcd
    MICRO_SERVICE_TARGET: etcd:2379
    MICRO_SERVICE_END_POINT: dtm:36790
```

### 3. 修改 Account RPC 注册到 etcd

需要在 Account RPC 中添加 etcd 注册逻辑（这需要修改 go-zero 的服务注册）。

**缺点**：需要同时维护 Consul 和 etcd 两套服务发现系统。

---

## 📊 三种方案对比

| 方案 | 复杂度 | 性能 | 灵活性 | 推荐度 |
|------|--------|------|--------|--------|
| **Docker 网络直连** (当前) | ⭐ 简单 | ⭐⭐⭐ 高 | ⭐ 低 | ✅ 推荐（单机部署） |
| **自定义 DTM + Consul** | ⭐⭐⭐ 复杂 | ⭐⭐ 中 | ⭐⭐⭐ 高 | ✅ 推荐（生产环境） |
| **DTM + etcd** | ⭐⭐ 中等 | ⭐⭐ 中 | ⭐⭐ 中 | ⚠️ 需要双服务发现 |

---

## 🎯 推荐实践

### 开发/测试环境
**使用当前方案**（Docker 网络直连）
- 配置简单
- 无需额外组件
- 满足功能验证需求

```go
accountRpcTarget := "account-rpc:9001"
```

### 生产环境
**使用方案一**（自定义 DTM + Consul）
- 统一服务发现机制（所有服务都用 Consul）
- 支持动态服务发现和负载均衡
- 架构更加规范

```go
accountRpcTarget := "account.rpc"
```

---

## 🔍 故障排查

### 问题：DTM 提示 "no dtm driver with name: consul"

**原因**：官方镜像不包含 Consul 驱动

**解决**：使用方案一构建自定义镜像

### 问题：DTM 无法连接到 Consul

**检查步骤**：

1. **验证 Consul 可访问**：
   ```bash
   docker exec zerobank_dtm ping -c 2 consul
   docker exec zerobank_dtm nc -zv consul 8500
   ```

2. **查看 DTM 配置**：
   ```bash
   docker-compose logs dtm | grep "MicroService"
   ```

3. **检查 Consul 服务注册**：
   ```bash
   curl http://localhost:8500/v1/catalog/services
   ```

### 问题：DTM 调用服务时地址解析失败

**检查步骤**：

1. **验证服务名格式**：
   - Consul 中注册的服务名：`account.rpc`
   - transferlogic.go 使用的服务名：`account.rpc`（必须一致）

2. **查看 DTM 日志**：
   ```bash
   docker-compose logs dtm | grep ERROR
   ```

---

## 📝 完整示例：自定义 DTM 镜像构建脚本

创建 `deploy/scripts/build-dtm-consul.sh`：

```bash
#!/bin/bash

set -e

echo "📦 开始构建包含 Consul 驱动的 DTM 镜像..."

# 创建临时构建目录
BUILD_DIR=$(mktemp -d)
cd "$BUILD_DIR"

# 克隆 DTM 源码
echo "📥 克隆 DTM 源码..."
git clone --depth 1 https://github.com/dtm-labs/dtm.git
cd dtm

# 创建包含 Consul 驱动的 main.go
echo "📝 创建自定义 main.go..."
cat > main_consul.go <<'EOF'
package main

import (
	_ "github.com/dtm-labs/dtmdriver-clients/driver_gozero"
	"github.com/dtm-labs/dtm/dtm"
)

func main() {
	dtm.Main()
}
EOF

# 下载依赖
echo "📦 下载依赖..."
go mod download
go get github.com/dtm-labs/dtmdriver-clients

# 编译
echo "🔨 编译 DTM..."
CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o dtm-consul main_consul.go

# 构建 Docker 镜像
echo "🐳 构建 Docker 镜像..."
cat > Dockerfile.consul <<'EOF'
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY dtm-consul /app/dtm
COPY admin /app/admin

EXPOSE 36789 36790

ENTRYPOINT ["/app/dtm"]
EOF

docker build -t zerobank-dtm:consul -f Dockerfile.consul .

# 清理
cd ..
rm -rf "$BUILD_DIR"

echo "✅ DTM Consul 镜像构建完成！"
echo "使用方式："
echo "  docker-compose.yml 中修改 dtm.image 为: zerobank-dtm:consul"
```

使用脚本：
```bash
chmod +x deploy/scripts/build-dtm-consul.sh
./deploy/scripts/build-dtm-consul.sh
```

---

## 🎓 总结

### 当前实现（推荐保持）

ZeroBank 项目当前使用 **Docker 网络直连** 方式，这对于开发和单机部署已经足够：

```yaml
# docker-compose.yml
dtm:
  image: yedf/dtm:latest  # 官方镜像
  environment:
    MICRO_SERVICE_DRIVER: default  # 不使用服务发现
```

```go
// transferlogic.go
accountRpcTarget := "account-rpc:9001"  // Docker 服务名
```

### 升级到 Consul 集成（生产环境推荐）

如果需要在生产环境部署，建议：
1. 使用**方案一**构建自定义 DTM 镜像
2. 统一使用 Consul 服务发现
3. 所有微服务（包括 DTM）都通过 Consul 进行服务注册和发现

### 关键要点

- ✅ DTM 官方镜像**不支持** Consul
- ✅ 需要自定义编译才能使用 Consul
- ✅ Docker 网络直连对单机部署已足够
- ✅ 生产环境建议统一服务发现机制
