package defaultcomponent

import (
	"github.com/thteam47/common-libs/confg"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/common/handler"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/component"
	customerservice "github.com/thteam47/go-agent-pc/pkg/helper/customer-service"
	authenservice "github.com/thteam47/go-agent-pc/pkg/helper/identity-authen-service"
	identityservice "github.com/thteam47/go-agent-pc/pkg/helper/identity-service"
	recommendservice "github.com/thteam47/go-agent-pc/pkg/helper/recommend-service"
	surveyservice "github.com/thteam47/go-agent-pc/pkg/helper/survey-service"
)

type ComponentFactory struct {
	properties confg.Confg
	handle     *handler.Handler
}

func NewComponentFactory(properties confg.Confg, handle *handler.Handler) (*ComponentFactory, error) {
	inst := &ComponentFactory{
		properties: properties,
		handle:     handle,
	}

	return inst, nil
}

func (inst *ComponentFactory) CreateAuthService() *grpcauth.AuthInterceptor {
	authService := grpcauth.NewAuthInterceptor(inst.handle)
	return authService
}

func (inst *ComponentFactory) CreateKeyInfoRepository() (component.KeyInfoRepository, error) {
	keyInfoRepository, err := NewKeyInfoRepositoryWithConfig(inst.properties.Sub("key-info-repository"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewKeyInfoRepositoryWithConfig")
	}
	return keyInfoRepository, nil
}

func (inst *ComponentFactory) CreateKeyInfoItemRepository() (component.KeyInfoItemRepository, error) {
	keyInfoItemRepository, err := NewKeyInfoItemRepositoryWithConfig(inst.properties.Sub("key-info-item-repository"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewKeyInfoItemRepositoryWithConfig")
	}
	return keyInfoItemRepository, nil
}

func (inst *ComponentFactory) CreateKeyInfoItemPhase3Repository() (component.KeyInfoItemPhase3Repository, error) {
	keyInfoItemPhase3Repository, err := NewKeyInfoItemPhase3RepositoryWithConfig(inst.properties.Sub("key-info-item-phase3-repository"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewKeyInfoItemPhase3RepositoryWithConfig")
	}
	return keyInfoItemPhase3Repository, nil
}

func (inst *ComponentFactory) CreateResultCardRepository() (component.ResultCardRepository, error) {
	resultCardRepository, err := NewResultCardRepositoryWithConfig(inst.properties.Sub("result-card-repository"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewResultCardRepositoryWithConfig")
	}
	return resultCardRepository, nil
}

func (inst *ComponentFactory) CreateProcessDataSurveyRepository() (component.ProcessDataSurveyRepository, error) {
	processDataSurveyRepository, err := NewProcessDataSurveyRepositoryWithConfig(inst.properties.Sub("process-data-survey-repository"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewProcessDataSurveyRepositoryWithConfig")
	}
	return processDataSurveyRepository, nil
}

func (inst *ComponentFactory) CreateIdentityAuthenService() (authenservice.AuthenService, error) {
	identityAuthenService, err := authenservice.NewDefaultAuthenServiceWithConfig(inst.properties.Sub("identity-authen-service"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewIdentityAuthenServiceWithConfig")
	}
	return identityAuthenService, nil
}

func (inst *ComponentFactory) CreateCustomerService() (customerservice.CustomerService, error) {
	customerService, err := customerservice.NewDefaultCustomerServiceWithConfig(inst.properties.Sub("customer-service"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewDefaultCustomerServiceWithConfig")
	}
	return customerService, nil
}

func (inst *ComponentFactory) CreateIdentityService() (identityservice.IdentityService, error) {
	identityService, err := identityservice.NewDefaultIdentityServiceWithConfig(inst.properties.Sub("identity-service"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewDefaultIdentityService")
	}
	return identityService, nil
}

func (inst *ComponentFactory) CreateRecommendService() (recommendservice.RecommendService, error) {
	recommendService, err := recommendservice.NewDefaultRecommendServiceWithConfig(inst.properties.Sub("recommend-service"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewDefaultRecommendServiceWithConfig")
	}
	return recommendService, nil
}

func (inst *ComponentFactory) CreateSurveyService() (surveyservice.SurveyService, error) {
	surveyService, err := surveyservice.NewDefaultSurveyServiceWithConfig(inst.properties.Sub("survey-service"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewDefaultSurveyServiceWithConfig")
	}
	return surveyService, nil
}
