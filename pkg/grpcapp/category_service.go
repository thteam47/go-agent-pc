package grpcapp

import (
	"context"

	"github.com/thteam47/common/entity"

	pb "github.com/thteam47/common/api/agent-pc"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/models"
	"github.com/thteam47/go-agent-pc/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func makeCategory(item *models.Category) (*pb.Category, error) {
	if item == nil {
		return nil, nil
	}
	category := &pb.Category{}
	err := util.ToMessage(item, category)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return category, nil
}

func makeCategories(items []models.Category) ([]*pb.Category, error) {
	categories := []*pb.Category{}
	for _, item := range items {
		category, err := makeCategory(&item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeCategory")
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func makeCategoryRecommend(item *models.Category) (*pb.CategoryRecommend, error) {
	if item == nil {
		return nil, nil
	}
	category := &pb.CategoryRecommend{}
	err := util.ToMessage(item, category)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return category, nil
}

func makeCategoriesRecommend(items []models.Category, resultRecommend map[int32]float32) ([]*pb.CategoryRecommend, error) {
	categories := []*pb.CategoryRecommend{}
	for _, item := range items {
		category, err := makeCategoryRecommend(&item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeCategory")
		}
		category.ProcessData = resultRecommend[item.Position]
		categories = append(categories, category)
	}
	return categories, nil
}


// func (inst *AgentpcService) GetSurveyByTenant(ctx context.Context, req *pb.StringRequest) (*pb.ListSurveyResponse, error) {
// 	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "@any", "@any", &grpcauth.AuthenOption{})
// 	if err != nil {
// 		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
// 	}
// 	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Value)
// 	if err != nil {
// 		return nil, errutil.Wrap(err, "CustomerService.GetTenantById")
// 	}
// 	if tenant == nil {
// 		return nil, status.Errorf(codes.NotFound, "Tenant not found")
// 	}

// 	surveys, count, err := inst.componentsContainer.SurveyService().
// 		GetSurveysByDomainId(userContext.Clone().EscalatePrivilege(), req.Value)
// 	if err != nil {
// 		return nil, errutil.Wrap(err, "SurveyService.GetSurveysByUserId")
// 	}
// 	items, err := makeSurveys(surveys)
// 	if err != nil {
// 		return nil, errutil.Wrap(err, "makeSurveys")
// 	}
// 	return &pb.ListSurveyResponse{
// 		Data:  items,
// 		Total: count,
// 	}, nil
// }

func (inst *AgentpcService) GetCategoriesByRecommendTenant(ctx context.Context, req *pb.Request) (*pb.ListCategoryResponse, error) {
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

	categories, total, err := inst.componentsContainer.SurveyService().GetCategoriesByRecommendTenant(userContext.SetDomainId(userInfo.DomainId).Clone().EscalatePrivilege(), userInfo.DomainId)
	if err != nil {
		return nil, errutil.Wrap(err, "SurveyService.GetSurveysByUserId")
	}
	items, err := makeCategories(categories)
	if err != nil {
		return nil, errutil.Wrap(err, "makeSurveys")
	}
	return &pb.ListCategoryResponse{
		Data:  items,
		Total: total,
	}, nil
}
