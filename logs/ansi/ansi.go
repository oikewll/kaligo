package ansi

import (
    "fmt"
    "strconv"
    "strings"
)

// 终端颜色值 https://zh.wikipedia.org/wiki/ANSI转义序列#颜色
type color int

const (
    black color = iota
    red
    green
    yellow
    blue
    magenta
    cyan
    white
    custom
    clear // 某些终端不支持
)

type colorAttr int

const (
    foreground colorAttr = 30
    background colorAttr = 40
)

type attr int

type attrs []attr

const (
    reset     attr = iota // 重置
    bold                  // 粗体
    dim                   // 弱化
    italic                // 斜体（未广泛支持）
    underline             // 下划线
    blink                 // 闪烁
)

const (
    start     string = "\033["
    separator string = ";"
    end       string = "m"
)

var (
    Black   = foreground.color(black).String()
    Red     = foreground.color(red).String()
    Green   = foreground.color(green).String()
    Yellow  = foreground.color(yellow).String()
    Blue    = foreground.color(blue).String()
    Magenta = foreground.color(magenta).String()
    Cyan    = foreground.color(cyan).String()
    White   = foreground.color(white).String()

    BlackLight   = foreground.light().color(black).String()
    RedLight     = foreground.light().color(red).String()
    GreenLight   = foreground.light().color(green).String()
    YellowLight  = foreground.light().color(yellow).String()
    BlueLight    = foreground.light().color(blue).String()
    MagentaLight = foreground.light().color(magenta).String()
    CyanLight    = foreground.light().color(cyan).String()
    WhiteLight   = foreground.light().color(white).String()

    Reset           = reset.String()                                                   // 重置所有属性
    Clear           = attrs{foreground.color(clear), background.color(clear)}.String() // 清除前景和背景
    ClearForeground = foreground.color(clear).String()
    ClearBackground = background.color(clear).String()
)

func (c colorAttr) light() colorAttr {
    if c > 60 {
        return colorAttr(c)
    }
    return colorAttr(c + 60)
}

func (c colorAttr) color(color color) attr {
    return attr(int(c) + int(color))
}

func (a attr) String() string {
    return attrs{a}.String()
}

func (a attrs) String() string {
    s := mapSlice(a, func(a attr) string { return strconv.Itoa(int(a)) })
    return fmt.Sprintf("%s%s%s", start, strings.Join(s, separator), end)
}

func mapSlice[T, U any](from []T, transformer func(T) U) []U {
    ret := make([]U, len(from))
    for i, v := range from {
        ret[i] = transformer(v)
    }
    return ret
}
