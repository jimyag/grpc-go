# grpc-go

这是学习 gRPC 而实现的一个Demo，这个Demo的主要功能就是通过gRPC实现对笔记本电脑参数的创建，储存等功能。

## 如何快速开始

1. 首先确保已经配置好protoc，protobuf等。
2. 确保配置好 make 环境

### clone 代码

```shell
git clone https://github.com/jimyag/grpc-go.git
```

### 下载依赖

```shell
go mod tidy 
go mod download
```

### 启动服务端

```shell
make server
```

### 启动客户端

```shell
make client
```
