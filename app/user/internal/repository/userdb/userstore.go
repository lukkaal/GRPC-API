package userdb

import (
	"context"
	"errors"

	"github.com/lukkaal/GRPC-API/app/user/internal/repository/usermodel"
	userpb "github.com/lukkaal/GRPC-API/idl/user"
	"gorm.io/gorm"
)

type UserStore struct {
	*gorm.DB
}

// get gorm.DB instance after InitDB (global real conn)
func NewUserStore(ctx context.Context) *UserStore {
	userstore := &UserStore{
		NewDBClient(ctx),
	}
	return userstore
}

// login
func (userstore *UserStore) GetuserInfo(req *userpb.
	LoginRequest) (r *usermodel.User, err error) {
	err = userstore.Model(&usermodel.User{}).
		Where("user_name=?", req.UserName).
		First(&r).Error
	return
}

// register
func (userstore *UserStore) CreateUser(
	req userpb.RegisterRequest) (err error) {
	var user usermodel.User
	var count int64

	userstore.Model(&usermodel.User{}).Where(
		"user_name=?", req.UserName).Count(&count)
	if count != 0 {
		return errors.New("Username already exist")
	}

	user = usermodel.User{
		UserName: req.UserName,
	}

	if err = user.EncryptPassword(req.Password); err != nil {
		return err
	}

	err = userstore.Model(&usermodel.User{}).Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

// no DB lever operation when logout
