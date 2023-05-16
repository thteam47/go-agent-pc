package grpcauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/thteam47/go-agent-pc/pkg/db"
	"github.com/thteam47/go-agent-pc/pkg/models"
	"github.com/thteam47/common/api/agent-pc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	handler *db.Handler
}

func NewAuthInterceptor(handler *db.Handler) *AuthInterceptor {
	return &AuthInterceptor{handler: handler}
}

func (interceptor *AuthInterceptor) Authentication(ctx context.Context, ctxRequest *pb.Context, privilege string, action string) (UserContext, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	accessToken := ""
	authorization := md["authorization"]
	if len(authorization) < 1 {
		if ctxRequest.AccessToken != "" {
			accessToken = ctxRequest.AccessToken
		} else {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}
	}
	if accessToken == "" {
		accessToken = strings.TrimPrefix(authorization[0], "Bearer ")
		if accessToken == "undefined" {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}
	}

	err := interceptor.handler.RedisRepository.GetValueCache(fmt.Sprintf("invalid-token-%s", accessToken), nil)
	if err == nil {
		return nil, status.Errorf(codes.Unauthenticated, "Token is expired")
	}

	claims, err := interceptor.VerifyToken(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	if claims.TokenInfo == nil {
		return nil, nil
	}
	userContext := NewUserContext(claims.TokenInfo)
	userContext.SetAccessToken(accessToken)

	if privilege == "@any" && action == "@any" {
		return userContext, nil
	}
	if !userContext.GetTokenInfo().AuthenticationDone {
		return nil, status.Error(codes.PermissionDenied, "Fobbiden!")
	}
	if claims.TokenInfo.PermissionAll {
		return userContext, nil
	}
	for _, permission := range claims.TokenInfo.Permissions {
		if permission.Privilege == privilege {
			return userContext, nil
			// if contains(permission.Actions, action) {
			// 	return userContext, nil
			// }
		}
	}
	return nil, status.Error(codes.PermissionDenied, "Fobbiden!")

}

func (interceptor *AuthInterceptor) VerifyToken(accessToken string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&models.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(interceptor.handler.JwtKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil

}
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
