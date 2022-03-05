package controller

import (
    "fmt"

    "github.com/owner888/kaligo/contex"
)

// Controller is a base controller struct
type Controller struct {
    Context   *contex.Context
    ChildName string
    H         map[string]interface{}
}

// New returns a new initialized Controller.
func New() *Controller {
    fmt.Printf("init Controller\n")
    return &Controller{}
}

// Init returns a new initialized Controller.
func (c *Controller) Init(contex *contex.Context, childName string) {

    c.Context = contex
    c.ChildName = childName

    fmt.Println("\n---------")
    fmt.Println("\nhello Init")

}

// Prepare is use for some processing before starting to execute.
func (c *Controller) Prepare() {
    fmt.Println("\nhello Prepare")
}

// Finish is use for processing after execution is complete.
func (c *Controller) Finish() {
    fmt.Println("\nhello Finish")
    fmt.Println("\n---------")
}

// JSON is use for return jsonresult
func (c *Controller) JSON() {
}
