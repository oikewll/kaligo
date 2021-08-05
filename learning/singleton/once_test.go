package once_test

import (
    "fmt"
    "sync"
    "testing"
    "unsafe"
)

type Singleton struct {

}

var singleInstance *Singleton
var once sync.Once    

func GetSingleObj() *Singleton {
    once.Do(func(){
        fmt.Println("Create Obj")
        singleInstance = new(Singleton)
    })
    return singleInstance
}

func TestGetSingletonObj(t *testing.T) {
    var wg sync.WaitGroup    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            obj := GetSingleObj()
            fmt.Printf("Object  Address: %x\n", unsafe.Pointer(obj))
            fmt.Printf("Pointer Address: %p\n", obj)
            wg.Done()
        } ()
    }
    wg.Wait()
}

