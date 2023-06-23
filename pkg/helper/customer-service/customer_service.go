package customerservice

import (
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type CustomerService interface {
	GetTenant(userContext entity.UserContext, domain string, tokenAgent string) (*models.Tenant, error)
	GetTenantById(userContext entity.UserContext, tenantId string) (*models.Tenant, error)
	GetAllTenantByCustomer(userContext entity.UserContext) ([]models.Tenant, error)
}
