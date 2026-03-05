# Easy-Chat Start All Script
# Usage: .\start-all.ps1

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Easy-Chat Start All Services" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$ScriptDir = $PSScriptRoot
if ([string]::IsNullOrEmpty($ScriptDir)) {
    $ScriptDir = Get-Location
}

# Start Docker Compose
Write-Host "[1/2] Starting Docker Compose services..." -ForegroundColor Yellow
Set-Location $ScriptDir
docker-compose up -d

if ($LASTEXITCODE -ne 0) {
    Write-Host "Docker Compose failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Waiting for services..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# Start Go services
Write-Host "[2/2] Starting Go services..." -ForegroundColor Yellow
Write-Host ""

$services = @(
    @{name="user-rpc";   path="apps\user\rpc\user.go";        config="etc/local/user.yaml"},
    @{name="user-api";   path="apps\user\api\user.go";        config="etc/local/user.yaml"},
    @{name="social-rpc"; path="apps\social\rpc\social.go";     config="etc/local/social.yaml"},
    @{name="social-api"; path="apps\social\api\social.go";     config="etc/local/social.yaml"},
    @{name="im-rpc";     path="apps\im\rpc\im.go";            config="etc/local/im.yaml"},
    @{name="im-api";     path="apps\im\api\im.go";            config="etc/local/im.yaml"},
    @{name="im-ws";     path="apps\im\ws\im.go";            config="etc/local/im.yaml"},
    @{name="task-mq";    path="apps\task\mq\task.go";         config="etc/local/mq.yaml"}
)

foreach ($svc in $services) {
    $servicePath = Join-Path $ScriptDir $svc.path
    $configPath = Join-Path (Split-Path $servicePath) $svc.config
    
    Write-Host "  Starting $($svc.name)..." -ForegroundColor Green
    Start-Process -FilePath "go" -ArgumentList "run", $servicePath, "-f", $configPath -WorkingDirectory (Split-Path $servicePath) -PassThru
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  All services started!" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Service URLs:" -ForegroundColor White
Write-Host "  - MySQL:      localhost:13306" -ForegroundColor White
Write-Host "  - Redis:      localhost:16379" -ForegroundColor White
Write-Host "  - MongoDB:    localhost:27017" -ForegroundColor White
Write-Host "  - etcd:       localhost:3379" -ForegroundColor White
Write-Host "  - Kafka:      localhost:9092" -ForegroundColor White
Write-Host "  - APISIX:     localhost:9080" -ForegroundColor White
Write-Host "  - APISIX Dashboard: localhost:9000" -ForegroundColor White
Write-Host "  - user-api:   localhost:8888" -ForegroundColor White
Write-Host "  - user-rpc:   localhost:10000" -ForegroundColor White
Write-Host "  - social-api: localhost:8889" -ForegroundColor White
Write-Host "  - social-rpc: localhost:10001" -ForegroundColor White
Write-Host "  - im-api:     localhost:8890" -ForegroundColor White
Write-Host "  - im-rpc:     localhost:10002" -ForegroundColor White
Write-Host "  - im-ws:      localhost:8891" -ForegroundColor White
Write-Host ""