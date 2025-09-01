##### automigrate

数据库初始化的时候

为了建表的安全性需要加上

Set("gorm:table\_options", "charset=utf8mb4").

AutoMigrate(\&model.Task{}，)

配合 db_init 使用