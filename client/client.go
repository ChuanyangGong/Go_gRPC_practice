package main

import (
	"bufio"
	"context"
	"fmt"
	"go_grpc_practice/proto"
	"os"
	"path"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	client := client{}
	client.init()

	// 测试错误的用户名
	client.Login("root1", "root123")

	// 测试密码错误
	client.Login("root", "root1234")

	// 测试成功登录
	client.Login("root", "root123")

	// 测试获得指定目录文件列表
	client.GetFileList("/notexist")

	client.GetFileList("/log.txt")

	client.GetFileList("../")

	client.GetFileList("./")

	// 测试发送文件
	client.UploadFile("client.log")

	client.close()
}

type client struct {
	conn  *grpc.ClientConn
	token string
	user  string
	c     proto.LoginerClient
	fc    proto.FileServerClient
}

func (c *client) init() {
	interceptor := grpc.WithUnaryInterceptor(
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			md := metadata.New(map[string]string{
				"user":  c.user,
				"token": c.token})
			ctx = metadata.NewOutgoingContext(ctx, md)
			err := invoker(ctx, method, req, reply, cc, opts...)
			return err
		},
	)
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure(), interceptor)
	if err != nil {
		panic(err)
	}

	c.c = proto.NewLoginerClient(conn)
	c.fc = proto.NewFileServerClient(conn)
	c.conn = conn
}

func (c *client) close() {
	c.conn.Close()
}

func (c *client) Login(username string, password string) {
	r, err := c.c.Login(context.Background(), &proto.LoginReqData{
		Password: password,
		Username: username,
	})
	if err != nil {
		panic(err)
	}
	if r.Success {
		c.token = r.Token
		c.user = username
		fmt.Println("Token: " + r.Token)
	}
	fmt.Println(r.Msg)
}

func (c *client) GetFileList(path string) {
	r, err := c.fc.ListDirectory(context.Background(), &proto.ListDirReq{
		Path: path,
	})
	if err != nil {
		panic(err)
	}

	if r.Success {
		fmt.Println(r.FileOrDirs)
	} else {
		fmt.Printf("目录路径 %s 不存在\n", path)
	}
}

func (c *client) UploadFile(tarPath string) {
	sender, err := c.fc.UploadFile(context.Background())
	if err != nil {
		fmt.Println("连接失败：", err)
		return
	}
	filepath, _ := os.Getwd()
	filepath = path.Join(filepath, "private", tarPath)
	file, err := os.OpenFile(filepath, os.O_RDONLY, 600)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	b := make([]byte, 20)
	for {
		len, err := r.Read(b)
		if err != nil {
			panic(err)
		}

		err = sender.Send(&proto.UploadFileReq{
			Filename: tarPath,
			File:     b,
		})
		if err != nil {
			fmt.Print(err)
			return
		}

		if len <= 20 {
			break
		}
	}
	fmt.Println("发送成功")
}
