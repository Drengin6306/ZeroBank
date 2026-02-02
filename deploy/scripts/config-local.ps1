# 将配置文件中的 Docker 服务名替换为 localhost（用于本地开发）

Write-Host "正在切换到本地配置..." -ForegroundColor Green

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
    if (Test-Path $config) {
        # 备份原始配置（如果还没有备份）
        $backup = "$config.docker"
        if (-not (Test-Path $backup)) {
            Copy-Item $config $backup
            Write-Host "  备份: $config -> $backup" -ForegroundColor Gray
        }

        # 替换为本地配置
        $content = Get-Content $config -Raw -Encoding UTF8
        $content = $content -replace 'mysql:3306', 'localhost:33306'
        $content = $content -replace 'consul:8500', 'localhost:8500'
        $content = $content -replace 'redis:6379', 'localhost:36379'
        Set-Content $config $content -Encoding UTF8 -NoNewline

        Write-Host "  ✓ $config" -ForegroundColor Green
    }
}

Write-Host "本地配置已生效！" -ForegroundColor Green
