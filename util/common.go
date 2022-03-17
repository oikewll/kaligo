package util

import (
    "io"
    "runtime"
)

var tab8s = "        "

// CheckErr is...
func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}

// CatchError is...
func CatchError(err *error) {
    if pv := recover(); pv != nil {
        switch e := pv.(type) {
        case runtime.Error:
            panic(pv)
        case error:
            if e == io.EOF {
                *err = io.ErrUnexpectedEOF
            } else {
                *err = e
            }
        default:
            panic(pv)
        }
    }
}
