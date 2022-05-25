package controller

import (
    "net/http"
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/sessions"
    "examples/model"
)

type User struct {
    kaligo.Controller
}

// @Summary 用户信息
// @Tags    User
// @Success 200 {object} map[string]string
// @Router  /user [GET]
func (c *User) Detail() {
    session := sessions.Default(c.Context)

    logs.Debug(c)

    c.JSON(200, kaligo.H{"hello": session.Get("hello")})
}

// @Summary 账户登陆
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
    if err := user.CheckUser(model.Accounts{ "username": username, "password": password, "validate": validate, "remember": remember }); err != nil {
        panic(err)
    }

    session := sessions.Default(c.Context)
    session.Set("user", user)

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

// @Summary 账户退出
// @Tags    User
// @Success 200 {object} map[string]string
// @Router  /user/logout [DELETE]
func (c *User) Logout() {

}
