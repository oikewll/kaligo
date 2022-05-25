package controller

import (
    "errors"
    "net/http"
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/sessions"
    "examples/model"
)

type User struct {
    kaligo.Controller
}

// @Summary List 分页获取用户信息
// @tags    User
// @Param   page       query integer false "当前页数, 1开始"
// @Param   size       query integer false "当前页数, 默认20"
// @Success 200 {object} []model.User
// @Router  /user [get]
func (c *User) List() {
    data, _, err := model.User{}.List()
    result(c.Context, data, err)
}

// @Summary Detail 用户信息
// @Tags    User
// @Param   id       path integer true "User ID"
// @Router  /user/{id} [GET]
// @Success 200 {object} map[string]string
func (c *User) Detail() {
    session := sessions.Default(c.Context)
    uid := session.Get("UID")

    logs.Debug("UID === ", uid)

    c.JSON(200, kaligo.H{"UID": uid})
}

// @Summary Update 更新单条或多条数据
// @tags    User
// @Param   user formData model.User true "User"
// @Success 200 {object} []model.User
// @Router  /user [PUT]
func (c *User) Update() {
    var user model.User
    err := c.JsonBodyValue(&user)
    if err != nil {
        result(c.Context, nil, err)
    }

    data, err := (&model.User{}).Update(user)
    result(c.Context, data, err)
}

// @Summary Delete 删除单条或多条数据
// @tags    User
// @Param   id       query integer false "User ID"
// @Success 200 {integer} integer
// @Router  /user [DELETE]
func (c *User) Delete() {
    id := c.QueryValue("id")
    if len(id) == 0 {
        result(c.Context, nil, errors.New("id is required"))
    }
    data, err := (&model.User{}).Delete(id)
    result(c.Context, data, err)
}

// @Summary Create 添加一条数据
// @tags    User
// @Param   user formData model.User true "User"
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

// @Summary Login 账户登陆
// @Tags    User
// @Success 200 {object} map[string]string
// @Router  /user/login [POST]
func (c *User) Login() {
    var err error    

    username := c.FormValue("username")
    password := c.FormValue("password")
    validate := c.FormValue("validate")
    remember := c.FormValue("remember")

    defer func() {
        if r := recover(); r != nil {
            c.JSON(http.StatusBadRequest, kaligo.H{
                "code": -1,
                "msg":  err.Error(),
            })
        }
    }()

    user := &model.User{}
    // CheckUser 里面调用的 GetUser 会把用户数据缓存到 context.Set("UserDefaultKey-UID", user)
    if err := user.CheckUser(model.Accounts{ "username": username, "password": password, "validate": validate, "remember": remember }); err != nil {
        panic(err)
    }

    session := sessions.Default(c.Context)
    // 要检查是否需要加UID为后缀
    session.Set("UID", user.UID)

    // if user.IsFirstLogin {
    //     session.Set("uid", user.UID)
    //     session.Save()
    //     c.Redirect(302, "/reset_pwd")
    // } else if config.GetBool("mfa_code") && user.OtpAuthCode == "" {
    //     // 启动强制MFA认证 并且 用户未绑定，进行绑定流程
    //     session.Set("otp_username", user.Username)
    //     session.Set("otp_authcode", googleAuth.CreateSecret())
    //     session.Save()
    //     c.Redirect(302, "/opt_enable/authentication")
    // } else {
    //     if config.GetBool("mfa_code") {
    //         session.Set("otp_uid", user.UID)
    //         session.Set("otp_remember", user.Remember)
    //         session.Set("otp_username", user.Username)
    //         session.Set("otp_authcode", user.OtpAuthCode)
    //         session.Save()
    //     }
    //     m.SetLoginInfo()
    //     c.Redirect(302, "/opt_enable/authentication")
    // }


    // c.JSON(http.StatusOK, kaligo.H{
    //     "code": 0,
    //     "msg" : "successful",
    //     "data": data,
    // })

    // c.SetCookie("access_token", username, 1000, "/", "", true, true)
}

// @Summary Logout 账户退出
// @Tags    User
// @Success 200 {object} map[string]string
// @Router  /user/logout [DELETE]
func (c *User) Logout() {

}
