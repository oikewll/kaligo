package util

import (
    "fmt"
    "testing"
    "time"
)

func TestQueue(t *testing.T) {
    q := NewQueue[int](3, 1, QueueFuncExecutor[int](func(i *QueueContext[int]) bool {
        time.Sleep(time.Millisecond * 200)
        fmt.Println(i.Data)
        return i.Data%2 == 0
    }))
    add := func(start int) {
        for i := start; i < start+10; i++ {
            // time.Sleep(time.Millisecond * 100)
            q.Add(i)
        }
    }
    go add(0)
    go add(10)
    q.Run(false)
}
