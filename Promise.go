package promise

import (
	"fmt"
	"sync"
)

func getSyncStateBool(initialState bool) (func() bool, func(bool), func(bool) bool) {
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

	readWrite := func(v bool) bool {
		lock.Lock()
		rt := init
		init = v
		lock.Unlock()
		return rt
	}
	return read, write, readWrite
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

func (obj *Promise) isFulfilled() bool {
	obj.valueLock.Lock()
	resolved := obj.resolved
	failed := obj.failed
	obj.valueLock.Unlock()
	return resolved && !failed
}

func (obj *Promise) isRejected() bool {
	obj.valueLock.Lock()
	resolved := obj.resolved
	failed := obj.failed
	obj.valueLock.Unlock()
	return resolved && failed
}

func (obj *Promise) isPending() bool {
	obj.valueLock.Lock()
	resolved := obj.resolved
	obj.valueLock.Unlock()
	return resolved
}

//Then ... for promise resolve
func (obj *Promise) Then(callback Callback) *Promise {
	return Create(func(resolve Callback, reject Callback) {

		obj.valueWaitLock.Wait()
		if obj.resolved && !obj.failed {
			value, err := callback(obj.value)

			/**
			* When return value is a Promise.Resolve statement or a new *Promise
			*
			***/
			vp, ok := value.(*Promise)
			if ok {
				vp.valueWaitLock.Wait()
				if vp.resolved && !vp.failed {
					resolve(vp.value)
					return
				}
				if vp.resolved && vp.failed {
					reject(vp.value)
				}
			}

			if err == nil {
				resolve(value)
				return
			}
			reject(err)
			return

		} else if obj.resolved && obj.failed {
			reject(obj.value)
			return
		}

	})
}

//Catch ... for promise fail
func (obj *Promise) Catch(callback Callback) *Promise {
	return Create(func(resolve Callback, reject Callback) {

		obj.valueWaitLock.Wait()
		if obj.resolved && obj.failed {
			value, err := callback(obj.value)

			/**
			* When return value is a Promise.Resolve statement or a new *Promise
			*
			***/
			vp, ok := value.(*Promise)
			if ok {
				vp.valueWaitLock.Wait()
				if vp.resolved && !vp.failed {
					resolve(vp.value)
					return
				}
				if vp.resolved && vp.failed {
					reject(vp.value)
				}
			}

			if err == nil {
				resolve(value)
				return
			}
			reject(err)
			return
		} else if obj.resolved && !obj.failed {
			resolve(obj.value)
			return
		}

	})
}

//Finally ... a synchronous call to finally do something at the end of Promise Chain
func (obj *Promise) Finally(callback func(interface{})) interface{} {
	obj.valueWaitLock.Wait()
	if callback != nil {
		callback(obj.value)
	}
	return obj.value
}

/**
* TO-DO : Props
* Promise.props(struct {
*  field1 : Promise.Resolve(1)
*  field2
* })
*
**/

//Map ... for the bluebird affiniadoes
func Map(promises []*Promise) *Promise {
	return All(promises)
}

//Reduce ... asynchronous reducer , does not waits for all promise to be resolve , it launches reduce callback as soon as first result is available
// func Reduce(promises []*Promise, reducer func(acc interface{}, value interface{}), start interface{}) *Promise {

// 	return Create(func(resolve Callback, reject Callback) {

// 	})

// }

//Race ... resolves to the very first promise, rejects if none of the promises resolves
func Race(promises []*Promise) *Promise {

	return Create(func(resolve Callback, reject Callback) {

		errList := make([]interface{}, len(promises))

		resolvedStateR, resolvedStateW, _ := getSyncStateBool(false)

		var wait sync.WaitGroup
		wait.Add(len(promises))
		for i, promise := range promises {
			index := i
			if promise == nil {
				continue
			}
			promise.Then(func(value interface{}) (interface{}, error) {

				resolve(value)
				// if !resolvedRW(true) {
				// 	message <- value
				// }
				wait.Done()
				resolvedStateW(true)
				return nil, nil
			}).Catch(func(value interface{}) (interface{}, error) {
				// if !resolvedR() {
				// 	errs <- value
				wait.Done()
				errList[index] = value
				// }
				return nil, nil
			})

			// this avoids launching extra go routings using then and catch
			if resolvedStateR() {
				return
			}
		}

		wait.Wait()
		if resolvedStateR() {
			return
		}
		reject(errList)

	})
}

//All ... resolves a Promise when all promises passed are resolved,
func All(promises []*Promise) *Promise {

	return Create(func(resolve Callback, reject Callback) {

		var w sync.WaitGroup

		successStateR, successStateW, _ := getSyncStateBool(true)

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
				return nil, nil
			})
			promise.Catch(func(value interface{}) (interface{}, error) {
				data[index] = value

				successStateW(false)
				w.Done()
				return nil, nil
			})
		}
		w.Wait()
		if successStateR() {
			resolve(data)
			return
		}

		reject(data)
		return
	})

}

//Resolve ... create a Promise with resolved value
func Resolve(value interface{}) *Promise {
	promise := new(Promise)
	promise.value = value
	promise.resolved = true
	promise.failed = false
	return promise
}

//Reject ... create a rejected promise
func Reject(value interface{}) *Promise {
	promise := new(Promise)
	promise.value = value
	promise.resolved = true
	promise.failed = true
	return promise
}

//Create ... creates a promise object
func Create(action func(resolve Callback, reject Callback)) *Promise {
	promise := new(Promise)
	promise.value = nil
	promise.resolved = false
	promise.valueWaitLock.Add(1)

	go action(func(data interface{}) (interface{}, error) {
		promise.valueLock.Lock()
		if promise.resolved == false {

			promise.value = data
			promise.resolved = true
			promise.failed = false
			promise.valueWaitLock.Done()
		}
		promise.valueLock.Unlock()

		return promise.value, nil

	}, func(err interface{}) (interface{}, error) {
		promise.valueLock.Lock()
		if promise.resolved == false {
			promise.value = err
			promise.resolved = true
			promise.failed = true
			promise.valueWaitLock.Done()
		}
		promise.valueLock.Unlock()
		v, _ := promise.value.(string)
		return promise.value, fmt.Errorf(v)
	})

	return promise
}
