package test

import (
	"fmt"
	"sync"
	"testing"

	Promise "github.com/rajatkb/go-promise"
)

func TestCreation(t *testing.T) {

	var w sync.WaitGroup
	var data int = 0
	w.Add(1)
	Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
		resolve(2)
	}).Then(func(value interface{}) (interface{}, error) {
		tmp, _ := value.(int)
		data = tmp
		w.Done()
		return nil, nil
	})
	w.Wait()
	if data == 0 {
		t.Errorf("expected data = %d found data = %d", 2, data)
	}
}

func TestThenCallback(t *testing.T) {

	var w sync.WaitGroup
	var data int = 0
	w.Add(2)
	promise := Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
		resolve(2)
	})
	promise.Then(func(value interface{}) (interface{}, error) {
		tmp, ok := value.(int)
		if ok {
			data = tmp
		}
		w.Done()
		return tmp, nil
	}).Then(func(value interface{}) (interface{}, error) {
		tmp, ok := value.(int)
		if ok {
			data = tmp
		}
		w.Done()
		return tmp, nil
	})

	w.Wait()
	if data == 0 {
		t.Errorf("expected data = %d found data = %d", 2, data)
	}
}

func TestChaining(t *testing.T) {

	var w sync.WaitGroup
	w.Add(3)

	Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
		fmt.Println("Promise called")
		resolve(2)
	}).Then(func(value interface{}) (interface{}, error) {
		fmt.Println("Then called")
		tmp, _ := value.(int)
		fmt.Println("found value", tmp)
		w.Done()
		return tmp, nil
	}).Then(func(value interface{}) (interface{}, error) {
		fmt.Println("Then called")
		tmp, _ := value.(int)
		fmt.Println("found value", tmp)
		w.Done()

		if tmp == 2 {

			return nil, fmt.Errorf("it worked")
		}
		return nil, nil
	}).Catch(func(value interface{}) (interface{}, error) {
		fmt.Println(value)
		w.Done()
		return nil, nil
	})
}
