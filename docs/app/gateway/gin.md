###### gin 和相关的 handler

###### 

###### 简单参数 → ctx.Query / ctx.Param / ctx.GetHeader

路由参数 /task/:id \&\& URL 查询参数 ?page=1



###### 复杂结构体 → ctx.Bind / ctx.ShouldBind

根据 Content-Type 自动选择绑定方式，会在失败时返回 400



###### 跨层信息 → ctx.Request.Context() + WithValue

###### 

###### 自定义返回/控制 → ctx.Writer / c.Status / c.Header





###### ctx.Request.Context() 底层 ctx 设值	

获取原生 context.Context（可传递超时、取消信号、Value/ 可以传递给 DAO 层、RPC 调用、异步 goroutine

**JWT 中间件里，可以用 WithValue 把 userID 放到 context，这样 DAO 或 gRPC 调用能拿到**

context.WithValue(ctx.Request.Context(), key, value) 





###### ctx.Set("key", value) Gin 自身存值 

**中间件之间共享状态、handler 内使用 / 中间件链内部使用（比如日志记录、校验标记）**

-> ctx.Get("key") / ctx.GetInt("key")

