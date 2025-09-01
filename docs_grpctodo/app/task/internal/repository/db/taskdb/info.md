##### db/taskdao

这个文件夹使用 gorm.DB 封装了

数据库的 init 和 gorm 的数据库 curd 操作



首先是数据库的初始化 mysql.New 完成驱动层配置

同时 gorm.Config 完成 gorm 映射层配置

设置 gorm.DB 的连接池后 完成 gorm.DB 句柄封装



初始化数据库 在初始化建表时 使用 automigration

保证初始化后 表存在



gorm 映射需要定义映射结构体 

定义了结构体之后 需要指定 tag

格式为　`gorm:"type:longtext"`



为 gorm.DB 句柄实现相关的数据库功能

首先使用 匿名嵌入 的方式

嵌入 gorm.DB 结构体



type TaskStore struct {

&nbsp;	\*gorm.DB 

}



然后为 TaskStore 实现具体的 db 方法

比如 CreateTask 等





##### 总结

总结起来 实现的思路就是 

初始化 db 句柄并配置 返回 instance

定义 db 中存储的结构体 配合 gorm 的 tag

为 db 实现 gorm 和 mysql 的业务逻辑 以便后面实现

rpc 服务的函数签名





