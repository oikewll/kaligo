## KaliGo v1.2.4 (2014-10-12)

KaliGo is an open-source, high-performance, modularity, full-stack web framework.

More info [doc.kaligo.com](https://doc.kaligo.com)

* [Changelog](https://github.com/owner888/kaligo/blob/master/README.md#changelog)
* [Installing](https://github.com/owner888/kaligo/blob/master/README.md#installing)
* [Testing](https://github.com/owner888/kaligo/blob/master/README.md#testing)
* [Examples](https://github.com/owner888/kaligo/blob/master/README.md#examples)
* [To do](https://github.com/owner888/kaligo/blob/master/README.md#to-do)
* [Known bugs](https://github.com/owner888/kaligo/blob/master/README.md#known-bugs)
* [Documentation](https://github.com/owner888/kaligo/blob/master/README.md#documentation)

## Changelog

v1.2.8: 重构config类，支持点式获取方式，config.Get[string]("database.mysql.charset")

v1.2.7: 定时器增加全局锁，修复同时设置多个定时器时协程并发写Timer map出现空指针异常Bug，增加用读写锁实现的安全map，排序map类改名

v1.2.6: 迁移util类到util namespace里面去，调用就不用每次都new那么麻烦了，修改CRUD数据库类插入、批量插入、修改为返回sql，以便记录日志;增加项目绝对路径，数据绝对路径util.PATH_ROOT、util.PATH_DATA

v1.2.5: 增加计时器，修复InsertBatch()方法Bug

v1.2.4: 抽取http监听地址和端口、mysql连接参数，redis连接参数 等公用资源到配置文件

v1.2.3: 采用channel实现的多功能连接池，只要是tcp连接方式都可以使用，并在此基础上实现Mysql连接池，每秒处理并发读请求接近2W，写请求8K，有没有感觉，生成静态HTML，生成缓存已经意义不大了？呵呵

v1.2.2: 封装MyMysql实现的CRUD数据库类，GetOne()、GetAll()、Insert()、InsertBatch()、Update()，采用Map，可直接利用post、get方式获取的Map进行数据库操作，比一张表就对应写一个类的ORM更轻量，更实用，让繁琐的ORM见鬼去吧

v1.2.1: 实现基本类库编写: util.go、配置文件读取: app.ini

v1.2: github线上项目，实现 package 远程调用，只需短短几行代码，即可启用一个高并发Web服务，妈妈再也不用担心我不会写 Hello World 了

v1.1: 采用redigo，成功调用redis连接池，应用在线上日志系统统计，每秒处理并发请求2W多

v1.0: 本地的MVC框架，实现控制器

## Installation

    go get -u github.com/owner888/kaligo@v1.2.14

## Examples

### 新建项目

    $ cd $GOPATH/src/
    $ mkdir project
    $ cd project
    $ tree
        ├── config
        │   └── app.go
        │   └── database.go
        ├── control
        │   └── user.go
        ├── model
        │   └── user.go
        ├── data
        │   ├── cache
        │   └── log
        ├── static
        │   ├── css
        │   ├── images
        │   └── js
        └── template
            └── index.tpl
        └── main.go

### 配置项目 - config/database.go
    
```go
package config

import (
    "fmt"
    "github.com/owner888/kaligo/config"
)

func init() {
    config.Add("database", config.StrMap{
        "mysql": map[string]interface{}{
            "host":     "127.0.0.1",
            "port":     "3306",
            "name":     "test",
            "user":     "root",
            "pass":     "root",
            "charset":  "utf8mb4",
            "loc":      "Asia/Shanghai",
            "table_prefix" : "",
            "crypt_key"    : "",
            "crypt_fields" : map[string][]string{
                "user"  : {"name", "age"},
                "player": {"nickname"},
            },
            "check_privilege" : true,
            // 连接池配置
            "max_idle_connections": 300,
            "max_open_connections": 25,
            "max_life_seconds":     5*60,
        },
    })

    user := config.Get[string]("database.mysql.user")
    pass := config.Get[string]("database.mysql.pass")
    host := config.Get[string]("database.mysql.host")
    port := config.Get[string]("database.mysql.port")
    name := config.Get[string]("database.mysql.name")
    dsn  := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", user, pass, host+":"+port, name, "utf8mb4")
    config.Set("database.mysql.dsn", dsn)
}
```

### Example 1 - 路由设置

```go
// main.go
package main

import (
    "project/control"
    "github.com/owner888/kaligo"
    "github.com/owner888/project/controller"
)

func main() {
    app := &kaligo.App{}
    AddRoutes(app)
    kaligo.Run(app)
}

// AddRoutes is use for add route to App struct
func AddRoutes(app *kaligo.App) {
    app.AddStaticRoute("/static", "static")

    app.AddRoute("/", map[string]string{
        "GET": "Index",
    }, &controller.Get{})

    app.AddRoute("/posts/:post_id([0-9]+)", map[string]string{
        "POST": "Show",
    }, &controller.Post{})
}
```

### Example 2 - Controller 编写

```go
// controller/get.go
package control

import (
    "github.com/owner888/kaligo/controller"
)

// Get is use for
type Get struct {
    controller.Controller
}

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

// Index is a Get Index method
func (c *Get) Index() {
    c.JSON(200, H{"message": "hello world"})
}
```

### Example 3 - Model 编写

```go
// model/get.go
package model

func GetString() string {
    return "Hi"
}
```

#### 在控制器 controller/get.go 中使用model 

```go
func (c *Get) Index() {
    str := model.GetString()
    c.JSON(200, H{"message": str})
}
```

### Example 4 - View 编写

#### 模板文件：index.tpl

```html
<!DOCTYPE html>
<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
        <title>KaliGo</title>
    </head>
    <body>
        <h1>Hello Kali Go!</h1>
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
```

#### 在控制器ctl_index.go 中使用模板

    import (
        "net/http"
        "github.com/owner888/kaligo"
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

### Example 5 - 静态文件处理
    
    // 如何把刚才创建的应用如下目录当做静态来访问：
    ├── static
    │   ├── css
    │   ├── images
    │   └── js

    // 在项目初始文件 main.go 里面设置静态文件处理目录
    kaligo.SetStaticPath("/static", "static")

    这样用户访问 URL http://localhost/static/123.txt 则会请求 static 目录下的 123.txt 文件

### Example 6 - 数据库操作

#### 原生SQL查询
    
    KaliGo框架数据库底层驱动使用的是mymysql，利用autorc的断开重连机制，加上channel实现了连接池，相比普通短连接方式，性能提升10倍以上
    经过测试，数据库每秒可并发读取数据库2W次，写入数据库8K次，读取速度接近Redis+连接池方式，可以说有了这个框架，对Mysql读取出来的数据进行缓存毫无意义
    原生用法请参考：
    https://github.com/ziutek/mymysql    

#### CRUD 操作

##### Select
    
    // 从连接池获取一个连接，不需要Close，Sql执行完框架会自动回收连接到池里
    db, err := kaligo.MysqlConn.Get().(*kaligo.DB)

    // 单条记录查询,GetOne方法会自动给sql加上Limit 1
    row := db.GetOne("Select `name`, `pass` From `user`")
    // row的结构如下：
    //row := map[string]string {
    //    "name": "name111",
    //    "pass": "pass111",
    //}

    // 多条记录查询
    rows, err := db.GetAll("Select `name`, `pass` From `user`")
    // rows的结构如下：
    //rows := []map[string]string {
    //    map[string]string {
    //        "name": "name111",
    //        "pass": "pass111",
    //    },
    //    map[string]string {
    //        "name": "name222",
    //        "pass": "pass222",
    //    },
    //}

    // map方式存储的数据好处，就是不用每张表对应的写一个类与之对应，再配合Go强大的模板引擎，你懂的，easy到吓人

##### Insert

    // 单条插入，相当于 Insert Into `user`(`name`,`pass`) Values ("name111", "pass111")
    row := map[string]string {
        "name": "name111",
        "pass": "pass111",
    }
    if sql, err := db.Insert("user", row); err != nil {
        fmt.Println("错误信息：", err)
    } else {
        // 记录一下sql？
        logfile := util.PATH_DATA+"/log/user_"+time.Now().Format("2006-01-02")+".sql";
        util.PutFile(logfile, sql+"\n", 1)
        fmt.Println("自增ID：", db.InsertId())
    }

##### InsertBatch

    // 批量插入，相当于 Insert Into `user`(`name`,`pass`) Values ("name111", "pass111"),("name222", "pass222"),...
    rows := []map[string]string {
        map[string]string {
            "name": "name111",
            "pass": "pass111",
        },
        map[string]string {
            "name": "name222",
            "pass": "pass222",
        },
    }
    // InsertBatch(table string, data []map[string]string) (string, error)
    if sql, err := db.InsertBatch("user", rows); err != nil {
        fmt.Println("错误信息：", err)
    } else {
        //要不要把sql记录起来？
        fmt.Println("影响条数：", db.AffectedRows())
    }

##### Update

    // 单条修改，相当于 Update `user` Set `name`="name111",`pass`="pass111" Where `id`=10
    row := map[string]string {
        "name": "name111",
        "pass": "pass111",
    }
    if sql, err := db.Update("user", rows, "`id`=10"); err != nil {
        fmt.Println("错误信息：", err)
    } else {
        //要不要把sql记录起来？
        fmt.Println("影响条数：", db.AffectedRows())
    }

##### UpdateBatch

    // 批量修改，相当于 
    Update `user` Set 
    `pass` = Case  
    When `name` = 'name111' Then 'pass111' 
    When `name` = 'name222' Then 'pass222' 
    Else `plat_name` End, 
    `age` = Case  
    When `name` = 'name111' Then '11' 
    When `name` = 'name222' Then '22' 
    Else `channel` End 
    Where `plat_user_name` In ('name111', 'name222')

    rows := []map[string]string {
        map[string]string {
            "name": "name111",
            "pass": "pass111",
            "age": "11",
        },
        map[string]string {
            "name": "name222",
            "pass": "pass222",
            "age": "22",
        },
    }
    // UpdateBatch(table string, data []map[string]string, index string) (string, error)
    if sql, err := db.UpdateBatch("user", rows, "name"); err != nil {
        fmt.Println("错误信息：", err)
    } else {
        //要不要把sql记录起来？
        fmt.Println("影响条数：", db.AffectedRows())
    }

    
### Example 7 - 定时器

    // 增加定时任务，设置时间小于当前时间则不执行，大于当前时间则当到达时间时执行
    kaligo.AddTasker("default", &control.Task{}, "import_database", "2014-10-15 15:33:00")
    // 删除定时任务
    kaligo.DelTasker("default")
    // 增加定时器，每5秒运行一次
    kaligo.AddTimer("default", &control.Timer{}, "Import_database_login_v2", 5000)
    // 删除定时器
    kaligo.DelTimer("default")

## To do

1. http header处理、cookie、session、isAjax
2. 数据库分页
3. 测试cpu、内存使用率
4. 完善使用文档
5. 基于此框架开发一个CMS
6. 一个服务自定义绑定多个端口
7. 服务器硬件监控：CPU、内存、网卡带宽、硬盘、RAID
8. redis访问图表、mysql访问图表、请求统计信息、定时任务(跑一个gotinue会不会影响性能？)

## Documentation
* [control](http://www.godoc.org/pkg/github.com/owner888/kaligo/control)
* [db](http://www.godoc.org/pkg/github.com/owner888/kaligo/db)
* [redis](http://www.godoc.org/pkg/github.com/owner888/kaligo/redis)

## LICENSE

KaliGo is licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).
