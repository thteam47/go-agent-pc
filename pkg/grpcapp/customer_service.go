package grpcapp

import (
	"context"
	pb "github.com/thteam47/common/api/agent-pc"
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/models"
	"github.com/thteam47/go-agent-pc/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func makeTenant(item *models.Tenant) (*pb.Tenant, error) {
	if item == nil {
		return nil, nil
	}
	survey := &pb.Tenant{}
	err := util.ToMessage(item, survey)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return survey, nil
}

func makeTenants(items []models.Tenant) ([]*pb.Tenant, error) {
	surveys := []*pb.Tenant{}
	for _, item := range items {
		survey, err := makeTenant(&item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeSurvey")
		}
		surveys = append(surveys, survey)
	}
	return surveys, nil
}
func (inst *AgentpcService) GetAllTenant(ctx context.Context, req *pb.Request) (*pb.ListTenantResponse, error) {
	if req.Ctx.TokenAgent == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}
	userContext := entity.NewUserContext("default")
	userInfo := inst.componentsContainer.IdentityAuthenService().GetUserInfo(userContext, req.Ctx.TokenAgent)
	if userInfo == nil {
		return nil, status.Errorf(codes.Unauthenticated, "userInfo Unauthenticated")
	}
	accessToken := inst.componentsContainer.IdentityAuthenService().AccessToken(userContext, req.Ctx.TokenAgent)

	if accessToken == "" {
		return nil, status.Errorf(codes.Unauthenticated, "accessToken Unauthenticated")
	}
	userContext.SetAccessToken(accessToken)
	tenants, err := inst.componentsContainer.CustomerService().GetAllTenantByCustomer(userContext)
	if err != nil {
		return nil, errutil.Wrap(err, "CustomerService.GetTenantById")
	}
	items, err := makeTenants(tenants)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	return &pb.ListTenantResponse{
		Data: items,
	}, nil
}
