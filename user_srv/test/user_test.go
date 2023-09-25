package test

import (
	"context"
	"fmt"
	"go-gin-rme/user_srv/proto"
	"google.golang.org/grpc"
	"testing"
)

var clent proto.UserClient

var conn *grpc.ClientConn
var err error

func Int() {
	conn, err = grpc.Dial("0.0.0.0:8088", grpc.WithInsecure())
	if err != nil {
		defer conn.Close()
	}
	clent = proto.NewUserClient(conn)
}
func GetUserList() {
	resp, err := clent.GetUserList(context.Background(), &proto.PageInfo{
		Page: 1,
		Size: 10,
	})
	if err != nil {
		panic(err)

	}
	for _, v := range resp.Data {
		fmt.Println(v)
		r, err := clent.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          "qwe123",
			EncryptedPassword: v.Password,
		})
		if err != nil {
			panic(err)

		}
		fmt.Println(r.Success)
	}
}
func TestGetAllUser(test *testing.T) {
	Int()
	GetUserList()
	conn.Close()
}
