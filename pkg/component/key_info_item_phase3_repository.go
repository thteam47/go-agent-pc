package component

import (
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type KeyInfoItemPhase3Repository interface {
	FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) (int64, []models.KeyInfoItemPhase3, error)
	FindAllProcessed(userContext entity.UserContext) (int64, []models.KeyInfoItemPhase3, error)
	Create(userContext entity.UserContext, item *models.KeyInfoItemPhase3) error
	FindLast(userContext entity.UserContext, findRequest *entity.FindRequest) (*models.KeyInfoItemPhase3, error)
	Delete(userContext entity.UserContext, filters map[string]interface{}) error
	DeleteByUUID(userContext entity.UserContext, item models.KeyInfoItemPhase3) error
	UpdateByUUID(userContext entity.UserContext, item *models.KeyInfoItemPhase3) error
}
