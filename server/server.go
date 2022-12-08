package main

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go_grpc_practice/proto"
)

type server struct {
	users  map[string]string
	tokens map[string]string
}

func (s *server) generateToken() string {
	length := 18
	rand.Seed(time.Now().UnixNano())
	rs := make([]string, length)
	for start := 0; start < length; start++ {
		t := rand.Intn(3)
		if t == 0 {
			rs = append(rs, strconv.Itoa(rand.Intn(10)))
		} else if t == 1 {
			rs = append(rs, string(rand.Intn(26)+65))
		} else {
			rs = append(rs, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(rs, "")
}

func (s *server) Login(ctx context.Context, req *proto.LoginReqData) (*proto.LoginResData, error) {
	var reply *proto.LoginResData = nil
	val, ok := s.users[req.Username]
	if !ok || req.Password != val {
		reply = &proto.LoginResData{
			Msg:     "用户名或密码错误",
			Success: false,
		}
	} else {
		token := s.generateToken()
		reply = &proto.LoginResData{
			Msg:     "登录成功",
			Success: true,
			Token:   token,
		}
		s.tokens[req.Username] = token
	}

	return reply, nil
}

func newServer() *server {
	return &server{
		users: map[string]string{
			"root": "root123",
			"user": "user123",
		},
		tokens: map[string]string{},
	}
}
