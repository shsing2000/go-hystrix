package example

import (
	"fmt"
	"testing"

	"github.com/shsing2000/go-hystrix/hystrix"
)

type helloWorldCommand struct {
	name string
}

func (h helloWorldCommand) Run() (interface{}, error) {
	return fmt.Sprintf("Hello %s!", h.name), nil
}

func TestSynchronous(t *testing.T) {
	helloWorld := helloWorldCommand{name: "World"}
	helloBob := helloWorldCommand{name: "Bob"}

	c := hystrix.NewRunCommand("ExampleGroup", helloWorld)
	result, err := c.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	s, ok := result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "Hello World!" {
		t.Errorf("expected result to be \"Hello World!\", but got \"%s\"", s)
	}

	c = hystrix.NewRunCommand("ExampleGroup", helloBob)
	result, err = c.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	s, ok = result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "Hello Bob!" {
		t.Errorf("expected result to be \"Hello Bob!\", but got \"%s\"", s)
	}
}

func TestAsynchronous(t *testing.T) {
	helloWorld := helloWorldCommand{name: "World"}
	helloBob := helloWorldCommand{name: "Bob"}

	cmdWorld := hystrix.NewRunCommand("ExampleGroup", helloWorld)
	chanWorld, err := cmdWorld.Queue()

	cmdBob := hystrix.NewRunCommand("ExampleGroup", helloBob)
	chanBob, err := cmdBob.Queue()
	if err != nil {
		t.Error("expected err to be nil")
	}

	result := <-chanWorld
	s, ok := result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "Hello World!" {
		t.Errorf("expected result to be \"Hello World!\", but got \"%s\"", s)
	}

	result = <-chanBob
	s, ok = result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "Hello Bob!" {
		t.Errorf("expected result to be \"Hello Bob!\", but got \"%s\"", s)
	}
}

func TestObservable(t *testing.T) {

}
