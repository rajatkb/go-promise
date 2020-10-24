package promise

import (
	"fmt"
	"sync"
)

//Callback ... type for callbacks
type Callback func(interface{}) (interface{}, error)

//Promise ... represents a promise task.
type Promise struct {
	resolved      bool
	failed        bool
	value         interface{}
	valueWaitLock sync.WaitGroup
	valueLock     sync.Mutex
}

//Then ... for promise resolve
func (obj *Promise) Then(callback Callback) *Promise {
	return Create(func(resolve Callback, reject Callback) {
		go func() {
			obj.valueWaitLock.Wait()
			if obj.resolved && !obj.failed {
				value, err := callback(obj.value)
				if err == nil {
					resolve(value)
				} else {
					reject(err)
				}
			}
		}()
	})
}

//Catch ... for promise fail
func (obj *Promise) Catch(callback Callback) *Promise {
	return Create(func(resolve Callback, reject Callback) {
		go func() {
			obj.valueWaitLock.Wait()
			if obj.resolved && obj.failed {
				value, err := callback(obj.value)
				if err == nil {
					resolve(value)
				} else {
					reject(err)
				}
			}
		}()
	})
}

//Create ... creates a promise object
func Create(action func(resolve Callback, reject Callback)) *Promise {
	promise := new(Promise)
	promise.resolved = false
	promise.valueWaitLock.Add(1)

	go action(func(data interface{}) (interface{}, error) {
		if promise.resolved == false {
			promise.valueLock.Lock()
			promise.value = data
			promise.resolved = true
			promise.failed = false
			promise.valueLock.Unlock()
			promise.valueWaitLock.Done()
		}
		return promise.value, nil

	}, func(err interface{}) (interface{}, error) {
		if promise.resolved == false {
			promise.valueLock.Lock()
			promise.value = err
			promise.resolved = true
			promise.failed = true
			promise.valueLock.Unlock()
			promise.valueWaitLock.Done()
		}
		v, _ := promise.value.(string)
		return promise.value, fmt.Errorf(v)
	})

	return promise
}
