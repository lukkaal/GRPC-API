package rpc

import (
	"fmt"

	"github.com/lukkaal/GRPC-API/config"
	taskpb "github.com/lukkaal/GRPC-API/idl/task"
	userpb "github.com/lukkaal/GRPC-API/idl/user"
	"github.com/lukkaal/GRPC-API/pkg/discovery"
	"github.com/lukkaal/GRPC-API/pkg/utils/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

// client (global)
var (
	Register   *discovery.Resolver // etcd resolver
	UserClient userpb.UserServiceClient
	TaskClient taskpb.TaskServiceClient

	// ctx: for Dial() when make conn
	// ctx        context.Context
	// CancelFunc context.CancelFunc
)

// using resolver to init grpc client
func Init() {
	// make etcd conn
	Register := discovery.NewResolver([]string{config.Conf.Etcd.Address}, logrus.New())
	resolver.Register(Register) // init to local
	defer Register.Close()

	// init Client(with)
	initClient(config.Conf.Domain["user"].Name, &UserClient)
	initClient(config.Conf.Domain["task"].Name, &TaskClient)
}

func initClient(serviceName string, client interface{}) {
	conn, err := connectServer(serviceName)
	if err != nil {
		logger.GinloggerObj.Panic(err.Error())
	}

	switch c := client.(type) {
	case *userpb.UserServiceClient:
		*c = userpb.NewUserServiceClient(conn)
	case *taskpb.TaskServiceClient:
		*c = taskpb.NewTaskServiceClient(conn)
	default:
		logger.GinloggerObj.Panic("unsupport type")
	}
}

// conn
func connectServer(serviceName string) (
	conn *grpc.ClientConn, err error) {
	// option only for dev
	// how we set up grpc service connection
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	opts = append(opts, TLSDial()...)

	// concat addr
	addr := fmt.Sprintf("%s:///%s", Register.Scheme(), serviceName)

	// check if loadbalance is true(local check:?)
	if config.Conf.Services[serviceName].LoadBalance {
		// logger.LogrusObj.Printf("load balance enabled for %s\n",
		// 	serviceName)
		opts = append(opts, grpc.
			WithDefaultServiceConfig(fmt.Sprintf(
				`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	}

	// lazy dial: only dial when first call rpc
	// or Dial()
	conn, err = grpc.NewClient(addr, opts...)
	return

}

func TLSDial() []grpc.DialOption {
	creds, err := credentials.NewClientTLSFromFile("cert.pem", "")
	if err != nil {
		logger.GinloggerObj.Error(err.Error())
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	return opts
}
