package example

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/shsing2000/go-hystrix/hystrix"
)

var prefixStoredOnRemoteDataStore = "ValueBeforeSet_"

type flusher interface {
	flushCache(id int)
}
type requestCacheGetterCommand struct {
	id int
}

//command key is GetterCommand
func (r requestCacheGetterCommand) Run() (interface{}, error) {
	return fmt.Sprintf("%s%d", prefixStoredOnRemoteDataStore, r.id), nil
}

func (r requestCacheGetterCommand) Fallback() (interface{}, error) {
	return "", nil
}

func (r requestCacheGetterCommand) GetCacheKey() string {
	return strconv.Itoa(r.id)
}

func (r *requestCacheGetterCommand) flushCache(id int) {
	//flush the cache
}

type requestCacheSetterCommand struct {
	id     int
	prefix string
	getter flusher
}

func (r *requestCacheSetterCommand) Run() (interface{}, error) {
	prefixStoredOnRemoteDataStore = r.prefix
	r.getter.flushCache(r.id)
	return "", nil
}

func (r *requestCacheSetterCommand) Fallback() (interface{}, error) {
	return "", nil
}

func TestCacheInvalidationGetSetGet(t *testing.T) {
	//initialize context
	getter1 := &requestCacheGetterCommand{id: 1}
	c1 := hystrix.NewCommand("GetSetGet", getter1)
	result, err := c1.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	s, ok := result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "ValueBeforeSet_1" {
		t.Errorf("expected result to be ValueBeforeSet_1, but got %s", s)
	}

	getter2 := &requestCacheGetterCommand{id: 1}
	c2 := hystrix.NewCommand("GetSetGet", getter2)
	result, err = c2.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	s, ok = result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "ValueBeforeSet_1" {
		t.Errorf("expected result to be ValueBeforeSet_1, but got %s", s)
	}
	if !c2.IsResponseFromCache() {
		t.Error("expected result to be from cache, but was executed")
	}

	//set the new value
	setter := &requestCacheSetterCommand{
		id:     1,
		prefix: "ValueAfterSet_",
		getter: getter2,
	}
	c3 := hystrix.NewCommand("GetSetGet", setter)
	c3.Execute()

	getter4 := &requestCacheGetterCommand{id: 1}
	c4 := hystrix.NewCommand("GetSetGet", getter4)
	result, err = c4.Execute()
	if err != nil {
		t.Error("expected err to be nil")
	}
	s, ok = result.(string)
	if !ok {
		t.Error("expected result to be a string type")
	}
	if s != "ValueAfterSet_1" {
		t.Errorf("expected result to be ValueAfterSet_1, but got %s", s)
	}
	if c4.IsResponseFromCache() {
		t.Error("expected result to be executed, but was from cache")
	}
	//shutdown context
}
