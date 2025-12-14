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