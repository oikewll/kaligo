// Package learning ...
// Golang 切片Slice原理详解和使用技巧
// https://zhuanlan.zhihu.com/p/88741532
package learning

import(
    "testing"
    "fmt"
    "time"
    "sync"
)

// 第一个协程
func TestGroutine(t *testing.T) {
    for i := 0; i < 10; i++ {
        // go 的方法传递的时候，都是值传递，这里的 i 是被复制了一份，每个协程里面所拥有的变量地址是不一样的，所以没有竞争关系
        go func(i int) {
            fmt.Printf("%v\n", i)
        }(i)
        // i 在这些协程里面被共享了，共享变量存在竞争条件，只能用锁的机制来完成
        //go func() {
            //fmt.Printf("%v\n", i)
        //}()
    }

    time.Sleep(time.Millisecond * 50)
}

// 协程计数，因为 i 是共享变量，所以存在多个协程同时对他 ++ 的情况下，他只 ++ 了一次，所以结果不准
func TestCounter(t *testing.T) {
    counter := 0
    for i := 0; i < 5000; i++ {
        go func() {
            counter++
        } ()
    }
    // 加sleep是为了避免协程还没执行完程序就Exit()掉了，结果就不准了
    time.Sleep(1 * time.Second)
    t.Logf("counter = %d", counter)
}

// 协程计数，每个协程要执行 ++ 时都需要等待其他协程完成，就不存在多个协程同时对 i 进行 ++ 的操作，计数就准了
func TestCounterThreadSafe(t *testing.T) {
    var mut sync.Mutex    
    counter := 0
    for i := 0; i < 5000; i++ {
        go func() {
            defer func() {
                mut.Unlock()
            }()
            mut.Lock()
            counter++
        } ()
    }
    // 加sleep是为了避免协程还没执行完程序就Exit()掉了，结果就不准了
    time.Sleep(1 * time.Second)
    t.Logf("counter = %d", counter)
}

// 协程计数，每个协程要执行 ++ 时都需要等待其他协程完成，就不存在多个协程同时对 i 进行 ++ 的操作，计数就准了
func TestCountetWaitGroup(t *testing.T) {
    var mut sync.Mutex    
    var wg sync.WaitGroup    
    counter := 0
    for i := 0; i < 5000; i++ {
        wg.Add(1)
        go func() {
            defer func() {
                mut.Unlock()
            }()
            mut.Lock()
            counter++
            wg.Done()
        } ()
    }
    // 加sleep是为了避免协程还没执行完程序就Exit()掉了，结果就不准了
    //time.Sleep(1 * time.Second)
    // 用 wg.Wait() 来替代 time.Sleep()，当所有任务完成以后，就退出
    // 因为你不知道所有协程啥时候会执行完，所以sleep就无法定时间了
    wg.Wait()
    t.Logf("counter = %d", counter)
}

// 发送者和接受者一对一关系，发送者 for i:= 0; i < 10; i++ {}，那么接受者也是循环10次
//   1、如果是多个接受者，它自己都不知道能处理多少数据
// 发送者发了N个任务，接受者并不知道N是多少个
//  发送者可以在最后发一个token，比如-1，接受者收到-1就说明任务没了，就把自己退出去
//  但是如果有多个接受者，发送者就不知道应该发多少个-1好了
// 发送者发了10个任务，但是他不知道有多少个接受者，所以他需要通过close channel来广播给所有接受者，告诉他们没有任务了
// close channel 是一种广播机制，所有接受者都会收到，就可以不用发起多少任务，要给所有接受者发 -1，因为发起者也不知道接受者有多少个


// select{} 多路可以用于规定任务超时时间
// close channel 可以用于广播给给所有runner来取消任务
