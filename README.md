# Go_gRPC_practice

这个项目为一个 gRPC 的入门实践项目，将通过实现一个简单的服务器和客户端，实践操作学习 gRPC 的基本使用、错误处理、流模式、拦截器、验证器，进行一个简单的入门。

## 使用方式

1. 编译 proto 文件

```proto
protoc -I .\proto\ login.proto --go_out=plugins=grpc:.\proto\
```

2. 运行服务器

```shell
go run .\server\server.go
```

3. 运行客户端

```shell
go run .\client\client.go
```