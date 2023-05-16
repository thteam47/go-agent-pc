package defaultcomponent

import (
	"github.com/thteam47/common-libs/confg"
	sqlutil "github.com/thteam47/common-libs/sqliteutil"
	"github.com/thteam47/common/entity"
	sqlrepository "github.com/thteam47/common/pkg/sqliterepository"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type ResultCardRepository struct {
	config         *ResultCardRepositoryConfig
	baseRepository *sqlrepository.BaseRepository
}

type ResultCardRepositoryConfig struct {
	SQLClientWrapper *sqlutil.SQLClientWrapper `mapstructure:"sqlite-client-wrapper"`
}

func NewResultCardRepositoryWithConfig(properties confg.Confg) (*ResultCardRepository, error) {
	config := ResultCardRepositoryConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}

	sqliteClientWrapper, err := sqlutil.NewSQLClientWrapperWithConfig(properties.Sub("sqlite-client-wrapper"))
	if err != nil {
		return nil, errutil.Wrap(err, "NewSQLClientWrapperWithConfig")
	}
	return NewResultCardRepository(&ResultCardRepositoryConfig{
		SQLClientWrapper: sqliteClientWrapper,
	})
}

func NewResultCardRepository(config *ResultCardRepositoryConfig) (*ResultCardRepository, error) {
	inst := &ResultCardRepository{
		config: config,
	}

	var err error
	inst.baseRepository, err = sqlrepository.NewBaseRepository(&sqlrepository.BaseRepositoryConfig{
		SqlClientWrapper: inst.config.SQLClientWrapper,
		Prototype:        models.ResultCard{},
	})
	if err != nil {
		return nil, errutil.Wrap(err, "mongorepository.NewBaseRepository")
	}

	return inst, nil
}
func (inst *ResultCardRepository) FindAllProcessed(userContext entity.UserContext) (int64, []models.ResultCard, error) {
	result := []models.ResultCard{}
	count, err := inst.baseRepository.FindAllWithAttributes(userContext, map[string]interface{}{
		"IsProcess": true,
	}, &result)
	if err != nil {
		return 0, nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return count, result, nil
}

func (inst *ResultCardRepository) FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) (int64, []models.ResultCard, error) {
	result := []models.ResultCard{}
	if findRequest == nil {
		findRequest = &entity.FindRequest{}
	}
	count, err := inst.baseRepository.FindAllByFindRequest(userContext, findRequest, &result)
	if err != nil {
		return 0, nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return count, result, nil
}

func (inst *ResultCardRepository) Create(userContext entity.UserContext, item *models.ResultCard) error {
	err := inst.baseRepository.Create(userContext, item)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.Create")
	}
	return nil
}

func (inst *ResultCardRepository) FindLast(userContext entity.UserContext, findRequest *entity.FindRequest) (*models.ResultCard, error) {
	result := &models.ResultCard{}
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
func (inst *ResultCardRepository) Delete(userContext entity.UserContext, filters map[string]interface{}) error {
	item := &models.ResultCard{}
	err := inst.baseRepository.DeleteWithAttributes(userContext, filters, item)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.DeleteWithFilters")
	}
	return nil
}

func (inst *ResultCardRepository) DeleteByUUID(userContext entity.UserContext, item models.ResultCard) error {
	baseItem := &models.ResultCard{}
	err := inst.baseRepository.DeleteWithAttributes(userContext, map[string]interface{}{
		"UUID": item.UUID,
	}, baseItem)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.Update")
	}
	return nil
}

func (inst *ResultCardRepository) UpdateByUUID(userContext entity.UserContext, item *models.ResultCard) error {
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
