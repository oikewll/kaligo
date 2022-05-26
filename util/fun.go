package util

import (
    "bytes"
    "compress/flate"
    "crypto/aes"
    "crypto/cipher"
    "crypto/md5"
    "crypto/sha256"
    "encoding/base64"
    "encoding/gob"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/pborman/uuid"
    "golang.org/x/net/publicsuffix"
    "io"
    "io/ioutil"
    "net"
    "log"
    "regexp"
    "math/rand"
    "net/url"
    "os"
    "os/exec"
    "path"
    "path/filepath"
    "reflect"
    "runtime"
    "strconv"
    "strings"
    "time"
    "golang.org/x/crypto/bcrypt"
)


func PasswordHash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func PasswordVerify(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// 去重
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
    a_len := len(a)
    for i := 0; i < a_len; i++ {
        if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
            continue
        }
        ret = append(ret, a[i])
    }
    return
}

func Md5(str string) string {
    data := []byte(str)
    has := md5.Sum(data)
    md5str1 := fmt.Sprintf("%x", has)
    return md5str1
}

func GenerateToken() string {
    return Md5(UUID())
}

// uuid
func UUID() string {
    return uuid.NewUUID().String()
}

// 随机生成测试IP
func RandIPAddr() string {
    ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
    log.Println("警告：这个是测试随机IP",ip)
    return ip
}

// 获取日期
func GetDay() string {
    return time.Now().Format("2006-01-02")
}

// 获取时间
func GetTimes() string {
    return time.Now().Format("2006-01-02 15:04:05")
}

// 获取文件 hash 值
func GetFileHash(f string) (string, error) {
    of, err := os.Open(f)
    if err != nil {
        return "", fmt.Errorf("Can't read binary (%s)", err)
    }
    hash := sha256.New()
    _, err = io.Copy(hash, of)
    if err != nil {
        return "", err
    }
    of.Close()
    return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// 随机数
func RandInt() int {
    return rand.Intn(100)
}

// 随机字符串
func RandomString(length int) string {
    str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
    b := []byte(str)
    var result []byte
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < length; i++ {
        result = append(result, b[r.Intn(len(b))])
    }
    return string(result)
}

// 字符串裁切
func StrLen(str string, strlen int) string {
    runes := []rune(str)
    if len(runes) > strlen {
        return string(runes[0:strlen]) + "..."
    }
    return str
}

// 随机
func Rands(len int) int {
    n := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(int32(len)))
    return ToInt(n)
}

// 生成订单 按日期
func OrderNum() string {
    return fmt.Sprintf("%s%d", time.Now().Format("20060102150101"), Rands(100000))
}

// 大小驼峰 转下划线 xorm 的 update 直接所有了mpa格式更新 要自己转换字段
func SnakeString(s string) string {
    data := make([]byte, 0, len(s)*2)
    j := false
    num := len(s)
    for i := 0; i < num; i++ {

        d := s[i]

        if i > 0 && d >= 'A' && d <= 'Z' && j {
            if string(s[i-1]) != "." {
                data = append(data, '_')
            }
        }
        if d != '_' {
            j = true
        }
        data = append(data, d)
    }
    return strings.ToLower(string(data[:]))
}

// 字符串加密 异或
func StringEncode(src string, key string) string {
    var result string
    j := 0
    s := ""
    bt := []rune(src)

    for i := 0; i < len(bt); i++ {
        s = strconv.FormatInt(int64(byte(bt[i])^key[j]), 16)
        if len(s) == 1 {
            s = "0" + s
        }
        result = result + (s)
        j = (j + 1) % len(key)
    }
    return result
}

// 字符串解密 异或
func StringDecode(src string, key string) string {
    var result string
    var s int64
    j := 0
    bt := []rune(src)
    for i := 0; i < len(src)/2; i++ {
        s, _ = strconv.ParseInt(string(bt[i*2:i*2+2]), 16, 0)
        result = result + string(byte(s)^key[j])
        j = (j + 1) % len(key)
    }
    return result
}

// json
func MapToJson(m interface{}) string {
    b, err := json.Marshal(m)
    if err == nil {
        return string(b)
    } else {

        return ""
    }
}

func MemStat() string {
    memStat := new(runtime.MemStats)
    runtime.ReadMemStats(memStat)
    return fmt.Sprintf("%dm", memStat.Sys/1024/1024)
}

