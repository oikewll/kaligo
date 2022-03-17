// Package learning ...
// Golang 切片Slice原理详解和使用技巧
// https://zhuanlan.zhihu.com/p/88741532
package learning

import (
    "fmt"
    "io"
    "regexp"
    "runtime"
    "strings"
    "testing"
    //"sync"
)

//const LENGTH int = 10
const (
    LENGTH int64 = 10
    WIDTH        = 20
)

var arrDetails [][]string

type Singleton struct {
    Name      string
    instances map[string]*Singleton
}

func GetInstance(name string) *Singleton {
    instance := &Singleton{
        Name: name,
    }
    if instance.instances == nil {
        instance.instances = make(map[string]*Singleton)
    }
    //return instance
    instance.instances[name] = instance
    return instance.instances[name]
}

func (g *Singleton) GetName() string {
    return g.Name
}

//var instance *singleton
//var once sync.Once
//func GetInstance() *singleton {
//once.Do(func() {
//instance = new(singleton)
//})
//return instance
//}

// CatchError is...
func catchError(err *error) {
    if pv := recover(); pv != nil {
        switch e := pv.(type) {
        case runtime.Error:
            panic(pv)
        case error:
            if e == io.EOF {
                *err = io.ErrUnexpectedEOF
            } else {
                *err = e
            }
        default:
            panic(pv)
        }
    }
}

//func TryFuncError() (err error) {
//defer catchError(&err)

//panic("hello world")
//}

//func TestTryFuncError(t *testing.T) {
//TryFuncError()
//}

func TryFuncArgsInterface(args any) bool {

    index := strings.Index("chicken", "cjj")
    fmt.Printf("TryFuncArgsInterface Index: %T=[%v]\n", index, index)

    //sqlStr := "hypertext, language, programming"
    sqlStr := "Select * From user where id = 10"
    a := regexp.MustCompile(`[\s]+`).Split(strings.TrimLeft(sqlStr[0:11], "("), 2)
    fmt.Printf("TryFuncArgsInterface sqlStr: %T=[%v]\n", a[0], a[0])
    parts := regexp.MustCompile(`\.`).Split("tablecolumn", 2)
    fmt.Printf("TryFuncArgsInterface sqlStr: %T=[%v]\n", parts, parts)

    //sqlStr := "Select * From user where id = 10"
    //sqlStr = sqlStr[0:11]
    //sqlStr = strings.TrimLeft(sqlStr[0:11], "(")
    //a := regexp.MustCompile(`[\s]+`)
    //sqlArr := a.Split(sqlStr, -1)
    //fmt.Printf("TryFuncArgsInterface sqlStr: %T=[%v]\n", sqlStr, sqlStr)
    //fmt.Printf("TryFuncArgsInterface sqlArr: %T=[%v]\n", sqlArr, sqlArr)

    g := GetInstance("111")
    name := g.GetName()
    fmt.Printf("TryFuncArgsInterface Name: %T=%v\n", name, name)

    g2 := GetInstance("222")
    name2 := g2.GetName()
    fmt.Printf("TryFuncArgsInterface Name: %T=%v\n", name2, name2)

    g3 := GetInstance("333")
    name3 := g3.GetName()

    fmt.Printf("TryFuncArgsInterface Name: %T=%v\n", name3, name3)

    switch vals := args.(type) {
    case []string:
        fmt.Println("[]string Type", vals)
        arrDetails = append(arrDetails, vals)
    case [][]string:
        fmt.Println("[][]string Type", vals)
        for _, v := range vals {
            arrDetails = append(arrDetails, v)
        }
    default:
        fmt.Println("Unknow Type")
    }

    fmt.Printf("TryFuncArgsInterface: %T=%v\n", arrDetails, arrDetails)
    arrDetails = nil
    fmt.Printf("TryFuncArgsInterface: %T=%v\n", arrDetails, arrDetails)
    return true
}

func TryFuncGetArgs(args ...[]string) bool {
    for k, v := range args {
        fmt.Printf("%d=%v\n", k, v)
    }
    return true
}

// 多维array、slice
func TestSliceMulti(t *testing.T) {

    value := "1, 2"
    valueArr := strings.Split(value, ",")
    min := strings.Trim(valueArr[0], " ")
    max := strings.Trim(valueArr[1], " ")
    fmt.Printf("TestSliceMulti min: %T=%v\n", min, min)
    fmt.Printf("TestSliceMulti max: %T=%v\n", max, max)

    var i int
    fmt.Printf("TestSliceMulti: %T=%v\n", i, i)
    //if WIDTH == 20 {
    //fmt.Printf("yes... const Test %T = %v\n", LENGTH, LENGTH)
    //} else {
    //fmt.Printf("no... const Test %T = %v\n", LENGTH, LENGTH)
    //}
    //fmt.Printf("const Test %T = %v\n", LENGTH, LENGTH)
    //fmt.Printf("const Test %T = %v\n", WIDTH, WIDTH)

    //var s1 []string
    //s1 = append(s1, "aaa", "bbb")
    //fmt.Printf("%T = %v\n", s1, s1)

    //var arrDetails [][3]string
    //a := [3]string{"111", "222", "333"}
    //arrDetails = append(arrDetails, a)

    //var s []string
    //s = []string{"111", "222", "333"}
    //arrDetails = append(arrDetails, s)
    //s = []string{"444", "555", "666"}
    //arrDetails = append(arrDetails, s)
    //s = []string{"777", "888", "999"}
    //arrDetails = append(arrDetails, s)

    TryFuncArgsInterface(arrDetails)
    //TryFuncArgsInterface(s)
    //fmt.Printf("%T = %v\n", arrDetails, arrDetails)

    //TryFuncGetArgs(s1, s2, s3)
}

