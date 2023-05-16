package recommendservice

import (
	"fmt"
	"strconv"

	"github.com/antihax/optional"
	"github.com/thteam47/go-agent-pc/pkg/models"

	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/apiclientutil"

	"github.com/thteam47/common-libs/confg"
	"github.com/thteam47/common-libs/reflectutil"
	apiclient "github.com/thteam47/common/pkg/apiswagger/recommend-api"
	"github.com/thteam47/go-agent-pc/errutil"
)

type DefaultRecommendService struct {
	config        *DefaultRecommendServiceConfig
	tenants       map[string]*models.Tenant
	apiClientInst *apiclient.APIClient
}

type DefaultRecommendServiceConfig struct {
	Port string `mapstructure:"port"`
}

func NewDefaultRecommendServiceWithConfig(properties confg.Confg) (*DefaultRecommendService, error) {
	config := DefaultRecommendServiceConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}
	return NewDefaultRecommendService(&config)
}

func NewDefaultRecommendService(config *DefaultRecommendServiceConfig) (*DefaultRecommendService, error) {
	inst := &DefaultRecommendService{
		config: config,
	}
	return inst, nil
}

func (inst *DefaultRecommendService) apiClient() *apiclient.APIClient {
	if inst.apiClientInst == nil {
		inst.apiClientInst = apiclient.NewAPIClient(&apiclient.Configuration{
			BasePath: fmt.Sprintf("http://127.0.0.1:%s", inst.config.Port),
			Scheme:   "http",
		})
	}
	return inst.apiClientInst
}

