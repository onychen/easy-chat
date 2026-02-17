#!/bin/bash

reso_addr='crpi-osk929019sdokpya.cn-guangzhou.personal.cr.aliyuncs.com/easy-ct/user-rpc-dev'
tag='latest'

echo "登录阿里云镜像仓库..."
docker login --username=onychen crpi-osk929019sdokpya.cn-guangzhou.personal.cr.aliyuncs.com

pod_ip="192.168.88.130"

container_name="easy-chat-user-rpc-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_easy-im -v /easy-im/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 10000:10000 -e POD_IP=${pod_ip}  --name=${container_name} -d ${reso_addr}:${tag}