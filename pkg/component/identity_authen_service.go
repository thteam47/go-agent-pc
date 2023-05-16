package component

import (
	"github.com/thteam47/common/entity"
)

type IdentityAuthenService interface {
	UpdatePassword(userContext entity.UserContext, userId string, password string) error
}
