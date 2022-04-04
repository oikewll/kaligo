// +----------------------------------------------------------------------
// | KaliGo [ A Golang Framework For Web ]
// +----------------------------------------------------------------------
// | Copyright (c) 2006-2014 https://doc.kaligo.org All rights reserved.
// +----------------------------------------------------------------------
// | Licensed ( http://www.apache.org/licenses/LICENSE-2.0 )
// +----------------------------------------------------------------------
// | Author: Seatle <seatle888@gmail.com>
// +----------------------------------------------------------------------

// +----------------------------------------------------------------------
// | GET 请求
// | r := requests.Get('http://www.test.com').Param("username", "test");
// | PHP SERVER
// | $_GET
// +----------------------------------------------------------------------
// | POST 请求
// | $data = array('name'=>'request');
// | requests.Post('http://www.test.com', $data);
// | PHP SERVER
// | $_POST
// +----------------------------------------------------------------------
// | POST RESTful请求
// | $data = array('name'=>'request');
// | $data_string = json_encode($data);
// | requests.SetHeader("Content-Type", "application/json");
// | requests.Post('http://www.test.com', $data_string);
// | PHP SERVER
// | file_get_contents('php://input')
// +----------------------------------------------------------------------
// | POST 文件上传
// | $data = array('file1'=>''./data/phpspider.log'');
// | requests.Post('http://www.test.com', null, $data);
// | PHP SERVER
// | $_FILES
// +----------------------------------------------------------------------
// | 代理
// | requests.SetProxy(array('223.153.69.150:42354'));
// | $html = requests.Get('https://www.test.com');
// +----------------------------------------------------------------------

//----------------------------------
// KaliGo 请求类
//----------------------------------

package requests

import (
    "bytes"
    "compress/gzip"
    "crypto/tls"
    "encoding/json"
    "encoding/xml"
    "io"
    "io/ioutil"
    // "log"
    "mime/multipart"
    "net"
    "net/http"
    "net/http/cookiejar"
    "net/http/httputil"
    "net/url"
    "os"
    "regexp"
    "strings"
    "sync"
    "time"
)

const (
    VERSION = "1.0.1"
)

var (
    defaultSetting = Settings{false, "requests/2.0.0", 60 * time.Second, 60 * time.Second, nil, nil, nil, false, true, true}
    defaultCookieJar http.CookieJar
    settingMutex sync.Mutex
)

// Settings
type Settings struct {
    ShowDebug        bool
    UserAgent        string
    ConnectTimeout   time.Duration
    ReadWriteTimeout time.Duration
    TlsClientConfig  *tls.Config
    Proxy            func(*http.Request) (*url.URL, error)
    Transport        http.RoundTripper
    EnableCookie     bool
    Gzip             bool
    DumpBody         bool
}

// Requests provides more useful methods for requesting one url than http.Request.
type Requests struct {
    setting Settings
    url     string
    req     *http.Request
    params  map[string]string
    files   map[string]string
    resp    *http.Response
    body    []byte
    dump    []byte
    errs    []error
}

// createDefaultCookie creates a global cookiejar to store cookies.
func createDefaultCookie() {
    settingMutex.Lock()
    defer settingMutex.Unlock()
    defaultCookieJar, _ = cookiejar.New(nil)
}

// Overwrite default settings
func SetDefaultSetting(setting Settings) {
    settingMutex.Lock()
    defer settingMutex.Unlock()
    defaultSetting = setting
    if defaultSetting.ConnectTimeout == 0 {
        defaultSetting.ConnectTimeout = 60 * time.Second
    }
    if defaultSetting.ReadWriteTimeout == 0 {
        defaultSetting.ReadWriteTimeout = 60 * time.Second
    }
}

// return *Requests with specific method
func NewRequests(rawurl, method string) *Requests {
    var resp http.Response
    u, err := url.Parse(rawurl)
    if err != nil {
        // log.Fatal(err)
        var errs []error    
        errs = append(errs, err)
        return &Requests{defaultSetting, rawurl, nil, map[string]string{}, map[string]string{}, &resp, nil, nil, errs}
    }

    req := http.Request{
        URL:        u,
        Method:     method,
        Header:     make(http.Header),
        Proto:      "HTTP/1.1",
        ProtoMajor: 1,
        ProtoMinor: 1,
    }
    return &Requests{defaultSetting, rawurl, &req, map[string]string{}, map[string]string{}, &resp, nil, nil, nil}
}

// Get returns *Requests with GET method.
func Get(url string) *Requests {
    return NewRequests(url, "GET")
}

// Post returns *Requests with POST method.
func Post(url string) *Requests {
    return NewRequests(url, "POST")
}

