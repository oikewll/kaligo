package controller

import "github.com/owner888/kaligo"

type User struct {
    kaligo.Controller
}

func (c *User) Login() {
    username := c.FormValue("username")
    c.SetCookie("access_token", username, 1000, "/", "", true, true)
}

func (c *User) Logout() {

}
