package kaligo

import (
    "errors"
    "reflect"
    "regexp"
)

// Route is a Router's route
type Route struct {
    Regex          *regexp.Regexp
    Methods        map[string]string
    Params         map[int]string
    ControllerType reflect.Type
}

func (r *Route) IsMethodsValid() (bool, error) {
    var err error
    c := reflect.New(r.ControllerType)
    for _, v := range r.Methods {
        m := c.MethodByName(v)
        if !m.IsValid() {
            if err == nil {
                err = errors.New("Routing method does not exist: ")
            }
            err = errors.New(err.Error() + v + ",")
        }
    }
    return err == nil, err
}

// StaticRoute is a Router's route
type StaticRoute struct {
    Prefix    string
    StaticDir string
}
