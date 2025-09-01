package main

import (
	"fmt"
	"net"

	"github.com/lukkaal/GRPC-API/app/user/internal/repository/userdb"
	"github.com/lukkaal/GRPC-API/app/user/internal/service"
	"github.com/lukkaal/GRPC-API/config"
	userpb "github.com/lukkaal/GRPC-API/idl/user"
	"github.com/lukkaal/GRPC-API/pkg/discovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// init viper config
	config.InitConfig()

	// init db conn
	userdb.InitDB()

	// etcd addrs
	etcdAddress := []string{config.Conf.Etcd.Address}

	// init register
	etcdRegister := discovery.NewRegister(
		etcdAddress, logrus.New())
	grpcAddress := config.Conf.Services["user"].Addr[0]
	defer etcdRegister.Stop() // like resolver.Close

	// new tasknode
	taskNode := discovery.ServerNode{
		Name: config.Conf.Domain["user"].Name,
		Addr: grpcAddress,
	}

	// init grpc server
	server := grpc.NewServer()
	defer server.Stop()

	// register server
	userpb.RegisterUserServiceServer(
		server, service.GetUserSrv())

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}

	// register a service to etcd
	if _, err := etcdRegister.Register(taskNode, 10); err != nil {
		panic(fmt.Sprintf("start server failed, err: %v", err))
	}
	logrus.Info("server started listen on ", grpcAddress)

	if err := server.Serve(lis); err != nil {
		panic(err)
	}

}