func (inst *DefaultRecommendService) KeyPublicUserSend(userContext entity.UserContext, domain string, tokenAgent string, keyInfo *models.KeyInfo) error {
	_, _, err := inst.apiClient().RecommendServiceApi.KeyPublicUserReceive(userContext.Context(), domain, apiclient.Body3{
		Ctx: &apiclient.V1recommendapictxDomainIdcombinedDataapproveCtx{
			AccessToken: userContext.AccessToken(),
		},
		Data: &apiclient.RecommendApiKeyPublicUser{
			KeyPublic:    keyInfo.KeyPublic,
			DomainId:     keyInfo.DomainId,
			PositionUser: keyInfo.PositionUserId,
			UserId:       keyInfo.UserId,
		},
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityRecommendServiceApi.Login")
	}

	return nil
}

func (inst *DefaultRecommendService) KeyPublicItemSend(userContext entity.UserContext, domain string, tokenAgent string, keyInfoItem *models.KeyInfoItem) error {
	_, _, err := inst.apiClient().RecommendServiceApi.KeyPublicItemReceive(userContext.Context(), domain, apiclient.Body2{
		Ctx: &apiclient.V1recommendapictxDomainIdcombinedDataapproveCtx{
			AccessToken: userContext.AccessToken(),
		},
		Data: &apiclient.RecommendApiKeyPublicItem{
			KeyPublic:    keyInfoItem.KeyPublic,
			DomainId:     keyInfoItem.DomainId,
			PositionUser: keyInfoItem.PositionUser,
			UserId:       keyInfoItem.UserId,
			Part:         keyInfoItem.Part,
			PositionItem: keyInfoItem.PositionItem,
		},
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityRecommendServiceApi.Login")
	}

	return nil
}

func (inst *DefaultRecommendService) GetCombinedData(userContext entity.UserContext, tokenAgent string, tenantId string) (*models.CombinedData, error) {
	response, _, err := inst.apiClient().RecommendServiceApi.CombinedDataGetByTenantId(userContext.Context(), "default", tenantId, &apiclient.CombinedDataGetByTenantIdOpts{
		CtxAccessToken: optional.NewString(userContext.AccessToken()),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityIdentityServiceApi.Login")
	}

	combinedData := &models.CombinedData{}
	err = reflectutil.Convert(response.Data, &combinedData)
	if err != nil {
		return nil, errutil.Wrap(err, "reflectutil.Convert")
	}

	return combinedData, nil
}

func (inst *DefaultRecommendService) KeyPublicUse(userContext entity.UserContext, tenantId string, tokenAgent string, positionItem int32, part int32) (*models.KeyPublicUse, error) {
	response, _, err := inst.apiClient().RecommendServiceApi.KeyPublicUseGet(userContext.Context(), tenantId, strconv.Itoa(int(positionItem)), &apiclient.KeyPublicUseGetOpts{
		CtxAccessToken: optional.NewString(userContext.AccessToken()),
		CtxPart:        optional.NewInt32(part),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, errutil.Wrap(apiclientutil.NormalizeError(err), "RecommendServiceApi.KeyPublicUseGet")
	}

	keyPublicUse := &models.KeyPublicUse{}
	err = reflectutil.Convert(response.Data, &keyPublicUse)
	if err != nil {
		return nil, errutil.Wrap(err, "reflectutil.Convert")
	}

	return keyPublicUse, nil
}

func (inst *DefaultRecommendService) ProcessedDataSendPart(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiResultCard) error {
	_, _, err := inst.apiClient().RecommendServiceApi.ResultCardCreate(userContext.Context(), domain, apiclient.Body12{
		Ctx: &apiclient.V1recommendapictxDomainIdcombinedDataapproveCtx{
			AccessToken: userContext.AccessToken(),
		},
		Data: data,
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityRecommendServiceApi.Login")
	}

	return nil
}

func (inst *DefaultRecommendService) ProcessedDataSendPhase3TwoPart(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiPhase3TwoPart) error {
	_, _, err := inst.apiClient().RecommendServiceApi.Phase3TwoPartCreate(userContext.Context(), domain, apiclient.Body4{
		Ctx: &apiclient.V1recommendapictxDomainIdcombinedDataapproveCtx{
			AccessToken: userContext.AccessToken(),
		},
		Data: data,
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityRecommendServiceApi.Login")
	}

	return nil
}

func (inst *DefaultRecommendService) ProcessedDataSendPhase3Get(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiPhase3TwoPart) (*apiclient.RecommendApiPhase3TwoPart, error) {
	response, _, err := inst.apiClient().RecommendServiceApi.Phase3TwoPartGet(userContext.Context(), domain, &apiclient.Phase3TwoPartGetOpts{
		CtxAccessToken:   optional.NewString(userContext.AccessToken()),
		CtxPart:          optional.NewInt32(2),
		DataPositionUser: optional.NewInt32(data.PositionUser),
		DataPositionItem: optional.NewInt32(data.PositionItem),
		DataPart:         optional.NewInt32(2),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityRecommendServiceApi.Login")
	}

	return response.Data, nil
}

func (inst *DefaultRecommendService) ProcessedDataSendPhase4Get(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiPhase4TwoPart) (*apiclient.RecommendApiPhase4TwoPart, error) {
	response, _, err := inst.apiClient().RecommendServiceApi.Phase4TwoPartGet(userContext.Context(), domain, &apiclient.Phase4TwoPartGetOpts{
		CtxAccessToken:   optional.NewString(userContext.AccessToken()),
		CtxPart:          optional.NewInt32(2),
		DataPositionUser: optional.NewInt32(data.PositionUser),
		DataPositionItem: optional.NewInt32(data.PositionItem),
		DataPart:         optional.NewInt32(2),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityRecommendServiceApi.Login")
	}

	return response.Data, nil
}

func (inst *DefaultRecommendService) KeyPublicUserGet(userContext entity.UserContext, domain string, tokenAgent string, positionUser int32) (*apiclient.RecommendApiKeyPublicUser, error) {
	response, _, err := inst.apiClient().RecommendServiceApi.KeyPublicUserGet(userContext.Context(), domain, &apiclient.KeyPublicUserGetOpts{
		CtxAccessToken:   optional.NewString(userContext.AccessToken()),
		DataPositionUser: optional.NewInt32(positionUser),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityRecommendServiceApi.Login")
	}

	return response.Data, nil
}

func (inst *DefaultRecommendService) ProcessedDataSendPhase4TwoPart(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiPhase4TwoPart) error {
	_, _, err := inst.apiClient().RecommendServiceApi.Phase4TwoPartCreate(userContext.Context(), domain, apiclient.Body5{
		Ctx: &apiclient.V1recommendapictxDomainIdcombinedDataapproveCtx{
			AccessToken: userContext.AccessToken(),
		},
		Data: data,
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityRecommendServiceApi.Login")
	}

	return nil
}

func (inst *DefaultRecommendService) ProcessedDataSend2(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiProcessDataSurvey2) error {
	_, _, err := inst.apiClient().RecommendServiceApi.ProcessDataSurveyCreate2(userContext.Context(), domain, apiclient.Body8{
		Ctx: &apiclient.V1recommendapictxDomainIdcombinedDataapproveCtx{
			AccessToken: userContext.AccessToken(),
		},
		Data: data,
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityRecommendServiceApi.Login")
	}

	return nil
}

func (inst *DefaultRecommendService) RecommendCbfGenarate(userContext entity.UserContext, domain string, tokenAgent string, userId string, data map[string]apiclient.RecommendApiRecommend) (map[string]apiclient.RecommendApiRecommendCbfResult12, *apiclient.RecommendApiRecommendCbfResult34, error) {
	response, _, err := inst.apiClient().RecommendServiceApi.RecommendCbfGenarate(userContext.Context(), domain, userId, apiclient.Body10{
		Ctx: &apiclient.V1recommendapictxDomainIdcombinedDataapproveCtx{
			AccessToken: userContext.AccessToken(),
		},
		Data: data,
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityRecommendServiceApi.Login")
	}

	return response.Data12, response.Data34, nil
}

func (inst *DefaultRecommendService) RecommendCfGenarate(userContext entity.UserContext, domain string, tokenAgent string, userId string, data map[string]apiclient.RecommendApiRecommend) (map[string]apiclient.RecommendApiRecommendCfResult910, *apiclient.RecommendApiRecommendCfResult1112, error) {
	response, _, err := inst.apiClient().RecommendServiceApi.RecommendCfGenarate(userContext.Context(), domain, userId, apiclient.Body11{
		Ctx: &apiclient.V1recommendapictxDomainIdcombinedDataapproveCtx{
			AccessToken: userContext.AccessToken(),
		},
		Data: data,
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityRecommendServiceApi.Login")
	}

	return response.Data910, response.Data1112, nil
}