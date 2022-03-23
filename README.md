## KaliGo v1.2.4 (2014-10-12)

KaliGo is an open-source, high-performance, modularity, full-stack web framework.

More info [doc.kaligo.com](https://doc.kaligo.com)

* [Changelog](https://github.com/owner888/kaligo/blob/master/README.md#changelog)
* [Installation](https://github.com/owner888/kaligo/blob/master/README.md#installation)
* [Testing](https://github.com/owner888/kaligo/blob/master/README.md#testing)
* [Examples](https://github.com/owner888/kaligo/blob/master/README.md#examples)
* [To do](https://github.com/owner888/kaligo/blob/master/README.md#to-do)
* [Documentation](https://github.com/owner888/kaligo/blob/master/README.md#documentation)

## Changelog

- v1.2.9: 重构config类，支持范型、支持点式获取方式: `config.Get[string]("database.mysql.charset")`
- v1.2.8: 重构DB类，和 kaliphp 操作保持一致
- v1.2.7: 定时器增加全局锁，修复同时设置多个定时器时协程并发写Timer map出现空指针异常Bug
- v1.2.6: 迁移util类到util namespace
- v1.2.5: 增加计时器，修复InsertBatch()方法Bug
- v1.2.4: 抽取http监听地址和端口、mysql连接参数，redis连接参数 等公用资源到配置文件
- v1.2.3: channel实现连接池，tcp连接方式都可使用，在此基础上实现Mysql连接池，每秒处理并发读请求接近2W，写请求8K
- v1.2.2: 封装MyMysql实现的CRUD数据库类，GetOne()、GetAll()、Insert()、InsertBatch()、Update()
- v1.2.1: 实现基本类库编写: util.go、配置文件读取: app.ini
- v1.2.0: github线上项目，实现 package 远程调用，只需短短几行代码，即可启用一个高并发Web服
- v1.1.0: 采用redigo，成功调用redis连接池，应用在线上日志系统统计，每秒处理并发请求2W多
- v1.0.0: 本地的MVC框架，实现控制器

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
        "mysql": map[string]any{
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
    "github.com/owner888/kaligo"
    "github.com/owner888/project/controller"
)

func main() {
    mux := kaligo.New()
    AddRoutes(mux)
    kaligo.Run(mux)
}

// AddRoutes is use for add route to Router struct
func AddRoutes(router *routes.Router) {
    router.AddStaticRoute("/static", "static")

    // RESTy routes for "articles" resource
    router.AddRoute("/article", map[string]string{
        "POST": "CreateArticle",
    }, &controller.ArticleController{})
    router.AddRoute("/article/search", map[string]string{
        "GET": "searchArticles",
    }, &controller.ArticleController{})

    // Subrouters:
    router.AddRoute("/article/:article_id([0-9]+)", map[string]string{
        "GET": "GetArticle",
        "PUT": "UpdateArticle",
        "DELETE": "DeleteArticle",
    }, &controller.ArticleController{})
}
```

### Example 2 - Controller 编写

```go
// controller/get.go
package control

import (
    "github.com/owner888/kaligo"
)

// Get is use for
type ArticleController struct {
    kaligo.Controller
}

type H map[string]any

func (c *ArticleController) Index() {
    // 通过范型获取路由数据
    articleID := Get[int](c.Params, "article_id", 20)
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
import (
    "github.com/owner888/project/model"
)

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

        args := map[string]any{
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

### Example 5 - 数据库操作

#### 原生SQL查询
    
    KaliGo框架数据库底层驱动使用的是mymysql，利用autorc的断开重连机制，加上channel实现了连接池，相比普通短连接方式，性能提升10倍以上
    经过测试，数据库每秒可并发读取数据库2W次，写入数据库8K次，读取速度接近Redis+连接池方式，可以说有了这个框架，对Mysql读取出来的数据进行缓存毫无意义
    原生用法请参考：
    https://github.com/ziutek/mymysql    

#### CRUD 操作

##### Select
```go
type User struct {
    ID   uint   `db:"id"`
    Name string `db:"name"`
    Age  uint   `db:"age"`
}
// 单条记录查询
var user User
db.Select("id", "name", "age").From("user").Where("id", "=", 1).Scan(&user).Execute()

// 多条记录查询
var users []User
db.Select("id", "name", "age").From("user").Limit(5).Scan(&users).Execute()
```

##### Insert

```go
// 单条插入，相当于 Insert Into `user`(`name`,`age`) Values ("test111", "20")
db.Insert("user", []string{"name", "age"}).Values([]string{"test111", "20"}).Execute()

// 批量插入，相当于 Insert Into `user`(`name`,`age`) Values ("test111", "20"),("test222", "25"),...
db.Insert("user", []string{"name", "age"}).Values([][]string{{"test111", "20"}, {"test222", "25"}}).Execute()
```

##### Update

```go
// 单条修改，相当于 Update `user` Set `name`="demo111", `age`="20" Where `id`=10
sets := map[string]string{"name":"demo111", "age":20}
db.Update("user").Set(sets).Where("id", "=", "1").Execute()
```

##### Delete

```go
db.Delete("user").Where("name", "=", "test").Execute()
```

##### UpdateBatch

```go
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
```
    
### Example 6 - 定时器

```go
// 增加定时任务，设置时间小于当前时间则不执行，大于当前时间则当到达时间时执行
mux.AddTasker("import_database", "2021-03-05 20:08:00", "ImportDatabase", &controller.Get{})
// 删除定时任务
mux.DelTasker("import_database")
// 增加定时器，每5秒运行一次
routkalier.AddTimer("import_database", 5000, "ImportDatabaseLoginV2", &controller.Get{})
// 删除定时器
mux.DelTimer("import_database")
```

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
