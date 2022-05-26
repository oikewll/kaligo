package controller

import (
    "errors"
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/logs"
    // "github.com/owner888/kaligo/util"
    "github.com/owner888/kaligo/sessions"
    "examples/model"
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
    data, _, err := model.User{}.List()
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
// @Param   emali     formData  string  false  "邮箱"       default(test@gmail.com)
// @Param   status    formData  integer false  "状态"       default(1)
// @Success 200 {object} []model.User
// @Router  /user/{id} [PUT]
func (c *User) Update() {
    // id := c.ParamValue("id")
    var user model.User
    // intVar, err := strconv.Atoi(id)
    err := c.JsonBodyValue(&user)
    if err != nil {
        result(c.Context, nil, err)
    }

    data, err := (&model.User{}).Update(user)
    result(c.Context, data, err)
}

// @Summary Create 添加一条数据
// @Tags    User
// @Param   username  formData  string  true   "账号"       default(test)
// @Param   password  formData  string  true   "密码"       default(test)
// @Param   realname  formData  string  true   "姓名"       default(test)
// @Param   groups    formData  string  true   "所属组IDs"  Enums(1, 2, 3)
// @Param   emali     formData  string  false  "邮箱"       default(test@gmail.com)
// @Param   status    formData  integer false  "状态"       default(1)
// @Success 200 {object} model.User
// @Router  /user [POST]
func (c *User) Create() {
    var user model.User
    err := c.JsonBodyValue(&user)
    if err != nil {
        result(c.Context, nil, err)
    }
    data, err := (&model.User{}).Create(user)
    result(c.Context, data, err)
}

// @Summary CreateForm 用户添加表单
// @Tags    User
// @Success 200 {object} model.Form
// @Router  /user/createform [GET]
func (c *User) CreateForm() {

    csrf := model.DefaultAuth(c.Context).MakeCsrfToken()
    form := model.Form{
        Name:   "用户添加",
        Path:   "/api/user",
        Method: "POST",
        Csrf:   csrf,
    }
    form.Fields = append(form.Fields, model.Field{
        Label: "账号",
        Field: "username",
        Type:  "text",
        Rules: "required",
    })
    form.Fields = append(form.Fields, model.Field{
        Label: "密码",
        Field: "password",
        Type:  "password",
        Rules: "required",
    })
    form.Fields = append(form.Fields, model.Field{
        Label: "姓名",
        Field: "realname",
        Type:  "text",
        Rules: "required",
    })
    form.Fields = append(form.Fields, model.Field{
        Label: "邮箱",
        Field: "email",
        Type:  "text",
        Rules: "required|email",
    })

    result(c.Context, form, nil)
}

// @Summary Delete 删除单条或多条数据
// @Tags    User
// @Param   id        path      integer false  "账号ID"     default(1)
// @Success 200 {integer} integer
// @Router  /user{id} [DELETE]
func (c *User) Delete() {
    id := c.QueryValue("id")
    if len(id) == 0 {
        result(c.Context, nil, errors.New("id is required"))
    }
    data, err := (&model.User{}).Delete(id)
    result(c.Context, data, err)
}

