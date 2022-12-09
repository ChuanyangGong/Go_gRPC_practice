package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"go_grpc_practice/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	ser := newServer()
	g := grpc.NewServer(
		grpc.UnaryInterceptor(
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				md, ok := metadata.FromIncomingContext(ctx)
				if !ok {
					fmt.Println("非法请求")
					return nil, status.Error(codes.PermissionDenied, "非法请求")
				}
				if info.FullMethod != "/Loginer/Login" {
					user, ok1 := md["user"]
					token, ok2 := md["token"]
					fmt.Println(ok1, ok2, user, token, ser.tokens)
					if !ok1 || !ok2 {
						fmt.Println("验证不通过")
						return nil, status.Error(codes.Unauthenticated, "验证不通过")
					}
					tok, ok := ser.tokens[user[0]]
					if !ok || tok != token[0] {
						fmt.Println("验证不通过")
						return nil, status.Error(codes.Unauthenticated, "验证不通过")
					}

				}
				res, err := handler(ctx, req)
				return res, err
			},
		),
	)
	proto.RegisterLoginerServer(g, ser)
	proto.RegisterFileServerServer(g, ser)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	err = g.Serve(listener)
	if err != nil {
		panic(err)
	}
}

type server struct {
	users  map[string]string
	tokens map[string]string
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

func (s *server) GetFilePath(tarPath string) (string, error) {
	// get file root path
	root, _ := os.Getwd()
	root = path.Join(root, "public")

	tarPath = path.Join(root, tarPath)
	if strings.Index(tarPath, root) != 0 {
		return "", fmt.Errorf("不存在该路径")
	}
	fmt.Println(tarPath)
	return tarPath, nil
}

func (s *server) ListDirectory(ctx context.Context, req *proto.ListDirReq) (*proto.ListDirRes, error) {
	var reply *proto.ListDirRes = nil

	path, err := s.GetFilePath(req.Path)
	if err != nil {
		reply = &proto.ListDirRes{
			Success: false,
		}
	} else {
		fileInfoList, err := ioutil.ReadDir(path)
		if err != nil {
			reply = &proto.ListDirRes{
				Success: false,
			}
		} else {
			reply = &proto.ListDirRes{
				Success:    true,
				FileOrDirs: make([]*proto.ListDirRes_FileOrDirItem, 0),
			}
			for _, info := range fileInfoList {
				reply.FileOrDirs = append(reply.FileOrDirs, &proto.ListDirRes_FileOrDirItem{
					IsFile: !info.IsDir(),
					Name:   info.Name(),
				})
			}
		}
	}

	return reply, nil
}

func (s *server) UploadFile(cliStream proto.FileServer_UploadFileServer) error {
	req, err := cliStream.Recv()
	target, err := s.GetFilePath(req.Filename)
	if err != nil {
		fmt.Println("该文件路径不合法1：" + target)
		cliStream.SendAndClose(&proto.UploadFileRes{
			Msg: "该文件路径不合法",
		})
		return nil
	}

	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 600)
	if err != nil {
		fmt.Println("该文件路径不合法2：" + target)
		cliStream.SendAndClose(&proto.UploadFileRes{
			Msg: "该文件路径不合法",
		})
		return nil
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()
	for err == nil {
		fmt.Println(string(req.File))
		w.Write(req.File)
		req, err = cliStream.Recv()
	}
	fmt.Println("上传完成")
	cliStream.SendAndClose(&proto.UploadFileRes{
		Msg: "上传完成",
	})

	return nil
}

func (s *server) DownloadFile(req *proto.DownloadFileReq, serStream proto.FileServer_DownloadFileServer) error {

	return nil
}
