package kaligo

type AuthUser interface {
    Name() string
    CheckUser(account, password string, remember ...bool) (err error)
    SaveUserSession(account, password string, remember ...bool) (err error) 
    PasswordHash(password string) (string, error)
    PasswordVerify(password, hash string) bool
    checkAccount(account string) (err error)
    checkLoginError24(account string) (err error)
    checkUserStatus(status int32) (err error)
}

/* vim: set expandtab: */
