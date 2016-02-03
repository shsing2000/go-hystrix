package example

import "testing"

type failFastCommand struct {
	err error
}

func (f *failFastCommand) Run() (string, error) {
	if f.err != nil {
		return "", f.err
	} else {
		return "success", nil
	}
}

func TestFailFastSuccess(t *testing.T) {
	//maybe use a separate command and fallback command
	//this test is specific to a command since it should
	//"fail fast" for failures, rejections, short-circuiting
}

func TestFailFastFailure(t *testing.T) {

}
