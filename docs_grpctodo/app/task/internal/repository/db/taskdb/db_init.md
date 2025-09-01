##### db\_init

##### 初始化数据库

用到了 mysql 和 gorm

gorm.Open() 会需要两个参数的传入



###### mysql.New(mysql.Config{})

用于配置 mysql 底层驱动层相关的设置



###### gorm.Config{}

配置和 gorm 映射相关的规则

比如 gorm 的 logger 

以及命名规则如 singulartable



###### 特别的 对于 gorm 的 logger

"gorm.io/gorm/logger"

配置的时候 可以根据 gin.Mode()选择

对应的模式 

声明 logger.Interface -> logger.Default



完成 gorm.Open() 之后 

对 \*gorm,.DB 调用 .DB() 获得驱动句柄

完成底层连接池的连接配置后 

即可返回 gorm 句柄 \*gorm.DB



###### 后续封装 \*gorm.DB 进结构体

可以为 \*gorm.DB 加上 ctx 后返回

使用 gorm 完成 db 相关 curd 函数的封装



