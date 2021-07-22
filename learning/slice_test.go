// Package learning ...
// Golang 切片Slice原理详解和使用技巧
// https://zhuanlan.zhihu.com/p/88741532
package learning

import(
    "testing"
    "strings"
)

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
func TestArray(t *testing.T) {
    var a [3]int // 声明并初始化为3个默认 0
    t.Logf("%v", a)
    // 把数组第一个元素改成 1
    a[0] = 1
    // 数组是无法使用append的，下面会报错
    //a = append(a, 3)
    t.Logf("%v", a)

    // 声明一个长度为3的数组，并且赋值为 1, 2, 3
    b := [3]int{1, 2, 3}
    // 声明一个二维数组，长度为3，并且赋值第一维：1,2；第二维：3,4
    c := [2][3]int{{1, 2}, {3, 4}}
    t.Log(len(b), cap(b))
    t.Log(len(c), cap(c))
}

// 切片不需要初始化给出长度，不断的append进去，它自己可以变长
func TestSliceInit(t *testing.T) {
    // 这样写死长度为3的叫数组
    //var s0 [3]int
    // 这样没有写死长度的叫切片
    var s0 []int
    t.Log(len(s0), cap(s0))
    // 为什么append以后要赋值给原来的变量，因为它要变长，内存并没有变
    s0 = append(s0, 1)
    t.Log(len(s0), cap(s0))

    // 初始化一个有3个元素的切片，这时它的长度会变成3，容量也是3
    s1 := []int{1, 2, 3}
    t.Log(len(s1), cap(s1))

    // 生成长度为3，容量为5的slice，注意array是无法用make的
    s2 := make([]int, 3, 5)
    t.Log(len(s2), cap(s2))
    t.Log(s2[0], s2[1], s2[2])
    s2 = append(s2, 1)
    t.Log(s2[0], s2[1], s2[2], s2[3])
    t.Log(len(s2), cap(s2))

    // Go对二维slice不太友好，只能手动去为第二维度分配空间，手动无奈：
    n := 5
    var s = make([][]int, n)
    for i := 0; i < n; i++ { 
        s[i] = make([]int, n) 
    }
    //dy := 5
    //dx := 8
    //a := make([][]int, dy) 
    //for i := range a { 
        //a[i] = make([]int, dx) 
    //}
}

// 切片 可变长度
func TestSliceGrowing(t *testing.T) {
    s := []int{}
    for i := 0; i < 10; i++ {
        s = append(s, i)
        t.Log(len(s), cap(s))
    }
}

// 共享内存
func TestSliceShareMemory(t *testing.T) {
    year := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep",
    "Oct", "Nov", "Dec"}
    Q2 := year[3:6]
    t.Log(Q2, len(Q2), cap(Q2))
    summer := year[5:8]
    t.Log(summer, len(summer), cap(summer))
    // 修改了summer的值，Q2、year 也收到影响，说明他们用的是一块内存
    summer[0] = "Unknow"
    t.Log(Q2)
    t.Log(year)
}

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

    //values  [][]string // 二维切片（变长数组）
    var slice [][]string

    slice = [][]string{{"ddd"}, {}}
    slice[0] = append(slice[0], "hi - 000")
    slice[1] = append(slice[1], "hi - 111")
    //slice[2] = append(slice[2], "hi - 111")

    t.Logf("%v", slice)

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


