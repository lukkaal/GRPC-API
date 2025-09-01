package discovery

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

const (
	schema = "etcd"
)

type Resolver struct {
	schema      string
	EtcdAddrs   []string
	DialTimeout int

	closeCh        chan struct{}
	watchCh        clientv3.WatchChan // monitor the key(Prefix)
	watchCtxCancel context.CancelFunc // cancel func for watch contex

	cli    *clientv3.Client
	cc     resolver.ClientConn // interface(literally *grpc.clientconn)
	logger *logrus.Logger

	keyPrefix   string
	srvAddrList []resolver.Address

	mu sync.RWMutex
}

func NewResolver(etcdAddrs []string,
	logger *logrus.Logger) *Resolver {
	return &Resolver{
		schema:      schema,
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// resolver.Builder
// resolver.Resolver
// resolver.clientConn

func (r *Resolver) Scheme() string {
	return r.schema
}

// grpc.Dial("etcd:///user/v1", ... ->
// target.Scheme → "etcd"
// target.Endpoint → "user/v1"

// create(complete) Resolver instance
func (r *Resolver) Build(target resolver.Target,
	cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	r.cc = cc

	// split the string into service/version
	parts := strings.Split(target.Endpoint(), "/")
	serviceName := parts[0] // "user"
	version := ""
	if len(parts) > 1 {
		version = parts[1] // "v1"
	}

	// exp: /user/v1/(need prefix only)
	r.keyPrefix = BuildPrefix(ServerNode{
		Name:    serviceName, // "user/v1" -> target.Endpoint
		Version: version})    // "/user/v1" -> target.URL.Path
	if _, err := r.start(); err != nil {
		return nil, err
	}
	return r, nil
}

// resolver.Resolver (interface)
func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) { _ = r.sync() }

// concurrency safety
// duty: close watchCh/ release etcd client
func (r *Resolver) Close() {
	select {
	// close watchCh(local, not etcd watch) through closeCh
	case <-r.closeCh:
		// has been closed
	default:
		close(r.closeCh)
	}

	// cancel watch context to stop (etcd client watch) goroutine
	if r.watchCtxCancel != nil {
		r.watchCtxCancel()
	}

	// close etcd client
	if r.cli != nil {
		r.cli.Close()
	}
}

// resolve the addr from etcd
func (r *Resolver) start() (chan<- struct{}, error) {
	var err error

	r.cli, err = clientv3.New(clientv3.Config{
		Endpoints: r.EtcdAddrs,
		DialTimeout: time.Duration(
			r.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}

	// register locally
	// resolver.Register(r)

	r.closeCh = make(chan struct{})

	if err = r.sync(); err != nil {
		return nil, err
	}

	go r.watch()

	return r.closeCh, nil
}

// watch keyprefix events change
func (r *Resolver) watch() {
	// ticker per min
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop() // stop when close watch goroutine

	// create cancelable context for watch
	ctx, cancel := context.WithCancel(context.Background())
	r.watchCtxCancel = cancel
	// ctx: when to cease watchCh
	r.watchCh = r.cli.Watch(ctx, r.keyPrefix, clientv3.WithPrefix())

	for {
		select {
		// events occur
		case resp, ok := <-r.watchCh:
			if ok {
				r.update(resp.Events)
			}

		case <-r.closeCh:
			return

		// sync per min
		case <-ticker.C:
			if err := r.sync(); err != nil {
				r.logger.Error("sync failed", err)
			}
		}
	}
}

// update
func (r *Resolver) update(events []*clientv3.Event) {
	// events from watchCh
	for _, ev := range events {
		var info ServerNode
		var err error

		switch ev.Type {
		case clientv3.EventTypePut:
			info, err = ParseValue(ev.Kv.Value)
			if err != nil {
				r.logger.Error("parse false occurs: ", err)
			}
			// concurrency safety
			r.mu.Lock()
			defer r.mu.Unlock()

			// init resolver.Address
			addr := resolver.Address{
				Addr:     info.Addr,
				Metadata: info.Weight,
			}
			if !Exist(r.srvAddrList, addr) {
				r.srvAddrList = append(r.srvAddrList, addr)
				r.cc.UpdateState(
					resolver.State{Addresses: r.srvAddrList},
				)
			}

		case clientv3.EventTypeDelete:
			info, err = SplitPath(string(ev.Kv.Key))
			if err != nil {
				r.logger.Error("split false occurs: ", err)
			}
			// concurrency safety
			r.mu.Lock()
			defer r.mu.Unlock()

			addr := resolver.Address{
				Addr: info.Addr,
			}
			if s, ok := Remove(r.srvAddrList, addr); ok {
				r.srvAddrList = s
				r.cc.UpdateState(
					resolver.State{Addresses: r.srvAddrList},
				)
			}
		}
	}
}

// get data with keyprefix
func (r *Resolver) sync() error {
	ctx, cancel := context.WithTimeout(
		context.Background(), 3*time.Second)
	defer cancel()

	res, err := r.cli.Get(ctx, r.keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	// concurrency safety
	r.mu.Lock()
	defer r.mu.Unlock()

	// initiate
	r.srvAddrList = []resolver.Address{}

	for _, v := range res.Kvs {
		// *mvccpb.KeyValue -> ServiceNode
		info, err := ParseValue(v.Value)
		if err != nil {
			continue
		}
		addr := resolver.Address{
			Addr:     info.Addr,
			Metadata: info.Weight}

		r.srvAddrList = append(r.srvAddrList, addr)
	}

	// inform grpc client the addrs of grpc server
	// Addresses → []resolver.Address (availble nodes)
	r.cc.UpdateState(resolver.State{
		Addresses: r.srvAddrList})

	return nil
}
