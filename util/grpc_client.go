package util

import (
	"context"
	"time"

	"github.com/thteam47/go-agent-pc/pkg/configs"
	"google.golang.org/grpc"
)

type GrpcClientConn struct {
	Config  configs.GrpcClientConnConfig
	Conn    *grpc.ClientConn
	Address string
}

var GrpcClientConnTimeout = 20 * time.Second

func NewGrpcClientConnWithConfig(config configs.GrpcClientConnConfig) (*GrpcClientConn, error) {
	if config.Timeout == 0 {
		config.Timeout = GrpcClientConnTimeout
	}

	return NewGrpcClientConn(config), nil
}

func NewGrpcClientConn(config configs.GrpcClientConnConfig) *GrpcClientConn {
	inst := &GrpcClientConn{
		Config: config,
	}

	inst.Address = inst.Config.Address

	return inst
}

func (inst *GrpcClientConn) Context() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), inst.Config.Timeout)
	defer cancel()
	return ctx
}

func (inst *GrpcClientConn) Stop() {
	if inst.Conn != nil {
		inst.Conn.Close()
		inst.Conn = nil
	}
}
