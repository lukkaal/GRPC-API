package usermodel

import (
	"github.com/lukkaal/GRPC-API/consts"
	"golang.org/x/crypto/bcrypt"
)

// user (table) for database
type User struct {
	UserId            int64  `gorm:"primarykey"`
	UserName          string `gorm:"unique;not null"`
	PasswordEncrypted string `gorm:"type:varchar(255);not null"`
}

func (user *User) TableName() string {
	return "user"
}

// encrypt password
func (user *User) EncryptPassword(password string) error {
	bcrypted_str, err := bcrypt.GenerateFromPassword(
		[]byte(password), consts.PassWordCost)
	if err != nil {
		return err
	}

	user.PasswordEncrypted = string(bcrypted_str)
	return nil
}

// check the pw
func (user *User) CheckPassword(password string) bool {
	// get/compare the password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(password), []byte(user.PasswordEncrypted)); err != nil {
		return false
	}
	return true
}
