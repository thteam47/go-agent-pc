package authenservice

import (
	"github.com/thteam47/common/entity"
)

type AuthenService interface {
	SetAccessToken(userContext entity.UserContext, tokenAgent string, accessToken string)
	AccessToken(userContext entity.UserContext, tokenAgent string) string
	PrepareLogin(userContext entity.UserContext, token string) (string, string, string, string, string, []string, string, *entity.ErrorInfo, error)
	Login(userContext entity.UserContext, loginType string, loginName string, password string, otp int32, token string, requestId string, userType string) (string, *entity.ErrorInfo, error)
	//LogOut(userContext entity.UserContext, token string) error
	SetAuthenInfo(userContext entity.UserContext, tokenAgent string, userInfo *entity.User)
	GetUserInfo(userContext entity.UserContext, tokenAgent string) *entity.User
}

type AuthenInfo struct {
	UserInfo  *entity.User      `json:"user_info,omitempty"`
	TokenInfo *entity.TokenInfo `json:"token_info,omitempty"`
}
