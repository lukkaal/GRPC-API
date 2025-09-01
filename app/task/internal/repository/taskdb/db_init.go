package taskdb

import (
	"context"
	"fmt"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukkaal/GRPC-API/config"
	"github.com/lukkaal/GRPC-API/pkg/utils/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var _MysqlDb_task *gorm.DB

func InitDB() {
	mConfig := config.Conf.MySQL
	host := mConfig.Host
	port := mConfig.Port
	database := mConfig.Database
	username := mConfig.UserName
	password := mConfig.Password
	charset := mConfig.Charset
	dsn := strings.Join([]string{
		username, ":",
		password, "@tcp(", host, ":", port, ")/",
		database, "?charset=" + charset + "&parseTime=true"}, "")
	err := Database(dsn)
	if err != nil {
		fmt.Println(err)
		logger.GinloggerObj.Error(err)
	}
}

func Database(dsn string) error {
	var ormLogger gorm_logger.Interface

	if gin.Mode() == "debug" {
		ormLogger = gorm_logger.Default.LogMode(gorm_logger.Info)
	} else {
		ormLogger = gorm_logger.Default
	}

	mysql_Conf := mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	})

	gorm_conf := gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	db, err := gorm.Open(mysql_Conf, &gorm_conf)
	if err != nil {
		return err
	}
	_MysqlDb_task = db

	sqldb, err := _MysqlDb_task.DB()
	if err != nil {
		return err
	}

	sqldb.SetConnMaxLifetime(time.Second * 30)
	sqldb.SetMaxIdleConns(20)
	sqldb.SetMaxOpenConns(100)

	migration()
	return err

}

func NewDBClient(ctx context.Context) *gorm.DB {
	return _MysqlDb_task.WithContext(ctx)
}
