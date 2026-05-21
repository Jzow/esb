#!/bin/bash
service_name="test"
port=29070
log_path=`pwd`
container_id=`docker ps -aqf "name=^$service_name$"`
echo "go build"
go mod tidy
go build -o main .
echo "[deploy]container_id=$container_id"
docker stop $container_id
echo "[deploy]the container has been stopped"
docker rm -f $container_id
echo "[deploy]the container has been deleted"
docker rmi -f $service_name:latest
echo "[deploy]the container's image has been deleted"
docker build -t $service_name .
echo "[deploy]have rebuilt the image"
docker run -d -p $port:$port --name $service_name --net=host -v $log_path/runtime/logs:/go/src/github.com/EDDYCJY/go-gin-example/runtime/logs --restart=always $service_name
echo "[deploy]restart a new container from the rebuilt image on port $port"