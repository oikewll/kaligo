package epooll

import (
)


type ConnPool struct {
	// Dial is an application supplied function for creating and configuring a
	// connection
	Dial func() (interface{}, error)
    //Dial func() (*autorc.Conn, error)

    // Maximum number of idle connections in the pool.
	MaxIdle int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActive int

	closed bool
	active int

    idle chan interface{}
}

// 批量生成连接，并把连接放到连接池channel里面
func (this *ConnPool)InitPool() error{
    this.idle = make(chan interface{}, this.MaxActive)
    for x := 0; x < this.MaxActive; x++ {
        conn, err := this.Dial()
        if err != nil {
            return err
        }
        this.idle <-conn
    }
    return nil
}

// 从连接池里取出连接
func (this *ConnPool)Get() interface{} {
    // 如果空闲连接为空，初始化连接池
    if this.idle == nil {
        this.InitPool()
    }

    // 赋值一下好给下面回收和返回
    conn := <-this.idle

    // 使用完把连接回收到连接池里
    defer this.Release(conn)

    // 因为channel是有锁的，所以就没必要借助sync.Mutex来进行读写锁定
    // container/list就需要锁住，不然并发就互抢出问题了
    return conn
}

// 回收连接到连接池
func (this *ConnPool)Release(conn interface{}) {
    this.idle <-conn
}

