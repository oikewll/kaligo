package kaligo

import (
    "fmt"
    "reflect"
    "testing"
)

type user struct {
    Controller
}

func (c *user) login() {
    fmt.Printf("Controller %v\n", c)
}

type user2 struct {
    *Controller
}

func (c user2) login() {
    fmt.Printf("Controller %v\n", c)
}

func login(c *Controller) {
    fmt.Printf("Controller %v\n", c)
}

func logout(c *user) {
    fmt.Printf("Controller %v\n", c)
}

func logout2(c user2) {
    fmt.Printf("Controller %v\n", c)
}

func addRoute[T Interface](path string, handler func(c T)) {
    var c T
    t := reflect.TypeOf(c)
    kind := t.Kind()
    if kind == reflect.Pointer {
        value := reflect.New(t.Elem())
        c = value.Interface().(T)
    } else {
        value := reflect.New(t)
        if kind == reflect.Struct {
            f := reflect.Indirect(value).FieldByName("Controller")
            if f.IsValid() && f.CanSet() {
                f.Set(reflect.ValueOf(&Controller{&Context{}, t.String()}))
            }
        }
        c = *value.Interface().(*T)
    }
    c.Init(&Context{}, t.String())
    c.Prepare()
    handler(c)
    c.Finish()
}

func TestType(t *testing.T) {
    addRoute("", (*user).login)
    addRoute("", user2.login)
    addRoute("", login)
    addRoute("", logout)
    addRoute("", logout2)
}
