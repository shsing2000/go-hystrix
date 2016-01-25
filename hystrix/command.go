package hystrix

import "log"

type Runner interface {
	Run() (interface{}, error)
}

type Fallbacker interface {
	Runner
	Fallback() (interface{}, error)
}

type Command struct {
	GroupKey               string
	IsFailedExecution      bool
	IsResponseFromFallback bool

	queueChan chan interface{}
}

type RunCommand struct {
	Command

	runner Runner
}

func (c *RunCommand) Execute() (interface{}, error) {
	return c.runner.Run()
}

func (c *RunCommand) Queue() (chan interface{}, error) {
	//queue for processing
	go func() {
		result, err := c.runner.Run()
		if err != nil {
			log.Print(err)
		}
		c.queueChan <- result
	}()

	return c.queueChan, nil
}

type FallbackCommand struct {
	Command

	fallbacker Fallbacker
}

func (c *FallbackCommand) Execute() (interface{}, error) {
	result, err := c.fallbacker.Run()
	if err == nil {
		return result, err
	}

	//failed execution
	c.IsFailedExecution = true
	result, err = c.fallbacker.Fallback()
	if err == nil {
		c.IsResponseFromFallback = true
	}

	return result, nil
}

func NewRunCommand(groupName string, runner Runner) *RunCommand {
	return &RunCommand{
		runner:  runner,
		Command: Command{queueChan: make(chan interface{})},
	}
}

func NewFallbackCommand(groupName string, fallbacker Fallbacker) *FallbackCommand {
	return &FallbackCommand{
		fallbacker: fallbacker,
		Command:    Command{queueChan: make(chan interface{})},
	}
}
