package kaligo

import (
    "fmt"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestTimer(t *testing.T) {
    timer := NewTimer(nil)
    assert.NotNil(t, timer)
}

func TestSchedule(t *testing.T) {
    schedule := NewSchedule("")
    for {
        select {
        case c := <-schedule.C:
            fmt.Println(c, time.Now())
            time.Sleep(time.Second * 2)
        }
    }
}
