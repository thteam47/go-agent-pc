package grpcapp

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/thteam47/common-libs/x509util"
	pb "github.com/thteam47/common/api/agent-pc"
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/models"
	"github.com/thteam47/go-agent-pc/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getUser(item *pb.User) (*entity.User, error) {
	if item == nil {
		return nil, nil
	}
	user := &entity.User{}
	err := util.FromMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return user, nil
}

func getUsers(items []*pb.User) ([]*entity.User, error) {
	users := []*entity.User{}
	for _, item := range items {
		user, err := getUser(item)
		if err != nil {
			return nil, errutil.Wrap(err, "getUser")
		}
		users = append(users, user)
	}
	return users, nil
}

func makeUser(item *entity.User) (*pb.User, error) {
	if item == nil {
		return nil, nil
	}
	user := &pb.User{}
	err := util.ToMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return user, nil
}

func makeUsers(items []entity.User) ([]*pb.User, error) {
	users := []*pb.User{}
	for _, item := range items {
		user, err := makeUser(&item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeUser")
		}
		users = append(users, user)
	}
	return users, nil
}

func (inst *AgentpcService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	userContext := entity.NewUserContext(req.Ctx.DomainId)
	tenantId := "default"
	userType := "customer"
	if req.Domain != "" {
		tenant, err := inst.componentsContainer.CustomerService().GetTenant(userContext, req.Domain, req.Ctx.TokenAgent)
		if err != nil {
			return nil, errutil.Wrap(err, "CustomerService.GetTenant")
		}
		if tenant == nil {
			return nil, status.Errorf(codes.NotFound, "Domain not found")
		}
		tenantId = tenant.TenantId
		userType = ""
	}

	token, errorInfo, err := inst.componentsContainer.IdentityAuthenService().Login(userContext.SetDomainId(tenantId), req.Type, req.Username, req.Password, req.Otp, req.Ctx.TokenAgent, req.RequestId, userType)
	if err != nil {
		return nil, errutil.Wrap(err, "IdentityAuthenService.Login")
	}
	return &pb.LoginResponse{
		Token:     token,
		ErrorCode: errorInfo.ErrorCode,
		Message:   errorInfo.Message,
	}, nil
}

func (inst *AgentpcService) PrepareLogin(ctx context.Context, req *pb.PrepareLoginRequest) (*pb.PrepareLoginResponse, error) {
	userContext := entity.NewUserContext(req.Ctx.DomainId)
	tokenAgent, accessToken, requestId, message, typeMfa, avaiabledMfa, url, errorInfo, err := inst.componentsContainer.IdentityAuthenService().PrepareLogin(userContext, req.Ctx.TokenAgent)
	if err != nil {
		return nil, errutil.Wrap(err, "IdentityAuthenService.Login")
	}
	errorCode := int32(0)
	if errorInfo != nil {
		errorCode = errorInfo.ErrorCode
	}
	if tokenAgent != "" && inst.componentsContainer.IdentityAuthenService().GetUserInfo(userContext, tokenAgent) != nil {
		userInfo := inst.componentsContainer.IdentityAuthenService().GetUserInfo(userContext, tokenAgent)
		keyInfo, err := inst.componentsContainer.KeyInfoRepository().FindLast(userContext, &entity.FindRequest{
			Filters: []entity.FindRequestFilter{
				entity.FindRequestFilter{
					Key:      "DomainId",
					Value:    userInfo.DomainId,
					Operator: entity.FindRequestFilterOperatorEqualTo,
				},
				entity.FindRequestFilter{
					Key:      "UserId",
					Value:    userInfo.UserId,
					Operator: entity.FindRequestFilterOperatorEqualTo,
				},
			},
		})
		if err != nil {
			return nil, errutil.Wrap(err, "KeyInfoRepository.FindLast")
		}
		if keyInfo == nil {
			curve := elliptic.P256()
			privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
			if err != nil {
				panic(err)
			}
			publicKey := &privateKey.PublicKey
			privateKeyGen, err := x509util.GenerateKeyPrivate(privateKey)
			if err != nil {
				panic(err)
			}
			publicKeyGen, err := x509util.GenerateKeyPublic(publicKey)
			if err != nil {
				panic(err)
			}
			err = inst.componentsContainer.KeyInfoRepository().Create(userContext, &models.KeyInfo{
				UserId:         userInfo.UserId,
				DomainId:       userInfo.DomainId,
				PositionUserId: userInfo.Position,
				KeyPrivate:     privateKeyGen,
				KeyPublic:      publicKeyGen,
			})
			if err != nil {
				return nil, errutil.Wrap(err, "KeyInfoRepository.Create")
			}
		}
	}
	return &pb.PrepareLoginResponse{
		TokenAgent:    tokenAgent,
		Token:         accessToken,
		ErrorCode:     errorCode,
		TypeMfa:       typeMfa,
		RequestId:     requestId,
		Message:       message,
		AvailableMfas: avaiabledMfa,
		Url:           url,
	}, nil
}
