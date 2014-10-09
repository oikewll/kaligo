package epooll

import(
    "testing"
)

func Test_Redis(t *testing.T) {
    for i := 0; i < 1000; i++ {
        //var _ = RedisConn
        NewRedisPool().Get()
    }
}
