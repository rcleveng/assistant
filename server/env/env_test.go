package env

import (
	"context"
	"testing"
)

func TestEmptyContext(t *testing.T) {
	if _, ok := FromContext(context.Background()); ok == true {
		t.Error("expected to not find env in empty context")
	}

}

func TestContext(t *testing.T) {
	environ, err := NewEnvironment(GOTEST)
	if err != nil {
		t.Fatal("failed to create Environment")
	}

	ctx := NewContext(context.Background(), environ)

	found := ctx.Value(EnvironmentKey)
	if found == nil {
		t.Error("unable to manually find server key in context")
	}

	if se, ok := found.(*Environment); ok == false {
		t.Errorf("wrong type for serverkey, got %#v", se)
	}

	if _, ok := FromContext(ctx); ok == false {
		t.Error("unable to find server key in context using FromContext")
	}

}

func TestNewServerEnv(t *testing.T) {
	if _, err := NewEnvironment(GOTEST); err != nil {
		t.Error("got error for test env")
	}
}
