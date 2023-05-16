package authenservice

import (
	"fmt"
	"math/rand"
	"time"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/thteam47/common-libs/jwtutil"
	"github.com/thteam47/common-libs/reflectutil"
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/pkg/apiclientutil"
	"github.com/thteam47/go-agent-pc/pkg/models"

	"github.com/thteam47/common-libs/confg"
	apiclient "github.com/thteam47/common/pkg/apiswagger/identity-authen-api"
	"github.com/thteam47/go-agent-pc/errutil"
)

type DefaultAuthenService struct {
	config        *DefaultAuthenServiceConfig
	accessToken   map[string]string
	authenInfo    map[string]*AuthenInfo
	apiClientInst *apiclient.APIClient
}

type DefaultAuthenServiceConfig struct {
	Port string `mapstructure:"port"`
}

func NewDefaultAuthenServiceWithConfig(properties confg.Confg) (*DefaultAuthenService, error) {
	config := DefaultAuthenServiceConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}
	return NewDefaultAuthenService(&config)
}

func NewDefaultAuthenService(config *DefaultAuthenServiceConfig) (*DefaultAuthenService, error) {
	inst := &DefaultAuthenService{
		config: config,
	}

	inst.authenInfo = map[string]*AuthenInfo{}

	return inst, nil
}

func (inst *DefaultAuthenService) apiClient() *apiclient.APIClient {
	if inst.apiClientInst == nil {
		inst.apiClientInst = apiclient.NewAPIClient(&apiclient.Configuration{
			BasePath: fmt.Sprintf("http://127.0.0.1:%s", inst.config.Port),
			Scheme:   "http",
		})
	}
	return inst.apiClientInst
}

func generateToken() string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 50)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (inst *DefaultAuthenService) PrepareLogin(userContext entity.UserContext, token string) (string, string, string, string, string, []string, string, *entity.ErrorInfo, error) {
	if token == "" {
		token = generateToken()
		return token, "", "", "", "", []string{}, "", nil, nil
	}
	if inst.authenInfo == nil {
		inst.authenInfo = map[string]*AuthenInfo{}
	}
	if inst.authenInfo[token] == nil {
		inst.authenInfo[token] = &AuthenInfo{}
	}
	if inst.authenInfo[token].TokenInfo == nil {
		inst.authenInfo[token].TokenInfo = &entity.TokenInfo{}
	}
	domainId := userContext.DomainId()
	if inst.authenInfo[token].UserInfo == nil {
		inst.authenInfo[token].UserInfo = &entity.User{}
	} else if inst.authenInfo[token].UserInfo.DomainId != "" {
		domainId = inst.authenInfo[token].UserInfo.DomainId
	}
	accessToken := ""
	if inst.accessToken != nil {
		accessToken = inst.AccessToken(userContext, token)
	}
	
	response, _, err := inst.apiClient().IdentityAuthenServiceApi.PrepareLogin(userContext.Context(), domainId, apiclient.Body3{
		Ctx: &apiclient.V1identityauthenapictxDomainIdforgotPasswordCtx{
			AccessToken: accessToken,
		},
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return "", "", "", "", "", []string{}, "", nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityAuthenServiceApi.PrepareLogin")
	}
	errorInfo := &entity.ErrorInfo{}
	if response.ErrorCode != 0 {
		errorInfo.ErrorCode = int32(response.ErrorCode)
		errorInfo.Message = response.Message
	}
	if len(response.AvailableMfas) == 0 {
		err = inst.parseAccessToken(userContext, token, response.Token)
		if err != nil {
			return "", "", "", "", "", []string{}, "", nil, errutil.Wrap(apiclientutil.NormalizeError(err), "parseAccessToken")
		}
		if response.ErrorCode == 0 {
			inst.SetAccessToken(userContext, token, response.Token)
		}
	}
	return token, response.Token, response.RequestId, response.Message, response.TypeMfa, response.AvailableMfas, response.Url, errorInfo, nil
}
func (inst *DefaultAuthenService) SetAccessToken(userContext entity.UserContext, tokenAgent string, accessToken string) {
	if inst.accessToken == nil {
		inst.accessToken = map[string]string{}
	}
	if tokenAgent != "" && accessToken != "" {
		inst.accessToken[tokenAgent] = accessToken
	}
}

func (inst *DefaultAuthenService) AccessToken(userContext entity.UserContext, tokenAgent string) string {
	return inst.accessToken[tokenAgent]
}

func (inst *DefaultAuthenService) SetAuthenInfo(userContext entity.UserContext, tokenAgent string, userInfo *entity.User) {
	if inst.authenInfo == nil {
		inst.authenInfo = map[string]*AuthenInfo{}
	}
	if tokenAgent != "" && userInfo != nil {
		inst.authenInfo[tokenAgent] = &AuthenInfo{
			UserInfo: userInfo,
		}
	}
}

func (inst *DefaultAuthenService) parseAccessToken(userContext entity.UserContext, tokenAgent string, accessToken string) error {
	if inst.authenInfo == nil {
		inst.authenInfo = map[string]*AuthenInfo{}
	}
	claims := &models.Claims{}
	err := jwtutil.ParseToken(accessToken, "", claims)
	if err != nil && !errutil.IsError(err, jwt.ErrSignatureInvalid) {
		return errutil.Wrap(err, "jwtutil.ParseToken")
	}
	tokenInfo := &entity.TokenInfo{}
	err = reflectutil.Convert(claims.TokenInfo, &tokenInfo)
	if err != nil {
		return errutil.Wrap(err, "reflectutil.Convert")
	}
	inst.authenInfo[tokenAgent].TokenInfo = tokenInfo
	if !inst.authenInfo[tokenAgent].TokenInfo.Get("UserInfo", &inst.authenInfo[tokenAgent].UserInfo) || inst.authenInfo[tokenAgent].UserInfo == nil {
		return errutil.NewWithMessage("Invalid token")
	}
	return nil
}

func (inst *DefaultAuthenService) GetUserInfo(userContext entity.UserContext, tokenAgent string) *entity.User {
	if inst.authenInfo[tokenAgent] != nil {
		return inst.authenInfo[tokenAgent].UserInfo
	}
	return nil
}

func (inst *DefaultAuthenService) Login(userContext entity.UserContext, loginType string, loginName string, password string, otp int32, token string, requestId string, userType string) (string, *entity.ErrorInfo, error) {
	accessToken := inst.accessToken[token]
	if loginType == "UsernamePassword" {
		accessToken = ""
	}
	response, _, err := inst.apiClient().IdentityAuthenServiceApi.Login(userContext.Context(), userContext.DomainId(), apiclient.Body2{
		Ctx: &apiclient.V1identityauthenapictxDomainIdforgotPasswordCtx{
			AccessToken: accessToken,
		},
		Type_:     loginType,
		Username:  loginName,
		Password:  password,
		Otp:       otp,
		RequestId: requestId,
		UserType:  userType,
	})
	if err != nil {
		if detailedError, ok := err.(apiclient.GenericSwaggerError); ok {
			err = apiclientutil.Error(detailedError.Body())
		}
		return "", nil, errutil.Wrap(apiclientutil.NormalizeError(err), "IdentityAuthenServiceApi.Login")
	}

	if response.ErrorCode == 0 {
		inst.SetAccessToken(userContext, token, response.Token)
	}

	return response.Token, &entity.ErrorInfo{
		Message:   response.Message,
		ErrorCode: response.ErrorCode,
	}, nil
}
