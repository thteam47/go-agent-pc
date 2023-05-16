package component

import (
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type KeyInfoItemRepository interface {
	FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) (int64, []models.KeyInfoItem, error)
	FindAllProcessed(userContext entity.UserContext) (int64, []models.KeyInfoItem, error)
	Create(userContext entity.UserContext, item *models.KeyInfoItem) error
	FindLast(userContext entity.UserContext, findRequest *entity.FindRequest) (*models.KeyInfoItem, error)
	Delete(userContext entity.UserContext, filters map[string]interface{}) error
	DeleteByUUID(userContext entity.UserContext, item models.KeyInfoItem) error
	UpdateByUUID(userContext entity.UserContext, item *models.KeyInfoItem) error
}
