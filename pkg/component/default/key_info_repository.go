package defaultcomponent

import (
	"github.com/thteam47/common-libs/confg"
	sqlutil "github.com/thteam47/common-libs/sqliteutil"
	"github.com/thteam47/common/entity"
	sqlrepository "github.com/thteam47/common/pkg/sqliterepository"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type KeyInfoRepository struct {
	config         *KeyInfoRepositoryConfig
	baseRepository *sqlrepository.BaseRepository
}

type KeyInfoRepositoryConfig struct {
	SQLClientWrapper *sqlutil.SQLClientWrapper `mapstructure:"sqlite-client-wrapper"`
}

func NewKeyInfoRepositoryWithConfig(properties confg.Confg) (*KeyInfoRepository, error) {
	config := KeyInfoRepositoryConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}

	sqliteClientWrapper, err := sqlutil.NewSQLClientWrapperWithConfig(properties.Sub("sqlite-client-wrapper"))
	if err != nil {
		return nil, errutil.Wrap(err, "NewSQLClientWrapperWithConfig")
	}
	return NewKeyInfoRepository(&KeyInfoRepositoryConfig{
		SQLClientWrapper: sqliteClientWrapper,
	})
}

func NewKeyInfoRepository(config *KeyInfoRepositoryConfig) (*KeyInfoRepository, error) {
	inst := &KeyInfoRepository{
		config: config,
	}

	var err error
	inst.baseRepository, err = sqlrepository.NewBaseRepository(&sqlrepository.BaseRepositoryConfig{
		SqlClientWrapper: inst.config.SQLClientWrapper,
		Prototype:        models.KeyInfo{},
	})
	if err != nil {
		return nil, errutil.Wrap(err, "mongorepository.NewBaseRepository")
	}

	return inst, nil
}
func (inst *KeyInfoRepository) FindAllProcessed(userContext entity.UserContext) (int64, []models.KeyInfo, error) {
	result := []models.KeyInfo{}
	count, err := inst.baseRepository.FindAllWithAttributes(userContext, map[string]interface{}{
		"IsProcess": true,
	}, &result)
	if err != nil {
		return 0, nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return count, result, nil
}

func (inst *KeyInfoRepository) FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) (int64, []models.KeyInfo, error) {
	result := []models.KeyInfo{}
	if findRequest == nil {
		findRequest = &entity.FindRequest{}
	}
	count, err := inst.baseRepository.FindAllByFindRequest(userContext, findRequest, &result)
	if err != nil {
		return 0, nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return count, result, nil
}

func (inst *KeyInfoRepository) Create(userContext entity.UserContext, item *models.KeyInfo) error {
	err := inst.baseRepository.Create(userContext, item)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.Create")
	}
	return nil
}

func (inst *KeyInfoRepository) FindLast(userContext entity.UserContext, findRequest *entity.FindRequest) (*models.KeyInfo, error) {
	result := &models.KeyInfo{}
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
func (inst *KeyInfoRepository) Delete(userContext entity.UserContext, filters map[string]interface{}) error {
	item := &models.KeyInfo{}
	err := inst.baseRepository.DeleteWithAttributes(userContext, filters, item)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.DeleteWithFilters")
	}
	return nil
}

func (inst *KeyInfoRepository) DeleteByUUID(userContext entity.UserContext, item models.KeyInfo) error {
	baseItem := &models.KeyInfo{}
	err := inst.baseRepository.DeleteWithAttributes(userContext, map[string]interface{}{
		"UUID": item.UUID,
	}, baseItem)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.Update")
	}
	return nil
}

func (inst *KeyInfoRepository) UpdateByUUID(userContext entity.UserContext, item *models.KeyInfo) error {
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
