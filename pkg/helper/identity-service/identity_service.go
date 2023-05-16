package identityservice

import (
	"github.com/thteam47/common/entity"
)

type IdentityService interface {
	GetUsers(userContext entity.UserContext, domain string) ([]entity.User, error)
}
