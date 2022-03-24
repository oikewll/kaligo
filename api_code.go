package kaligo

import (
    "fmt"
)

type Codes struct {
    SUCCESS uint
    FAILED uint
    CnMessage map[uint]string
    EnMessage map[uint]string
    LANG string
}

var ApiCode = &Codes{
    SUCCESS: 200,
    FAILED: 400,
    LANG: "en",
}

func init() {
       ApiCode.CnMessage = map[uint]string{
        ApiCode.SUCCESS: "操作成功",
        ApiCode.FAILED:  "操作失败",
    }
    ApiCode.EnMessage = map[uint]string{
        ApiCode.SUCCESS: "successful",
        ApiCode.FAILED:  "failed",
    }
}

func (c *Codes) GetMessage(code uint, a ...any) string {
    if c.LANG == "en" {
        message, ok := c.EnMessage[code]
        if !ok {
            return ""
        }
        return fmt.Sprintf(message, a...)
    }

    message, ok := c.CnMessage[code]
    if !ok {
        return ""
    }
    return fmt.Sprintf(message, a...)
}


