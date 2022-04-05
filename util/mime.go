package util

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strings"
)

// MIME http header 数据类型
type MIME string

const (
    MIMEJson              MIME = "application/json"
    MIMEPostForm          MIME = "application/x-www-form-urlencoded"
    MIMEMultipartPOSTForm MIME = "multipart/form-data"
)

// MIMEParser 数据解析器
type MIMEParser interface {
    Type() MIME
    ParseValues(request *http.Request) UrlValues
}

type MIMEParsers []MIMEParser

// DefaultMIMEParsers 默认的一组解析器，Context 用于解析 Request form 数据
var DefaultMIMEParsers = MIMEParsers{&MIMEJsonParser{}, &MIMEMultipartPOSTFormParser{maxMemory: 32 << 20}, &MIMEPostFormParser{}}

func (m MIMEParsers) ParseValues(request *http.Request) UrlValues {
    contentType := request.Header.Get("Content-Type")
    for _, v := range m {
        if strings.Contains(contentType, string(v.Type())) {
            return v.ParseValues(request)
        }
    }
    return UrlValues{}
}

type MIMEJsonParser struct{}

func (m *MIMEJsonParser) Type() MIME {
    return MIMEJson
}

func (m *MIMEJsonParser) ParseValues(request *http.Request) UrlValues {
    values := UrlValues{}
    var form map[string]any
    json.NewDecoder(request.Body).Decode(&form)
    defer request.Body.Close()
    for k, v := range form {
        values[k] = []string{fmt.Sprint(v)}
    }
    return values
}

type MIMEPostFormParser struct{}

func (m *MIMEPostFormParser) Type() MIME {
    return MIMEPostForm
}

func (m *MIMEPostFormParser) ParseValues(request *http.Request) UrlValues {
    if err := request.ParseForm(); err != nil {

    }
    return UrlValues(request.Form)
}

type MIMEMultipartPOSTFormParser struct {
    maxMemory int64
}

func (m *MIMEMultipartPOSTFormParser) Type() MIME {
    return MIMEMultipartPOSTForm
}

func (m *MIMEMultipartPOSTFormParser) ParseValues(request *http.Request) UrlValues {
    if err := request.ParseMultipartForm(m.maxMemory); err != nil {
        if !errors.Is(err, http.ErrNotMultipart) {
            // debugPrint("error on parse multipart form array: %v", err)
        }
    }
    return UrlValues(request.PostForm)
}
