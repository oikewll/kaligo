package util

import (
	"log"
    "os"
    "io/ioutil"
)

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

/**
 *  使用方法
    if ok, err := util.WriteLog("/data/golang/log/go.txt", "Just a test\n"); !ok {
        log.Print(err)
    }
 */
func WriteLog(file string, format string) (bool, error) {

    f, err := os.OpenFile(file, os.O_RDWR | os.O_APPEND |  os.O_CREATE, 0777)
    if err != nil {
        return false, err
    }
    defer f.Close() 
    logger := log.New(f, "", log.Ldate | log.Ltime | log.Lshortfile)
    logger.Print(format)
    return true, err
}

/**
 *  使用方法
    if ok, err := util.PutFile("/data/golang/log/go.txt", "Just a test\n", 1); !ok {
        log.Print(err)
    }
 */
func PutFile(file string, format string, args ...interface{}) (bool, error) {

    f, err := os.OpenFile(file, os.O_RDWR | os.O_APPEND |  os.O_CREATE, 0777)
    // 上面的0777并不起作用
    os.Chmod(file, 0777)
    // 如果没有传参数，重新新建文件
    if args == nil {
        f, err = os.Create(file)
    }
    for _, arg := range args {
        // 参数为0，也重新创建文件
        if arg == 0 {
            f, err = os.Create(file)
        }
    }
    defer f.Close()

    if err != nil {
        return false, err
    }

    f.WriteString(format)
    return true, err
    //f.Write([]byte("Just a test!\r\n"))
}

func GetFile(file string) (string, error) {
    
    f, err := os.Open(file)
    if err != nil {
        // 抛出异常
        //panic(err)
        return "", err
    }
    defer f.Close() 
    // 这里不用处理错误了，如果是文件不存在或者没有读权限，上面都直接抛异常了，这里还可能有错误么？
    fd, _  := ioutil.ReadAll(f)
    return string(fd), err
}

