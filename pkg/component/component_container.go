package component

import (
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/db"
	customerservice "github.com/thteam47/go-agent-pc/pkg/helper/customer-service"
	authenservice "github.com/thteam47/go-agent-pc/pkg/helper/identity-authen-service"
	identityservice "github.com/thteam47/go-agent-pc/pkg/helper/identity-service"
	recommendservice "github.com/thteam47/go-agent-pc/pkg/helper/recommend-service"
	surveyservice "github.com/thteam47/go-agent-pc/pkg/helper/survey-service"
)

type ComponentsContainer struct {
	keyInfoRepository           KeyInfoRepository
	authService                 *grpcauth.AuthInterceptor
	handler                     *db.Handler
	identityAuthenService       authenservice.AuthenService
	customerService             customerservice.CustomerService
	identityService             identityservice.IdentityService
	recommendService            recommendservice.RecommendService
	surveyService               surveyservice.SurveyService
	keyInfoItemRepository       KeyInfoItemRepository
	keyInfoItemPhase3Repository KeyInfoItemPhase3Repository
	resultCardRepository        ResultCardRepository
	processDataSurveyRepository ProcessDataSurveyRepository
	accessToken                 map[string]string
}

func NewComponentsContainer(componentFactory ComponentFactory) (*ComponentsContainer, error) {
	inst := &ComponentsContainer{}

	var err error
	inst.authService = componentFactory.CreateAuthService()
	inst.keyInfoRepository, err = componentFactory.CreateKeyInfoRepository()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateKeyInfoRepository")
	}
	inst.resultCardRepository, err = componentFactory.CreateResultCardRepository()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateResultCardRepository")
	}
	inst.keyInfoItemRepository, err = componentFactory.CreateKeyInfoItemRepository()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateKeyInfoItemRepository")
	}
	inst.keyInfoItemPhase3Repository, err = componentFactory.CreateKeyInfoItemPhase3Repository()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateKeyInfoItemPhase3Repository")
	}
	inst.processDataSurveyRepository, err = componentFactory.CreateProcessDataSurveyRepository()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateProcessDataSurveyRepository")
	}
	inst.identityAuthenService, err = componentFactory.CreateIdentityAuthenService()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateIdentityAuthenService")
	}
	inst.customerService, err = componentFactory.CreateCustomerService()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateCustomerService")
	}
	inst.identityService, err = componentFactory.CreateIdentityService()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateIdentityService")
	}
	inst.recommendService, err = componentFactory.CreateRecommendService()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateRecommendService")
	}
	inst.surveyService, err = componentFactory.CreateSurveyService()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateSurveyService")
	}
	return inst, nil
}

func (inst *ComponentsContainer) AuthService() *grpcauth.AuthInterceptor {
	return inst.authService
}

func (inst *ComponentsContainer) KeyInfoRepository() KeyInfoRepository {
	return inst.keyInfoRepository
}

func (inst *ComponentsContainer) KeyInfoItemRepository() KeyInfoItemRepository {
	return inst.keyInfoItemRepository
}

func (inst *ComponentsContainer) KeyInfoItemPhase3Repository() KeyInfoItemPhase3Repository {
	return inst.keyInfoItemPhase3Repository
}

func (inst *ComponentsContainer) ResultCardRepository() ResultCardRepository {
	return inst.resultCardRepository
}

func (inst *ComponentsContainer) ProcessDataSurveyRepository() ProcessDataSurveyRepository {
	return inst.processDataSurveyRepository
}

func (inst *ComponentsContainer) IdentityAuthenService() authenservice.AuthenService {
	return inst.identityAuthenService
}

func (inst *ComponentsContainer) CustomerService() customerservice.CustomerService {
	return inst.customerService
}

func (inst *ComponentsContainer) IdentityService() identityservice.IdentityService {
	return inst.identityService
}

func (inst *ComponentsContainer) RecommendService() recommendservice.RecommendService {
	return inst.recommendService
}

func (inst *ComponentsContainer) SurveyService() surveyservice.SurveyService {
	return inst.surveyService
}

var errorCodeBadRequest = 400
