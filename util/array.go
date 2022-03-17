package util

import (
    "fmt"
    "reflect"
)

// Arr 数组（array）、可变长数组（slice）、map 操作类
type Arr struct {
}

// MapKeys 用于获取 map 所有key
func (a *Arr) MapKeys(m map[string]string) []string {
    // 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率较高
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}

// SliceRemoveDuplicate 函数用于在slice中去除重复的元素，其中a必须是已经排序的序列。
// params:
//   a: slice对象，如[]string, []int, []float64, ...
// return:
//   []any: 已经去除重复元素的新的slice对象
func (a *Arr) SliceRemoveDuplicate(sli any) (ret []any) {
    if reflect.TypeOf(sli).Kind() != reflect.Slice {
        fmt.Printf("<SliceRemoveDuplicate> <a> is not slice but %T\n", sli)
        return ret
    }

    va := reflect.ValueOf(sli)
    for i := 0; i < va.Len(); i++ {
        if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
            continue
        }
        ret = append(ret, va.Index(i).Interface())
    }

    return ret
}
