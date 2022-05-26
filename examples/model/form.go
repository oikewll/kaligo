package model

// @Description 表单验证
type Validate struct {
    Required bool       `json:"required"`       // 是否必选
    Type     string     `json:"type"`           // 内建校验类型: date、datetimerange、emali、array
    Enum     string     `json:"enum"`           // 枚举类型
    Len      int        `json:"len"`            // 字段长度
    Min      int        `json:"min"`            // 最小长度
    Max      int        `json:"max"`            // 最大长度
    Pattern  string     `json:"pattern"`        // 正则表达式校验
    Message  string     `json:"message"`        // 校验文案
}

// @Description 表单组件属性 Property
type Props struct {
    Type string     `json:"type"`   // 类型: text、password
}

// @Description 表单组件
// http://form-create.com/v2/element-ui/components/input.html
type Component struct {
    Type  string        `json:"type"`           // 类型: input、editor、textarea、file、image
    Field string        `json:"field"`          // 字段
    Title string        `json:"title"`          // 标题
    Props Props         `json:"props"`          // 属性
    Value string        `json:"value"`          // 默认值
    Validate Validate   `json:"validate"`       // 验证: required、numeric..., 详细看本页最下面
    Tips  string        `json:"tips"`           // 格式:bm://open:com.xxx.xxx
    Placeholder string  `json:"placeholder"`    // 输入框占位文本
}

// @Description From表单
type Form struct {
    Name   string   `json:"name"`
    Path   string   `json:"path"`
    Method string   `json:"method"`
    Csrf   string   `json:"csrf"` 
    
    // 表单组件
    Components []Component
}

// @Description 表单按钮
type TableButton struct {
    Name string
    Path string
}

// @Description 表格公共操作按钮
type TableGlobalOperate struct {
    CreateButton     TableButton    // 添加
    DeleteButton     TableButton    // 删除
    EnableButton     TableButton    // 启用
    DisableButton    TableButton    // 禁用
    RefreshButton    TableButton    // 刷新
}

// @Description 表格列表每一行的操作按钮
type TableListOperate struct {
    UpdateButton     TableButton    // 修改按钮
    ResetMFAButton   TableButton    // 重置MFA
    TerminateAButton TableButton    // 终止session
}

// @Description Table 表格
type Table struct {
    Name   string   `json:"name"`
    Path   string   `json:"path"`
    Method string   `json:"method"`
    Csrf   string   `json:"csrf"` 
    
    // 查询表单组件
    SearchComponents []Component
    TableGlobalOperate TableGlobalOperate   // 公共按钮
    TableListOperate TableListOperate       // 列表最右边的按钮
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
