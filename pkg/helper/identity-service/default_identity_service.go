package identityservice

import (
	"encoding/json"
	"fmt"

	"github.com/antihax/optional"

	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/apiclientutil"

	"github.com/thteam47/common-libs/confg"
	"github.com/thteam47/common-libs/reflectutil"
	apiclient "github.com/thteam47/common/pkg/apiswagger/identity-api"
	"github.com/thteam47/go-agent-pc/errutil"
)

type DefaultIdentityService struct {
	config        *DefaultIdentityServiceConfig
	apiClientInst *apiclient.APIClient
}

type DefaultIdentityServiceConfig struct {
	Port string `mapstructure:"port"`
}

func NewDefaultIdentityServiceWithConfig(properties confg.Confg) (*DefaultIdentityService, error) {
	config := DefaultIdentityServiceConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}
	return NewDefaultIdentityService(&config)
}

func NewDefaultIdentityService(config *DefaultIdentityServiceConfig) (*DefaultIdentityService, error) {
	inst := &DefaultIdentityService{
		config: config,
	}
	return inst, nil
}

func (inst *DefaultIdentityService) apiClient() *apiclient.APIClient {
	if inst.apiClientInst == nil {
		inst.apiClientInst = apiclient.NewAPIClient(&apiclient.Configuration{
			BasePath: fmt.Sprintf("http://127.0.0.1:%s", inst.config.Port),
			Scheme:   "http",
		})
	}
	return inst.apiClientInst
}

func (inst *DefaultIdentityService) GetUsers(userContext entity.UserContext, domain string) ([]entity.User, error) {
	findRequest := &entity.FindRequest{
		Limit: -1,
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "DomainId",
				Value:    domain,
				Operator: entity.FindRequestFilterOperatorEqualTo,
			},
		},
	}
	findRequestData, _ := json.Marshal(findRequest)
	response, _, err := inst.apiClient().IdentityServiceApi.GetAll(userContext.Context(), "default", &apiclient.GetAllOpts{
		CtxAccessToken: optional.NewString(userContext.AccessToken()),
		RequestPayload: optional.NewString(string(findRequestData)),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityIdentityServiceApi.Login")
	}

	users := []entity.User{}
	for _, user := range response.Data {
		userTmp := entity.User{}
		err = reflectutil.Convert(user, &userTmp)
		if err != nil {
			return nil, errutil.Wrap(err, "reflectutil.Convert")
		}
		users = append(users, userTmp)
	}

	return users, nil
}