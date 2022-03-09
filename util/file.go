package util

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
    "runtime"
    "syscall"
	"path/filepath"
	"regexp"
)

// FileInit 初始化文件目录
func FileInit(file string) (err error) {
	if !FileExists(file) {
		// 检查是否结尾是否存在点, 如果有点则吧最后一个当做文件处理
		dir := filepath.Dir(file)
		if dir != "." {
			// 递归创建目录
			if runtime.GOOS == "linux" {
				mask := syscall.Umask(0)
				defer syscall.Umask(mask)
				err = os.MkdirAll(dir, 0766)
			} else {
				err = os.MkdirAll(dir, 0766)
			}
		}
	}
	return
}

//var PATH_ROOT = SelfDir()
//var PATH_DATA = SelfDir()+"/data"

// SelfPath gets compiled executable file absolute path
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// SelfDir gets compiled executable file directory
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// WriteLog is use for
// if ok, err := util.WriteLog("Just a test\n"); !ok {
//     log.Print(err)
// }
// func WriteLog(filename string, format string) (bool, error) {
//     logfile := config.PathData+"/log/"+filename+".log"
//     f, err := os.OpenFile(logfile, os.O_RDWR | os.O_APPEND |  os.O_CREATE, 0777)
//     if err != nil {
//         return false, err
//     }
//     defer f.Close()
//     logger := log.New(f, "", log.Ldate | log.Ltime | log.Lshortfile)
//     logger.Print(format)
//     return true, err
// }

// PutFile is use for
// 使用方法
// if ok, err := util.PutFile("/data/golang/log/go.txt", "Just a test\n", 1); !ok {
//     log.Print(err)
// }
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

    // f的类型是*os.File，所以即使上面因为权限问题等问题导致f为nil了，一样可以Close()，因为*os.File是有Close()方法的，可以Close()多少次都行
    // 如果err不为空，说明f本来就是Close()的，虽然Close()也不会报错，但是直接返回还是省点资源的，所以这里直接就return吧
    // 如果是http抓取网页的 res.Body.Close() 就不同了，res.Body为空是不能Close()的，会报空指针异常，因为本来就是空指针
    // http://stackoverflow.com/questions/16280176/go-panic-runtime-error-invalid-memory-address-or-nil-pointer-dereference
    if err != nil {
        return false, err
    }
    defer f.Close()

    f.WriteString(format)
    return true, err
}

// GetFile is use for get file content
func GetFile(file string) (string, error) {
    f, err := os.Open(file)
    if err != nil {
        return "", err
    }
    defer f.Close() 
    fd, _ := ioutil.ReadAll(f)
    return string(fd), err
}

// SearchFile is use for Search a file in paths.
// this is often used in search config file in /etc ~/
func SearchFile(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, filename); FileExists(fullpath) {
			return
		}
	}
	err = errors.New(fullpath + " not found in paths")
	return
}

// GrepFile like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(patten string, filename string) (lines []string, err error) {
	re, err := regexp.Compile(patten)
	if err != nil {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)
	prefix := ""
	isLongLine := false
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		if er == io.EOF {
			break
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		} else {
			isLongLine = true
		}

		line = prefix + line
		if isLongLine {
			prefix = ""
		}
		if re.MatchString(line) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
