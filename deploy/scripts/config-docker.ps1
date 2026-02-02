# 恢复 Docker 配置（从备份文件）

Write-Host "正在恢复 Docker 配置..." -ForegroundColor Yellow

# 定义服务配置文件路径
$configs = @(
    "service/account/api/etc/account-api.yaml",
    "service/account/rpc/etc/account.yaml",
    "service/transaction/api/etc/transaction-api.yaml",
    "service/transaction/rpc/etc/transaction.yaml",
    "service/riskcontrol/rpc/etc/riskcontrol.yaml",
    "service/report/api/etc/report-api.yaml"
)

foreach ($config in $configs) {
    $backup = "$config.docker"
    if (Test-Path $backup) {
        Copy-Item $backup $config -Force
        Write-Host "  ✓ $config" -ForegroundColor Yellow

        # 删除备份文件
        Remove-Item $backup -Force
        Write-Host "  删除备份: $backup" -ForegroundColor Gray
    }
}

Write-Host "Docker 配置已恢复！" -ForegroundColor Yellow
