#### proto

proto 文件是grpc自行生成 pb.go 之前需要自己定义的文件

格式上需要注意三个地方



##### message：

会在 .go 文件当中被映射成为结构体 struct

同时  // @inject\_tag: json:"nick\_name" form:"nick\_name"

代表了会在生成的结构体中自动加入这里的各 tag

其中 repeated 关键字是在结构体中生成 \[] 切片



##### service：

service 会被映射成为 interface

同时会为出现的 rpc 关键词后面的函数生成相关的函数签名



##### option go\_package：

遵循“路径; 包名”的整体格式

路径代表了其他文件 import 的时候应该如何获取包

包名代表的是 proto 生成的 pb.go 的包是属于哪个包的



使用 protoc 对 proto 文件进行处理的时候 应该满足

protoc --go\_out=. --go-grpc\_out=. task.proto

即 --go-out  --go-grpc\_out 代表两个文件分别输出的地址

task.proto 代表了 proto 文件的地址

-I 代表寻找 .proto 文件的地址



注意 option go\_package 和 以上两指令不要重合