// 打印内存信息
func PrintMemStats() {
    /**
    HeapSys：程序向应用程序申请的内存
    HeapAlloc：堆上目前分配的内存
    HeapIdle：堆上目前没有使用的内存
    HeapReleased：回收到操作系统的内存
    */
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    fmt.Printf("HeapAlloc = %v HeapIdel= %v HeapSys = %v  HeapReleased = %v\n", m.HeapAlloc/1024, m.HeapIdle/1024, m.HeapSys/1024, m.HeapReleased/1024)
}

// 打印结构体
func StructPrint(i any) {
    var kv = make(map[string]any)
    vValue := reflect.ValueOf(i)
    vType := reflect.TypeOf(i)
    for i := 0; i < vValue.NumField(); i++ {
        kv[vType.Field(i).Name] = vValue.Field(i)
    }

    fmt.Println("Struct: ")
    for k, v := range kv {
        fmt.Print(k)
        fmt.Print(":")
        fmt.Print(v)
        fmt.Println()
    }
}

func CopyBytes(b []byte) (copiedBytes []byte) {
    if b == nil {
        return nil
    }
    copiedBytes = make([]byte, len(b))
    copy(copiedBytes, b)

    return
}

// 转字符串数组
func IntArrayToString(values []int, delimiter string) string {
    var strBuffer bytes.Buffer
    for index, val := range values {
        strBuffer.WriteString(ToString(val))
        if index < (len(values) - 1) {
            strBuffer.WriteString(delimiter)
        }
    }
    return strBuffer.String()
}

// 以毫秒为单位返回当前unix时间戳
func GetTimestampString() string {
    return strconv.FormatInt(GetTimestamp(), 10)
}

func GetTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// 获取字符串键和值的映射 url.Values。
func MapToQueryParams(m map[string]string) url.Values {
    params := url.Values{}
    for key, value := range m {
        params.Add(key, value)
    }
    return params
}

func SerializeMap(m map[string]string) []byte {
    b := new(bytes.Buffer)
    e := gob.NewEncoder(b)
    e.Encode(m)
    return b.Bytes()
}

func DeserializeMap(b []byte) (map[string]string, error) {
    var decodedMap map[string]string
    d := gob.NewDecoder(bytes.NewBuffer(b))

    err := d.Decode(&decodedMap)
    if err != nil {
        return nil, err
    } else {
        return decodedMap, nil
    }
}

func AesEncrypt(orig string, key string) string {
    // 转成字节数组
    origData := []byte(orig)
    k := []byte(key)
    // 分组秘钥
    block, err := aes.NewCipher(k)
    if err != nil {
        return ""
    }
    // 获取秘钥块的长度
    blockSize := block.BlockSize()
    // 补全码
    origData = PKCS7Padding(origData, blockSize)
    // 加密模式
    blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
    // 创建数组
    cryted := make([]byte, len(origData))
    // 加密
    blockMode.CryptBlocks(cryted, origData)
    s := base64.StdEncoding.EncodeToString(cryted)
    s = strings.Replace(s, "+", "__", -1)
    return s
}

func AesDecrypt(cryted string, key string) (s string, err error) {
    cryted = strings.Replace(cryted, "__", "+", -1)
    // 转成字节数组
    crytedByte, err := base64.StdEncoding.DecodeString(cryted)
    if err != nil {
        return "", err
    }
    k := []byte(key)
    // 分组秘钥
    block, err := aes.NewCipher(k)
    if err != nil {
        return "", err
    }
    // 获取秘钥块的长度
    blockSize := block.BlockSize()

    // 加密模式
    blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
    // 创建数组
    orig := make([]byte, len(crytedByte))

    if len(crytedByte)%blockSize != 0 {
        return "", errors.New("rypto/cipher: input not full blocks")
    }

    if len(orig) < len(crytedByte) {
        return "", errors.New("crypto/cipher: output smaller than input")
    }

    // 解密
    blockMode.CryptBlocks(orig, crytedByte)
    // 去补全码
    orig = PKCS7UnPadding(orig)
    return string(orig), nil
}

// 补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
    padding := blocksize - len(ciphertext)%blocksize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

// 去码
func PKCS7UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

