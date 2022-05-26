package controller

import (
    "errors"
    "examples/model"

    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/form"
    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/sessions"
)

type User struct {
    kaligo.Controller
}

// @Summary List 分页获取用户信息
// @Tags    User
// @Param   page     query integer false "当前页数" default(1)
// @Param   size     query integer false "每页记录" default(20)
// @Success 200 {array} model.User
// @Router  /user [get]
func (c *User) List() {
    page := c.QueryInt64("page", 1)
    size := c.QueryInt64("size", 20)
    // 排序字段和排序方式: 只支持 ID、创建时间
    orderBy   := map[string]string{c.QueryValue("order_name", "id"): c.QueryValue("order_by", "desc")}
    keywords  := c.FormValue("keywords")
    status    := c.FormValue("status")
    createdAt := c.FormValue("created_at") // 2022-05-06 - 2022-06-08

    data, err := model.User{}.List(page, size, orderBy, keywords, status, createdAt)

    table := form.Table{
        Name:   "用户列表",
        Path:   "/api/user",
        Method: "GET",
        Csrf:   model.DefaultAuth(c.Context).MakeCsrfToken(),
    }
    table.SearchComponents = append(table.SearchComponents, form.Component{
        Title: "关键字",
        Placeholder: "账号/昵称",
        Field: "keywords",
        Type:  "input",
    })
    // table.SearchComponents = append(table.SearchComponents, form.Component{
    //     Title: "是否启用",
    //     Placeholder: "是否启用",
    //     Field: "status",
    //     Type:  "select",
    // })
    table.TableGlobalOperate = form.TableGlobalOperate{
        CreateButton:  form.TableButton{Name:"添加", Path:"/api/user", Method:"POST", Form:"/api/user/createform"},
        DeleteButton:  form.TableButton{Name:"删除", Path:"/api/user/:id", Method:"DELETE"},
        EnableButton:  form.TableButton{Name:"激活", Path:"/api/user/status/:id", Method:"PUT"},
        DisableButton: form.TableButton{Name:"禁用", Path:"/api/user/status/:id", Method:"DELETE"},
        RefreshButton: form.TableButton{Name:"刷新", Path:"/api/user/refresh", Method:"GET"},
    }
    table.TableListOperate = form.TableListOperate{
        UpdateButton:  form.TableButton{Name:"修改", Path:"/api/user/:id", Method:"PUT", Form:"/api/user/updateform/:id"},
    }

    data.Table = table
    result(c.Context, data, err)
}

// @Summary Detail 用户信息
// @Tags    User
// @Param   id        path      integer false  "账号ID"     default(1)
// @Router  /user/{id} [GET]
// @Success 200 {object} map[string]string
func (c *User) Detail() {
    session := sessions.Default(c.Context)
    uid := session.Get("UID")

    logs.Debug("UID === ", uid)

    c.JSON(200, kaligo.H{"UID": uid})
}

// @Summary Update 更新单条或多条数据
// @Tags    User
// @Param   id        path      integer false  "账号ID"     default(1)
// @Param   username  formData  string  true   "账号"       default(test)
// @Param   password  formData  string  true   "密码"       default(test)
// @Param   realname  formData  string  true   "姓名"       default(test)
// @Param   groups    formData  string  true   "所属组IDs"  Enums(1, 2, 3)
// @Param   email     formData  string  false  "邮箱"       default(test@gmail.com)
// @Param   status    formData  integer false  "状态"       default(1)
// @Success 200 {object} []model.User
// @Router  /user/{id} [PUT]
func (c *User) Update() {
    var user model.User
    user.Id       = c.ParamInt("id")
    user.Username = c.FormValue("username")
    user.Password = c.FormValue("password")
    user.Realname = c.FormValue("realname")
    user.Groups   = c.FormValue("groups")
    user.Email    = c.FormValue("email")

    data, err := (&model.User{}).Update(user)
    result(c.Context, data, err)
}