// 数组无法Join，slice才可以，所以要用Join，要先把array转化为slice
func TestSliceStringJoin(t *testing.T) {
    // 不写死长度的定义，实际上还是3，因为后面就3个元素，相当于 [3]string{}
    arr := [...]string{"a", "b", "c"}
    // 数组无法append
    //arr = append(arr, "d")
    // 下面这样会报错，因为数组无法Join，slice才可以
    //t.Log(strings.Join(arr, ","))
    // arr[:] 相当于把 array 拷贝成 slice
    t.Log(strings.Join(arr[:], ","))
}

// 数组需要初始化长度，并且长度不可变
//func TestArray(t *testing.T) {
//var a [3]int // 声明并初始化为3个默认 0
//t.Logf("%v", a)
//// 把数组第一个元素改成 1
//a[0] = 1
//// 数组是无法使用append的，下面会报错
////a = append(a, 3)
//t.Logf("%v", a)

//// 声明一个长度为3的数组，并且赋值为 1, 2, 3
//b := [3]int{1, 2, 3}
//// 声明一个二维数组，长度为3，并且赋值第一维：1,2；第二维：3,4
//c := [2][3]int{{1, 2}, {3, 4}}
//t.Log(len(b), cap(b))
//t.Log(len(c), cap(c))
//}

// 切片不需要初始化给出长度，不断的append进去，它自己可以变长
//func TestSliceInit(t *testing.T) {
//// 这样写死长度为3的叫数组
////var s0 [3]int
//// 这样没有写死长度的叫切片
//var s0 []int
//t.Log(len(s0), cap(s0))
//// 为什么append以后要赋值给原来的变量，因为它要变长，内存并没有变
//s0 = append(s0, 1)
//t.Log(len(s0), cap(s0))

//// 初始化一个有3个元素的切片，这时它的长度会变成3，容量也是3
//s1 := []int{1, 2, 3}
//t.Log(len(s1), cap(s1))

//// 生成长度为3，容量为5的slice，注意array是无法用make的
//s2 := make([]int, 3, 5)
//t.Log(len(s2), cap(s2))
//t.Log(s2[0], s2[1], s2[2])
//s2 = append(s2, 1)
//t.Log(s2[0], s2[1], s2[2], s2[3])
//t.Log(len(s2), cap(s2))

//// Go对二维slice不太友好，只能手动去为第二维度分配空间，手动无奈：
//n := 5
//var s = make([][]int, n)
//for i := 0; i < n; i++ {
//s[i] = make([]int, n)
//}
////dy := 5
////dx := 8
////a := make([][]int, dy)
////for i := range a {
////a[i] = make([]int, dx)
////}
//}

// 切片 可变长度
//func TestSliceGrowing(t *testing.T) {
//s := []int{}
//for i := 0; i < 10; i++ {
//s = append(s, i)
//t.Log(len(s), cap(s))
//}
//}

// 共享内存
//func TestSliceShareMemory(t *testing.T) {
//year := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep",
//"Oct", "Nov", "Dec"}
//Q2 := year[3:6]
//t.Log(Q2, len(Q2), cap(Q2))
//summer := year[5:8]
//t.Log(summer, len(summer), cap(summer))
//// 修改了summer的值，Q2、year 也收到影响，说明他们用的是一块内存
//summer[0] = "Unknow"
//t.Log(Q2)
//t.Log(year)
//}

// 数组支持比较，那么切片是否支持比较？
func TestSliceComparing(t *testing.T) {
    //a := []int{1, 2, 3, 4}
    //b := []int{1, 2, 3, 4}
    //// 下面会报错，slice只能和nil比较
    //if a == b {
    //t.Log("equal")
    //}
}

// 切片容量可伸缩，每次超出长度增加为原来的2倍，比如初始化为2，当超出2个时变成4，超出4时变成8
func TestSlice(t *testing.T) {
    //s1 := []string{}
    //s1 = append(s1, "s1")
    //s1 = append(s1, "s2")
    //s1 = append(s1, "s3")
    //s1 = append(s1, "s4")
    //s1 = append(s1, "s5")
    //s1 = append(s1, "s6")
    //s1 = append(s1, "s7")
    //s1 = append(s1, "s8")
    //t.Log(s1)
    //t.Log(len(s1), cap(s1))

    //s2 := s1[2:4]
    //t.Log(s2)
    //t.Log(len(s2), cap(s2))
    //s2 = append(s2, "s9")

    //t.Log(s1)
    //t.Log(len(s1), cap(s1))
    //t.Log(s2)
    //t.Log(len(s2), cap(s2))
}
