# go-promise
Go Promise aims to be a Promise/Future alternative of bluebird in golang. It will give all the various functionalities provided by any standard promise library along with something more.


### How to use ?
```golang

promise := Promise.Create(func(resolve Promise.Callback, reject Promise.Callback) {
		resolve(2)
    })
    

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

### Test
```
$ go test -v ./test
```