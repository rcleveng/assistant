package env

import (
	"context"
	"fmt"
	"os"
)

type ExecutionEnvironment int

const (
	GOTEST ExecutionEnvironment = iota
	TEST
	STAGING
	PRODUCTION_CLOUDRUN
)

type ServerEnvironment struct {
	PalmApiKey           func() (string, error)
	ExecutionEnvironment ExecutionEnvironment
}

type ServerEnvironmentKeyType int

const ServerEnvironmentKey ServerEnvironmentKeyType = 0

func NewServerEnvironment(env ExecutionEnvironment) (*ServerEnvironment, error) {
	switch env {
	case GOTEST:
		return &ServerEnvironment{
			PalmApiKey: func() (string, error) {
				return os.Getenv("PALM_KEY"), nil
			},
			ExecutionEnvironment: env,
		}, nil
	default:
		return nil, fmt.Errorf("invalid environment specified")
	}
}

func FromContext(ctx context.Context) (*ServerEnvironment, bool) {
	o := ctx.Value(ServerEnvironmentKey)
	e, ok := o.(*ServerEnvironment)
	return e, ok
}

func NewContext(ctx context.Context, environment *ServerEnvironment) context.Context {
	rctx := context.WithValue(ctx, ServerEnvironmentKey, environment)
	return rctx
}
