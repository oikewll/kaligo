package main

import "time"
import "fmt"

//type comparable interface {
    //type int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64
//}

//func max[T comparable](a []T) T {
    //m := a[0]
    //for _, v := range a {
        //if m < v {
            //m = v
        //}
    //}
    //return m
//}

func main() {
    c1 := make(chan string, 1)
    go func() {
        time.Sleep(time.Second * 2)
        c1 <- "result 1"
    }()
    select {
    case res := <-c1:
        fmt.Println(res)
        fmt.Printf("%T = %v\n", res, res)
    case <-time.After(time.Second * 3):
        fmt.Println("timeout 1")
    }
}
