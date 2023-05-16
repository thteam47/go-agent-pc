package recommendservice

import (
	"github.com/thteam47/common/entity"
	apiclient "github.com/thteam47/common/pkg/apiswagger/recommend-api"
	"github.com/thteam47/go-agent-pc/pkg/models"
)

type RecommendService interface {
	KeyPublicUserSend(userContext entity.UserContext, domain string, tokenAgent string, keyInfo *models.KeyInfo) error
	GetCombinedData(userContext entity.UserContext, tokenAgent string, tenantId string) (*models.CombinedData, error)
	KeyPublicItemSend(userContext entity.UserContext, domain string, tokenAgent string, keyInfo *models.KeyInfoItem) error
	KeyPublicUse(userContext entity.UserContext, tenantId string, tokenAgent string, positionItem int32, part int32) (*models.KeyPublicUse, error)
	ProcessedDataSendPart(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiResultCard) error
	ProcessedDataSendPhase3TwoPart(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiPhase3TwoPart) error
	ProcessedDataSendPhase3Get(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiPhase3TwoPart) (*apiclient.RecommendApiPhase3TwoPart, error)
	ProcessedDataSendPhase4TwoPart(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiPhase4TwoPart) error
	ProcessedDataSendPhase4Get(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiPhase4TwoPart) (*apiclient.RecommendApiPhase4TwoPart, error)
	KeyPublicUserGet(userContext entity.UserContext, domain string, tokenAgent string, positionUser int32) (*apiclient.RecommendApiKeyPublicUser, error)
	ProcessedDataSend2(userContext entity.UserContext, domain string, tokenAgent string, data *apiclient.RecommendApiProcessDataSurvey2) error
	RecommendCbfGenarate(userContext entity.UserContext, domain string, tokenAgent string, userId string, data map[string]apiclient.RecommendApiRecommend) (map[string]apiclient.RecommendApiRecommendCbfResult12, *apiclient.RecommendApiRecommendCbfResult34, error)
	RecommendCfGenarate(userContext entity.UserContext, domain string, tokenAgent string, userId string, data map[string]apiclient.RecommendApiRecommend) (map[string]apiclient.RecommendApiRecommendCfResult910, *apiclient.RecommendApiRecommendCfResult1112, error)
}
