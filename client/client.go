package main

import (
	"context"
	"fmt"
	"go_grpc_practice/proto"

	"google.golang.org/grpc"
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

	client.close()
}

type client struct {
	conn  *grpc.ClientConn
	token string
	c     proto.LoginerClient
}

func (c *client) init() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	c.c = proto.NewLoginerClient(conn)
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
		fmt.Println("Token: " + r.Token)
	}
	fmt.Println(r.Msg)
}
