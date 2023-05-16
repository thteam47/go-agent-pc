package defaultcomponent

import (
	"github.com/thteam47/common-libs/confg"
	sqlutil "github.com/thteam47/common-libs/sqliteutil"
	"github.com/thteam47/common/entity"
	sqlrepository "github.com/thteam47/common/pkg/sqliterepository"
	"github.com/thteam47/go-agent-pc/errutil"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type ProcessDataSurveyRepository struct {
	config         *ProcessDataSurveyRepositoryConfig
	baseRepository *sqlrepository.BaseRepository
}

type ProcessDataSurveyRepositoryConfig struct {
	SQLClientWrapper *sqlutil.SQLClientWrapper `mapstructure:"sqlite-client-wrapper"`
}

func NewProcessDataSurveyRepositoryWithConfig(properties confg.Confg) (*ProcessDataSurveyRepository, error) {
	config := ProcessDataSurveyRepositoryConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}

	sqliteClientWrapper, err := sqlutil.NewSQLClientWrapperWithConfig(properties.Sub("sqlite-client-wrapper"))
	if err != nil {
		return nil, errutil.Wrap(err, "NewSQLClientWrapperWithConfig")
	}
	return NewProcessDataSurveyRepository(&ProcessDataSurveyRepositoryConfig{
		SQLClientWrapper: sqliteClientWrapper,
	})
}

func NewProcessDataSurveyRepository(config *ProcessDataSurveyRepositoryConfig) (*ProcessDataSurveyRepository, error) {
	inst := &ProcessDataSurveyRepository{
		config: config,
	}

	var err error
	inst.baseRepository, err = sqlrepository.NewBaseRepository(&sqlrepository.BaseRepositoryConfig{
		SqlClientWrapper: inst.config.SQLClientWrapper,
		Prototype:        models.ProcessDataSurvey{},
	})
	if err != nil {
		return nil, errutil.Wrap(err, "mongorepository.NewBaseRepository")
	}

	return inst, nil
}
func (inst *ProcessDataSurveyRepository) FindAllProcessed(userContext entity.UserContext) (int64, []models.ProcessDataSurvey, error) {
	result := []models.ProcessDataSurvey{}
	count, err := inst.baseRepository.FindAllWithAttributes(userContext, map[string]interface{}{
		"IsProcess": true,
	}, &result)
	if err != nil {
		return 0, nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return count, result, nil
}

func (inst *ProcessDataSurveyRepository) FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) (int64, []models.ProcessDataSurvey, error) {
	result := []models.ProcessDataSurvey{}
	if findRequest == nil {
		findRequest = &entity.FindRequest{}
	}
	count, err := inst.baseRepository.FindAllByFindRequest(userContext, findRequest, &result)
	if err != nil {
		return 0, nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return count, result, nil
}

func (inst *ProcessDataSurveyRepository) Create(userContext entity.UserContext, item *models.ProcessDataSurvey) error {
	err := inst.baseRepository.Create(userContext, item)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.Create")
	}
	return nil
}

func (inst *ProcessDataSurveyRepository) FindLast(userContext entity.UserContext, findRequest *entity.FindRequest) (*models.ProcessDataSurvey, error) {
	result := &models.ProcessDataSurvey{}
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
func (inst *ProcessDataSurveyRepository) Delete(userContext entity.UserContext, filters map[string]interface{}) error {
	item := &models.ProcessDataSurvey{}
	err := inst.baseRepository.DeleteWithAttributes(userContext, filters, item)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.DeleteWithFilters")
	}
	return nil
}

func (inst *ProcessDataSurveyRepository) DeleteByUUID(userContext entity.UserContext, item models.ProcessDataSurvey) error {
	baseItem := &models.ProcessDataSurvey{}
	err := inst.baseRepository.DeleteWithAttributes(userContext, map[string]interface{}{
		"UUID": item.UUID,
	}, baseItem)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.Update")
	}
	return nil
}

func (inst *ProcessDataSurveyRepository) UpdateByUUID(userContext entity.UserContext, item *models.ProcessDataSurvey) error {
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

func (inst *ProcessDataSurveyRepository) CreateAndUpdate(userContext entity.UserContext, item *models.ProcessDataSurvey) error {
	processDataSurvey := &models.ProcessDataSurvey{}
	findRequest := &entity.FindRequest{
		Filters: []entity.FindRequestFilter{
			entity.FindRequestFilter{
				Key:      "DomainId",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    item.DomainId,
			},
			entity.FindRequestFilter{
				Key:      "UserId",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    item.UserId,
			},
			entity.FindRequestFilter{
				Key:      "PositionItem",
				Operator: entity.FindRequestFilterOperatorEqualTo,
				Value:    item.PositionItem,
			},
		},
	}
	if findRequest.Orders == nil {
		findRequest.Orders = map[string]entity.FindRequestOrderType{}
	}
	findRequest.Orders["CreatedAt"] = entity.FindRequestOrderTypeDesc
	err := inst.baseRepository.FindFirstByFindRequest(userContext, findRequest, &processDataSurvey)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.FindFirstByFindRequest")
	}
	if processDataSurvey == nil {
		err := inst.baseRepository.Create(userContext, item)
		if err != nil {
			return errutil.Wrap(err, "baseRepository.Create")
		}
	} else {
		processDataSurvey.ProcessedData = item.ProcessedData
		err := inst.baseRepository.UpdateByFindRequest(userContext, &entity.FindRequest{
			Filters: []entity.FindRequestFilter{
				entity.FindRequestFilter{
					Key:      "UUID",
					Value:    item.UUID,
					Operator: entity.FindRequestFilterOperatorEqualTo,
				},
			},
		}, processDataSurvey)
		if err != nil {
			return errutil.Wrap(err, "baseRepository.UpdateByFindRequest")
		}
	}

	return nil
}
