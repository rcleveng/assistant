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
	COMMANDLINE
)

type ServerEnvironment struct {
	PalmApiKey           string
	DatabaseConnection   string
	DatabaseUserName     string
	DatabasePassword     string
	ExecutionEnvironment ExecutionEnvironment
}

type ServerEnvironmentKeyType int

const ServerEnvironmentKey ServerEnvironmentKeyType = 0

func NewServerEnvironment(env ExecutionEnvironment) (*ServerEnvironment, error) {
	switch env {
	case GOTEST:
		return &ServerEnvironment{
			PalmApiKey:           os.Getenv("PALM_KEY"),
			DatabaseConnection:   os.Getenv("PG_URL"),
			DatabaseUserName:     os.Getenv("PG_USERNAME"),
			DatabasePassword:     os.Getenv("PG_PASSWORD"),
			ExecutionEnvironment: env,
		}, nil

	case COMMANDLINE:
		return &ServerEnvironment{
			PalmApiKey:           os.Getenv("PALM_KEY"),
			DatabaseConnection:   os.Getenv("PG_URL"),
			DatabaseUserName:     os.Getenv("PG_USERNAME"),
			DatabasePassword:     os.Getenv("PG_PASSWORD"),
			ExecutionEnvironment: COMMANDLINE,
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
