FROM golang:alpine AS builder
ENV GOPROXY https://goproxy.cn,direct

# 更新Alpine的软件源为国内镜像站点，提高下载速度
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add gcc g++ libc-dev librdkafka-dev pkgconf

RUN mkdir /app
ADD . /app/

WORKDIR /app
RUN go build -tags musl -o main .

FROM alpine

# 更新Alpine的软件源为国内镜像站点
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add gcc g++ libc-dev librdkafka-dev pkgconf

ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 29070
ENTRYPOINT ["./main"]

