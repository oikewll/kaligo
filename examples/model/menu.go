package model

import (
    "encoding/xml"
    "io"
    "io/ioutil"
    "os"
)

const defaultMenuPath = "menu.xml"

// 菜单
type Menu struct {
    ID     int    `json:"id"`
    Path   string `json:"path" xml:"path,attr"`
    Show   bool   `json:"show" xml:"show,attr"`
    Top    bool   `json:"top"  xml:"top,attr"`
    Reload bool   `json:"reload" xml:"reload,attr"`
    Method string `json:"method" xml:"method,attr"` // http method
    Meta
    Children []Menu `json:"children" xml:"menu"`
}

// 菜单显示信息
type Meta struct {
    Name string `json:"name" xml:"name,attr"`
    Icon string `json:"icon" xml:"icon,attr"`
}

// xmlRoot 菜单根节点
type xmlRoot struct {
    XMLName xml.Name `xml:"xml"`
    Menus   []Menu   `xml:"menu"`
}

func (m Menu) LoadDefault(filters ...string) ([]Menu, error) {
    file, err := os.Open(defaultMenuPath)
    if err != nil {
        return []Menu{}, err
    }
    defer file.Close()
    return m.Load(file)
}

func (m Menu) Load(reader io.Reader) ([]Menu, error) {
    data, err := ioutil.ReadAll(reader)
    if err != nil {
        return nil, err
    }
    v := xmlRoot{}
    err = xml.Unmarshal(data, &v)
    if err != nil {
        return []Menu{}, err
    }
    return v.Menus, nil
}

func (m Menu) getChildren() []Menu   { return m.Children }
func (m *Menu) setChildren(c []Menu) { m.Children = c }

// Permission 通过权限列表过滤菜单
func (m *Menu) Permission(p []Permission) *Menu {
    return filterTree(m, Menu.getChildren, (*Menu).setChildren, func(node Menu) bool {
        if len(node.Children) > 0 {
            return true
        }
        for _, v := range p {
            if v.Method == node.Method && v.Path == node.Path {
                return true
            }
        }
        return false
    })
}

// filterTree 树形结构过滤，深度优先，先过滤叶子结点
func filterTree[T any](root *T, getter func(node T) []T, setter func(node *T, children []T), filter func(node T) bool) *T {
    var children []T
    for _, v := range getter(*root) {
        n := filterTree(&v, getter, setter, filter)
        if n != nil {
            children = append(children, *n)
        }
    }
    setter(root, children)
    ok := filter(*root)
    if ok {
        return root
    }
    return nil
}