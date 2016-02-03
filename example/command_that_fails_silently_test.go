package example

import (
	"errors"
	"testing"

	"github.com/shsing2000/go-hystrix/hystrix"
)

type failSilentCommand struct {
	err error
}

func (f *failSilentCommand) Run() (interface{}, error) {
	if f.err != nil {
		return []string{}, f.err
	} else {
		return []string{"success"}, nil
	}
}

func (f *failSilentCommand) Fallback() (interface{}, error) {
	return []string{}, nil
}

func TestFailSilentSuccess(t *testing.T) {
	f := &failSilentCommand{err: nil}
	c := hystrix.NewCommand("ExampleGroup", f)
	result, err := c.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	s, ok := result.([]string)
	if !ok {
		t.Error("expected result to be a []string type")
	}
	if s[0] != "success" {
		t.Errorf("expected result to be \"success\", but got \"%s\"", s[0])
	}
}

func TestFailSilentFailure(t *testing.T) {
	f := &failSilentCommand{err: errors.New("some silent error")}
	c := hystrix.NewCommand("ExampleGroup", f)
	result, err := c.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	s, ok := result.([]string)
	if !ok {
		t.Error("expected result to be a []string type")
	}
	if len(s) != 0 {
		t.Errorf("expected result to be empty, but got \"%d\" elements", len(s))
	}
}
