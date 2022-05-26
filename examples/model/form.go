package model

// 表单字段
type Field struct {
    Label string    `json:"label"`  // 标签
    Field string    `json:"field"`  // 字段
    Type  string    `json:"type"`   // 类型: text、password、number、editor、textarea、file、image
    Rules string    `json:"rules"`  // 验证: required、numeric、integer、decimal、url、email、date
    Tips  string    `json:"tips"`   // 格式:bm://open:com.xxx.xxx
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
    CreateButton     TableButton    // 添加
    UpdateButton     TableButton    // 修改
    DeleteButton     TableButton    // 删除
    RefreshButton    TableButton    // 刷新
    EnableButton     TableButton    // 启用
    DisableButton    TableButton    // 禁用
    ResetMFAButton   TableButton    // 重置MFA
    TerminateAButton TableButton    // 终止session
}

// Table 表格
type Table struct {
    Name   string   `json:"name"`
    Path   string   `json:"path"`
    Method string   `json:"method"`
    Csrf   string   `json:"csrf"` 
    
    Fields []Field
}
