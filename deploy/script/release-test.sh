#!/bin/bash

echo "登录阿里云镜像仓库..."
docker login --username=onychen crpi-osk929019sdokpya.cn-guangzhou.personal.cr.aliyuncs.com

need_start_server_shell=(
    "user-rpc-test.sh"
    "user-api-test.sh"
    "social-rpc-test.sh"
    "social-api-test.sh"
)

for i in ${need_start_server_shell[*]} ; do
    chmod +x $i
    ./$i
done


docker ps

docker exec -it etcd etcdctl get --prefix ""
