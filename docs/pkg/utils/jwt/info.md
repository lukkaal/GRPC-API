##### JWT

Token（尤其是 JWT）是临时的、可随时失效的

没必要常驻数据库。

Token 是认证信息，不属于用户的核心业务数据

（用户名、密码、昵称等才是）。

如果用 JWT，它本身已经包含过期时间、签名等信息

后端只需验证即可，无需存储。

如果把 token 存在 user 表里，每次用户登录都要更新它

会增加数据库写压力。



方案 A：JWT（无状态）

不存 token

登录时生成 JWT → 返回给客户端

之后客户端每次请求携带 token

（放在 HTTP Header 的 Authorization）

后端用密钥验证 token 的有效性

优点：数据库无感知；扩展性好；分布式支持更好



方案 B：需要手动让 token 失效（有状态）

在 Redis 中维护 token（或 sessionID）→ 可快速失效

用户登录后生成 token，保存到 Redis（设置 TTL）

退出登录或需要强制下线时，直接删除 Redis 里的 token

优点：可控性强，支持手动踢人；缺点是需要额外维护 Redis

