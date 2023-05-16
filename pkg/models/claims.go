package models

import (
	"github.com/dgrijalva/jwt-go/v4"
)

type Claims struct {
	TokenInfo *TokenInfo `json:"token_info,omitempty" bson:"token_info,omitempty"`
}

func (inst *Claims) Valid(validationHelper *jwt.ValidationHelper) error {
	return nil
}
