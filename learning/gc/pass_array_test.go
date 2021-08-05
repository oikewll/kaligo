package gc
// go test -bench=.
// GODEBUG=gctrace=1 go test -bench=BenchmarkPassingArrayWithValue 
// GODEBUG=gctrace=1 go test -bench=BenchmarkPassingArrayWithReference 
// go test -bench=BenchmarkPassingArrayWithReference -trace=trace_val.out
// go tool trace trace_val.out

import(
    //"fmt"
    //"time"
    "testing"
)

const NumOfElems = 1000

type Content struct {
    Detail [10000]int
}

// 传递值
func withValue(arr [NumOfElems]Content) int{
    //fmt.Println(&arr[2])
    return 0
}

// 传递引用
func withReference(arr *[NumOfElems]Content) int{
    //b := *arr    
    //fmt.Println(&arr[2])
    return 0
}

func TestFn(t *testing.T) {
    var arr [NumOfElems]Content    
    //fmt.Println(&arr[2])
    withValue(arr)
    withReference(&arr)
}

func BenchmarkPassingArrayWithValue(b *testing.B) {
    var arr [NumOfElems]Content    

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        withValue(arr)
    }
    b.StopTimer()
}

func BenchmarkPassingArrayWithReference(b *testing.B) {
    var arr [NumOfElems]Content    

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        withReference(&arr)
    }
    b.StopTimer()
}

//// 第一个协程
//func TestGroutine(t *testing.T) {
    //for i := 0; i < 10; i++ {
        //// go 的方法传递的时候，都是值传递，这里的 i 是被复制了一份，每个协程里面所拥有的变量地址是不一样的，所以没有竞争关系
        //go func(i int) {
            //fmt.Printf("%v\n", i)
        //}(i)
        //// i 在这些协程里面被共享了，共享变量存在竞争条件，只能用锁的机制来完成
        ////go func() {
            ////fmt.Printf("%v\n", i)
        ////}()
    //}

    //time.Sleep(time.Millisecond * 50)
//}

