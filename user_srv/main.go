package main

import (
	"flag"
	"fmt"
	"go-gin-rme/user_srv/api"
	"go-gin-rme/user_srv/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 8088, "端口")
	flag.Parse()
	log.Printf("当前地址:%s:%d", *IP, *Port)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &api.UserServer{})
	//lis, err := net.Listen("tcp", "0.0.0.0:8088")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())

	}
	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
