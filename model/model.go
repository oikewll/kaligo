package model

import (
    "os"
    "log"
)


type common struct {
}
func (this *common) PutFile(file string, content string, flag int) bool {
    flag = 0
    if flag == 0 {
        f, err := os.OpenFile(file, os.O_RDWR | os.O_CREATE, 0777)
    } else {
        f, err := os.OpenFile(file, os.O_RDWR | os.O_APPEND |  os.O_CREATE, 0777)
    }
    // 操作完关闭
    defer f.Close()

    if err != nil {
        log.Printf("open file error=%s\r\n", err.Error())
        return false
    }

    //logger := log.New(f, "", log.Ldate | log.Ltime | log.Lshortfile)
    //logger.Print("normal log 1")

    f.WriteString("Just a test!\r\n")

    return true
}
