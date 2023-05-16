package component

import (
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type KeyInfoRepository interface {
	FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) (int64, []models.KeyInfo, error)
	FindAllProcessed(userContext entity.UserContext) (int64, []models.KeyInfo, error)
	Create(userContext entity.UserContext, item *models.KeyInfo) error
	FindLast(userContext entity.UserContext, findRequest *entity.FindRequest) (*models.KeyInfo, error)
	Delete(userContext entity.UserContext, filters map[string]interface{}) error
	DeleteByUUID(userContext entity.UserContext, item models.KeyInfo) error
	UpdateByUUID(userContext entity.UserContext, item *models.KeyInfo) error
}
