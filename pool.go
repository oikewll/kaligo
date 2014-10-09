package epooll

import (
    //"fmt"
)

type InitFunction func() (interface{}, error)

type ConnPool struct {
    size int
    conn chan interface{}
}

func (this *ConnPool)InitPool(size int, initfn InitFunction) error{
    this.conn = make(chan interface{}, size)
    for x := 0; x < size; x++ {
        conn, err := initfn()
        if err != nil {
            return err
        }
        this.conn <-conn
    }
    this.size = size
    return nil
}

// 从连接池里取出连接
func (this *ConnPool)Get() interface{} {
    // 因为channel是有锁的，所以就没必要借助sync.Mutex来进行读写锁定
    // container/list就需要锁住，不然并发就互抢出问题了
    return <-this.conn
}

// 回收连接到连接池
func (this *ConnPool)Release(conn interface{}) {
    this.conn <-conn
}

