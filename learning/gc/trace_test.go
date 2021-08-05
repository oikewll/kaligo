package gc
// 测试程序输出trace信息
// go test -trace trace.out
// 可视化 trace 信息，会自动打开浏览器
// go tool trace trace.out
// 注意：
//   1、初始化至合适的大小
//     自动扩容是有代价的
//   2、复用内存
// https://github.com/easierway/service_decorators

import(
    "os"
    "runtime/trace"
    "testing"
)

func TestTrace(t *testing.T) {
    f, err := os.Create("trace.out")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    err = trace.Start(f)
    if err != nil {
        panic(err)
    }

    defer trace.Stop()
    // Your program here
}

