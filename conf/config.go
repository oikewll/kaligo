/**
 * Read the configuration file
 *
 * @copyright           (C) 2014  seatle
 * @lastmodify          2014-12-01
 * @website		http://www.epooll.com
 *
 */

package conf

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
    "fmt"
)

var (
    PATH_ROOT string
    PATH_DATA string
	confFile string                         //your ini file path directory+file
	confList []map[string]map[string]string //configuration information slice
)

// 一导入conf package就初始化变量
func init() {
    confFile := "conf/app.ini"
    if len(os.Args) > 1 {
        confFile = os.Args[1]
    }
    InitConfig(confFile)
    PATH_ROOT = GetValue("base", "path_root")
    fmt.Println(PATH_ROOT)
    PATH_DATA = PATH_ROOT + "/data"
}
//Create an empty configuration file
func InitConfig(filename string) {
	confFile = filename
    err := ReadList()
    if err != nil {

        SetValue("http", "addr", "0.0.0.0")
        SetValue("http", "port", "9527")
        path, _ := filepath.Abs(os.Args[0])
        //dir := filepath.Dir(path)
        pathArr := strings.Split(path, "/")
        SetValue("base", "basename", pathArr[len(pathArr)-1])
    }
}

//To obtain corresponding value of the key values
func GetValue(section, name string) string {

	for _, v := range confList {
		for key, value := range v {
			if key == section {
				return value[name]
			}
		}
	}
	return ""
}

//Set the corresponding value of the key value, if not add, if there is a key change
func SetValue(section, key, value string) bool {

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
	} else {
		conf[section] = make(map[string]string)
		conf[section][key] = value
		confList = append(confList, conf)
		return true
	}

	return false
}

//Delete the corresponding key values
func DeleteValue(section, name string) bool {

	for i, v := range confList {
		for key, _ := range v {
			if key == section {
				delete(confList[i][key], name)
				return true
			}
		}
	}
	return false
}

//List all the configuration file
func ReadList() (err error) {

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
			value := strings.TrimSpace(line[i+1 : len(line)])
			data[section][strings.TrimSpace(line[0:i])] = value
			if uniquappend(section) == true {
				confList = append(confList, data)
			}
		}

	}

	return err
}

//Ban repeated appended to the slice method
func uniquappend(conf string) bool {
	for _, v := range confList {
		for k, _ := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}

