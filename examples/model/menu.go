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
    // XMLName xml.Name `xml:"menu"`
    ID     int
    Path   string `xml:"path,attr"`
    Show   bool   `xml:"show,attr"`
    Method string `xml:"method,attr"` // http method
    Meta
    Children []Menu `xml:"menu"`
}

// 菜单显示信息
type Meta struct {
    Name string `xml:"name,attr"`
    Icon string `xml:"icon,attr"`
}

// xmlRoot 菜单根节点
type xmlRoot struct {
    XMLName xml.Name `xml:"xml"`
    Menus   []Menu   `xml:"menu"`
}

func (m Menu) LoadDefault() ([]Menu, error) {
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
