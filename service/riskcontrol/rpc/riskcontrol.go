package main

import (
	"flag"
	"fmt"

	"github.com/Drengin6306/ZeroBank/service/riskcontrol/rpc/internal/config"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/rpc/internal/server"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/rpc/internal/svc"
	"github.com/Drengin6306/ZeroBank/service/riskcontrol/rpc/proto"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/riskcontrol.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		proto.RegisterRiskControlServer(grpcServer, server.NewRiskControlServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	consul.RegisterService(c.ListenOn, c.Consul)
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
