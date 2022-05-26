package controller

import (
    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/sessions"
    "examples/model"
)

type Auth struct {
    kaligo.Controller
}

// @Summary Login 账户登陆
// @Tags    Auth
// @Param   username  formData  string  true   "账号"       default(test)
// @Param   password  formData  string  true   "密码"       default(test)
// @Param   remember  formData  boolean true   "记住密码"   default(true)
// @Success 200 {object} map[string]string
// @Router  /auth/login [POST]
func (c *Auth) Login() {
    username := c.FormValue("username")
    password := c.FormValue("password")
    remember := c.FormValue("remember")

    user := &model.User{}
    // // CheckUser 里面调用的 GetUser 会把用户数据缓存到 context.Set("UserDefaultKey-UID", user)
    // if err := user.CheckUser(model.Accounts{ "username": username, "password": password, "remember": remember }); err != nil {
    //     c.JSON(http.StatusBadRequest, kaligo.H{
    //         "code": -1,
    //         "msg":  err.Error(),
    //     })
    // }

    user.UID = "abcdefg"
    session := sessions.Default(c.Context)
    // 要检查是否需要加UID为后缀
    session.Set("UID", user.UID)
    session.Save()

    c.JSON(200, map[string]string{
        "username": username,
        "password": password,
        "remember": remember,
    })

    // c.String(200, "Login successful")

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
// @Tags    Auth
// @Success 200 {object} map[string]string
// @Router  /auth/logout [DELETE]
func (c *Auth) Logout() {
    session := sessions.Default(c.Context)
    session.Delete("UID")
    session.Save()

    c.String(200, "Logout successful")
}

// @Summary Token 产生一个 CSRF Token
// @Tags    Auth
// @Success 200 {object} string
// @Router  /auth/token [GET]
func (c *Auth) Token() {
    csrf := model.DefaultAuth(c.Context).MakeCsrfToken()
    result(c.Context, map[string]string{"csrf": csrf}, nil)
}

// @Summary CheckToken 检查 CSRF Token 是否存在
// @Tags    Auth
// @Param   csrf  formData  string  true   "CSRF Token"
// @Success 200 {object} string
// @Router  /auth/check_token [POST]
func (c *Auth) CheckToken() {
    csrf := c.FormValue("csrf")
    if model.DefaultAuth(c.Context).CheckCsrfToken(csrf) {
        c.String(200, "验证通过")
    } else {
        c.String(404, "请勿重复提交")
    }
}
