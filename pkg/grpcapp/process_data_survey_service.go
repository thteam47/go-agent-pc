package grpcapp

import (
	"context"

	"github.com/thteam47/go-agent-pc/pkg/models"

	"github.com/thteam47/common/entity"

	pb "github.com/thteam47/common/api/agent-pc"
	"github.com/thteam47/go-agent-pc/errutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (inst *AgentpcService) ProcessDataSurvey(ctx context.Context, req *pb.StringRequest) (*pb.MessageResponse, error) {
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
	tenant, err := inst.componentsContainer.CustomerService().GetTenantById(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "CustomerService.GetTenantById")
	}
	if tenant == nil {
		return nil, status.Errorf(codes.NotFound, "Tenant not found")
	}
	if userInfo.UserId != tenant.CustomerId {
		return nil, status.Errorf(codes.PermissionDenied, "PermissionDenied")
	}
	users, err := inst.componentsContainer.IdentityService().GetUsers(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "CustomerService.GetTenant")
	}
	if len(users) == 0 {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	combinedData, err := inst.componentsContainer.RecommendService().GetCombinedData(userContext, req.Ctx.TokenAgent, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "RecommendService.GetCombinedData")
	}
	if combinedData == nil {
		return nil, status.Errorf(codes.NotFound, "CombinedData not found")
	}
	countItem := combinedData.NumberItem1
	if req.Value == combinedData.TenantId2 {
		countItem = combinedData.NumberItem2
	}
	for _, user := range users {
		for j := 1; j <= int(countItem); j++ {
			resultCard, err := inst.componentsContainer.ResultCardRepository().FindLast(userContext, &entity.FindRequest{
				Filters: []entity.FindRequestFilter{
					entity.FindRequestFilter{
						Key:      "DomainId",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    user.DomainId,
					},
					entity.FindRequestFilter{
						Key:      "UserId",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    user.UserId,
					},
					entity.FindRequestFilter{
						Key:      "PositionItem",
						Operator: entity.FindRequestFilterOperatorEqualTo,
						Value:    j,
					},
				},
			})
			if err != nil {
				return nil, errutil.Wrap(err, "ResultCardRepository.FindLast")
			}
			processedData := int32(0)
			if resultCard != nil {
				processedData = resultCard.PositionOption
			}
			err = inst.componentsContainer.ProcessDataSurveyRepository().CreateAndUpdate(userContext, &models.ProcessDataSurvey{
				DomainId:             user.DomainId,
				UserId:               user.UserId,
				PositionUser:         user.Position,
				PositionItem:         int32(j),
				ProcessedData:        processedData,
				PositionItemOriginal: int32(j),
			})
			if err != nil {
				return nil, errutil.Wrap(err, "ProcessDataSurveyRepository.Create")
			}
			processedDataIsRating := int32(0)
			if resultCard != nil {
				processedDataIsRating = 1
			}
			err = inst.componentsContainer.ProcessDataSurveyRepository().CreateAndUpdate(userContext, &models.ProcessDataSurvey{
				DomainId:             user.DomainId,
				UserId:               user.UserId,
				PositionUser:         user.Position,
				PositionItem:         countItem + int32(j),
				ProcessedData:        processedDataIsRating,
				PositionItemOriginal: int32(j),
			})
			if err != nil {
				return nil, errutil.Wrap(err, "ProcessDataSurveyRepository.Create")
			}
			err = inst.componentsContainer.ProcessDataSurveyRepository().CreateAndUpdate(userContext, &models.ProcessDataSurvey{
				DomainId:             user.DomainId,
				UserId:               user.UserId,
				PositionUser:         user.Position,
				PositionItem:         2*countItem + int32(j),
				ProcessedData:        processedData * processedData,
				PositionItemOriginal: int32(j),
			})
			if err != nil {
				return nil, errutil.Wrap(err, "ProcessDataSurveyRepository.Create")
			}
		}
		count := int32(1)
		for j := 1; j <= int(countItem)-1; j++ {
			for k := j + 1; k <= int(countItem); k++ {
				resultCardItem1, err := inst.componentsContainer.ResultCardRepository().FindLast(userContext, &entity.FindRequest{
					Filters: []entity.FindRequestFilter{
						entity.FindRequestFilter{
							Key:      "DomainId",
							Operator: entity.FindRequestFilterOperatorEqualTo,
							Value:    user.DomainId,
						},
						entity.FindRequestFilter{
							Key:      "UserId",
							Operator: entity.FindRequestFilterOperatorEqualTo,
							Value:    user.UserId,
						},
						entity.FindRequestFilter{
							Key:      "PositionItem",
							Operator: entity.FindRequestFilterOperatorEqualTo,
							Value:    j,
						},
					},
				})
				if err != nil {
					return nil, errutil.Wrap(err, "ResultCardRepository.FindLast")
				}
				processedDataItem1 := int32(0)
				if resultCardItem1 != nil {
					processedDataItem1 = resultCardItem1.PositionOption
				}
				resultCardItem2, err := inst.componentsContainer.ResultCardRepository().FindLast(userContext, &entity.FindRequest{
					Filters: []entity.FindRequestFilter{
						entity.FindRequestFilter{
							Key:      "DomainId",
							Operator: entity.FindRequestFilterOperatorEqualTo,
							Value:    user.DomainId,
						},
						entity.FindRequestFilter{
							Key:      "UserId",
							Operator: entity.FindRequestFilterOperatorEqualTo,
							Value:    user.UserId,
						},
						entity.FindRequestFilter{
							Key:      "PositionItem",
							Operator: entity.FindRequestFilterOperatorEqualTo,
							Value:    k,
						},
					},
				})
				if err != nil {
					return nil, errutil.Wrap(err, "ResultCardRepository.FindLast")
				}
				processedDataItem2 := int32(0)
				if resultCardItem2 != nil {
					processedDataItem2 = resultCardItem2.PositionOption
				}
				err = inst.componentsContainer.ProcessDataSurveyRepository().CreateAndUpdate(userContext, &models.ProcessDataSurvey{
					DomainId:              user.DomainId,
					UserId:                user.UserId,
					PositionUser:          user.Position,
					PositionItem:          3*countItem + count,
					ProcessedData:         processedDataItem2 * processedDataItem1,
					PositionItemOriginal1: int32(j),
					PositionItemOriginal2: int32(k),
				})
				if err != nil {
					return nil, errutil.Wrap(err, "ProcessDataSurveyRepository.Create")
				}
				count++
			}
		}
	}

	return &pb.MessageResponse{
		Ok: true,
	}, nil
}
