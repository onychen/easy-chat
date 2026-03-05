Write-Host "开始部署本地服务..."

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

& "$scriptDir\user-rpc-local.sh"
& "$scriptDir\user-api-local.sh"
& "$scriptDir\social-rpc-local.sh"
& "$scriptDir\social-api-local.sh"
& "$scriptDir\im-rpc-local.sh"
& "$scriptDir\im-api-local.sh"
& "$scriptDir\im-ws-local.sh"
& "$scriptDir\task-mq-local.sh"

Write-Host "`n部署完成！当前运行的容器："
docker ps
