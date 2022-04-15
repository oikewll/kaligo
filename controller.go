package kaligo

import (
    // "fmt"

    "fmt"
    "reflect"
)

// H is a shortcut for map[string]any
type H map[string]any

// Interface is use for
type Interface interface {
    Init(contex *Context, childName string)
    Prepare()
    Finish()
}

// Controller is a base controller struct
type Controller struct {
    *Context
    ChildName string
}

// New returns a new initialized Controller.
func New() *Controller {
    return &Controller{}
}

func runController(controllerType reflect.Type, m string, ctx *Context, params Params) (ret []reflect.Value, err error) {
    // Invoke the request handler
    vc := reflect.New(controllerType)

    // Init callback
    method := vc.MethodByName("Init")

    args := make([]reflect.Value, 2)
    args[0] = reflect.ValueOf(ctx)
    args[1] = reflect.ValueOf(controllerType.Name())
    method.Call(args)

    args = make([]reflect.Value, 0)

    // Prepare callback
    method = vc.MethodByName("Prepare")
    method.Call(args)

    // Request callback
    method = vc.MethodByName(m)
    if !method.IsValid() {
        err = fmt.Errorf("Controller Method not exist")
        return
    }
    ret = method.Call(args)

    // Finish callback
    method = vc.MethodByName("Finish")
    method.Call(args)

    return ret, err
}

// Init returns a new initialized Controller.
func (c *Controller) Init(ctx *Context, cn string) {
    c.Context = ctx
    c.ChildName = cn
}

// Prepare is use for some processing before starting to execute.
func (c *Controller) Prepare() {
}

// Finish is use for processing after execution is complete.
func (c *Controller) Finish() {
}

/* vim: set expandtab: */
