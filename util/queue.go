package util

type JobData any

type QueueExecutor[T JobData] interface {
    Exec(*QueueContext[T]) bool
}

type QueueFuncExecutor[T JobData] func(*QueueContext[T]) bool

func (f QueueFuncExecutor[T]) Exec(ctx *QueueContext[T]) bool {
    return f(ctx)
}

type QueueContext[T JobData] struct {
    Data  T
    Retry int
}

// Queue 队列，消费订阅模式，支持并发，重试
type Queue[T JobData] struct {
    channel  chan struct{}         // 并发任务队列
    jobs     chan *QueueContext[T] // 所有任务队列
    Retry    int                   // 任务失败重试次数，0 不重试
    Executor QueueExecutor[T]      // 任务执行
}

// NewQueue 初始化 concurrency 并发数，retry 重试次数（任务执行次数 = 1 + retry）
func NewQueue[T JobData](concurrency int, retry int, executor QueueExecutor[T]) *Queue[T] {
    return &Queue[T]{
        Retry:    retry,
        Executor: executor,
        jobs:     make(chan *QueueContext[T], 1000000),
        channel:  make(chan struct{}, concurrency),
    }
}

// Add 添加任务
func (w *Queue[T]) Add(data ...T) {
    for _, v := range data {
        w.add(&QueueContext[T]{Data: v})
    }
}

func (w *Queue[T]) add(job *QueueContext[T]) {
    w.jobs <- job
}

// Run 执行队列，exit 没有任务时退出
func (w *Queue[T]) Run(exit bool) {
    // TODO: exit 逻辑
    for v := range w.jobs {
        w.exec(v)
    }
}

func (w *Queue[T]) exec(j *QueueContext[T]) {
    w.channel <- struct{}{}
    go func(j *QueueContext[T]) {
        ret := w.Executor.Exec(j)
        if !ret {
            if j.Retry < w.Retry {
                j.Retry++
                w.add(j)
            }
        }
        <-w.channel
    }(j)
}
