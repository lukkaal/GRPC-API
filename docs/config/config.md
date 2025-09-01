#### config

config 文件是 yaml 格式的 存在分级

同时在定义结构体的时候 使用相关的 tag `yaml:"server"` 



对于文件中存在同一词条下多个对象的场景：

比如：

services:

&nbsp; gateway:

&nbsp;   name: gateway

&nbsp;	...

&nbsp; user:

&nbsp;   name: user

&nbsp;	...

在定义顶层结构体的时候 应该使用 map

如 Services map\[string]\*Service 

同时也定义出 Service 结构体

