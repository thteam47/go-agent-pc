package defaultcomponent

import (
	"github.com/thteam47/common-libs/confg"
	sqlutil "github.com/thteam47/common-libs/sqliteutil"
	"github.com/thteam47/common/entity"
	sqlrepository "github.com/thteam47/common/pkg/sqliterepository"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type KeyInfoItemPhase3Repository struct {
	config         *KeyInfoItemPhase3RepositoryConfig
	baseRepository *sqlrepository.BaseRepository
}

type KeyInfoItemPhase3RepositoryConfig struct {
	SQLClientWrapper *sqlutil.SQLClientWrapper `mapstructure:"sqlite-client-wrapper"`
}

func NewKeyInfoItemPhase3RepositoryWithConfig(properties confg.Confg) (*KeyInfoItemPhase3Repository, error) {
	config := KeyInfoItemPhase3RepositoryConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}

	sqliteClientWrapper, err := sqlutil.NewSQLClientWrapperWithConfig(properties.Sub("sqlite-client-wrapper"))
	if err != nil {
		return nil, errutil.Wrap(err, "NewSQLClientWrapperWithConfig")
	}
	return NewKeyInfoItemPhase3Repository(&KeyInfoItemPhase3RepositoryConfig{
		SQLClientWrapper: sqliteClientWrapper,
	})
}

func NewKeyInfoItemPhase3Repository(config *KeyInfoItemPhase3RepositoryConfig) (*KeyInfoItemPhase3Repository, error) {
	inst := &KeyInfoItemPhase3Repository{
		config: config,
	}

	var err error
	inst.baseRepository, err = sqlrepository.NewBaseRepository(&sqlrepository.BaseRepositoryConfig{
		SqlClientWrapper: inst.config.SQLClientWrapper,
		Prototype:        models.KeyInfoItemPhase3{},
	})
	if err != nil {
		return nil, errutil.Wrap(err, "mongorepository.NewBaseRepository")
	}

	return inst, nil
}
func (inst *KeyInfoItemPhase3Repository) FindAllProcessed(userContext entity.UserContext) (int64, []models.KeyInfoItemPhase3, error) {
	result := []models.KeyInfoItemPhase3{}
	count, err := inst.baseRepository.FindAllWithAttributes(userContext, map[string]interface{}{
		"IsProcess": true,
	}, &result)
	if err != nil {
		return 0, nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return count, result, nil
}

func (inst *KeyInfoItemPhase3Repository) FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) (int64, []models.KeyInfoItemPhase3, error) {
	result := []models.KeyInfoItemPhase3{}
	if findRequest == nil {
		findRequest = &entity.FindRequest{}
	}
	count, err := inst.baseRepository.FindAllByFindRequest(userContext, findRequest, &result)
	if err != nil {
		return 0, nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return count, result, nil
}

func (inst *KeyInfoItemPhase3Repository) Create(userContext entity.UserContext, item *models.KeyInfoItemPhase3) error {
	err := inst.baseRepository.Create(userContext, item)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.Create")
	}
	return nil
}

func (inst *KeyInfoItemPhase3Repository) FindLast(userContext entity.UserContext, findRequest *entity.FindRequest) (*models.KeyInfoItemPhase3, error) {
	result := &models.KeyInfoItemPhase3{}
	if findRequest == nil {
		findRequest = &entity.FindRequest{}
	}
	if findRequest.Orders == nil {
		findRequest.Orders = map[string]entity.FindRequestOrderType{}
	}
	findRequest.Orders["CreatedAt"] = entity.FindRequestOrderTypeDesc
	err := inst.baseRepository.FindFirstByFindRequest(userContext, findRequest, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindFirstByFindRequest")
	}
	return result, nil
}
func (inst *KeyInfoItemPhase3Repository) Delete(userContext entity.UserContext, filters map[string]interface{}) error {
	item := &models.KeyInfoItemPhase3{}
	err := inst.baseRepository.DeleteWithAttributes(userContext, filters, item)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.DeleteWithFilters")
	}
	return nil
}

func (inst *KeyInfoItemPhase3Repository) DeleteByUUID(userContext entity.UserContext, item models.KeyInfoItemPhase3) error {
	baseItem := &models.KeyInfoItemPhase3{}
	err := inst.baseRepository.DeleteWithAttributes(userContext, map[string]interface{}{
		"UUID": item.UUID,
	}, baseItem)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.Update")
	}
	return nil
}

func (inst *KeyInfoItemPhase3Repository) UpdateByUUID(userContext entity.UserContext, item *models.KeyInfoItemPhase3) error {
	err := inst.baseRepository.UpdateByFindRequest(userContext, &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "UUID",
				Value:    item.UUID,
				Operator: entity.FindRequestFilterOperatorEqualTo,
			},
		},
	}, item)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.Update")
	}
	return nil
}
