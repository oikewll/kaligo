package model

import (
    "encoding/xml"
    "io"
    "io/ioutil"
    "os"

    "github.com/owner888/kaligo/util"
)

const defaultMenuPath = "menu.xml"

// 菜单
type Menu struct {
    ID     int    `json:"id"`
    Permit bool   `json:"permit"`
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

// LoadDefault 加载默认配置，filters 过滤权限，removeDeny 移除无权限的结点
func (m Menu) LoadDefault(filters string, removeDeny bool) ([]Menu, error) {
    file, err := os.Open(defaultMenuPath)
    if err != nil {
        return []Menu{}, err
    }
    defer file.Close()
    menu, err := m.Load(file)
    root := &Menu{Children: menu}
    root = root.Permission(Permission{}.Parse(filters), removeDeny)
    if root != nil {
        return root.Children, err
    }
    return []Menu{}, err
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
    // 给 Menu 带上自增 ID
    id := 0
    root := &Menu{Children: v.Menus}
    util.ForTree(root, (*Menu).getChildren, (*Menu).setChildren, func(node *Menu) {
        node.ID = id
        id++
    })
    return root.Children, nil
}

func (m *Menu) getChildren() []Menu  { return m.Children }
func (m *Menu) setChildren(c []Menu) { m.Children = c }

// Permission 通过权限列表过滤菜单
func (m *Menu) Permission(p []Permission, removeDeny bool) *Menu {
    if removeDeny {
        util.ForTreeChild(m, (*Menu).getChildren, (*Menu).setChildren, func(node *Menu) {
            node.updatePermission(p)
        })
        return m
    } else {
        return util.FilterTree(m, (*Menu).getChildren, (*Menu).setChildren, func(node *Menu) bool {
            return node.updatePermission(p)
        })
    }
}

func (m *Menu) updatePermission(permissions []Permission) bool {
    if len(m.Children) > 0 {
        m.Permit = util.ReduceSlice(m.Children, true, func(x bool, n Menu) bool { return x && n.Permit })
    } else {
        for _, p := range permissions {
            if (m.Method == p.Method || p.Method == PermissionAll) &&
                (m.Path == p.Path || p.Path == PermissionAll) {
                m.Permit = true
                break
            }
        }
    }
    return m.Permit
}
