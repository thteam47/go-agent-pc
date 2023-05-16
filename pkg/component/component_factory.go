package component

import (
	grpcauth "github.com/thteam47/common/grpcutil"
	customerservice "github.com/thteam47/go-agent-pc/pkg/helper/customer-service"
	authenservice "github.com/thteam47/go-agent-pc/pkg/helper/identity-authen-service"
	identityservice "github.com/thteam47/go-agent-pc/pkg/helper/identity-service"
	recommendservice "github.com/thteam47/go-agent-pc/pkg/helper/recommend-service"
	surveyservice "github.com/thteam47/go-agent-pc/pkg/helper/survey-service"
)

type ComponentFactory interface {
	CreateAuthService() *grpcauth.AuthInterceptor
	CreateKeyInfoRepository() (KeyInfoRepository, error)
	CreateKeyInfoItemRepository() (KeyInfoItemRepository, error)
	CreateKeyInfoItemPhase3Repository() (KeyInfoItemPhase3Repository, error)
	CreateResultCardRepository() (ResultCardRepository, error)
	CreateProcessDataSurveyRepository() (ProcessDataSurveyRepository, error)
	CreateIdentityAuthenService() (authenservice.AuthenService, error)
	CreateCustomerService() (customerservice.CustomerService, error)
	CreateIdentityService() (identityservice.IdentityService, error)
	CreateRecommendService() (recommendservice.RecommendService, error)
	CreateSurveyService() (surveyservice.SurveyService, error)
}