// Put returns *Requests with PUT method.
func Put(url string) *Requests {
    return NewRequests(url, "PUT")
}

// Delete returns *Requests DELETE method.
func Delete(url string) *Requests {
    return NewRequests(url, "DELETE")
}

// Head returns *Requests with HEAD method.
func Head(url string) *Requests {
    return NewRequests(url, "HEAD")
}

// Head returns *Requests with OPTIONS method.
func Options(url string) *Requests {
    return NewRequests(url, "OPTIONS")
}

// Head returns *Requests with PATCH method.
func Patch(url string) *Requests {
    return NewRequests(url, "PATCH")
}

// Change request settings
func (r *Requests) Setting(setting Settings) *Requests {
    r.setting = setting
    return r
}

// GetErrors return all error
func (r *Requests) GetErrors() []error {
    return r.errs
}

// SetBasicAuth sets the request's Authorization header to use HTTP Basic Authentication with the provided username and password.
func (r *Requests) SetBasicAuth(username, password string) *Requests {
    r.req.SetBasicAuth(username, password)
    return r
}

// SetEnableCookie sets enable/disable cookiejar
func (r *Requests) SetEnableCookie(enable bool) *Requests {
    r.setting.EnableCookie = enable
    return r
}

// Debug sets show debug or not when executing request.
func (r *Requests) Debug(isdebug bool) *Requests {
    r.setting.ShowDebug = isdebug
    return r
}

// Dump Body.
func (r *Requests) DumpBody(isdump bool) *Requests {
    r.setting.DumpBody = isdump
    return r
}

// return the DumpRequest
func (r *Requests) DumpRequest() []byte {
    return r.dump
}

// SetTimeout sets connect time out and read-write time out for Requests.
func (r *Requests) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *Requests {
    r.setting.ConnectTimeout = connectTimeout
    r.setting.ReadWriteTimeout = readWriteTimeout
    return r
}

// SetTLSClientConfig sets tls connection configurations if visiting https url.
func (r *Requests) SetTLSClientConfig(config *tls.Config) *Requests {
    r.setting.TlsClientConfig = config
    return r
}

// 设置语言请求，默认中文
func (r *Requests) SetAcceptLanguage(lang ...string) *Requests {
    if len(lang) == 0 {
        r.req.Header.Set("Accept-Language", "zh-CN")
    } else {
        r.req.Header.Set("Accept-Language", lang[0])
    }
    return r
}

func (r *Requests) SetClientIP(ip string) *Requests {
    r.req.Header.Set("CLIENT-IP", ip)
    r.req.Header.Set("X-FORWARDED-FOR", ip)
    return r
}

func (r *Requests) SetReferer(referer string) *Requests {
    r.req.Header.Set("Referer", referer)
    return r
}

// Header add header item string in request.
func (r *Requests) Header(key, value string) *Requests {
    r.req.Header.Set(key, value)
    return r
}

// SetUserAgent sets User-Agent header field
// r.SetUserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
func (r *Requests) SetUserAgent(useragent string) *Requests {
    r.setting.UserAgent = useragent
    return r
}

// Set HOST
// 负载均衡到不同的服务器，如果对方使用CDN，采用这个是最好的了
func (r *Requests) SetHost(host string) *Requests {
    r.req.Host = host
    return r
}

// Set the protocol version for incoming requests.
// Client requests always use HTTP/1.1.
func (r *Requests) SetProtocolVersion(vers string) *Requests {
    if len(vers) == 0 {
        vers = "HTTP/1.1"
    }

    major, minor, ok := http.ParseHTTPVersion(vers)
    if ok {
        r.req.Proto = vers
        r.req.ProtoMajor = major
        r.req.ProtoMinor = minor
    }

    return r
}

// SetCookie add cookie into request.
func (r *Requests) SetCookie(cookie *http.Cookie) *Requests {
    r.req.Header.Add("Cookie", cookie.String())
    return r
}

// Set transport to
func (r *Requests) SetTransport(transport http.RoundTripper) *Requests {
    r.setting.Transport = transport
    return r
}

// Set http proxy
// example:
//
//	func(req *http.Request) (*url.URL, error) {
// 		u, _ := url.ParseRequestURI("http://127.0.0.1:8118")
// 		return u, nil
// 	}
func (r *Requests) SetProxy(proxy func(*http.Request) (*url.URL, error)) *Requests {
    r.setting.Proxy = proxy
    return r
}

func (r *Requests) SetProxyString(proxy string) *Requests {
    r.setting.Proxy = func(req *http.Request) (*url.URL, error) {
        u, _ := url.ParseRequestURI(proxy)
        return u, nil
    }

    return r
}

