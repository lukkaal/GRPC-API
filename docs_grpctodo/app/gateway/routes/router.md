###### sessions.Sessions 在 Gin 里创建并注册一个 Session 管理中间件



###### ginRouter.Use(sessions.Sessions("CookieSession", store))

为每个请求上加了一个自动识别和解密 Cookie 的中间件

检查请求里是否有 Cookie 名 "CookieSession"

如果有，就读取 Cookie 的值



###### cookie.NewStore(\[]byte(secretKey))：

定义 Session 的存储方式为存到客户端 Cookie 并加密

后续网关会使用 secretKey 解密 Cookie 中的 Session 数据

验证完整性，确保数据没有被篡改



###### 注册之后

每次请求自动解析 Cookie，获取 Session 数据

在处理请求过程中可以通过 sessions.Default(c) 访问、修改 Session 数据

请求结束后自动把修改后的 Session 数据写回 Cookie（加密后）

