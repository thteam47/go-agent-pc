package main

import (
	"fmt"
	"log"
	"net"

	"github.com/thteam47/common-libs/confg"

	"github.com/thteam47/common/entity"
	"github.com/thteam47/common/handler"

	"github.com/soheilhy/cmux"
	"github.com/thteam47/common/pkg/confgutil"
	clienthttp "github.com/thteam47/go-agent-pc/cmd/client"
	servergrpc "github.com/thteam47/go-agent-pc/cmd/server"
	"github.com/thteam47/go-agent-pc/errutil"
)

func client(lis net.Listener, grpc_port string, http_port string) {
	err := clienthttp.Run(lis, grpc_port, http_port)
	if err != nil {
		fmt.Println(err)
	}
}

func serverGrpc(lis net.Listener, properties confg.Confg, handler *handler.Handler) {
	err := servergrpc.Run(lis, properties, handler)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	properties, err := confgutil.GetPropertiesFromArgs()
	if err != nil {
		fmt.Println(errutil.Wrap(err, "confgutil.GetPropertiesFromArgs"))
	}
	cf := entity.Config{}
	err = properties.Unmarshal(&cf)
	if err != nil {
		log.Fatalln("Unmarshal", err)
	}
	handle, err := handler.NewHandlerWithConfig(&cf)
	if err != nil {
		log.Fatalln("NewHandlerWithConfig", err)
	}

	lis, err := net.Listen("tcp", cf.GrpcPort)
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Server run on", cf.GrpcPort)

	m := cmux.New(lis)
	// a different listener for HTTP1
	httpL := m.Match(cmux.HTTP1Fast())
	grpcL := m.Match(cmux.HTTP2())
	go serverGrpc(grpcL, properties, handle)
	go client(httpL, cf.GrpcPort, cf.HttpPort)
	err = m.Serve()
	if err != nil {
		log.Fatalln("Serve", err)
	}
}