// Param adds query param in to request.
// params build query string as ?key1=value1&key2=value2...
func (r *Requests) Param(key, value string) *Requests {
    // 不用加锁，因为一个请求一个 Requests struct
    r.params[key] = value
    return r
}

func (r *Requests) Params(params map[string]string) *Requests {
    for key, value := range params {
        r.params[key] = value
    }
    return r
}

func (r *Requests) PostFile(formname, filename string) *Requests {
    r.files[formname] = filename
    return r
}

// Body adds request raw body.
// it supports string and []byte.
func (r *Requests) Body(data any) *Requests {
    switch t := data.(type) {
    case string:
        buf := bytes.NewBufferString(t)
        r.req.Body = ioutil.NopCloser(buf)
        r.req.ContentLength = int64(len(t))
    case []byte:
        buf := bytes.NewBuffer(t)
        r.req.Body = ioutil.NopCloser(buf)
        r.req.ContentLength = int64(len(t))
    }
    return r
}

// JsonBody adds request raw body encoding by JSON.
func (r *Requests) JsonBody(obj any) (*Requests, error) {
    if r.req.Body == nil && obj != nil {
        buf := bytes.NewBuffer(nil)
        enc := json.NewEncoder(buf)
        if err := enc.Encode(obj); err != nil {
            return r, err
        }
        r.req.Body = ioutil.NopCloser(buf)
        r.req.ContentLength = int64(buf.Len())
        r.req.Header.Set("Content-Type", "application/json")
    }
    return r, nil
}

func (r *Requests) buildUrl(paramBody string) {
    // build GET url with query string
    if r.req.Method == "GET" && len(paramBody) > 0 {
        if strings.Index(r.url, "?") != -1 {
            r.url += "&" + paramBody
        } else {
            r.url = r.url + "?" + paramBody
        }
        return
    }

    // build POST/PUT/PATCH url and body
    if (r.req.Method == "POST" || r.req.Method == "PUT" || r.req.Method == "PATCH") && r.req.Body == nil {
        // with files
        if len(r.files) > 0 {
            pr, pw := io.Pipe()
            bodyWriter := multipart.NewWriter(pw)
            go func() {

                for formname, filename := range r.files {
                    fileWriter, err := bodyWriter.CreateFormFile(formname, filename)
                    if err != nil {
                        r.errs = append(r.errs, err)
                    }
                    fh, err := os.Open(filename)
                    if err != nil {
                        r.errs = append(r.errs, err)
                    }
                    // iocopy
                    _, err = io.Copy(fileWriter, fh)
                    fh.Close()
                    if err != nil {
                        r.errs = append(r.errs, err)
                    }
                }
                for k, v := range r.params {
                    bodyWriter.WriteField(k, v)
                }
                bodyWriter.Close()
                pw.Close()
            }()
            r.Header("Content-Type", bodyWriter.FormDataContentType())
            r.req.Body = ioutil.NopCloser(pr)
            return
        }

        // with params
        if len(paramBody) > 0 {
            r.Header("Content-Type", "application/x-www-form-urlencoded")
            r.Body(paramBody)
        }
    }
}

func (r *Requests) getResponse() (*http.Response, error) {
    if r.resp.StatusCode != 0 {
        return r.resp, nil
    }
    resp, err := r.SendOut()
    if err != nil {
        return nil, err
    }
    r.resp = resp
    return resp, nil
}

func (r *Requests) SendOut() (*http.Response, error) {
    var paramBody string
    if len(r.params) > 0 {
        var buf bytes.Buffer
        for k, v := range r.params {
            buf.WriteString(url.QueryEscape(k))
            buf.WriteByte('=')
            buf.WriteString(url.QueryEscape(v))
            buf.WriteByte('&')
        }
        paramBody = buf.String()
        paramBody = paramBody[0 : len(paramBody)-1]
    }

    r.buildUrl(paramBody)
    url, err := url.Parse(r.url)
    if err != nil {
        return nil, err
    }

    r.req.URL = url

    trans := r.setting.Transport

    if trans == nil {
        // create default transport
        trans = &http.Transport{
            TLSClientConfig: r.setting.TlsClientConfig,
            Proxy:           r.setting.Proxy,
            Dial:            TimeoutDialer(r.setting.ConnectTimeout, r.setting.ReadWriteTimeout),
        }
    } else {
        // if b.transport is *http.Transport then set the settings.
        if t, ok := trans.(*http.Transport); ok {
            if t.TLSClientConfig == nil {
                t.TLSClientConfig = r.setting.TlsClientConfig
            }
            if t.Proxy == nil {
                t.Proxy = r.setting.Proxy
            }
            if t.Dial == nil {
                t.Dial = TimeoutDialer(r.setting.ConnectTimeout, r.setting.ReadWriteTimeout)
            }
        }
    }

    var jar http.CookieJar = nil
    if r.setting.EnableCookie {
        if defaultCookieJar == nil {
            createDefaultCookie()
        }
        jar = defaultCookieJar
    }

    client := &http.Client{
        Transport: trans,
        Jar:       jar,
    }

    if r.setting.UserAgent != "" && r.req.Header.Get("User-Agent") == "" {
        r.req.Header.Set("User-Agent", r.setting.UserAgent)
    }

    if r.setting.ShowDebug {
        dump, err := httputil.DumpRequest(r.req, r.setting.DumpBody)
        if err != nil {
            println(err.Error())
        }
        r.dump = dump
    }
    return client.Do(r.req)
}

