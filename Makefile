.PHONY: build run stop

# Build the Go application
# service/account/api/
# service/account/rpc/
# service/transaction/api/
# service/riskcontrol/rpc/

build:
	cd service/account/api && go build
	cd service/account/rpc && go build
	cd service/transaction/api && go build
	cd service/riskcontrol/rpc && go build

run:
	powershell -NoProfile -Command "Start-Process -WorkingDirectory 'service\account\api' -FilePath '.\api.exe'"
	powershell -NoProfile -Command "Start-Process -WorkingDirectory 'service\account\rpc' -FilePath '.\rpc.exe'"
	powershell -NoProfile -Command "Start-Process -WorkingDirectory 'service\transaction\api' -FilePath '.\api.exe'"
	powershell -NoProfile -Command "Start-Process -WorkingDirectory 'service\riskcontrol\rpc' -FilePath '.\rpc.exe'"

stop:
	powershell -NoProfile -Command "Get-Process -Name api,rpc -ErrorAction SilentlyContinue | Stop-Process -Force"