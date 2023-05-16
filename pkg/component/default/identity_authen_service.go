package defaultcomponent

import (
	"context"
	"time"

	"github.com/thteam47/common-libs/confg"
	v1 "github.com/thteam47/common/api/identity-authen-api"
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-agent-pc/errutil"
	"google.golang.org/grpc"
)

type IdentityAuthenService struct {
	config *IdentityAuthenServiceConfig
	client v1.IdentityAuthenServiceClient
}

type IdentityAuthenServiceConfig struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	AccessToken string        `mapstructure:"access_token"`
}

func NewIdentityAuthenServiceWithConfig(properties confg.Confg) (*IdentityAuthenService, error) {
	config := IdentityAuthenServiceConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}
	return NewIdentityAuthenService(&config)
}

func NewIdentityAuthenService(config *IdentityAuthenServiceConfig) (*IdentityAuthenService, error) {
	inst := &IdentityAuthenService{
		config: config,
	}
	conn, err := grpc.Dial(config.Address, grpc.WithInsecure())
	if err != nil {
		return nil, errutil.Wrapf(err, "grpc.Dial")
	}
	client := v1.NewIdentityAuthenServiceClient(conn)
	inst.client = client
	return inst, nil
}

func (inst *IdentityAuthenService) requestCtx(userContext entity.UserContext) *v1.Context {
	return &v1.Context{
		AccessToken: inst.config.AccessToken,
		DomainId:    userContext.DomainId(),
	}
}

func (inst *IdentityAuthenService) UpdatePassword(userContext entity.UserContext, userId string, password string) error {
	_, err := inst.client.UpdatePassword(context.Background(), &v1.UpdatePasswordRequest{
		Ctx:      inst.requestCtx(userContext),
		UserId:   userId,
		Password: password,
	})
	if err != nil {
		return errutil.Wrapf(err, "client.UpdatePassword")
	}
	return nil
}
