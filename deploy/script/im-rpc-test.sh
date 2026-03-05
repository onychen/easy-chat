#!/bin/bash

reso_addr='crpi-osk929019sdokpya.cn-guangzhou.personal.cr.aliyuncs.com/easy-ct/im-rpc-dev'
tag='latest'

pod_ip="192.168.88.130"

container_name="easy-chat-im-rpc-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


docker run -p 8080:8080 -e POD_IP=${pod_ip}  --name=${container_name} -d ${reso_addr}:${tag}
