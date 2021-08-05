package signal

import (
    "testing"
    "time"
    "os"
    "os/signal"
)

func TestSignal(t *testing.T) {
    // 必须是 block channel，否则下面sleep会让信号丢失，
    // 如果是block channel，你按下ctrl+c的时候，信号已经存在channel里面了，这样sleep完，下面就可以收到信号直接退出
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)

    // 模拟一个 5分钟的执行代码
    time.Sleep(5 * time.Second)

    s := <-c
    t.Logf("Got signal: %p", s)
    t.Logf("Got signal: %x", s)
    t.Log("Got signal:", s)
}


