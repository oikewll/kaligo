package model

// 表单字段
type Field struct {
    Label string    `json:"label" comment:"标签"`
    Field string    `json:"field"`
    Type  string    `json:"type"`
    Rules string    `json:"rules"`
}

// From 表单
type Form struct {
    Name   string   `json:"name"`
    Path   string   `json:"path"`
    Method string   `json:"method"`
    Csrf   string   `json:"csrf"` 
    
    Fields []Field
}

// 表单按钮
type TableButton struct {
    Name string
    Path string
}

// 表格操作
type TableOperate struct {
    CreateButton TableButton
    UpdateButton TableButton
    DeleteButton TableButton
}

// Table 表格
type Table struct {
    Name   string   `json:"name"`
    Path   string   `json:"path"`
    Method string   `json:"method"`
    Csrf   string   `json:"csrf"` 
    
    Fields []Field
}