// @Summary Create 添加一条数据
// @Tags    User
// @Param   username  formData  string  true   "账号"       default(test)
// @Param   password  formData  string  true   "密码"       default(test)
// @Param   realname  formData  string  true   "姓名"       default(test)
// @Param   groups    formData  string  true   "所属组IDs"  Enums(1, 2, 3)
// @Param   email     formData  string  false  "邮箱"       default(test@gmail.com)
// @Param   status    formData  integer false  "状态"       default(1)
// @Success 200 {object} model.User
// @Router  /user [POST]
func (c *User) Create() {
    var user model.User
    user.Username = c.FormValue("username")
    user.Password = c.FormValue("password")
    user.Realname = c.FormValue("realname")
    user.Groups   = c.FormValue("groups")
    user.Email    = c.FormValue("email")
    data, err := (&model.User{}).Create(user)
    result(c.Context, data, err)
}

// @Summary CreateForm 用户添加表单
// @Tags    User
// @Router  /user/createform [GET]
func (c *User) CreateForm() {
    table := form.Form{
        Name:   "用户添加",
        Path:   "/api/user",
        Method: "POST",
        Csrf:   model.DefaultAuth(c.Context).MakeCsrfToken(),
    }
    table.Components = append(table.Components, form.Component{
        Title: "账号",
        Field: "username",
        Type:  "text",
        Validate: form.Validate{Required: true, Message: "请输入账号", Trigger: "change"},
    })
    table.Components = append(table.Components, form.Component{
        Title: "密码",
        Field: "password",
        Type:  "password",
        Validate: form.Validate{Required: true, Message: "请输入密码", Trigger: "change"},
    })
    table.Components = append(table.Components, form.Component{
        Title: "姓名",
        Field: "realname",
        Type:  "text",
        Validate: form.Validate{Required: true, Message: "请输入姓名", Trigger: "change"},
    })
    table.Components = append(table.Components, form.Component{
        Title: "邮箱",
        Field: "email",
        Type:  "text",
        Validate: form.Validate{Required: true, Message: "请输入邮箱", Trigger: "change", Type: "email"},
    })

    result(c.Context, table, nil)
}

// @Summary CreateForm 用户添加表单
// @Tags    User
// @Param   id       path      integer false  "账号ID"     default(1)
// @Router  /user/updateform/{id} [GET]
func (c *User) UpdateForm() {
    table := form.Form{
        Name:   "用户修改",
        Path:   "/api/user/:id",
        Method: "PUT",
        Csrf:   model.DefaultAuth(c.Context).MakeCsrfToken(),
    }
    table.Components = append(table.Components, form.Component{
        Title: "账号",
        Field: "username",
        Value: "username",
        Type:  "text",
        Validate: form.Validate{Required: true, Message: "请输入账号", Trigger: "change"},
    })
    table.Components = append(table.Components, form.Component{
        Title: "密码",
        Field: "password",
        Value: "password",
        Type:  "password",
        Validate: form.Validate{Required: true, Message: "请输入密码", Trigger: "change"},
    })
    table.Components = append(table.Components, form.Component{
        Title: "姓名",
        Field: "realname",
        Value: "realname",
        Type:  "text",
        Validate: form.Validate{Required: true, Message: "请输入姓名", Trigger: "change"},
    })
    table.Components = append(table.Components, form.Component{
        Title: "邮箱",
        Field: "email",
        Value: "email",
        Type:  "text",
        Validate: form.Validate{Required: true, Message: "请输入邮箱", Trigger: "change", Type: "email"},
    })

    result(c.Context, table, nil)
}

// @Summary Delete 删除单条或多条数据
// @Tags    User
// @Param   id       path      integer false  "账号ID"     default(1)
// @Success 200 {integer} integer
// @Router  /user/{id} [DELETE]
func (c *User) Delete() {
    id := c.ParamValue("id")
    if len(id) == 0 {
        result(c.Context, nil, errors.New("id is required"))
    }
    data, err := (&model.User{}).Delete(id)
    result(c.Context, data, err)
}
