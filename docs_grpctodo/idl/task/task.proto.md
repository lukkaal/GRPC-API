##### task.proto

message:

task 本身的作用是记录任务 因此rpc服务包含有 request 和 response

对于 request 而言会记录任务创建的具体信息

对于 response 来说 重要的是记录 rpc 方法调用(处理 task)的状态信息



详细一些的 需要展示某个 userid 下的所有 task

因此还需要定义出 taskmodel 和 对应的 detail response

会包含 repeated TaskModel -> 即 \[]\*TaskModel



service

TaskService: 

rpc TaskCreate(TaskRequest) returns(TaskCommonResponse);

rpc TaskUpdate(TaskRequest) returns(TaskCommonResponse);

rpc TaskShow(TaskRequest) returns(TasksDetailResponse);

rpc TaskDelete(TaskRequest) returns(TaskCommonResponse);

编译之后生成对应的接口 interface{} 之后需要自行实现这些 rpc 函数签名

这些方法的本质其实是包含了 gorm 相关的 db 操作

会在 app/task/internal 当中实现具体的 db 和 rpc 句柄



