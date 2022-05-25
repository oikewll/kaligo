package model

import (
    "errors"
    "fmt"
    "strings"
    // "time"

    "github.com/owner888/kaligo"
    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/util"
    "github.com/go-playground/validator/v10"
)

const (
    UserDefaultKeyFormat    = "UserDefaultKey-%s"       // 根据UID获取用户信息
    PurviewDefaultKeyFormat = "PurviewDefaultKey-%s"    // 更加PurViewID获取用户权限信息
    errorFormat             = "[sessions] ERROR! %s\n"
    timeFormat              = "2006-01-02 15:04:05"
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

    UID string              `db:"uid" json:"uid"`
    Groups []int            `db:"groups" json:"groups"`
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
    logs.Debugf("UID === %s", c.UID)

    u := &User{ctx: c}
    // 从缓存或数据库中获取数据填充到 User struct
    u.GetUser(c.UID, "uid", true)

    return u
}

type Accounts map[string]string

func (m *User) Table() string { return "user" }

// List 分页获取数据列表
func (m User) List() ([]User, int, error) {
    var users []User
    _, err := DB.Select("id", "username", "realname", "status").From(m.Table()).Scan(&users).Execute()
    return users, 0, err
}

// 检查用户权限
// GET-/api/todo,GET-/api/todo/:id,POST-/api/todo,PUT-/api/todo/:id,DELETE-/api/todo/:id
func (m *User) CheckPurviews(url, method string) bool {
    return true
}

// 获取用户私有权限(非组权限)
// 用户权限 = 用户私有权限 + 组权限
func (m *User) GetPurviews() string {
    return "GET-/api/todo,GET-/api/todo/:id,POST-/api/todo,PUT-/api/todo/:id,DELETE-/api/todo/:id"
}

// 获取组权限
func (m *User) GetGroupsPurviews() string {
    return "GET-/api/todo,GET-/api/todo/:id,POST-/api/todo,PUT-/api/todo/:id,DELETE-/api/todo/:id"
}

// 检测用户
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

    if !util.PasswordVerify(password, m.Password) {
        return errors.New("用户名或密码无效！")
    }

    // 用户被禁用
    if m.Status <= 0 {
        return errors.New("用户禁用！")
    }

    // 缓存用户信息 GetUser() 方法已经做了
    // m.ctx.Set(fmt.Sprintf(UserDefaultKeyFormat, uid), m)
    return
}

// 用户是否24小时内连续登录失败超过3次
func (m *User) getLoginError24(username string) bool {
    // ErrorNum = 3
    // loc, _ := time.LoadLocation("Asia/Shanghai")
    // t, _ := time.ParseInLocation(timeFormat, "")
    // _, err = DB.Select("uid").From(m.table()).Where(ftype, "=", account).Scan(&uid).Execute()
    return true
}

// 获得用户信息, 如果缓存中存在直接返回, 否则查数据库并且缓存起来
func (m *User) GetUser(account string, ftype string, useCache bool) (err error) {
    // account 为空字符串直接返回, 避免下面数据库查询不到把当前对象设置为 nil 导致空指针异常
    if account == "" {
        return
    }

    var uid string = account
    // 非uid, 通过 account 查询 uid 先
    if ftype != "uid" {
        _, err = DB.Select("uid").From(m.Table()).Where(ftype, "=", account).Scan(&uid).Execute()
    }

    var exists bool
    if useCache {
        var u any
        // 获取缓存中的用户信息
        u, exists = m.ctx.Get(fmt.Sprintf(UserDefaultKeyFormat, uid))
        m = u.(*User)
    }

    if !exists {
        _, err = DB.Select("*").From(m.Table()).Where("uid", "=", uid).Scan(m).Execute()
        // 缓存用户信息
        m.ctx.Set(fmt.Sprintf(UserDefaultKeyFormat, uid), m)
    }

    // 处理头像
    if m.Avatar != "" {
        m.Avatar = m.getUserAvatar(m.Avatar)
    }

    return
}

// 获取用户头像
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
    _, err = DB.Select("*").From(m.Table()).Where("id", "=", id).Scan(&u).Execute()
    return
}

// Create 添加一条数据
func (m *User) Create(u User) (ID, error) {
    q, err := DB.Insert(m.Table(), []string{"username", "password"}).Values([]any{u.Username, u.Password}).Execute()
    return ID(q.LastInsertId), err
}

// Update 更新单条或多条数据
func (m *User) Update(u User) (ID, error) {
    password, err := util.PasswordHash(u.Password)

    q, err := DB.Update(m.Table()).Set(map[string]string{
        "username": u.Username, 
        "password": password,
    }).Where("id", "=", u.Id).Execute()
    return ID(q.LastInsertId), err
}

// Delete 删除单条或多条数据
func (m *User) Delete(ids string) (bool, error) {
    _, err := DB.Delete(m.Table()).Where("id", "in", strings.Split(ids, ",")).Execute()
    return err != nil, err
}
