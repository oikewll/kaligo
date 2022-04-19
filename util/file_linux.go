package util

import "syscall"

func init() {
    umask = syscall.Umask
}
