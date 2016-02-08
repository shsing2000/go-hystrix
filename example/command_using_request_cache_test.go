package example

import (
	"strconv"
	"testing"

	"github.com/shsing2000/go-hystrix/hystrix"
)

type requestCacheCommand struct {
	value int
}

func (r requestCacheCommand) Run() (interface{}, error) {
	result := r.value == 0 || r.value%2 == 0
	return result, nil
}

func (r requestCacheCommand) Fallback() (interface{}, error) {
	return false, nil
}

func (r requestCacheCommand) GetCacheKey() string {
	return strconv.Itoa(r.value)
}

func TestRequestCacheWithoutCacheHits(t *testing.T) {
	//hystrix initialize context

	c1 := hystrix.NewCommand("ExampleGroup", requestCacheCommand{value: 2})
	result, err := c1.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	b, ok := result.(bool)
	if !ok {
		t.Error("expected result to be a bool type")
	}
	if !b {
		t.Error("expected result to be true, but got false")
	}
	//hystrix shutdown context
}

func TestRequestCacheWithCacheHits(t *testing.T) {
	//initialize context
	command2a := hystrix.NewCommand("ExampleGroup", requestCacheCommand{value: 2})
	command2b := hystrix.NewCommand("ExampleGroup", requestCacheCommand{value: 2})

	result, err := command2a.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	b, ok := result.(bool)
	if !ok {
		t.Error("expected result to be a bool type")
	}
	if !b {
		t.Error("expected result to be true, but got false")
	}
	if command2a.IsResponseFromCache() {
		t.Error("expected result to be executed, but was from cache")
	}

	result, err = command2b.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	b, ok = result.(bool)
	if !ok {
		t.Error("expected result to be a bool type")
	}
	if !b {
		t.Error("expected result to be true, but got false")
	}
	if !command2b.IsResponseFromCache() {
		t.Error("expected result to be from cache, but was executed")
	}
	//shutdown context

	//start new context
	command2c := hystrix.NewCommand("ExampleGroup", requestCacheCommand{value: 2})
	result, err = command2c.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	b, ok = result.(bool)
	if !ok {
		t.Error("expected result to be a bool type")
	}
	if !b {
		t.Error("expected result to be true, but got false")
	}
	if command2c.IsResponseFromCache() {
		t.Error("expected result to be executed, but was from cache")
	}
	//shutdown context
}
