package api

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"go-gin-rme/user_srv/global"
	"go-gin-rme/user_srv/model"
	"go-gin-rme/user_srv/proto"
	"go-gin-rme/user_srv/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
)

type UserServer struct{}

var (
	db = global.DB
)

func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

func ModelToResponse(user model.User) (userInfoResp proto.UserInfoResponse) {

	userInfoResp = proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.NickName,
		Sex:      int32(user.Sex),
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoResp.Birthday = uint64(user.Birthday.Unix())
	}
	return
}

func (u *UserServer) GetUserList(c context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	//获取用户列表
	var users []model.User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	resp := &proto.UserListResponse{
		Total: result.RowsAffected,
	}
	db.Scopes(Paginate(int(req.Page), int(req.Size))).Find(&users)
	for _, us := range users {
		userInfoResp := ModelToResponse(us)
		resp.Data = append(resp.Data, &userInfoResp)
	}
	return resp, nil
}
func (u *UserServer) GetUserMobile(c context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	//手机号查询用户
	var user model.User
	result := db.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "该用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	resp := ModelToResponse(user)

	return &resp, nil
}
func (u *UserServer) GetUserId(c context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	//id查询用户
	var user model.User
	result := db.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "该用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	resp := ModelToResponse(user)

	return &resp, nil
}
func (u *UserServer) CreateUserInfo(c context.Context, req *proto.CreateUserRequest) (*proto.UserInfoResponse, error) {
	//创建用户
	var user model.User
	check := db.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if check.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "用户已存在")
	}
	user.Mobile = req.Mobile
	user.NickName = req.NickName
	user.Password = util.BcryptHash(req.Password)
	result := db.Create(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}
	resp := ModelToResponse(user)
	return &resp, nil
}
func (u *UserServer) EditUserInfo(c context.Context, req *proto.EditUserRequest) (*empty.Empty, error) {
	//更新信息
	var user model.User
	check := db.First(&user, req.Id)
	if check.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "该用户不存在")
	}
	birthday := time.Unix(int64(req.Birthday), 0)
	user.NickName = req.NickName
	user.Sex = int(req.Sex)
	user.Birthday = &birthday
	result := db.Save(user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}
func (u *UserServer) CheckPassword(c context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	check := util.BcryptCheck(req.Password, req.EncryptedPassword)
	return &proto.CheckResponse{Success: check}, nil
}
