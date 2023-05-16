package component

import (
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type ProcessDataSurveyRepository interface {
	FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) (int64, []models.ProcessDataSurvey, error)
	FindAllProcessed(userContext entity.UserContext) (int64, []models.ProcessDataSurvey, error)
	Create(userContext entity.UserContext, item *models.ProcessDataSurvey) error
	FindLast(userContext entity.UserContext, findRequest *entity.FindRequest) (*models.ProcessDataSurvey, error)
	Delete(userContext entity.UserContext, filters map[string]interface{}) error
	DeleteByUUID(userContext entity.UserContext, item models.ProcessDataSurvey) error
	UpdateByUUID(userContext entity.UserContext, item *models.ProcessDataSurvey) error
	CreateAndUpdate(userContext entity.UserContext, item *models.ProcessDataSurvey) error
}
