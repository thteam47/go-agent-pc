package surveyservice

import (
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type SurveyService interface {
	GetCategories(userContext entity.UserContext, domain string) ([]models.Category, error)
	GetSurveysByUserId(userContext entity.UserContext, domain string, userId string) ([]models.Survey, error)
	GetSurveysByDomainId(userContext entity.UserContext, domain string) ([]models.Survey, int32, error)
	GetCategoriesByRecommendTenant(userContext entity.UserContext, domain string) ([]models.Category, int32, error)
}
