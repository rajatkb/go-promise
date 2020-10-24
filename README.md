# go-promise 
Go Promise aims to be a Promise/Future alternative of bluebird in golang. It will give all the various functionalities provided by any standard promise library along with something more.

### Want to try ? 🧐

```
$ go get github.com/rajatkb/go-promise
```
 
### How to use ? 🤨

* Create a single promise
```golang

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


* Want bunch of Promises executed at once 👀
```golang

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
        fmt.Println(v)
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



### Test
```
$ go test -v ./test
```