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

// LoadDefault("GET-/api/todo,GET-/api/todo/:id,POST-/api/todo,PUT-/api/todo/:id,DELETE-/api/todo/:id")
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
