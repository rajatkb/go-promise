package promise

import (
	"fmt"
	"sync"
)

func getSyncStateBool(initialState bool) (func() bool, func(bool)) {
	var lock sync.Mutex
	init := initialState
	read := func() bool {
		lock.Lock()
		v := init
		lock.Unlock()
		return v
	}
	write := func(v bool) {
		lock.Lock()
		init = v
		lock.Unlock()
	}
	return read, write
}

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
			} else if obj.resolved && obj.failed {
				reject(obj.value)
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
			} else if obj.resolved && !obj.failed {
				resolve(obj.value)
			}
		}()
	})
}

//All ... resolves a Promise when all promises passed are resolved,
func All(promises []*Promise) *Promise {

	return Create(func(resolve Callback, reject Callback) {
		go func() {
			var w sync.WaitGroup

			resolveStateR, resolveStateW := getSyncStateBool(true)

			data := make([]interface{}, len(promises))
			w.Add(len(promises))
			for i, promise := range promises {
				index := i // because go catches current i value not the ones that was encountered when loop was at this loop state
				if promise == nil {
					data[index] = nil
					w.Done()
					continue
				}
				promise.Then(func(value interface{}) (interface{}, error) {
					data[index] = value
					w.Done()
					resolveStateW(true && resolveStateR())
					return nil, nil
				})
				promise.Catch(func(value interface{}) (interface{}, error) {

					data[index] = value
					resolveStateW(false && resolveStateR())
					w.Done()
					return nil, nil
				})
			}
			w.Wait()
			if resolveStateR() {
				resolve(data)
			} else {
				reject(data)
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
