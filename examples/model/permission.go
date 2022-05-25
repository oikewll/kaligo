package model

import (
    "strings"

    "github.com/owner888/kaligo/util"
)

type Permission struct {
    Method string
    Path   string
}

// GET-/api/todo,GET-/api/todo/:id,POST-/api/todo,PUT-/api/todo/:id,DELETE-/api/todo/:id
const (
    permissionsSeperator      = ","
    permissionStructSeperator = "-"
)

// Parse 解析 Permission 数组，不符合规范的数据直接被移除
func (m Permission) Parse(permission string) []Permission {
    p, _ := util.CompactMapSlice(strings.Split(permission, permissionsSeperator), func(p string) (Permission, bool) {
        data := strings.Split(p, permissionStructSeperator)
        if len(data) == 2 {
            return Permission{Method: data[0], Path: data[1]}, true
        }
        return Permission{}, false
    })
    return p
}
