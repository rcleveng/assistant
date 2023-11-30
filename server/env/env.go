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

type Environment struct {
	PalmApiKey string
	// Databse hostname
	DatabaseHostname     string
	DatabaseUserName     string
	DatabasePassword     string
	DatabaseDatabase     string
	ExecutionEnvironment ExecutionEnvironment
}

type EnvironmentKeyType int

const EnvironmentKey EnvironmentKeyType = 0

func NewEnvironment(env ExecutionEnvironment) (*Environment, error) {
	switch env {
	case GOTEST:
		return &Environment{
			PalmApiKey:           os.Getenv("PALM_KEY"),
			DatabaseHostname:     os.Getenv("PG_HOSTNAME"),
			DatabaseUserName:     os.Getenv("PG_USERNAME"),
			DatabasePassword:     os.Getenv("PG_PASSWORD"),
			DatabaseDatabase:     os.Getenv("PG_DATABASE"),
			ExecutionEnvironment: env,
		}, nil

	case COMMANDLINE:
		return &Environment{
			PalmApiKey:           os.Getenv("PALM_KEY"),
			DatabaseHostname:     os.Getenv("PG_HOSTNAME"),
			DatabaseUserName:     os.Getenv("PG_USERNAME"),
			DatabasePassword:     os.Getenv("PG_PASSWORD"),
			DatabaseDatabase:     os.Getenv("PG_DATABASE"),
			ExecutionEnvironment: COMMANDLINE,
		}, nil

	default:
		return nil, fmt.Errorf("invalid environment specified")
	}
}

func FromContext(ctx context.Context) (*Environment, bool) {
	o := ctx.Value(EnvironmentKey)
	e, ok := o.(*Environment)
	return e, ok
}

func NewContext(ctx context.Context, environment *Environment) context.Context {
	rctx := context.WithValue(ctx, EnvironmentKey, environment)
	return rctx
}
