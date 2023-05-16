package servergrpc

import (
	"net"

	defaultcomponent "github.com/thteam47/go-agent-pc/pkg/component/default"

	"github.com/thteam47/common-libs/confg"
	pb "github.com/thteam47/common/api/agent-pc"
	"github.com/thteam47/common/handler"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/component"
	"github.com/thteam47/go-agent-pc/pkg/grpcapp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(lis net.Listener, properties confg.Confg, handler *handler.Handler) error {
	componentFactory, err := defaultcomponent.NewComponentFactory(properties.Sub("components"), handler)
	if err != nil {
		return errutil.Wrap(err, "NewComponentFactory")
	}
	componentsContainer, err := component.NewComponentsContainer(componentFactory)
	if err != nil {
		return errutil.Wrap(err, "NewComponentsContainer")
	}
	serverOptions := []grpc.ServerOption{}
	s := grpc.NewServer(serverOptions...)
	pb.RegisterAgentpcServiceServer(s, grpcapp.NewAgentpcService(componentsContainer))
	reflection.Register(s)
	return s.Serve(lis)
}
