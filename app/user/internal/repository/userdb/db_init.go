package userdb

import (
	"context"
	"fmt"
	"time"

	"github.com/lukkaal/GRPC-API/config"
	"github.com/lukkaal/GRPC-API/pkg/utils/logger"

	"gorm.io/driver/mysql"
	gorm_logger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var _MySQLDB_user *gorm.DB

func InitDB() {
	// concat destination
	mConfig := config.Conf.MySQL
	host := mConfig.Host
	port := mConfig.Port
	database := mConfig.Database
	username := mConfig.UserName
	password := mConfig.Password
	charset := mConfig.Charset
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username, password, host, port, database, charset)

	if err := Database(dsn); err != nil {
		fmt.Println(err)
		// temporarily using the gin logger(logrus)
		logger.GinloggerObj.Fatal(err)
	}
}

// make connection to db
func Database(connString string) error {
	var ormLogger gorm_logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = gorm_logger.Default.LogMode(gorm_logger.Info)
	} else {
		ormLogger = gorm_logger.Default
	}

	// set db connection
	mysqlDial := mysql.New(mysql.Config{
		DSN:                       connString,
		DefaultStringSize:         256,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	})

	// gorm config
	gormConf := gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	// make conn
	db, err := gorm.Open(mysqlDial, &gormConf)
	if err != nil {
		return err
	}
	_MySQLDB_user = db

	sqldb, err := _MySQLDB_user.DB()
	if err != nil {
		return err
	}

	sqldb.SetConnMaxLifetime(time.Second * 30)
	sqldb.SetMaxIdleConns(20)
	sqldb.SetMaxOpenConns(100)

	// automigrate
	migration()
	return nil
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _MySQLDB_user
	return db.WithContext(ctx)
}