// 时间是加一天
func TimeAdd1(t string) string {
    now, err := time.ParseInLocation("2006-01-02", t, time.Local)
    if err != nil {
        return ""
    }
    pd, _ := time.ParseDuration("24h")
    now = now.Add(pd)
    return time.Time(now).Format("2006-01-02")
}

// 时间是减几天
func TimeReduce(n int) string {
    return time.Now().AddDate(0, 0, -n).Format("2006-01-02")
}

func TimeAddYears(n int) string {
    return time.Now().AddDate(n, 0, 0).Format("2006-01-02")
}

func TimeCurrentMonthDay() string {
    currentYear, currentMonth, _ := time.Now().Date()
    currentLocation := time.Now().Location()

    firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
    //lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
    return firstOfMonth.Format("2006-01-02")
}

func TimeDayNum(startDate string, endDate string) int64 {
    t1, err := time.Parse("2006-01-02", startDate)
    t2, err := time.Parse("2006-01-02", endDate)
    if err == nil && t1.Before(t2) { //判断时间t是否在时间ｕ的前面
        return (t2.Unix() - t1.Unix()) / 86400
    }
    return 0
}

func TimeSub(t1, t2 time.Time) int {
    t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
    t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
    return int(t1.Sub(t2).Hours() / 24)
}

func TimeDayList(startDate string, endDate string) []string {
    dayNum := int(TimeDayNum(startDate, endDate))

    var d []string
    if dayNum > 0 {
        var dlist string
        for i := 0; i <= dayNum; i++ {
            if i == 0 {
                d = append(d, startDate)
            } else {
                dlist = TimeAdd1(d[len(d)-1])
                d = append(d, dlist)
            }
        }
    } else {
        d = append(d, startDate)
    }

    return d
}

// 获取应用系统的当前路径
func GetBasePath() (string, error) {
    var dir, currentPath string
    cwd, err := exec.LookPath(os.Args[0])
    if err == nil {
        //cwd = strings.Replace(cwd, `\`, `/`, -1)

        if runtime.GOOS == "windows" {
            cwd = filepath.ToSlash(cwd)
        }

        dir, _ = path.Split(cwd)
        os.Chdir(dir)
        currentPath, err = os.Getwd()
    }
    return currentPath, err
}

// 判断文件是否存在
func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

// 获取用户自增ID 16 位MD5
func GetKNS(KNS int64) string {
    return Md5(Md5(ToString(KNS)))[7:23]
}

func DeWriter(data string) (string, error) {
    buf := bytes.NewBuffer(nil)
    flateWrite, err := flate.NewWriter(buf, flate.BestCompression)
    if err != nil {
        return "", err
    }
    _, err = flateWrite.Write([]byte(data))
    if err != nil {
        return "", err
    }
    err = flateWrite.Flush()
    if err != nil {
        return "", err
    }
    err = flateWrite.Close()
    if err != nil {
        return "", err
    }
    return buf.String(), nil
}

func DeReader(data string) (string, error) {
    reader := flate.NewReader(strings.NewReader(string(data)))
    defer reader.Close()
    out, err := ioutil.ReadAll(reader)
    if err != nil {
        return "", err
    }
    return string(out), nil
}

// 获取域名信息
func GetRootDomain(urls string) (root, tld string, err error) {
    ra, err := regexp.Compile("^(http:|https:)")
    if err != nil {
        return "", "", err
    }
    isAddScheme := len(ra.FindIndex([]byte(urls))) == 0

    if isAddScheme {
        urls = "http://" + urls
    }
    var host string
    r, err := url.Parse(urls)
    if err != nil {
        return "", "", err
    }

    host = r.Hostname()
    //port = r.Port()

    isIP := net.ParseIP(host)
    if isIP != nil {
        return host, host, err
    }

    root, _ = publicsuffix.EffectiveTLDPlusOne(host)
    tld, _  = publicsuffix.PublicSuffix(host)
    //fmt.Printf("根域名%s 域名后缀%s 后缀是否有效 %v", root, tld, isok)
    return
}

func CreateFolder(folder string) error {
    _, err := os.Stat(folder)
    if err != nil {
        if os.IsNotExist(err) {
            err = os.MkdirAll(folder, os.ModePerm)
            return err
        }
    }
    return err
}
