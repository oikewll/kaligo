package kaligo

import (
    "sync"
    "time"

    "github.com/owner888/kaligo/logs"
)

// Timer 是所有定时任务管理中心
type Timer struct {
    storeTimers sync.Map
    mux         *Mux
}

// New 新建 Timer，不要重复创建
func NewTimer(mux *Mux) *Timer {
    return &Timer{mux: mux}
}

// DelTasker is the function for delete tasker
func (t *Timer) DelTasker(name string) bool {
    tasker, ok := t.storeTimers.Load(name)
    if ok {
        tasker.(*time.Ticker).Stop()
        return true
    }
    return false
}

// DelTimer is the function for delete timer
func (t *Timer) DelTimer(name string) bool {
    tasker, ok := t.storeTimers.Load(name)
    if ok {
        tasker.(*time.Timer).Stop()
        return true
    }
    return false
}

// AddTasker is the function for add tasker
// c.Timer.AddTasker("import_database", "2014-10-15 15:33:00", "ImportDatabase", &controller.Get{})
func (t *Timer) AddTasker(name, taskTime, method string, runner Interface, params Params) {
    go func() {
        then, _ := time.ParseInLocation("2006-01-02 15:04:05", taskTime, time.Local)
        dura := then.Sub(time.Now())
        //fmt.Println(dura)
        if dura > 0 {
            timeTasker := time.AfterFunc(dura, func() {
                t.mux.CallController(runner, method, params)
            })
            t.storeTimers.Store(name, timeTasker)
        } else {
            logs.Error("定时任务 --- [ " + name + " ] --- 小于当前时间，将不会被执行")
        }
    }()
}

// AddTimer is the function for add timer, The interval is in microseconds
// c.Timer.AddTimer("import_database", time.Second * 3, "ImportDatabase", &controller.Get{})
func (t *Timer) AddTimer(name string, duration time.Duration, method string, runner Interface, params Params) {
    go func() {
        timeTicker := time.NewTicker(duration)
        t.storeTimers.Store(name, timeTicker)
        for {
            select {
            case <-timeTicker.C:
                t.mux.CallController(runner, method, params)
            }
        }
    }()
}

func (t *Timer) AddSchedule(name string, s string, method string, runner Interface, params Params) {
    go func() {
        timeTicker := NewSchedule(s)
        t.storeTimers.Store(name, timeTicker)
        for {
            select {
            case <-timeTicker.C:
                t.mux.CallController(runner, method, params)
            }
        }
    }()
}

// Schedule 基于时间的循环任务
type Schedule struct {
    C      <-chan time.Time
    ticker *time.Ticker
}

func NewSchedule(s string) *Schedule {
    first, period := scheduleParse(s)
    c := make(chan time.Time, 1)
    schedule := &Schedule{C: c}
    time.AfterFunc(first, func() {
        schedule.ticker = time.NewTicker(period)
        c <- time.Now()
        for {
            select {
            case t := <-schedule.ticker.C:
                c <- t
            }
        }
    })
    return schedule
}

func scheduleParse(s string) (first, period time.Duration) {
    // TODO: 解析规则
    return time.Hour, time.Hour
}
