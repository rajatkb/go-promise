package test

import (
	"fmt"
	"sync"
	"testing"
	"time"

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

func TestPromiseAllResolve(t *testing.T) {

	var d bool = true
	var w sync.WaitGroup
	w.Add(1)

	promises := make([]*Promise.Promise, 10)
	for i := 0; i < 10; i++ {
		promises[i] = Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
			time.Sleep(time.Duration(i*100) * time.Millisecond)
			resolve(2)
		})
	}

	Promise.All(promises).Then(func(value interface{}) (interface{}, error) {
		v, _ := value.([]interface{})
		array := make([]int, 10)
		for i, _v := range v {
			array[i] = _v.(int)
		}

		if len(array) != 10 {
			d = false
			w.Done()
			return nil, nil
		}
		for _, r := range array {
			if r != 2 {
				d = false
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

func TestPromiseAll(t *testing.T) {
	var d bool = true
	var w sync.WaitGroup
	w.Add(1)

	promises := make([]*Promise.Promise, 20)
	for i := 0; i < 10; i++ {
		promises[i] = Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
			time.Sleep(time.Duration(i*100) * time.Millisecond)
			resolve(2)
		})
	}

	for i := 10; i < 20; i++ {
		promises[i] = Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
			time.Sleep(time.Duration(i*100) * time.Millisecond)
			reject(3)
		})
	}

	Promise.All(promises).Then(func(value interface{}) (interface{}, error) {
		w.Done()
		return nil, nil
	}).Catch(func(value interface{}) (interface{}, error) {

		v, _ := value.([]interface{})
		array := make([]int, 20)
		for i, _v := range v {
			array[i] = _v.(int)
		}

		if len(array) != 20 {
			d = false
			w.Done()
			return nil, nil
		}
		for i, r := range array {
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

func TestChainAndFinally(t *testing.T) {
	var data int = 0
	Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
		resolve(2)
	}).Then(func(value interface{}) (interface{}, error) {
		tmp, _ := value.(int)
		return tmp + 1, nil
	}).Then(func(value interface{}) (interface{}, error) {
		tmp, _ := value.(int)
		return tmp + 1, nil
	}).Finally(func(value interface{}) {
		tmp, _ := value.(int)
		data = tmp
	})

	if data != 4 {
		t.Errorf("expected data = %d found data = %d , Finally did not work", 4, data)
	}
}

func TestResolveAndFinally(t *testing.T) {
	var data int = 0
	Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
		resolve(2)
	}).Then(func(value interface{}) (interface{}, error) {
		tmp, _ := value.(int)
		return Promise.Resolve(tmp + 1), nil
	}).Then(func(value interface{}) (interface{}, error) {
		tmp, _ := value.(int)
		return Promise.Resolve(tmp + 1), nil
	}).Finally(func(value interface{}) {
		tmp, _ := value.(int)
		data = tmp
	})

	if data != 4 {
		t.Errorf("expected data = %d found data = %d , Resolve passing did not work", 4, data)
	}
}

func TestJustFinally(t *testing.T) {

	dt, _ := Promise.Resolve(3).Finally(func(value interface{}) {}).(int)
	if dt != 3 {
		t.Errorf("expected data = %d found data = %d , Synchronous wait on Finally failed", 3, dt)
	}
}

func TestResolve(t *testing.T) {
	var data int = 0
	Promise.Resolve(4).Finally(func(value interface{}) {
		tmp, _ := value.(int)
		data = tmp
	})

	if data != 4 {
		t.Errorf("expected data = %d found data = %d , Resolve passing did not work", 4, data)
	}
}

func TestReject(t *testing.T) {
	var data int = 0
	Promise.Reject(4).Finally(func(value interface{}) {
		tmp, _ := value.(int)
		data = tmp
	})

	if data != 4 {
		t.Errorf("expected data = %d found data = %d , Resolve passing did not work", 4, data)
	}
}

