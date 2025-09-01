##### rpc 服务

实现的时候 务必包含



taskpb.UnimplementedTaskServiceServer



本质上是一个空结构体 但是实现了 rpc 的接口

所以需要使用匿名嵌套的方法



type TaskSrv struct {

 	taskpb.UnimplementedTaskServiceServer

}



进行嵌套
相当于这里进行了第二层的封装：第一层是 init db 并且完成 db 的逻辑操作
第二层封装是这里：实现了server side 的接口函数

##### 在实际上线过程中 会将 TaskSrv 注册为 grpc server：
server := grpc.NewServer()
taskpb.RegisterTaskServiceServer(server, &TaskSrv)
