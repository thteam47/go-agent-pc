package customerservice

import (
	"fmt"

	"github.com/antihax/optional"
	"github.com/thteam47/common-libs/reflectutil"
	"github.com/thteam47/go-agent-pc/pkg/models"

	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/apiclientutil"

	"github.com/thteam47/common-libs/confg"
	apiclient "github.com/thteam47/common/pkg/apiswagger/customer-api"
	"github.com/thteam47/go-agent-pc/errutil"
)

type DefaultCustomerService struct {
	config        *DefaultCustomerServiceConfig
	tenants       map[string]*models.Tenant
	apiClientInst *apiclient.APIClient
}

type DefaultCustomerServiceConfig struct {
	Port string `mapstructure:"port"`
}

func NewDefaultCustomerServiceWithConfig(properties confg.Confg) (*DefaultCustomerService, error) {
	config := DefaultCustomerServiceConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}
	return NewDefaultCustomerService(&config)
}

func NewDefaultCustomerService(config *DefaultCustomerServiceConfig) (*DefaultCustomerService, error) {
	inst := &DefaultCustomerService{
		config: config,
	}
	return inst, nil
}

func (inst *DefaultCustomerService) apiClient() *apiclient.APIClient {
	if inst.apiClientInst == nil {
		inst.apiClientInst = apiclient.NewAPIClient(&apiclient.Configuration{
			BasePath: fmt.Sprintf("http://127.0.0.1:%s", inst.config.Port),
			Scheme:   "http",
		})
	}
	return inst.apiClientInst
}

func (inst *DefaultCustomerService) GetTenant(userContext entity.UserContext, domain string, tokenAgent string) (*models.Tenant, error) {
	if inst.tenants != nil && inst.tenants[tokenAgent] != nil {
		return inst.tenants[tokenAgent], nil
	}
	response, _, err := inst.apiClient().CustomerServiceApi.GetByDomain(userContext.Context(), "default", domain, apiclient.Body1{
		Ctx: &apiclient.V1customerapictxDomainIdtenantsCtx{},
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityCustomerServiceApi.Login")
	}

	tenant := &models.Tenant{}
	err = reflectutil.Convert(response.Data, &tenant)
	if err != nil {
		return nil, errutil.Wrap(err, "reflectutil.Convert")
	}

	if tenant != nil {
		if inst.tenants == nil {
			inst.tenants = map[string]*models.Tenant{}
		}
		inst.tenants[tokenAgent] = tenant
	}

	return tenant, nil
}

func (inst *DefaultCustomerService) GetAllTenantByCustomer(userContext entity.UserContext) ([]models.Tenant, error) {
	response, _, err := inst.apiClient().CustomerServiceApi.GetAll(userContext.Context(), "default", &apiclient.GetAllOpts{
		CtxAccessToken: optional.NewString(userContext.AccessToken()),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityCustomerServiceApi.Login")
	}

	tenants := []models.Tenant{}

	for _, tenant := range response.Data {
		tenantTmp := models.Tenant{}
		err = reflectutil.Convert(tenant, &tenantTmp)
		if err != nil {
			return nil, errutil.Wrap(err, "reflectutil.Convert")
		}
		tenants = append(tenants, tenantTmp)
	}

	return tenants, nil
}

func (inst *DefaultCustomerService) GetTenantById(userContext entity.UserContext, tenantId string) (*models.Tenant, error) {
	response, _, err := inst.apiClient().CustomerServiceApi.GetById(userContext.Context(), "default", tenantId, &apiclient.GetByIdOpts{
		CtxAccessToken: optional.NewString(userContext.AccessToken()),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityCustomerServiceApi.Login")
	}

	tenant := &models.Tenant{}
	err = reflectutil.Convert(response.Data, &tenant)
	if err != nil {
		return nil, errutil.Wrap(err, "reflectutil.Convert")
	}

	return tenant, nil
}
