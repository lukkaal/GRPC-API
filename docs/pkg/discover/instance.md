重点关注出现的结构体: ServerNode

etcd 注册的时候 是遵循 kv 存储的



BuildRegisterPath:　

拼接 string 的 key

exp: /user-service/v1/127.0.0.1:8080

**这是作为 key 进行存储的**

**ServerNode 作为 value**



ParseValue:

将从 etcd 获取的 value 的 \[]byte(stream)

转换成 ServerNode 

