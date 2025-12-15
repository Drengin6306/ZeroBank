.PHONY: build run stop quickstart

# 定义服务目录
SERVICES := service/account/api service/account/rpc service/transaction/api service/transaction/rpc service/riskcontrol/rpc service/report/api

# Build the Go application
build:
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $$service && go build && cd - > /dev/null || exit 1; \
	done

run:
	@for service in $(SERVICES); do \
		name=$${service##*/}; \
		exe="$$name.exe"; \
		echo "Starting $$service..."; \
		cd $$service && start "" $$exe && cd - > /dev/null; \
	done

stop:
	- powershell -NoProfile -Command "Get-Process -Name api,rpc -ErrorAction SilentlyContinue | Stop-Process -Force"

quickstart: stop build run

docker:
	docker build --build-arg SERVICE_PATH=service/account/api --build-arg SERVICE_NAME=account-api --build-arg SERVICE_PORT=8001 -t account-api:1.0 -f Dockerfile .
	docker build --build-arg SERVICE_PATH=service/account/rpc --build-arg SERVICE_NAME=account-rpc --build-arg SERVICE_PORT=9001 -t account-rpc:1.0 -f Dockerfile .
	docker build --build-arg SERVICE_PATH=service/transaction/api --build-arg SERVICE_NAME=transaction-api --build-arg SERVICE_PORT=8002 -t transaction-api:1.0 -f Dockerfile .
	docker build --build-arg SERVICE_PATH=service/transaction/rpc --build-arg SERVICE_NAME=transaction-rpc --build-arg SERVICE_PORT=9002 -t transaction-rpc:1.0 -f Dockerfile .
	docker build --build-arg SERVICE_PATH=service/riskcontrol/rpc --build-arg SERVICE_NAME=riskcontrol-rpc --build-arg SERVICE_PORT=9003 -t riskcontrol-rpc:1.0 -f Dockerfile .
	docker build --build-arg SERVICE_PATH=service/report/api --build-arg SERVICE_NAME=report-api --build-arg SERVICE_PORT=8003 -t report-api:1.0 -f Dockerfile .