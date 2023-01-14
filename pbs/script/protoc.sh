#!/usr/bin/env bash

# 使用姿势
# pb文件必须以服务名字开头，里面所有数据结构以及rpc都带上服务名称，防止跟别的服务冲突
#/vas/apps/core/internal/pbs
#zander@macos pbs % ./goout.sh ./proto/core_order.proto

docker run --name go --rm -v $(pwd):$(pwd) -w $(pwd) -it registry.cn-hangzhou.aliyuncs.com/zander84/golang:1.19 protoc.sh $1
