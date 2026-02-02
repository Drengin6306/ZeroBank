.PHONY: build run stop quickstart config-local config-docker

# 定义服务目录
SERVICES := service/account/api service/account/rpc service/transaction/api service/transaction/rpc service/riskcontrol/rpc service/report/api

# Build the Go application
build:
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $$service && go build && cd - > /dev/null || exit 1; \
	done

# 切换到本地配置
config-local:
	@powershell -NoProfile -ExecutionPolicy Bypass -File deploy/scripts/config-local.ps1

# 恢复 Docker 配置
config-docker:
	@powershell -NoProfile -ExecutionPolicy Bypass -File deploy/scripts/config-docker.ps1

# 本地运行（自动切换配置）
run: config-local
	@for service in $(SERVICES); do \
		name=$${service##*/}; \
		exe="$$name.exe"; \
		echo "Starting $$service..."; \
		cd $$service && start "" $$exe && cd - > /dev/null; \
	done

# 停止服务（自动恢复 Docker 配置）
stop:
	- powershell -NoProfile -Command "Get-Process -Name api,rpc -ErrorAction SilentlyContinue | Stop-Process -Force"
	@$(MAKE) config-docker

quickstart: stop build run

docker:
	docker build --build-arg SERVICE_PATH=service/account/api --build-arg SERVICE_NAME=account-api --build-arg SERVICE_PORT=8001 -t account-api:1.0 -f Dockerfile .
	docker build --build-arg SERVICE_PATH=service/account/rpc --build-arg SERVICE_NAME=account-rpc --build-arg SERVICE_PORT=9001 -t account-rpc:1.0 -f Dockerfile .
	docker build --build-arg SERVICE_PATH=service/transaction/api --build-arg SERVICE_NAME=transaction-api --build-arg SERVICE_PORT=8002 -t transaction-api:1.0 -f Dockerfile .
	docker build --build-arg SERVICE_PATH=service/transaction/rpc --build-arg SERVICE_NAME=transaction-rpc --build-arg SERVICE_PORT=9002 -t transaction-rpc:1.0 -f Dockerfile .
	docker build --build-arg SERVICE_PATH=service/riskcontrol/rpc --build-arg SERVICE_NAME=riskcontrol-rpc --build-arg SERVICE_PORT=9003 -t riskcontrol-rpc:1.0 -f Dockerfile .
	docker build --build-arg SERVICE_PATH=service/report/api --build-arg SERVICE_NAME=report-api --build-arg SERVICE_PORT=8003 -t report-api:1.0 -f Dockerfile .