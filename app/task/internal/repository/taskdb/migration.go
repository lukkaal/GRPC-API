package taskdb

import (
	"os"

	"github.com/lukkaal/GRPC-API/app/task/internal/repository/taskmodel"
	"github.com/lukkaal/GRPC-API/pkg/utils/logger"
)

func migration() {
	err := _MysqlDb_task.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&taskmodel.Task{},
		)
	if err != nil {
		logger.GinloggerObj.Infoln("register table fail")
		os.Exit(0) // successfully exit
	}
	logger.GinloggerObj.Infoln("register table success")
}
