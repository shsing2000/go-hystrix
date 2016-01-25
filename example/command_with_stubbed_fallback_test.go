package example

import (
	"errors"
	"testing"

	"github.com/shsing2000/go-hystrix/hystrix"
)

//start of tests
type user struct {
	customerId          int
	name                string
	countryCode         string
	isFeatureXPermitted bool
	isFeatureYPermitted bool
	isFeatureZPermitted bool
}

type mockFallbacker struct {
	customerId  int
	countryCode string
}

func (m *mockFallbacker) Run() (interface{}, error) {
	return user{}, errors.New("forced failure for example")
}

func (m *mockFallbacker) Fallback() (interface{}, error) {
	return user{
		customerId:          m.customerId,
		countryCode:         m.countryCode,
		isFeatureXPermitted: true,
		isFeatureYPermitted: true,
		isFeatureZPermitted: false,
	}, nil
}

func TestFallbackCommand(t *testing.T) {
	fallbacker := &mockFallbacker{
		customerId:  1234,
		countryCode: "ca",
	}
	c := hystrix.NewFallbackCommand("ExampleGroup", fallbacker)
	result, err := c.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	u, ok := result.(user)
	if !ok {
		t.Error("expected result to be a user type")
	}
	if !c.IsFailedExecution {
		t.Error("expected failed execution to be true")
	}
	if !c.IsResponseFromFallback {
		t.Error("expected response from fallback to be true")
	}
	if u.customerId != 1234 {
		t.Errorf("expected customerId to be 1234, but got %d", u.customerId)
	}
	if u.countryCode != "ca" {
		t.Errorf("expected countryCode to be \"ca\", but got %s", u.countryCode)
	}
	if !u.isFeatureXPermitted {
		t.Error("expected isFeatureXPermitted to be true")
	}
	if !u.isFeatureYPermitted {
		t.Error("expected isFeatureYPermitted to be true")
	}
	if u.isFeatureZPermitted {
		t.Error("expected isFeatureZPermitted to be false")
	}
}
