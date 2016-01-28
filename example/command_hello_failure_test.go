package example

import (
	"errors"
	"fmt"
	"testing"

	"github.com/shsing2000/go-hystrix/hystrix"
)

type helloFailureCommand struct {
	name string
}

func (h helloFailureCommand) Run() (interface{}, error) {
	return fmt.Sprintf("Hello %s!", h.name), errors.New("run error")
}

func (h helloFailureCommand) Fallback() (interface{}, error) {
	return fmt.Sprintf("Hello Failure %s!", h.name), nil
}

func TestHelloFailureSynchronous(t *testing.T) {
	helloWorld := helloFailureCommand{name: "World"}
	helloBob := helloFailureCommand{name: "Bob"}

	c := hystrix.NewCommand("ExampleGroup", helloWorld)
	result, err := c.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	s, ok := result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "Hello Failure World!" {
		t.Errorf("expected result to be \"Hello Failure World!\", but got \"%s\"", s)
	}

	c = hystrix.NewCommand("ExampleGroup", helloBob)
	result, err = c.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	s, ok = result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "Hello Failure Bob!" {
		t.Errorf("expected result to be \"Hello Failure Bob!\", but got \"%s\"", s)
	}
}

func TestHelloFailureAsynchronous1(t *testing.T) {
	helloWorld := helloFailureCommand{name: "World"}
	helloBob := helloFailureCommand{name: "Bob"}

	cmdWorld := hystrix.NewCommand("ExampleGroup", helloWorld)
	chanWorld, err := cmdWorld.Queue()
	if err != nil {
		t.Error("expected err to be nil")
	}

	cmdBob := hystrix.NewCommand("ExampleGroup", helloBob)
	chanBob, err := cmdBob.Queue()
	if err != nil {
		t.Error("expected err to be nil")
	}

	result := <-chanWorld
	s, ok := result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "Hello Failure World!" {
		t.Errorf("expected result to be \"Hello Failure World!\", but got \"%s\"", s)
	}

	result = <-chanBob
	s, ok = result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "Hello Failure Bob!" {
		t.Errorf("expected result to be \"Hello Failure Bob!\", but got \"%s\"", s)
	}
}

func TestHelloFailureAsynchronous2(t *testing.T) {

}
