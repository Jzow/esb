#!/bin/bash
service_name="crpi-5zbmuq0cr4awkhv6.cn-guangzhou.personal.cr.aliyuncs.com/iwork-club/test:latest"
port=29070
container_id=`docker ps -aqf "name=^$service_name$"`
echo "删除容器：$service_name"
docker rmi -f $service_name
echo "创建容器：$service_name"
docker build -t $service_name .
echo "推送容器：$service_name"
echo "docker login --username=aliyun7596916190 crpi-5zbmuq0cr4awkhv6.cn-guangzhou.personal.cr.aliyuncs.com"
echo "Chattry@2024" | sudo 命令
docker push $service_name 
