package tpl

import (
    "bytes"
    "fmt"
    "html/template"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"

    "github.com/owner888/kaligo/util"
)

type Tpl struct {
    *template.Template
    dir        string
    ext        string
    debug      bool
    funcs      template.FuncMap
    loadedAt   time.Time
    mu         sync.Mutex
    reloadTime time.Duration
}

func NewTmpl(dir, ext string, debug bool, funcs template.FuncMap, reloadTime time.Duration) (tpl *Tpl, err error) {
    if dir, err = filepath.Abs(dir); err != nil {
        return
    }
    tpl = new(Tpl)
    tpl.dir = dir
    tpl.ext = ext
    tpl.debug = debug
    tpl.funcs = funcs
    tpl.reloadTime = reloadTime
    if err = tpl.Load(); err != nil {
        tpl = nil
        return
    }
    if tpl.debug {
        var buf bytes.Buffer
        for _, tmpl := range tpl.Templates() {
            buf.WriteString("\t- ")
            buf.WriteString(tmpl.Name())
            buf.WriteString("\n")
        }
        //log.Printf("[Omux-debug] Loaded HTML Templates (%d): \n%s\n", len(tpl.Templates()), buf.String())
    }

    reload_ticker := time.NewTicker(reloadTime)
    go func() {
        for range reload_ticker.C {
            err := tpl.Load()
            if err != nil {
                log.Printf("tpl reload error:%v", err.Error())
            }
        }
    }()

    return
}

func (t *Tpl) Load() (err error) {

    t.loadedAt = time.Now()
    var root = template.New("").Funcs(t.funcs)
    var walkFunc = func(path string, info os.FileInfo, err error) (_ error) {
        if err != nil {
            return err
        }
        if !info.Mode().IsRegular() {
            return
        }
        if filepath.Ext(path) != t.ext {
            return
        }
        var rel string
        if rel, err = filepath.Rel(t.dir, path); err != nil {
            return err
        }
        rel = strings.TrimSuffix(rel, t.ext)
        var (
            nt = root.New(filepath.ToSlash(rel))
            b  []byte
        )
        if b, err = ioutil.ReadFile(path); err != nil {
            return err
        }
        _, err = nt.Parse(string(b))
        return err
    }

    if err = filepath.Walk(t.dir, walkFunc); err != nil {
        return
    }

    t.Template = root
    return
}

func (t *Tpl) Render(w io.Writer, name string, data interface{}) (err error) {
    t.mu.Lock()
    defer t.mu.Unlock()
    if t.debug == true {
        if err = t.Load(); err != nil {
            return
        }
    }
    tpf := filepath.Join(t.dir, name+t.ext)
    if !util.FileExists(tpf) {
        msg := "Template html paht not :%s"
        data = fmt.Sprintf(msg, name+t.ext)
        name = "error/404"
        log.Printf(msg, tpf)
    }
    err = t.ExecuteTemplate(w, name, data)
    return
}

func (t *Tpl) ReLoad() (err error) {
    t.mu.Lock()
    defer t.mu.Unlock()
    if err = t.Load(); err != nil {
        return
    }
    return
}
