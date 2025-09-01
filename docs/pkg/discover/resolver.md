resolver 是用来在 gateway 处注册一个 etcd 服务发现的结构体

在本地注册(register)之后 每次调用 grpc.NewClient(addr, opts...) 

会根据 scheme 主动触发对应 resolver.Build 创建 resolver 实例对象

进行 etcd client 的初始化和远程服务发现  





**配合 ./app/gateway/rpc 和  ./app/gateway/cmd/main 来进行讲解:**



对于一个 resolver 而言 需要实现三个接口:

resolver.Builder 创建 resolver

resolver.Resolver 执行解析和后台的 watch 更新

resolver.clientConn 抽象的 grpc 连接 用于更新



type Resolver struct {

&nbsp;	schema      string

&nbsp;	EtcdAddrs   \[]string

&nbsp;	DialTimeout int



&nbsp;	closeCh        chan struct{}

&nbsp;	watchCh        clientv3.WatchChan // monitor the key(Prefix)

&nbsp;	watchCtxCancel context.CancelFunc // cancel func for watch contex



&nbsp;	cli    \*clientv3.Client

&nbsp;	cc     resolver.ClientConn // interface(literally \*grpc.clientconn)

&nbsp;	logger \*logrus.Logger



&nbsp;	keyPrefix   string

&nbsp;	srvAddrList \[]resolver.Address



&nbsp;	mu sync.RWMutex

}





首先是 resolver.Builder 需要实现两个函数签名: scheme 和 build

1）前者是返回一个 scheme 比如说 etcd

也就是 <scheme>://<authority>/<endpoint> 中的第一个

当使用 grpc.Dial("etcd:///user/v1", ... 会根据这个 scheme 指定执行对应的 resolver

2）后者是一个 resolver 实例对象具体初始化的执行流程 需要自行实现

包括连接客户端和开启 watch 

&nbsp;



对于 resolver.Resolver 需要实现 ResolveNow 和 close 函数

前者是立刻执行一次 resolve 后者是关闭或者说释放掉这个 resolver





resolver.clientConn 是一个接口 实际上是 grpc.Dial 所返回的 \*grpc.clientconn

也就是说 每一次的对服务地址更新的操作 执行的 UpdateState

都是直接更新到 \*grpc.clientconn 内部维护的 grpcaddrs 当中





对于 build 而言：

target resolver.Target 和 cc resolver.ClientConn 作为参数传入

**执行 grpc.Dial("etcd:///user/" ... 的时候 会将 etcd:///user/ 保存为 target**

**同时返回值 \*grpc.clientconn 本身作为 cc resolver.ClientConn 接口的具体实现**

1）target.Endpoint() 进行解析: 一般不会有 v1 

所以将解析得到的 user 转换为 ServerNode 后设置为 r.keyPrefix

r.keyPrefix 的 exp: /user/

2）执行 resolver 的 start 函数



start 函数的具体组成：
1）初始化 etcd client(和 register 的操作一样 使用 clientv3.New)

2）初始化 r.cc = cc 在结构体内部保存连接句柄

3）执行 r.sync() 传入 keyPrefix 解析到一次 grpc addrs

4）开启 goroutine 执行 r.watch()



重点说明 sync 和 watch

1）sync: (mutex 线程安全)

r.cli.Get(ctx, r.keyPrefix, clientv3.WithPrefix()) 传入 keyPrefix 获取到和这个 prefix 相关的所有 kv

for \_, v := range res.Kvs 使用一个循环 获得到所有的 v (var v \*mvccpb.KeyValue)

获取 v.Value (\[]byte) 并转换为 ServerNode (json.Unmarshal - 对应 register 当中的 json.Marshal)



提取 ServerNode 当中的服务地址 存放到 resolver.Address.Addr 当中 并且 append 到 list 当中

最后 r.cc.UpdateState(resolver.State{Addresses: r.srvAddrList}) 传入这些 grpc 的 addrs





2）watch:

开一个额外的协程用来开启 watch ch 并且持续监听指定前缀的 kvs 变化

r.watchCh = r.cli.Watch(ctx, r.keyPrefix, clientv3.WithPrefix())

也就是说 会监听指定前缀的所有 kv 的变化  r.keyPrefix

本地保存这个 Ch (read only)



陷入 for\_select 的开启之后三种情况:

case resp, ok := <-r.watchCh 也就是有更新发生 可能是删除或者更新(del/ update)

这时候需要重新访问包含前缀的 etcd server 端 得到更新之后的地址 (详见 .update)





case <-r.closeCh 说明此时要释放掉这个 resolver 了

此时退出循环 执行完毕 goroutine 避免资源泄露



case <-ticker.C:

定时执行一次 sync() 更新一下列表





update:

r.update(resp.Events) 

本质还是 for \_, ev := range events 加上一个

switch ev.Type {

&nbsp;		case clientv3.EventTypePut:

&nbsp;		case clientv3.EventTypeDelete:

然后分别触发一次 

r.cc.UpdateState(resolver.State{Addresses: r.srvAddrList},)



只不过一个是 ev.Kv.Value 一个是(sync)  \*mvccpb.KeyValue.Value

前者是来自 watchCh 后者是 r.cli.Get





Close

释放资源和 register 类似

case <-r.closeCh 关闭本地的 watch 轮询

r.watchCtxCancel() 释放 etcd client 创建的 watchCh 的后台协程

r.cli.Close() 关闭 etcd client (也会清理 watchCh 但是可能不够及时 所以要使用 cancel func)



执行 main 的时候 在 rpc.Init 当中 也会使用

Register := discovery.NewResolver(\[]string{config.Conf.Etcd.Address}, logrus.New())

resolver.Register(Register) // init to local

defer Register.Close()



在本地执行 register 以便后面执行 dial 可以被发现 

同时 defer 清理结束后的资源





