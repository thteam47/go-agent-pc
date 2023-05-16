package component

import (
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type ResultCardRepository interface {
	FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) (int64, []models.ResultCard, error)
	FindAllProcessed(userContext entity.UserContext) (int64, []models.ResultCard, error)
	Create(userContext entity.UserContext, item *models.ResultCard) error
	FindLast(userContext entity.UserContext, findRequest *entity.FindRequest) (*models.ResultCard, error)
	Delete(userContext entity.UserContext, filters map[string]interface{}) error
	DeleteByUUID(userContext entity.UserContext, item models.ResultCard) error
	UpdateByUUID(userContext entity.UserContext, item *models.ResultCard) error
}
