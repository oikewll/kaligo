/**
 * Read the configuration file
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 */

package conf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
    // AppPath is ...
    AppPath string
    // PathRoot is PathRoot
    PathRoot string
    // PathData is PathData
    PathData string
	//confFile string                       // Your ini file path directory+file
	confList []map[string]map[string]string // Configuration information slice
)

// 一导入conf package 就初始化变量
func init() {
    //confFile := "../conf/app.ini"
    ////confFile := AppPath + "conf/app.ini"
    ////if len(os.Args) > 1 {
        ////confFile = os.Args[1]
    ////}
    //InitConfig(confFile)
    //PathRoot = Get("base", "path_root")
    //fmt.Println(PathRoot)
    //PathData = PathRoot + "/data"
}

// InitConfig is the function for create an empty configuration file
func InitConfig(confFile string) {
    fmt.Printf("confFile: [ %v ]", confFile)

    err := ReadList(confFile)
    if err != nil {
        Set("http", "addr", "0.0.0.0")
        Set("http", "port", "9527")
        path, _ := filepath.Abs(os.Args[0])
        //dir := filepath.Dir(path)
        pathArr := strings.Split(path, "/")
        Set("base", "basename", pathArr[len(pathArr)-1])
    }
}

// Get is the function for obtain corresponding value of the key values
func Get(section, name string) string {

	for _, v := range confList {
		for key, value := range v {
			if key == section {
				return value[name]
			}
		}
	}
	return ""
}

// Set is the function for set the corresponding value of the key value, if not add, if there is a key change
func Set(section, key, value string) bool {

	var ok bool
	var index = make(map[int]bool)
	var conf = make(map[string]map[string]string)
	for i, v := range confList {
		_, ok = v[section]
		index[i] = ok
	}

	i, ok := func(m map[int]bool) (i int, v bool) {
		for i, v := range m {
			if v == true {
				return i, true
			}
		}
		return 0, false
	}(index)

	if ok {
		confList[i][section][key] = value
        return true
    }

    conf[section] = make(map[string]string)
    conf[section][key] = value
    confList = append(confList, conf)
    return true
}

// Delete is the function for delete the corresponding key values
func Delete(section, name string) bool {

	for i, v := range confList {
		for key := range v {
			if key == section {
				delete(confList[i][key], name)
				return true
			}
		}
	}
	return false
}

// ReadList is the function for list all the configuration file
func ReadList(confFile string) (err error) {

	file, err := os.Open(confFile)
	if err != nil {
        return err
	}

	defer file.Close()
	var data map[string]map[string]string
	var section string
	buf := bufio.NewReader(file)
    for {
        l, err := buf.ReadString('\n')
        line := strings.TrimSpace(l)
        if err != nil {
            if err != io.EOF {
                return err
            }
            if len(line) == 0 {
                break
            }
        }
        switch {
        case strings.Index(line, "#") >= 0 || strings.Index(line, ";") >= 0:
        case len(line) == 0:
        case line[0] == '[' && line[len(line)-1] == ']':
            section = strings.TrimSpace(line[1 : len(line)-1])
            data = make(map[string]map[string]string)
            data[section] = make(map[string]string)
        default:
            i := strings.IndexAny(line, "=")
            value := strings.TrimSpace(line[i+1:])
            data[section][strings.TrimSpace(line[0:i])] = value
            if uniquappend(section) == true {
                confList = append(confList, data)
            }
        }
    }

	return err
}

// uniquappend is the function for ban repeated appended to the slice method
func uniquappend(conf string) bool {
	for _, v := range confList {
		for k := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}
