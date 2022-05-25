package middlewares

import (
    // "errors"
    // "sync"
    // "golang.org/x/crypto/bcrypt"

    // "github.com/owner888/kaligo"
    // "github.com/owner888/kaligo/logs"
    // "github.com/owner888/kaligo/sessions"
)

// AuthUserKey is the cookie name for user credential in basic auth.
const AuthUserKey = "user"

// func Default(c *kaligo.Context) {
//     if req := c.Request; req != nil {
//         logs.Debug(req.Method, " ", req.URL)
//     }
//
//     session := sessions.Default(c)
//     c.UID = session.Get(AuthUserKey).(string)
//
//     c.Next()
// }

// // 后台管理员表
// const (
//     TableAdmin             = "#PB#_admin"
//     TableAdminGroup        = "#PB#_admin_group"
//     TableAdminLoginLog     = "#PB#_admin_login_log"
//     TableAdminOperationLog = "#PB#_admin_operation_log"
//     TableAdminPurview      = "#PB#_admin_purview"
//     TableAdminSession      = "#PB#_admin_session"
// )
//
// // 后台管理员字段
// var (
//     TableAdminField = []any{
//         "id",
//         "groups",
//         "account",
//         "password",
//         "realname",
//         "avatar",
//         "email",
//         "session",
//         "session_expired",
//         "status",
//     }
//
//     TableAdminSessionField = []string{
//         "id",
//         "token",
//         "utma",
//         "ip",
//     }
//
//     // Author = &Auth{}
// )

// type UserParam struct {
//     key   string
//     value string
// }
//
// type AuthUser struct {
//     ID uint64 `db:"id"`
//     // Groups           []uint64    `db:"groups"`
//     Account        string `db:"account"`
//     Password       string `db:"password"`
//     Realname       string `db:"realname"`
//     Avatar         string `db:"avatar"`
//     Email          string `db:"email"`
//     Session        string `db:"session"`
//     SessionExpired uint64 `db:"session_expired"`
//     Status         int32  `db:"status"`
// }

// type Auth struct {
//     User       *AuthUser
//     ctx        *Context
//     cacheStore *sync.Map
// }
//
// func BasicAuthForRealm(accounts Accounts, realm string) HandlerFunc {
//     if realm == "" {
//         realm = "Authorization Required"
//     }
//     // Basic realm=Authorization Required
//     realm = "Basic realm=" + strconv.Quote(realm)
//     pairs := processAccounts(accounts)
//     return func(c *Context) {
//         // Search user in the slice of allowed credentials
//         user, found := pairs.searchCredential(c.requestHeader("Authorization"))
//         if !found {
//             // Credentials doesn't match, we return 401 and abort handlers chain.
//             c.Header("WWW-Authenticate", realm)
//             c.AbortWithStatus(http.StatusUnauthorized)
//             return
//         }
//
//         session := sessions.Default(c)
//         session.Set("hello", "world")
//         session.Save()
//
//         // The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
//         // c.MustGet(AuthUserKey).
//         c.Set(AuthUserKey, user)
//     }
// }

// func init() {
//     // 只初始化一次
//     if Author == nil {
//         NewAuth()
//     }
// }
//
// func NewAuth() *Auth {
//     Author := &Auth{}
//     // cache, err := cache.New()
//     // if err != nil {
//     //     panic(err)
//     // }
//     // mux.Cache = cache
//     // mux.Timer = NewTimer()
//     // mux.pool.New = func() any {
//     //     return &Context{DB: mux.DB, Cache: mux.Cache, Timer: mux.Timer}
//     // }
//     return Author
// }

// @description    检查账号
// @author         seatle            时间（2022/03/30   08:02 ）
// @param          account           string           "账号"
// @param          password          string           "密码"
// @param          remember          string           "是否记住登陆状态"
// @return         err               error            "错误信息"
// func (a *Auth) CheckUser(account, password string, remember ...bool) (err error) {
//     if account == "" || password == "" {
//         return errors.New("请输入会员名密码！")
//     }
//
//     // if err = a.checkAccount(account); err != nil {
//     //     return
//     // }
//     // if err = a.checkLoginError24(account); err != nil {
//     //     return
//     // }
//
//     var user AuthUser
//     _, err = a.ctx.DB.Select(TableAdminField...).From("user").Where("id", "=", "1").Scan(&user).Execute()
//     if err != nil {
//         return errors.New("账号不存在！")
//     }
//
//     if err = a.checkUserStatus(user.Status); err != nil {
//         return
//     }
//
//     // 存到全局去
//     Author.User = &user
//     return
// }
//
// func (a *Auth) SaveUserSession(account, password string, remember ...bool) (err error) {
//     return
// }

// func (a *Auth) ListUserSession(account, password string, remember ...bool) (sess []*Session) {
//     return
// }

// func (a *Auth) PasswordHash(password string) (string, error) {
//     bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//     return string(bytes), err
// }
//
// func (a *Auth) PasswordVerify(password, hash string) bool {
//     err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//     return err == nil
// }
//
// // 检测账号格式合法性
// func (a *Auth) checkAccount(account string) (err error) {
//     return errors.New("账号格式不合法！")
// }
//
// // 同 IP 同帐号 24 小时内连续错误次数检测(默认：3次)
// func (a *Auth) checkLoginError24(account string) (err error) {
//     return errors.New("连续登录失败超过3次，暂时禁止登录！")
// }
//
// // 检测用户状态
// func (a *Auth) checkUserStatus(status int32) (err error) {
//     if status <= 0 {
//         return errors.New("用户禁用！")
//     }
//     return
// }
