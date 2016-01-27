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

	fallbacker Fallbacker
	queueChan  chan interface{}
}

func (c *Command) Queue() (chan interface{}, error) {
	//queue for processing
	go func() {
		result, err := c.fallbacker.Run()
		if err != nil {
			log.Print(err)
		}
		c.queueChan <- result
	}()

	return c.queueChan, nil
}

func (c *Command) Execute() (interface{}, error) {
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

func NewCommand(groupName string, fallbacker Fallbacker) *Command {
	return &Command{
		fallbacker: fallbacker,
		queueChan:  make(chan interface{}),
	}
}
