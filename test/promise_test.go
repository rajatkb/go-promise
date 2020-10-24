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

func TestThenAndCatchTogether(t *testing.T) {
	var data int = 0
	var w sync.WaitGroup
	w.Add(1)

	Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
		reject(3)
	}).Then(func(value interface{}) (interface{}, error) {
		return nil, nil
	}).Catch(func(value interface{}) (interface{}, error) {
		tmp, _ := value.(int)
		data = tmp
		w.Done()
		return nil, nil
	})

	w.Wait()
	if data != 3 {
		t.Errorf("expected data = %d found data = %d", 3, data)
	}
}

func TestPromiseAll(t *testing.T) {
	var d bool = true
	var w sync.WaitGroup
	w.Add(1)

	promises := make([]*Promise.Promise, 20)
	for i := 0; i < 10; i++ {
		promises[i] = Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
			resolve(2)
		})
	}

	for i := 10; i < 20; i++ {
		promises[i] = Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
			reject(3)
		})
	}

	Promise.All(promises).Then(func(value interface{}) (interface{}, error) {
		w.Done()
		return nil, nil
	}).Catch(func(value interface{}) (interface{}, error) {
		v, _ := value.([]int)
		for i, r := range v {
			if i < 10 {
				if r != 2 {
					d = false
				}
			}
			if i > 10 {
				if r != 3 {
					d = false
				}
			}

		}
		w.Done()
		return nil, nil
	})

	w.Wait()

	if !d {
		t.Errorf("Promise.All failed")
	}
}

func TestChaining(t *testing.T) {
	var data int = 0
	var w sync.WaitGroup
	w.Add(3)

	Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
		resolve(2)
	}).Then(func(value interface{}) (interface{}, error) {
		tmp, _ := value.(int)
		w.Done()
		return tmp + 1, nil
	}).Then(func(value interface{}) (interface{}, error) {
		tmp, _ := value.(int)
		w.Done()

		if tmp == 3 {

			return nil, fmt.Errorf("it worked")
		}
		return nil, nil
	}).Catch(func(value interface{}) (interface{}, error) {
		data = 3
		w.Done()
		return nil, nil
	})

	w.Wait()
	if data != 3 {
		t.Errorf("expected data = %d found data = %d", 3, data)
	}

}
