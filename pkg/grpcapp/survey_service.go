package grpcapp

import (
	"github.com/thteam47/common/entity"
	"context"

	pb "github.com/thteam47/common/api/agent-pc"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/models"
	"github.com/thteam47/go-agent-pc/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func makeSurvey(item *models.Survey) (*pb.Survey, error) {
	if item == nil {
		return nil, nil
	}
	survey := &pb.Survey{}
	err := util.ToMessage(item, survey)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	if survey != nil {
		survey.UserIdJoin = []string{}
		survey.UserIdCreate = ""
		survey.UserIdVerify = ""
	}
	return survey, nil
}

func makeSurveys(items []models.Survey) ([]*pb.Survey, error) {
	surveys := []*pb.Survey{}
	for _, item := range items {
		survey, err := makeSurvey(&item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeSurvey")
		}
		surveys = append(surveys, survey)
	}
	return surveys, nil
}
func (inst *AgentpcService) GetSurveyByTenant(ctx context.Context, req *pb.StringRequest) (*pb.ListSurveyResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "@any", "@any", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "CustomerService.GetTenantById")
	}
	if tenant == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant not found")
	}

	surveys, count, err := inst.componentsContainer.SurveyService().
		GetSurveysByDomainId(userContext.Clone().EscalatePrivilege(), req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyService.GetSurveysByUserId")
	}
	items, err := makeSurveys(surveys)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	return &pb.ListSurveyResponse{
		Data:  items,
		Total: count,
	}, nil
}

func (inst *AgentpcService) GetSurveyByUser(ctx context.Context, req *pb.Request) (*pb.ListSurveyResponse, error) {
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
	
	surveys, err := inst.componentsContainer.SurveyService().GetSurveysByUserId(userContext.SetDomainId(userInfo.DomainId).Clone().EscalatePrivilege(), userInfo.DomainId, userInfo.UserId)
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyService.GetSurveysByUserId")
	}
	items, err := makeSurveys(surveys)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	return &pb.ListSurveyResponse{
		Data:  items,
		Total: int32(len(surveys)),
	}, nil
}

func (inst *AgentpcService) GetSurveyByTenantId(ctx context.Context, req *pb.Request) (*pb.ListSurveyResponse, error) {
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

	surveys, count, err := inst.componentsContainer.SurveyService().
		GetSurveysByDomainId(userContext.Clone().EscalatePrivilege(), userInfo.DomainId)
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyService.GetSurveysByUserId")
	}
	items, err := makeSurveys(surveys)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	return &pb.ListSurveyResponse{
		Data:  items,
		Total: count,
	}, nil
}