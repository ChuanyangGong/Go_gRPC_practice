# Go_gRPC_practice

这个项目为一个 gRPC 的入门实践项目，将通过实现一个简单的服务器和客户端，实践操作学习 gRPC 的基本使用、错误处理、流模式、拦截器、验证器，进行一个简单的入门。

## 使用方式

1. 编译 proto 文件

```proto
protoc -I .\proto\ login.proto --go_out=plugins=grpc:.\proto\
protoc -I .\proto\ fileserver.proto --go_out=plugins=grpc:.\proto\
```

2. 运行服务器

```shell
go run .\server\server.go
```

3. 运行客户端

```shell
go run .\client\client.go
```

## 涵盖内容

- [x] 登录功能，进行 token 的记录和返回

- [x] 获取目录列表功能，进行 gRPC 基本使用的实践

- [x] 上传文件功能，实践流模式

- [x] 客户端请求自动携带 token，服务器非登录请求自动验证 token，实现拦截器

- [x] 在所有操作中进行错误处理

- [ ] 对请求列表参数添加目录地址验证功能，实践验证器