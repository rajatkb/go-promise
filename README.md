# go-promise 
Go Promise aims to be a Promise/Future alternative of bluebird in golang. It will give all the various functionalities provided by any standard promise library along with something more.

### Want to try ? 🧐

```
$ go get github.com/rajatkb/go-promise
```
 
### How to use ? 🤨

* Create a single promise
```golang

import "github.com/rajatkb/go-promise"


// Create a single 
promise := Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
		resolve(2)
    })
    
```

* Then & Catch in for same promise 🔥
```golang
w.Add(1)

Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
    reject(3)
}).Then(func(value interface{}) (interface{}, error) {
    return nil, nil
}).Catch(func(value interface{}) (interface{}, error) {
    tmp, _ := value.(int)
    fmt.Println("got value:",tmp)
    w.Done()
    return nil, nil
})

w.Wait()
```




* Chaing your Then and Catch
```golang

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

w.Wait()
```

* Chain , but want to wait at the end ? 🕔️
```golang
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

fmt.Println(data)
```


* I want a promise right now ? 👊
```golang

var data int = 0
Promise.Resolve(4).Finally(func(value interface{}) {
    tmp, _ := value.(int)
    data = tmp
})

fmt.Println(data)

Promise.Reject(5).Finally(func(value interface{}) {
    tmp, _ := value.(int)
    data = tmp
})

fmt.Println(data)
```



Ya Go is async already with GoRoutines but how do I get a Synchronous wait ? 👐
```golang

dt, _ := Promise.Resolve(3).Finally(func(value interface{}) {}).(int)
fmt.Println(dt)
//
// as a pattern you can orchestrate many jobs together by 
// combining promises only to collect one final result at one s
// ingle place using Finally 
//
// You can pass nil to Finally also , instead of a function

```




* Want bunch of Promises executed at once 👀 Note: 
  All is same as Map i.e position aware results
```golang

var w sync.WaitGroup
	w.Add(1)

	promises := make([]*Promise.Promise, 20)
	for i := 0; i < 10; i++ {
		promises[i] = Promise.Create(
            func(resolve Promise.Callback, reject Promise.Callback) {
		    	resolve(2)
		})
	}

	for i := 10; i < 20; i++ {
		promises[i] = Promise.Create(
            func(resolve Promise.Callback, reject Promise.Callback) {
			    reject(3)
		})
	}

	Promise.All(promises).Then(
        func(value interface{}) (interface{}, error) {
		    w.Done()
		    return nil, nil
	}).Catch(func(value interface{}) (interface{}, error) {
        fmt.Println(value)
		w.Done()
		return nil, nil
	})

	w.Wait()

```

* What about Go Routines that race to finish and you just want the winner ? 😎
```golang


promises := make([]*Promise.Promise, 3)
	for i := 0; i < len(promises); i++ {
		index := i
		promises[i] = Promise.Create(
            func(resolve Promise.Callback, reject Promise.Callback) {
			    time.Sleep(time.Duration(index)*time.Second + time.Duration(10))
			    resolve(index + 1)
		})
	}

val, _ := Promise.Race(promises).Finally(nil).(int)

// Cactch statements in case of collection operation returns a list of error for all promises that provided error

fmt.Println("Winner is :",val)

```



### Development

#### Test
```
$ go test -v ./test
```