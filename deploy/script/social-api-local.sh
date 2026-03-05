#!/bin/bash

reso_addr='crpi-osk929019sdokpya.cn-guangzhou.personal.cr.aliyuncs.com/easy-ct/social-api-local'
tag='latest'

pod_ip="192.168.88.130"

container_name="easy-chat-social-api-local"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


docker run -p 8881:8881 -e POD_IP=${pod_ip}  --name=${container_name} -d ${reso_addr}:${tag}
