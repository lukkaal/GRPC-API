##### config

config 文件夹存储了 config.yaml 和 config.go

前者是用来存储相关信息的

后者存在一个 var 并且使用函数将 yaml 中的变量

使用 viper 的 unmarshal 之后可以访问变量成员

从而获取到 yaml 中各信息

如 viper.GetString("server.jwtSecret")