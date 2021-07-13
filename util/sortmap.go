package util

import (
)

// SortMap is ...
type SortMap []Item
                 
// Item is ...
type Item struct {
    Key string
    Val string
}

// NewSortMap is ...
func NewSortMap(m map[string]string) SortMap {
    ms := make(SortMap, 0, len(m))
                 
    for k, v := range m {
        ms = append(ms, Item{k, v})
    }
    return ms
}

func (ms SortMap) Len() int {
    return len(ms)
}

func (ms SortMap) Less(i, j int) bool {
    //return ms[i].Val < ms[j].Val // 按值排序
    return ms[i].Key < ms[j].Key // 按键排序,入数据库，一定要用这个，也就是根据字段排序了
}

func (ms SortMap) Swap(i, j int) {
    ms[i], ms[j] = ms[j], ms[i]
}

