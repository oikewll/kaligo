package requests

import (
    "fmt"
    "io/ioutil"
    "math/rand"
    "net/http"
    "net/url"
    "sync"
    "time"
)

type Url struct {
    Method string
    UrlStr string
    Params map[string]string
    Extras any // 自定义数据
}

func (u *Url) String() string {
    urlStr := u.UrlStr
    if len(u.Params) > 0 && u.Method == http.MethodGet {
        param := url.Values{}
        for k, v := range u.Params {
            param.Set(k, v)
        }
        urlStr = urlStr + "?" + param.Encode()
    }
    return urlStr
}

type Worker struct {
    count    int           // 协程数量
    channel  chan struct{} // 空结构体变量的内存占用大小为 0，而 bool 类型内存占用大小为 1
    urls     []*Url        // 要采集的 URLs
    Callback func(*Url)    // 采集成功回调
}

// Worker struct 初始化
func NewWorker(count int) *Worker {
    return &Worker{
        count:   count,
        channel: make(chan struct{}, count),
    }
}

func (w *Worker) AddUrl(method, urlStr string, params map[string]string, extras any) {
    w.urls = append(w.urls, &Url{
        Method: method,
        UrlStr: urlStr,
        Params: params,
        Extras: extras,
    })
}

// Run 方法：创建有限的 go callBack 函数的 goroutine
func (w *Worker) Run() {
    for _, v := range w.urls {
        w.channel <- struct{}{}
        go func(u *Url) {
            w.Callback(u)
            <-w.channel
        }(v)
    }
}

// WaitGroup 对象内部有一个计数器，从 0 开始
// 有三个方法：Add(), Done(), Wait() 用来控制计数器的数量
var wg = sync.WaitGroup{}

func main() {
    start := time.Now()
    worker := NewWorker(5)

    // 接口请求URL
    //max := int(math.Pow10(8))                 // 模拟一千万数据
    max := 15 // 先测试 5 次吧

    for i := 0; i < max; i++ {
        wg.Add(1)

        // 随机手机号码参数
        p := map[string]string{
            "phone": RandMobile(),
        }
        worker.AddUrl("GET", "http://apis.juhe.cn/mobile/get", p, nil)
        worker.Callback = func(u *Url) {
            param := url.Values{}
            param.Set("key", "您申请的KEY") // 接口请求Key
            for key, value := range u.Params {
                param.Set(key, value)
            }

            // 发送请求
            data, err := HttpGet(u.UrlStr, param)
            if err != nil {
                fmt.Println(err)
                return
            }

            // 其它逻辑代码...

            fmt.Println(string(data))
            // time.Sleep(time.Second * 2)
            wg.Done()
        }
    }

    worker.Run()

    // 阻塞代码防止退出
    wg.Wait()

    fmt.Printf("耗时: %fs\n", time.Now().Sub(start).Seconds())
}

// Get 方式发起网络请求
func HttpGet(apiURL string, params url.Values) (rs []byte, err error) {
    var Url *url.URL
    Url, err = url.Parse(apiURL)
    if err != nil {
        return nil, err
    }
    // 如果参数中有中文参数,这个方法会进行URLEncode
    Url.RawQuery = params.Encode()
    req, err := http.NewRequest(http.MethodGet, Url.String(), nil)
    if err != nil {
        return nil, err
    }
    req.Header.Add("Pragma", "no-cache")
    req.Header.Add("Accept", "*/*")
    req.Header.Add("Accept-Language", "zh-CN,zh-Hans;q=0.9")
    // req.Header.Add("Accept-Encoding", "gzip, deflate, br")
    req.Header.Add("Cache-Control", "no-cache")
    req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)
}

var MobilePrefix = [...]string{"130", "131", "132", "133", "134", "135", "136", "137", "138", "139", "145", "147", "150", "151", "152", "153", "155", "156", "157", "158", "159", "170", "176", "177", "178", "180", "181", "182", "183", "184", "185", "186", "187", "188", "189"}

// GeneratorPhone 生成手机号码
func RandMobile() string {
    return MobilePrefix[RandInt(0, len(MobilePrefix))] + fmt.Sprintf("%0*d", 8, RandInt(0, 100000000))
}

// 指定范围随机 int
func RandInt(min, max int) int {
    rand.Seed(time.Now().UnixNano())
    return min + rand.Intn(max-min)
}
