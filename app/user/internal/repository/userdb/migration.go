package userdb

import (
	"os"

	"github.com/lukkaal/GRPC-API/app/user/internal/repository/usermodel"
	"github.com/lukkaal/GRPC-API/pkg/utils/logger"
)

func migration() {
	// AutoMigrate:modify/create when init db
	err := _MySQLDB_user.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&usermodel.User{},
		)
	if err != nil {
		logger.GinloggerObj.Infoln("register table fail")
		os.Exit(0)
	}
	logger.GinloggerObj.Infoln("register table success")
}
