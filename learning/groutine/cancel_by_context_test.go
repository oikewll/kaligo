package learning

import(
    "testing"
    "time"
    "context"
)

func isChanneled(ctx context.Context) bool{
    select {
    case <- ctx.Done():
        return true
    default:
        return false
    }
}

// 测试取消任务
func TestCancel(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    for i := 0; i < 5; i++ {
        go func(i int, ctx context.Context) {
            for {
                if isChanneled(ctx) {
                    break;
                }
                // 模拟任务执行时间
                time.Sleep(time.Millisecond * 5)
            }
        } (i, ctx)
    }
    cancel()
    time.Sleep(time.Second * 1)
}

// 根 Context：通过 context.Background() 创建
// 子 context：context。WithCancel(parentContext) 创建
//   ctx, cancel := context.WithCancel(context.Backgroud())
// 当前 Context 被取消时，基于它的子 context 都会被取消
// 接受取消通知 <-ctx.Done()
