package timer

import (
    "fmt"
    "reflect"
    "sync"
    "time"

    "github.com/astaxie/beego/logs"
)

// Timer 是所有定时任务管理中心
type Timer struct {
    storeTimers sync.Map
}

// New 新建 Timer，不要重复创建
func New() *Timer {
    return &Timer{}
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
// AddTasker("default", &control.Task{}, "import_database", "2014-10-15 15:33:00")
// func AddTasker(name string, control any, action string, taskTime string) {
func (t *Timer) AddTasker(name, taskTime, method string, runner any, params any) {
    go func() {
        then, _ := time.ParseInLocation("2006-01-02 15:04:05", taskTime, time.Local)
        dura := then.Sub(time.Now())
        //fmt.Println(dura)
        if dura > 0 {
            timeTasker := time.AfterFunc(dura, func() {
                run(runner, method, params)
            })
            t.storeTimers.Store(name, timeTasker)
        } else {
            logs.Error("定时任务 --- [ " + name + " ] --- 小于当前时间，将不会被执行")
        }
    }()
}

// AddTimer is the function for add timer, The interval is in microseconds
// router.AddTimer("import_database", 3000, "ImportDatabase", &controller.Get{})
func (t *Timer) AddTimer(name string, duration time.Duration, method string, runner any, params any) {
    go func() {
        timeTicker := time.NewTicker(duration * time.Millisecond)
        t.storeTimers.Store(name, timeTicker)
        for {
            select {
            case <-timeTicker.C:
                run(runner, method, params)
            }
        }
    }()
}

func run(runner any, method string, params any) error {
    r := reflect.New(reflect.Indirect(reflect.ValueOf(runner)).Type())
    m := r.MethodByName(method)
    if !m.IsValid() {
        return fmt.Errorf("Controller method \"%s\" not exist", m)
    }
    m.Call([]reflect.Value{reflect.ValueOf(params)})
    return nil
}