// String returns the body string in response.
// it calls Response inner.
func (r *Requests) String() (string, error) {
    data, err := r.Bytes()
    if err != nil {
        return "", err
    }

    return string(data), nil
}

// Bytes returns the body []byte in response.
// it calls Response inner.
func (r *Requests) Bytes() ([]byte, error) {
    if r.body != nil {
        return r.body, nil
    }
    resp, err := r.getResponse()
    if err != nil {
        return nil, err
    }
    if resp.Body == nil {
        return nil, nil
    }
    defer resp.Body.Close()
    if r.setting.Gzip && resp.Header.Get("Content-Encoding") == "gzip" {
        reader, err := gzip.NewReader(resp.Body)
        if err != nil {
            return nil, err
        }
        r.body, err = ioutil.ReadAll(reader)
    } else {
        r.body, err = ioutil.ReadAll(resp.Body)
    }
    if err != nil {
        return nil, err
    }
    return r.body, nil
}

// ToFile saves the body data in response to one file.
// it calls Response inner.
func (r *Requests) ToFile(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    resp, err := r.getResponse()
    if err != nil {
        return err
    }
    if resp.Body == nil {
        return nil
    }
    defer resp.Body.Close()
    _, err = io.Copy(f, resp.Body)
    return err
}

// ToJson returns the map that marshals from the body bytes as json in response .
// it calls Response inner.
func (r *Requests) ToJson(v any) error {
    data, err := r.Bytes()
    if err != nil {
        return err
    }
    return json.Unmarshal(data, v)
}

// ToXml returns the map that marshals from the body bytes as xml in response .
// it calls Response inner.
func (r *Requests) ToXml(v any) error {
    data, err := r.Bytes()
    if err != nil {
        return err
    }
    return xml.Unmarshal(data, v)
}

// Response executes request client gets response mannually.
func (r *Requests) Response() (*http.Response, error) {
    return r.getResponse()
}

// TimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
    return func(netw, addr string) (net.Conn, error) {
        conn, err := net.DialTimeout(netw, addr, cTimeout)
        if err != nil {
            return nil, err
        }
        conn.SetDeadline(time.Now().Add(rwTimeout))
        return conn, nil
    }
}

// var request *Requests

// func init() {
//     request = New()
// }
//
// func New() *Requests {
//     return &Requests{}
// }

// type Requestsss struct {
//     Timeout         int
//     Encoding        string      // Xpath only supports utf-8, other encodings need to be transcoded
//     InputEncoding   string
//     OutputEncoding  string
//     Cookies         *sync.Map   // map of cookies to pass
//     RawHeaders      *sync.Map   // map of raw headers to send
//     DomainCookies   *sync.Map   // map of cookies for domain to pass
//     Hosts           *sync.Map   // random host binding for make request faster
//     Headers         *sync.Map   // headers returned from server sent here
//     UserAgents      []string    // random agent we masquerade as
//     ClientIPs       []string    // random ip we masquerade as
//     Proxies         []string    // random proxy ip
//     Raw             string      // head + body content returned from server sent here
//     Header          string      // head content returned from server sent here
//     Content         string      // The body before encoding
//     Text            string      // The body after encoding
//     Info            *sync.Map   // curl info like: http_code、header_size
//     History         int         // http request status before redirect. ex:302
//     Status          int         // http request status
//     Err             error       // error message sent here
// }

func (r *Requests) Request(method, url string, fields any, files []string, allowRedirects bool, cert string) string {
    return ""
}

/**
* 简单的判断一下参数是否为一个URL链接
* @param  string  $str 
* @return boolean      
*/
func (r *Requests) ISUrl(url string) bool {
    // $pattern = "/\b(([\w-]+:\/\/?|www[.])[^\s()<>]+(?:\([\w\d]+\)|([^[:punct:]\s]|\/)))/";
    if match, _ := regexp.MatchString(`^(\w+):\/\/([^/:]+)\/([^/:]+)\/(\d*)?\/(\d*)$`, url); !match {
        return false
    } else {
        return true
    }
}
