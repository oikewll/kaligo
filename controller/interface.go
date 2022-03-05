package controller

import (
    "github.com/owner888/kaligo/contex"
)

// Interface is use for 
type Interface interface {
    Init(contex *contex.Context, childName string)
    Prepare()
    Finish()
}
