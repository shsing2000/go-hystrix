package example

import (
	"errors"
	"testing"

	"github.com/shsing2000/go-hystrix/hystrix"
)

type commandWithFallbackViaNetworkCommand struct {
	id         int
	commandKey string
}

func (c commandWithFallbackViaNetworkCommand) Run() (interface{}, error) {
	return "", errors.New("force failure for example")
}

func (c commandWithFallbackViaNetworkCommand) Fallback() (interface{}, error) {
	fallbacker := &fallbackViaNetworkCommand{
		id:         c.id,
		commandKey: "GetValueFallbackCommand",
	}
	c1 := hystrix.NewCommand("RemoteServiceXFallback", fallbacker)
	return c1.Execute()
}

type fallbackViaNetworkCommand struct {
	id         int
	commandKey string
}

func (f fallbackViaNetworkCommand) Run() (interface{}, error) {
	return "", errors.New("the fallback also failed")
}

func (f fallbackViaNetworkCommand) Fallback() (interface{}, error) {
	// the fallback also failed
	// so we will always return an error
	// and let the caller react
	return "", errors.New("fallback-of-a-fallback failed")
}

func TestFallbackViaNetwork(t *testing.T) {
	//context = hystrix initialize context
	fallbacker := &commandWithFallbackViaNetworkCommand{
		id:         1,
		commandKey: "GetValueCommand",
	}
	c := hystrix.NewCommand("RemoteServiceX", fallbacker)
	result, err := c.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	s, ok := result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "" {
		t.Errorf("expected result to be empty, but got %s", s)
	}

	// verify request log for 2 failures
	// failure on command "GetValueCommand"
	// failure on command "GetValueFallbackCommand"
}
