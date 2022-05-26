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