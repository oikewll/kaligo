## Epooll v1.0.0 (2014-10-01)

epooll is an open-source, high-performance, modularity, full-stack web framework.

More info [epooll.com](http://www.epooll.com)

* [Changelog](https://github.com/owner888/epooll/blob/master/README.md#changelog)
* [Installing](https://github.com/owner888/epooll/blob/master/README.md#installing)
* [Testing](https://github.com/owner888/epooll/blob/master/README.md#testing)
* [Examples](https://github.com/owner888/epooll/blob/master/README.md#examples)
* [To do](https://github.com/owner888/epooll/blob/master/README.md#to-do)
* [Known bugs](https://github.com/owner888/epooll/blob/master/README.md#known-bugs)
* [Documentation](https://github.com/owner888/epooll/blob/master/README.md#documentation)

## Changelog

v1.0.0: 初始化类库

## Installation

    $ go get github.com/owner888/epooll

## Examples

### 新建项目

    $ cd $GOPATH/src/
    $ mkdir epoollprojects
    $ cd epoollprojects
    $ tree
        ├── conf
        │   └── app.ini
        ├── control
        │   └── ctl_index.go
        ├── data
        │   ├── cache
        │   └── log
        ├── main.go
        ├── model
        │   └── mod_common.go
        ├── static
        │   ├── css
        │   ├── images
        │   └── js
        └── template
            └── index.tpl

### Example 1: 路由设置

    // main.go
    package main

    import (
        "epoollprojects/control"
        "github.com/owner888/epooll"
    )

    func main() {
        // 设置路由
        // 当用户访问 /?ct=index&ac=login 的时候就是调用了 control/ctl_index.go 里面的login方法
        epooll.Router("index", &control.Index{})
        // 设置静态路径，当用户访问 /static 的时候，就访问 static 目录下面的静态文件
        epooll.SetStaticPath("/static", "static")
        // 解析配置文件、编译模板、启动epooll模块、监听服务端口
        epooll.Run()
    }

### Example 2: Controller 编写

    package control

    import (
        "github.com/owner888/epooll"
        "net/http"
        "io"
    )

    type Index struct {}

    func (this *Index) Index(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "Hello World!")
    }


### Example 3: Model 编写

#### model 编写 mod_common.go

    package model

    import (
    )

    func GetString() string {
        return "Hi"
    }

#### 在控制器ctl_index.go 中使用model 

    import (
        "epoollprojects/model"
        "net/http"
        "io"
    )

    type Index struct {}

    func (this *Index) Index(w http.ResponseWriter, r *http.Request) {
        str := model.GetString()
        io.WriteString(w, str)
    }

### Example 4: View 编写

#### 模板文件：index.tpl

    <!DOCTYPE html>
    <html>
        <head>
            <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
            <title>Epooll</title>
        </head>
        <body>
            <h1>Hello Epooll Go!</h1>
            <p>Username: {{.username}}</p>
            <p>
                <ul>
                    {{range .data}}
                    <li>
                        id: {{.id}} --- email: {{.email}}
                    </li>
                    {{end}}
                </ul>
            </p>
        </body>
    </html>

#### 在控制器ctl_index.go 中使用模板

    import (
        "net/http"
        "github.com/owner888/epooll"
        "io"
        "html/template"
    )

    func (this *Demo) Index(w http.ResponseWriter, r *http.Request) {

        args := map[string]interface{}{
            "username":"yangzetao",
            "data":[]map[string]string{
                map[string]string{
                    "id":"1",
                    "email":"seatle@foxmail.com",
                },
                map[string]string{
                    "id":"2",
                    "email":"seatle@163.com",
                },
            },
        }

        t, err := template.ParseFiles("template/index.tpl")
        if err != nil {
            io.WriteString(w, fmt.Sprintf("%s", err));
            return
        }
        t.Execute(w, args)
    }

### Example 5: 静态文件处理

    ├── static
    │   ├── css
    │   ├── images
    │   └── js

    // 设置静态文件处理目录
    epooll.SetStaticPath("/static", "static")

    这样用户访问 URL http://localhost/static/123.txt 则会请求 static 目录下的 123.txt 文件

## To do

1. http header处理、cookie、session、isAjax
2. 数据库分页
3. 测试cpu、内存使用率
4. 完善使用文档
5. 基于此框架开发一个CMS
6. 一个服务自定义绑定多个端口

## Known bugs


## Documentation
* [control](http://www.godoc.org/pkg/github.com/owner888/epooll/control)
* [db](http://www.godoc.org/pkg/github.com/owner888/epooll/db)
* [redis](http://www.godoc.org/pkg/github.com/owner888/epooll/redis)

## Contact US
QQ:525773145

## LICENSE

epooll is licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).
