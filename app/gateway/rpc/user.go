package rpc

import (
	"context"
	"errors"

	userpb "github.com/lukkaal/GRPC-API/idl/user"
	"github.com/lukkaal/GRPC-API/pkg/errcode"
)

func UserLogin(ctx context.Context,
	req *userpb.LoginRequest) (
	resp *userpb.UserResponse, err error) {

	// the return would be UserDetailResponse
	resdetail, err := UserClient.UserLogin(ctx, req)
	if err != nil {
		return
	}

	if resdetail.Code != errcode.SUCCESS {
		err = errors.New("login error")
		resp = nil
		return
	}

	resp = resdetail.UserDetail
	return
}

func UserRegister(ctx context.Context,
	req *userpb.RegisterRequest) (
	resp *userpb.UserCommonResponse, err error) {
	resp, err = UserClient.UserRegister(ctx, req)
	if err != nil {
		return
	}

	if resp.Code != errcode.SUCCESS {
		err = errors.New(resp.Msg)
		return
	}
	return

}

// no execution for rpc
func UserLogout(ctx context.Context,
	req *userpb.LogoutRequest) (
	resp *userpb.UserCommonResponse, err error) {
	return
}
