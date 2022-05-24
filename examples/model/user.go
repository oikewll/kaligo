package model

import (
    "errors"
    "strings"

    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/logs"
    "golang.org/x/crypto/bcrypt"
    "github.com/go-playground/validator/v10"
)

const (
    DefaultKey  = "UserDefaultKey"
    errorFormat = "[sessions] ERROR! %s\n"
)

type Purview struct {
    url    string
    method string
}

type UserOptions struct {
    OS string
    Version string
    Utma string
    DeviceType string
    Device string
    OSVersion string
}

type User struct {
    Base

    // ID int                  `db:"id" json:"id"`
    UID string              `db:"uid" json:"uid"`
    Groups []int            `db:"ugroupsid" json:"groups"`
    Username string         `db:"username" json:"username" validate:"required"`
    Password string         `db:"validate" validate:"required"`
    Realname string         `db:"realname" json:"realname"`
    Avatar string           `db:"avatar" json:"avatar"`
    Email string            `db:"email" json:"email" validate:"required|email"`
    SessionID string        `db:"session_id"`
    SessionExpire string    `db:"session_expire"`
    Status int              `db:"status"`
    IsFirstLogin bool       `db:"is_first_login"`

    ctx *kaligo.Context
    Purviews []Purview
}

func DefaultUser(c *kaligo.Context) *User {
    logs.Debug("%v", c.UID)

    u := &User{ctx: c}
    u.GetUser(c.UID, "uid", true)

    return u
}

type Accounts map[string]string

func (m *User) table() string { return "user" }

func (m *User) CheckPurviews(url, method string) bool {
    return true
}

// 获取用户私有权限(非组权限)
// func (m *User) GetPurviews() string {
// }

func (m *User) CheckUser(accounts Accounts) (err error) {
    var username, password string    
    var ok bool

    if username, ok = accounts["username"]; !ok {
        return errors.New("请输入会员名密码")
    }

    if password, ok = accounts["password"]; !ok {
        return errors.New("请输入会员名密码")
    }

    user := &User{
        Username: username,
        Password: password,
    }
    validate := validator.New()
    // 注册自定义函数
    // _ = validate.RegisterValidation("CustomValidationErrors", CustomValidationErrors)
    err = validate.Struct(user)
    if err != nil {
        return err
        // 打印每条错误信息
        // for _, err := range err.(validator.ValidationErrors) {
            // fmt.Println(err)//Key: 'Users.Passwd' Error:Field validation for 'Passwd' failed on the 'min' tag
            // return
        // }
    }

    if err = m.GetUser(username, "username", false); err != nil {
        return err
    }

    // 同一IP使用某帐号连续错误次数检测
    if m.getLoginError24(username) {
        return errors.New("连续登录失败超过3次，暂时禁止登录！")
    }

    if !m.PasswordVerify(password, m.Password) {
        return errors.New("用户名或密码无效！")
    }

    // 用户被禁用
    if m.Status <= 0 {
        return errors.New("用户禁用！")
    }

    // 缓存用户信息 GetUser() 方法已经做了
    // m.ctx.Set(DefaultKey, m)
    return
}

// 验证登录成功后对用户进行授权
// 登录接口才会到这里来
//
// @param array $user   用户信息 
// @param int $remember 是否自动登陆 
// @param int $seclogin 是否私密登陆
// func (m *User) AuthUser(username string) bool {
// }

func (m *User) getLoginError24(username string) bool {
    return true
}

func (m *User) PasswordHash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func (m *User) PasswordVerify(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// 获得用户信息, 如果缓存中存在直接返回, 否则查数据库并且缓存起来
func (m *User) GetUser(account string, ftype string, useCache bool) (err error) {
    var uid string = account
    // 非uid, 通过 account 查询 uid 先
    if ftype != "uid" {
        _, err = DB.Select("uid").From(m.table()).Where(ftype, "=", account).Scan(&uid).Execute()
    }

    var exists bool
    if useCache {
        var u any
        // 获取缓存中的用户信息
        u, exists = m.ctx.Get(DefaultKey)
        m = u.(*User)
    }

    if !exists {
        _, err = DB.Select("*").From(m.table()).Where("uid", "=", uid).Scan(m).Execute()
        // 缓存用户信息
        m.ctx.Set(DefaultKey, m)
    }

    // 处理头像
    if m.Avatar != "" {
        m.Avatar = m.getUserAvatar(m.Avatar)
    }

    return
}

func (m *User) getUserAvatar(avatarURL string) string {
    return avatarURL
}

// 获取随机头像
func (m *User) getRandomAvatar(uid string, width, height int, avatarDIR string) string {
    // $filepath = config::instance('upload')->get('filepath');
    // $avatar_path = $filepath . '/' . $avatar_dir;
    // util::path_exists($avatar_path);

    //获取随机头像
    // avatar  = avatarDIR + '/' + md5(uid) + '.jpg'
    // $imgdata = file_get_contents("https://picsum.photos/{$width}/{$height}?random=1&?blur");
    // file_put_contents($filepath. '/' . $avatar, $imgdata);

    return ""
}

// Detail 获取单条数据详情
func (m *User) Detail(id string) (u User, err error) {
    _, err = DB.Select("*").From(m.table()).Where("id", "=", id).Scan(&u).Execute()
    return
}

// Create 添加一条数据
func (m *User) Create(u User) (ID, error) {
    q, err := DB.Insert(m.table(), []string{"title", "date"}).Values([]any{u.Username, u.Password}).Execute()
    return ID(q.LastInsertId), err
}

// Update 更新单条或多条数据
func (m *User) Update(u User) (ID, error) {
    password, err := u.PasswordHash(u.Password)

    q, err := DB.Update(m.table()).Set(map[string]string{
        "username": u.Username, 
        "password": password,
    }).Where("id", "=", u.Id).Execute()
    return ID(q.LastInsertId), err
}

// Delete 删除单条或多条数据
func (m *User) Delete(ids string) (bool, error) {
    _, err := DB.Delete(m.table()).Where("id", "in", strings.Split(ids, ",")).Execute()
    return err != nil, err
}