func TestRaceSomePass(t *testing.T) {
	var d bool = false
	var w sync.WaitGroup
	w.Add(1)

	promises := make([]*Promise.Promise, 3)
	for i := 0; i < len(promises); i++ {
		index := i
		promises[i] = Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
			time.Sleep(time.Duration(index)*time.Second + time.Duration(10))
			resolve(index + 1)
		})
	}

	Promise.Race(promises).Then(func(value interface{}) (interface{}, error) {
		v, _ := value.(int)
		if v == 1 {
			d = true
		}
		w.Done()
		return nil, nil
	})

	w.Wait()
	if !d {
		t.Errorf("Promise.Race failed ")
	}
}

func TestRaceSomePassWithReject(t *testing.T) {
	var d bool = false
	var w sync.WaitGroup
	w.Add(1)

	promises := make([]*Promise.Promise, 20)
	for i := 0; i < 10; i++ {
		index := i
		promises[i] = Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
			time.Sleep(time.Duration(index)*time.Second + time.Duration(10))
			resolve(index + 1)
		})
	}

	for i := 10; i < 20; i++ {
		index := i
		promises[i] = Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
			time.Sleep(time.Duration(index)*time.Second + time.Duration(10))
			reject(index + 1)
		})
	}

	Promise.Race(promises).Then(func(value interface{}) (interface{}, error) {
		v, _ := value.(int)
		if v == 1 {
			d = true
		}
		w.Done()
		return nil, nil
	})

	w.Wait()
	if !d {
		t.Errorf("Promise.Race failed ")
	}
}

func TestRaceWithReject(t *testing.T) {
	var d bool = false
	var w sync.WaitGroup
	w.Add(1)

	promises := make([]*Promise.Promise, 10)

	for i := 0; i < 10; i++ {
		index := i
		promises[i] = Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
			time.Sleep(time.Duration(index)*time.Second + time.Duration(10))
			reject(index + 1)
		})
	}

	Promise.Race(promises).Catch(func(value interface{}) (interface{}, error) {
		array, _ := value.([]interface{})
		first := array[0].(int)

		if len(array) == 10 && first == 1 {
			d = true
		}
		w.Done()
		return nil, nil
	})

	w.Wait()
	if !d {
		t.Errorf("Promise.Race failed ")
	}
}

func TestRaceWithFinally(t *testing.T) {
	promises := make([]*Promise.Promise, 3)
	for i := 0; i < len(promises); i++ {
		index := i
		promises[i] = Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
			time.Sleep(time.Duration(index)*time.Second + time.Duration(10))
			resolve(index + 1)
		})
	}

	val, _ := Promise.Race(promises).Finally(nil).(int)

	if val != 1 {
		t.Errorf("Race with finally failed")
	}

}

func TestIsPending(t *testing.T) {
	p := Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
		time.Sleep(time.Duration(3) * time.Second)
		resolve(3)
	})

	p2 := Promise.Resolve(3)

	if !p.IsPending() || p2.IsPending() {
		t.Errorf("IsPending not working")
	}
}

func TestIsFullFilled(t *testing.T) {
	p := Promise.Resolve(3)
	if !p.IsFulfilled() {
		t.Errorf("IsFullFilled not working")
	}
}

func TestTimeout(t *testing.T) {
	s := false
	p := Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
		time.Sleep(time.Duration(1) * time.Second)
		resolve(false)
	}).Timeout(500).Catch(func(value interface{}) (interface{}, error) {
		return true, nil
	}).Finally(nil)
	s = p.(bool)
	if !s {
		t.Errorf("Failed to TimeOUt promise")
	}

}

func TestAsyncGenerator(t *testing.T) {
	promises := make([]*Promise.Promise, 3)
	for i := 0; i < len(promises); i++ {
		index := i
		promises[i] = Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
			time.Sleep(time.Duration(index)*time.Second + time.Duration(10))
			resolve(index + 1)
		})
	}
	count := 0
	i := 1
	for value := range Promise.AsyncGenerator(promises) {
		v, _ := value.(int)
		if i == v {
			count++
			i++
		}
	}
	if count != 3 {
		t.Errorf("Failed async generator")
	}
}
