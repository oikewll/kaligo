// Package learning ...
// 空接口
package learning

import(
    "testing"
    "fmt"
)

func DoSomething(p interface{}) {
    //if i, ok := p.(int); ok {
        //fmt.Println("Interger", i)
    //}
    //if s, ok := p.(string); ok {
        //fmt.Println("String", s)
    //}
    //fmt.Println("Unknow Type")

    switch v := p.(type) {
    case int:
        fmt.Println("Interger", v)
    case string:
        fmt.Println("String", v)
    default:
        fmt.Println("Unknow Type")
    }
}

func TestEmptyInterfaceAssertion(t *testing.T) {
    DoSomething(10);
    DoSomething("10");
}

// Go 接口最佳实践
// 倾向于使用小的接口定义，很多接口只包含一个方法
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
// 较大的接口定义，可以由多个小接口定义组合而成
type ReadWriter interface {
    Reader
    Writer
}

// 只依赖于必要功能的最小接口
//func StoreData(reader Reader) error {
    //return error.Error()
//}
