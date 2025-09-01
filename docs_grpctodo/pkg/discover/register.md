register 是用来向 etcd 集群注册一个服务的实例对象的

注册前提：即将传入的 etcdaddrs 的 etcd servers 都已经正确初始化并暴露端口 



type Register struct {

&nbsp;	leaseID   clientv3.LeaseID

&nbsp;	cli       \*clientv3.Client

&nbsp;	keepalive <-chan \*clientv3.LeaseKeepAliveResponse

&nbsp;	srvInfo   ServerNode

&nbsp;	srvTTL    int64



&nbsp;	logger \*logrus.Logger



&nbsp;	EtcdAddrs   \[]string

&nbsp;	DialTimeout int



&nbsp;	closeCh chan struct{}

}



register 是作为 etcd client 向 etcd 注册服务

同时保持 keepalive 使服务对象定时在 etcd 续约 保证发现





**配合 ./app/task/cmd/main 来进行讲解:**

1）初始化 Register 结构体：

从 yaml 文件中读取获得 etcdaddrs 集群的地址 并使用 logrus 作为 logger

作为参数传入 Register 后返回实例



**main 函数中 也是先初始化 Register 实例对象 然后对 ServerNode** 

**也就是需要注册的 service 细节进行初始化:
从 yaml 文件中获取 rpc 服务本机的地址 和 服务名称**



2）使用自定义的 Register 初始化 etcd client 

&nbsp;clientv3.New(clientv3.Config{Endpoints，DialTimeout}) 

结构体内部保存 client 连接句柄

同时传入 ServerNode 实例对象 保存为结构体内部变量

(其余变量如 ttl 等也被初始化)



然后执行 r.register() 以及 go r.keepAlive() 

对于 register 函数而言:

1）先使用 r.cli.Grant(ctx, r.srvTTL) 获得 lease 也就是 "租约"

后续需要不断 "续约" 来保证服务存活

2）获得租约后 开启 r.cli.KeepAlive (leaseid) 自动续约

返回一个 read only 的 channel 会有租约的响应: TTL/ ID

如果这个 channel 关闭的话 则说明失效

3）将 srvInfo (ServerNode) 进行 JSON 序列化

同时使用 BuildRegisterPath 为 srvInfo 构造 key

使用 r.cli.Put 对服务进行注册



对于 go r.keepAlive() 新开了一个协程:

for\_select 进行两个部分的轮询操作: 

case <-r.closeCh 和 case resp, ok := <-r.keepalive



第一个是用来查看是不是这个 register 收到了关闭的信号 

如果是的那么要停掉这个 goroutine 避免出现协程的资源泄露



第二个是来查看 keepalive chan 当中是不是出现问题了

如果是的话 要先撤销一下原来的租约 然后重新执行以下 Register



如果是正常的情况 那么要看一下当前本机的情况 (cpu/ memory)

可以根据情况动态更新一下 srvInfo 的 weight 参数

然后重新使用 r.cli.Put 覆盖上传一下这个 key 的 value



**main 函数中 在初始化了结构体之后 就立刻 defer etcdRegister.Stop()**

**然后使用 Register 进行 etcd 的注册**



**在释放资源的时候 其实 Register 本身不需要考虑 因为 golang 有 GC** 

**但是开启的 协程/ 远程连接的 etcd client** 

**以及  r.cli.KeepAlive 的底层后台协程 和 租约 都是需要手动释放的**



**所以 .Stop() 会** 

**close(r.closeCh) 关闭 channel 触发 keepalive 协程释放**

**r.cli.Delete 删除 key**

**以及 r.cli.Revoke 撤销 租约 也会释放 r.cli.KeepAlive 的底层后台协程**

**最后 r.cli.Close() 释放 client 连接**









