package util

import (
)


func Colorize(text string, status string) string {
    out := ""
    switch status {
        case "succ":
            out = "\033[32;1m"    // Blue
        case "fail":
            out = "\033[31;1m"    // Red
        case "warn":
            out = "\033[33;1m"    // Yellow
        case "note":
            out = "\033[34;1m"    // Green
        default:
            out = "\033[0m"    // Default
    }
	return out+text+"\033[0m"
}

