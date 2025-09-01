生成文件的内容说明

以 task 和 user 为例 分为 client 和 server 两个部分

分别会在 rpc server （也就是service）和 rpc client（gin 网关分发执行的 handler）、



###### 对于 server 端 如./user 和 ./task

因为 pb 文件夹中 service 是以接口的形式实现的 包裹了函数签名

所以对于server 端 比如 task 就需要手动实现这些接口的函数



###### type TaskSrv struct {

###### &nbsp;	taskpb.UnimplementedTaskServiceServer

###### }



1）其中 type UnimplementedTaskServiceServer struct{} 是一个空的结构体

但是这个结构体实现了 task 的 service 接口



2）所以这里通过 嵌入匿名结构体 的方式 会自动提升方法和变量

使得 TaskSrv 也成为一个 service 的具体实现

然后就可以为 TaskSrv 封装方法



(对于匿名结构体嵌入来说 方法提升的效果 

等价于 编译器在外层类型上 自动帮你写了一层 wrapper )



(如果外层类型既能用值又能用指针完整实现接口 通常会选择 指针嵌入)



###### 对于 client 端来说 比如说 ./app/gateway 当中的 user 和 task

###### 

###### UserClient userpb.UserServiceClient

###### TaskClient taskpb.TaskServiceClient 



是具体的服务连接句柄 连接服务端之后 可以直接 TaskClient.() 调用对应的函数

1）首先使用 grpc 包中的函数 grpc.NewClient 传入地址后获得 \*grpc.ClientConn(conn)

也就是底层的连接句柄指针

2）对于 taskpb 来说 实用 taskpb.NewTaskServiceClient 函数 传入 conn 

返回 taskpb.TaskServiceClient 传入到全局的变量 TaskClient (指针接收)

