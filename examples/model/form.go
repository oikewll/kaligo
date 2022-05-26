package model

// 表单字段
type Field struct {
    Label string    `json:"label"`  // 标签
    Field string    `json:"field"`  // 字段
    Type  string    `json:"type"`   // 类型: text、password、number、editor、textarea、file、image
    Rules string    `json:"rules"`  // 验证: required、numeric..., 详细看本页最下面
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

// 表单验证类
//
// required: "必选字段",
// remote: "请修正该字段",
// email: "请输入正确格式的电子邮件",
// url: "请输入合法的网址",
// date: "请输入合法的日期",
// numeric: "请输入合法的数字",
// integer: "只能输入整数",
// decimal: "只能输入小数",
// idcard: "请输入合法的身份证号",
// creditcard: "请输入合法的信用卡号",
// matches[param]: "请再次输入相同的值",
// accept: "请输入拥有合法后缀名的字符串",
// maxlength[param]: "长度不能大于 {param} 位",
// minlength[param]: "长度不能小于 {param} 位",
// exactlength[param]: "长度只能等于 {param} 位",
// rangelength[minlen:maxlen]: "长度介于 {minlen} 和 {maxlen} 之间",
// max[param]: "请输入一个最大为 {param} 的值",
// min[param]: "请输入一个最小为 {param} 的值"
// range[minnum:maxnum]: "请输入一个介于 {minnum} 和 {maxnum} 之间的值",
