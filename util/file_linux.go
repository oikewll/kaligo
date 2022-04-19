package util

import "syscall"

func init() {
    umask = syscall.Umask
}

var _i = 0
