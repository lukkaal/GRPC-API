package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// grpc service(server side) registered as etcd client
// GC: 存放在栈上 (Go GC/ Go 垃圾回收器)
type Register struct {
	leaseID   clientv3.LeaseID
	cli       *clientv3.Client
	keepalive <-chan *clientv3.LeaseKeepAliveResponse
	srvInfo   ServerNode
	srvTTL    int64

	logger *logrus.Logger

	EtcdAddrs   []string
	DialTimeout int

	closeCh chan struct{}
}

func NewRegister(etcdaddrs []string,
	logger *logrus.Logger) *Register {
	return &Register{
		logger:      logger,
		DialTimeout: 3,
		EtcdAddrs:   etcdaddrs,
	}
}

// register a service to etcd
func (r *Register) Register(srvinfo ServerNode,
	ttl int64) (chan<- struct{}, error) {
	var err error

	// IP:PORT
	if strings.Split(srvinfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip")
	}

	// init etcd client(preserve in Register)
	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints: r.EtcdAddrs,
		DialTimeout: time.Duration(
			r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	// set servernode & ttl
	r.srvInfo = srvinfo
	r.srvTTL = ttl

	// register servernode
	if err = r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})

	// new go routine to monitor keepalive
	go r.keepAlive()

	return r.closeCh, nil
}

// register service through etcd client
func (r *Register) register() error {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	// Grant creates a new lease(id)
	// expire after ttl
	leaseResp, err := r.cli.Grant(ctx, r.srvTTL)
	if err != nil {
		return err
	}

	r.leaseID = leaseResp.ID

	// r.keepAliveCh : read only
	if r.keepalive, err = r.cli.KeepAlive(
		context.Background(), r.leaseID); err != nil {
		return err
	}

	// transform struct to []byte -> register service
	data, err := json.Marshal(r.srvInfo)

	// register ServiceNode
	// key: exp: /user-service/v1/127.0.0.1:8080(unique) kv <= 1
	_, err = r.cli.Put(context.Background(),
		BuildRegisterPath(r.srvInfo), string(data),
		clientv3.WithLease(r.leaseID))

	return err
}

func (r *Register) keepAlive() {
	for {
		select {
		// check closeCh
		case <-r.closeCh:
			r.logger.Info("register stopped, exiting keepalive loop")
			return

		// loop check the chan
		case resp, ok := <-r.keepalive:
			if !ok {
				r.logger.Warn("keepalive channel closed, try re-register")

				// clean
				if _, err := r.cli.Revoke(context.Background(),
					r.leaseID); err != nil {
					r.logger.Error("revoke old lease failed: ", err)
				}

				if err := r.register(); err != nil {
					r.logger.Error("register failed: ", err)
					time.Sleep(time.Second)
				}
				continue
			}

			// dynamicly set weight through cpu/mem usage
			r.srvInfo.Weight = getDynamicWeight()
			data, err := json.Marshal(r.srvInfo)
			if err != nil {
				r.logger.Error("marshal server info failed: ", err)
				continue
			}

			// dynamicly change the weight per keepalive
			_, err = r.cli.Put(context.Background(), BuildRegisterPath(r.srvInfo),
				string(data), clientv3.WithLease(r.leaseID))
			if err != nil {
				r.logger.Error("update server info in etcd failed: ", err)
				continue
			}

			r.logger.Debugf("lease %d refreshed, TTL=%d, updated weight=%d",
				resp.ID, resp.TTL, r.srvInfo.Weight)

		}
	}
}

// revoke lease and clear registeration
// duty: clean goroutine/ etcd client/ lease
func (r *Register) Stop() {
	// r.closeCh <- struct{}{}
	// delete key in etcd
	if r.closeCh != nil {
		close(r.closeCh)
	}

	if err := r.unregister(); err != nil {
		r.logger.Error("unregister failed, error: ", err)
	}

	// revoke the lease
	_, err := r.cli.Revoke(context.Background(), r.leaseID)
	if err != nil {
		r.logger.Error("revoke failed, error: ", err)
	}

	if r.cli != nil {
		r.cli.Close()
	}
}

func (r *Register) unregister() error {
	_, err := r.cli.Delete(
		context.Background(), BuildRegisterPath(r.srvInfo))
	return err
}

// get service
func (r *Register) GetServerInfo() (ServerNode, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), 3*time.Second)
	defer cancel()

	serviceInfo := &ServerNode{}
	resp, err := r.cli.Get(ctx, BuildRegisterPath(r.srvInfo))
	if err != nil {
		// empty struct
		return *serviceInfo, err
	}

	// unlikely to have resp.Count > 1
	if resp.Count >= 1 {
		// resp.Kvs: key and value
		if err := json.Unmarshal(resp.Kvs[0].Value, serviceInfo); err != nil {
			return *serviceInfo, err
		}
	}

	return *serviceInfo, nil
}

// 手动修改 weight 的逻辑 待修改
// change weight through http handler func
func (r *Register) UpdateHandler() http.HandlerFunc {
	handler := func(w http.ResponseWriter, req *http.Request) {

	}

	return handler
}
