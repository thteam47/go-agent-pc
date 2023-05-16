package surveyservice

import (
	"encoding/json"
	"fmt"

	"github.com/antihax/optional"

	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/apiclientutil"
	"github.com/thteam47/go-agent-pc/pkg/models"

	"github.com/thteam47/common-libs/confg"
	"github.com/thteam47/common-libs/reflectutil"
	apiclient "github.com/thteam47/common/pkg/apiswagger/survey-api"
	"github.com/thteam47/go-agent-pc/errutil"
)

type DefaultSurveyService struct {
	config        *DefaultSurveyServiceConfig
	apiClientInst *apiclient.APIClient
}

type DefaultSurveyServiceConfig struct {
	Port string `mapstructure:"port"`
}

func NewDefaultSurveyServiceWithConfig(properties confg.Confg) (*DefaultSurveyService, error) {
	config := DefaultSurveyServiceConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}
	return NewDefaultSurveyService(&config)
}

func NewDefaultSurveyService(config *DefaultSurveyServiceConfig) (*DefaultSurveyService, error) {
	inst := &DefaultSurveyService{
		config: config,
	}
	return inst, nil
}

func (inst *DefaultSurveyService) apiClient() *apiclient.APIClient {
	if inst.apiClientInst == nil {
		inst.apiClientInst = apiclient.NewAPIClient(&apiclient.Configuration{
			BasePath: fmt.Sprintf("http://127.0.0.1:%s", inst.config.Port),
			Scheme:   "http",
		})
	}
	return inst.apiClientInst
}

func (inst *DefaultSurveyService) GetCategories(userContext entity.UserContext, domain string) ([]models.Category, error) {
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
	response, _, err := inst.apiClient().SurveyServiceApi.GetAllCategory(userContext.Context(), "default", &apiclient.GetAllCategoryOpts{
		CtxAccessToken: optional.NewString(userContext.AccessToken()),
		RequestPayload: optional.NewString(string(findRequestData)),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentitySurveyServiceApi.Login")
	}

	categories := []models.Category{}
	for _, category := range response.Data {
		categoryTmp := models.Category{}
		err = reflectutil.Convert(category, &categoryTmp)
		if err != nil {
			return nil, errutil.Wrap(err, "reflectutil.Convert")
		}
		categories = append(categories, categoryTmp)
	}

	return categories, nil
}

func (inst *DefaultSurveyService) GetSurveysByUserId(userContext entity.UserContext, domain string, userId string) ([]models.Survey, error) {
	response, _, err := inst.apiClient().SurveyServiceApi.GetSurveyByUserJoin(userContext.Context(), domain, userId, &apiclient.GetSurveyByUserJoinOpts{
		CtxAccessToken: optional.NewString(userContext.AccessToken()),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, errutil.Wrap(apiclientutil.NormalizeError(err), "SurveyServiceApi.GetSurveyByUserJoin")
	}

	surveys := []models.Survey{}
	for _, survey := range response.Data {
		surveyTmp := models.Survey{}
		err = reflectutil.Convert(survey, &surveyTmp)
		if err != nil {
			return nil, errutil.Wrap(err, "reflectutil.Convert")
		}
		surveys = append(surveys, surveyTmp)
	}

	return surveys, nil
}

func (inst *DefaultSurveyService) GetSurveysByDomainId(userContext entity.UserContext, domain string) ([]models.Survey, int32, error) {
	response, _, err := inst.apiClient().SurveyServiceApi.GetSurveyByTenant(userContext.Context(), domain, domain, &apiclient.GetSurveyByTenantOpts{
		CtxAccessToken: optional.NewString(userContext.AccessToken()),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, 0, errutil.Wrap(apiclientutil.NormalizeError(err), "SurveyServiceApi.GetSurveyByUserJoin")
	}

	surveys := []models.Survey{}
	for _, survey := range response.Data {
		surveyTmp := models.Survey{}
		err = reflectutil.Convert(survey, &surveyTmp)
		if err != nil {
			return nil, 0, errutil.Wrap(err, "reflectutil.Convert")
		}
		surveys = append(surveys, surveyTmp)
	}

	return surveys, response.Total, nil
}

func (inst *DefaultSurveyService) GetCategoriesByRecommendTenant(userContext entity.UserContext, domain string) ([]models.Category, int32, error) {
	response, _, err := inst.apiClient().SurveyServiceApi.GetCategoriesByRecommendTenant(userContext.Context(), domain, domain, &apiclient.GetCategoriesByRecommendTenantOpts{
		CtxAccessToken: optional.NewString(userContext.AccessToken()),
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return nil, 0, errutil.Wrap(apiclientutil.NormalizeError(err), "SurveyServiceApi.GetSurveyByUserJoin")
	}

	categories := []models.Category{}
	for _, category := range response.Data {
		categoryTmp := models.Category{}
		err = reflectutil.Convert(category, &categoryTmp)
		if err != nil {
			return nil, 0, errutil.Wrap(err, "reflectutil.Convert")
		}
		categories = append(categories, categoryTmp)
	}

	return categories, response.Total, nil
}
