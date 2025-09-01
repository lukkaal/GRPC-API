package main

import (
	"fmt"
	"net"

	"github.com/lukkaal/GRPC-API/app/task/internal/repository/taskdb"
	"github.com/lukkaal/GRPC-API/app/task/internal/service"
	"github.com/lukkaal/GRPC-API/config"
	taskpb "github.com/lukkaal/GRPC-API/idl/task"
	"github.com/lukkaal/GRPC-API/pkg/discovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// init viper config
	config.InitConfig()

	// init db conn
	taskdb.InitDB()

	// etcd addrs
	etcdAddress := []string{config.Conf.Etcd.Address}

	// init register
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := config.Conf.Services["task"].Addr[0]
	defer etcdRegister.Stop() // like resolver.Close

	// new tasknode
	taskNode := discovery.ServerNode{
		Name: config.Conf.Domain["task"].Name,
		Addr: grpcAddress,
	}

	// init grpc server
	server := grpc.NewServer()
	defer server.Stop()

	// register server
	taskpb.RegisterTaskServiceServer(server, service.GetTaskSrv())

	// start to listen on ip:port
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}

	// register a service to etcd
	if _, err := etcdRegister.Register(taskNode, 10); err != nil {
		panic(fmt.Sprintf("start server failed, err: %v", err))
	}
	logrus.Info("server started listen on ", grpcAddress)

	// start to listen and accept service
	if err := server.Serve(lis); err != nil {
		panic(err)
	}

}
