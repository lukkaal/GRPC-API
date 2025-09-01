// define the grpc funcs with db embeded
package service

import (
	"context"
	"sync"

	"github.com/lukkaal/GRPC-API/app/user/internal/repository/userdb"
	"github.com/lukkaal/GRPC-API/app/user/internal/repository/usermodel"
	userpb "github.com/lukkaal/GRPC-API/idl/user"
	"github.com/lukkaal/GRPC-API/pkg/errcode"
)

// grpc server side struct
type UserSrv struct {
	userpb.UnimplementedUserServiceServer
}

// global funcs
var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})

	return UserSrvIns
}

// implement rpc interface

// userlogin
func (u *UserSrv) UserLogin(
	ctx context.Context, req *userpb.LoginRequest) (
	resp *userpb.UserDetailResponse, err error) {
	resp = new(userpb.UserDetailResponse)
	resp.Code = errcode.SUCCESS

	var user_info *usermodel.User

	user_info, err = userdb.NewUserStore(ctx).GetuserInfo(req)
	if err != nil {
		resp.Code = errcode.ERROR
		return
	}

	// check pw
	if !user_info.CheckPassword(req.Password) {
		resp.Code = errcode.ErrorNotCompare
	}

	// set resp if verified
	resp.UserDetail.UserId = user_info.UserId
	resp.UserDetail.UserName = user_info.UserName

	return
}

// userregister
func (u *UserSrv) UserRegister(ctx context.Context,
	req *userpb.RegisterRequest) (
	resp *userpb.UserCommonResponse, err error) {

	resp = new(userpb.UserCommonResponse)
	resp.Code = errcode.SUCCESS

	// a db conn pool
	err = userdb.NewUserStore(ctx).CreateUser(*req)
	if err != nil {
		resp.Code = errcode.ERROR
	}
	resp.Msg = errcode.GetMsg(int(resp.Code))

	return
}

// logout
func (u *UserSrv) UserLogout(ctx context.Context,
	req *userpb.LogoutRequest) (
	resp *userpb.UserCommonResponse, err error) {
	return
}
