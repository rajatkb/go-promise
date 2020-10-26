package promise

import (
	"fmt"
	"sync"
	"time"
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

//IsFulfilled ... checks whether promiseis fullfilled or not
func (obj *Promise) IsFulfilled() bool {
	obj.valueLock.Lock()
	resolved := obj.resolved
	failed := obj.failed
	obj.valueLock.Unlock()
	return resolved && !failed
}

//IsRejected ... checks whether promise is rejected or not
func (obj *Promise) IsRejected() bool {
	obj.valueLock.Lock()
	resolved := obj.resolved
	failed := obj.failed
	obj.valueLock.Unlock()
	return resolved && failed
}

//IsPending ... checks whether promise is pending or not
func (obj *Promise) IsPending() bool {
	obj.valueLock.Lock()
	resolved := obj.resolved
	obj.valueLock.Unlock()
	return !resolved
}

//Cancel ... cancels a promise
func (obj *Promise) Cancel() bool {
	wasCancelled := false
	obj.valueLock.Lock()
	if !obj.resolved {
		obj.resolved = true
		obj.failed = true
		obj.value = fmt.Errorf("Error: Promise Canecelled")
		obj.valueWaitLock.Done()
		wasCancelled = true
	}
	obj.valueLock.Unlock()
	return wasCancelled
}

//Timeout ... cancels a promise
func (obj *Promise) Timeout(ms int) *Promise {
	go func() {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		obj.Cancel()
	}()
	return obj
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
	value := obj.value
	vp, ok := value.(*Promise)
	if ok {
		vp.valueWaitLock.Wait()
		value = vp.value
	}
	if callback != nil {
		callback(value)
	}

	return value
}

/**
* TO-DO : Props
* Promise.props(struct {
*  field1 : Promise.Resolve(1)
*  field2
* })
*
**/

//// WE GO FUNCTIONALL BRRRRRRRRRRRRRR................. ///////////////

//Map ... for the bluebird affiniadoes , maps your array of promise with a new common then
func Map(promises []*Promise, cb Callback) []*Promise {
	promisesT := make([]*Promise, len(promises))
	for i, promise := range promises {
		promisesT[i] = promise.Then(cb)
	}
	return promisesT
}

//Reduce ... asynchronous reducer , does not waits for all promise to be resolve , it launches reduce callback as soon as first result is available
// It will be used to process both errors and values. So reduces should account for that.
func Reduce(promises []*Promise, reducer func(index int, acc interface{}, value interface{}) interface{}, acc interface{}) *Promise {
	return Create(func(resolve Callback, reject Callback) {
		promise := Resolve(acc)
		count := 0
		for asyncValue := range AsyncGenerator(promises) {
			index := count
			promise = promise.Then(func(acc interface{}) (interface{}, error) {
				return reducer(index, acc, asyncValue), nil
			})
			count++
		}
		resolve(promise)
	})
}

//// FUNCTIONAL DONE /////////////////////

//AsyncGenerator ... returns the results & errors of promises , without any ordering to the caller
func AsyncGenerator(promises []*Promise) <-chan interface{} {
	messages := make(chan interface{}, len(promises))
	var w sync.WaitGroup
	w.Add(len(promises))
	go func() {
		for _, promise := range promises {
			promise.Then(func(value interface{}) (interface{}, error) {
				messages <- value
				w.Done()
				return nil, nil
			}).Catch(func(err interface{}) (interface{}, error) {
				messages <- err
				w.Done()
				return nil, nil
			})
		}
	}()

	go func() {
		w.Wait()
		close(messages)
	}()
	return messages
}

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
