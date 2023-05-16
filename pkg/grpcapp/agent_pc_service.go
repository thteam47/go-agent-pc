package grpcapp

import (
	"context"

	pb "github.com/thteam47/common/api/agent-pc"
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/component"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AgentpcService struct {
	pb.AgentpcServiceServer
	componentsContainer *component.ComponentsContainer
}

func NewAgentpcService(componentsContainer *component.ComponentsContainer) *AgentpcService {
	return &AgentpcService{
		componentsContainer: componentsContainer,
	}
}

func (inst *AgentpcService) GetById(ctx context.Context, req *pb.Request) (*pb.UserResponse, error) {
	if req.Ctx.TokenAgent == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}
	userContext := entity.NewUserContext("default")
	userInfo := inst.componentsContainer.IdentityAuthenService().GetUserInfo(userContext, req.Ctx.TokenAgent)
	if userInfo == nil {
		return nil, status.Errorf(codes.Unauthenticated, "userInfo Unauthenticated")
	}
	userApi, err := makeUser(userInfo)

	if err != nil {
		return nil, errutil.Wrap(err, "makeUser")
	}
	return &pb.UserResponse{
		Data: userApi,
	}, nil
}
