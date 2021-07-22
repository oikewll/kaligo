package learning

import(
    "testing"
)

// make 有很多作用，可以用来创建slice、channel、map
// map 和 slice 一样是一种自增长度的存储
func TestMapInit(t *testing.T) {
    //m := map[string]int{"one": 1, "two": 2, "three": 3}
    // 先声明，再赋值
    //m1 := map[string]int{}
    //m1["one"] = 1
    // 对比一下slice用make方法生成，第一个参数是长度，第二个参数是容量
    //s2 := make([]int, 3, 5)
    // map 用 make 生成时只有一个参数，就是容量，为什么没有长度？
    // 因为slice len所指的这些单元格都会默认初始化为0，但是map不行
    // 下面这种比较少用，但是也很有用，可以用来提高性能
    // 因为map和slice一样在自增长的时候会分配新的内存空间，同时把数据进行拷贝，如果我们能在初始化的时候就把它初始化为我们需要的大小，就可以避免
    //m2 := make(map[string]int, 10)

    m1 := map[int]int{1: 1, 2: 4, 3: 9}
    // 访问其中一个元素
    t.Log(m1[2])
    t.Logf("len m1=%d", len(m1))
    m2 := map[int]int{}
    m2[4] = 16
    t.Logf("len m2=%d", len(m2))
    m3 := make(map[int]int, 10)
    t.Logf("len m3=%d", len(m3))
    // 下面会报错，因为map没有cap
    //t.Logf("len m3=%d %d", len(m3), cap(m3))
}

// 访问一个不存在的key
func TestMapAccessNotExistingKey(t *testing.T) {
    m1 := map[int]int{}
    // 在访问的Key不存在时，仍会返回零值，不过通过返回nil来判断元素是否存在
    t.Log(m1[1])
    // 设置一个key对应的值为0，输出也是0
    m1[2] = 0
    t.Log(m1[2])
    // 那么怎么判断到底是不存在还是0值？
    m1[3] = 0
    if v, ok := m1[3]; ok {
        t.Logf("key 3's value is %d", v)
    } else {
        t.Log("key 3 is not existing")
    }
}

// map 遍历，和 array、slice 一样，通过foreach，也就是for range来遍历
func TestMapTravel(t *testing.T) {
    m1 := map[int]int{1: 1, 2: 4, 3: 9}
    for k, v := range m1 {
        t.Log(k, v)
    }
}

