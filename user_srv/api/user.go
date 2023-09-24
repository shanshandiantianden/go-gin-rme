package api

import (
	"context"
	"go-gin-rme/user_srv/global"
	"go-gin-rme/user_srv/model"
	"go-gin-rme/user_srv/proto"
	"gorm.io/gorm"
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
