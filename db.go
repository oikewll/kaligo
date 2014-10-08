package epooll

import (
    "fmt"
    "log"
    "regexp"
    "time"
    "strings"
    //"reflect"
    "github.com/owner888/epooll/util" 
    "github.com/ziutek/mymysql/mysql" 
    _ "github.com/ziutek/mymysql/native" // Native engine 
    // _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine 
)

type DB struct {
    logSlowQuery bool
    logSlowTime int64
    conn mysql.Conn
    row []mysql.Row
    res mysql.Result
    err error
}

// 程序一启动就会执行这一句，也就是初始化一个单例的数据库连接
var DBConn = InitDB()

// (读+写)连接数据库+选择数据库
func InitDB() *DB{
    fmt.Println("InitDB")
    conf := InitConfig()
    host := conf.GetValue("db", "host")+":"+conf.GetValue("db", "port")
	user := conf.GetValue("db", "user")
	pass := conf.GetValue("db", "pass")
	name := conf.GetValue("db", "name")

    db := new(DB)
    conn := mysql.New("tcp", "", host, user, pass, name) 

    err := conn.Connect() 
    if err != nil { 
        panic(err) 
    }

    conn.Query("set names utf8") 

    db.conn = conn
    return db
}

// 记录慢查询日志
func (this *DB) slowQueryLog(sql string, queryTime int64) {
    msg  := "Time: "+fmt.Sprintf("%s",queryTime)+" -- "+time.Now().Format("2006-01-02 15:04:05")+" -- "+sql+"\n";
    if ok, err := util.WriteLog("/data/golang/log/slow_query.log", msg); !ok {
        log.Print(err)
    }
}

// 执行一条语句(读 + 写)
func (this *DB) Query(sql string) ([]mysql.Row, mysql.Result, error){
    startTime := time.Now().UnixNano()
    rows, res, err := this.conn.Query(sql) 
    //endTime := time.Now().Unix() - startTime
    //endTime := (time.Now().UnixNano() - 1412524713953787006) / 1000000000
    queryTime := (time.Now().UnixNano() - startTime) / 1000000000

    if queryTime > this.logSlowTime && this.logSlowQuery {
        this.slowQueryLog(sql, queryTime)
    }

    return rows, res, err
}

// 提取数据表字段名称
//func (this *DB) getFieldList(str string) ([]string) {
    //reg, _ := regexp.Compile(`map\[(.*?)\]`)
    //arr := reg.FindAllString(str, 2)
    //str = fmt.Sprintf("%s", arr[1])
    //reg = regexp.MustCompile(`:%!s\(.*?\)`)
    //str = reg.ReplaceAllString(str, "")
    //str = strings.Replace(str, "map[", "", 1)
    //str = strings.Replace(str, "]", "", 1)
    //fieldList := strings.Split(str, " ")
    //return fieldList
//}

// (读)直接从一个sql语句返回一条记录数据
func (this *DB) GetOne(sql string) (map[string]string, error) {

    reg, _ := regexp.Compile(`(?i:limit)`)
    if (!reg.MatchString(sql)) {
        sql = strings.Trim(sql, " ") 
        reg, _ = regexp.Compile(`(?i:[,;])$`)
        sql = reg.ReplaceAllString(sql, "")
    }
    sql = fmt.Sprintf("%s Limit 1", sql)

    rows, res, err := this.Query(sql)

    //fields := this.getFieldList(fmt.Sprintf("%s", res))
    fields := []string{}
    for _, field := range res.Fields() {
        fields = append(fields, field.Name)
	}

    if err != nil { 
        panic(err) 
    } 

    results := map[string]string {}
    for _, row := range rows { 
        for _, field := range fields {
            results[field] = row.Str(res.Map(field)) 
        }
    }

    return results, err
}

// (读)直接从一个sql语句返回多条记录数据
func (this *DB) GetAll(sql string) ([]map[string]string, error) {

    rows, res, err := this.Query(sql)

    //fields := this.getFieldList(fmt.Sprintf("%s", res))
    fields := []string{}
    for _, field := range res.Fields() {
        fields = append(fields, field.Name)
	}

    if err != nil { 
        panic(err) 
    } 

    results := []map[string]string {}
    for _, row := range rows { 
        v := map[string]string {}
        for _, field := range fields {
            v[field] = row.Str(res.Map(field)) 
        }
        results = append(results, v)
    }

    return results, err
}

// 转义：引号、双引号添加反斜杠
func (this *DB) AddSlashes(val string) (string) {
    val = strings.Replace(val, "\"", "\\\"", -1)
    val = strings.Replace(val, "'", "\\'", -1)
    return val
}

// 反转义：引号、双引号去除反斜杠
func (this *DB) StripSlashes(val string) (string) {
    val = strings.Replace(val, "\\\"", "\"", -1)
    val = strings.Replace(val, "\\'", "'", -1)
    return val
}

// 防止XSS跨站攻击
func (this *DB) GetSafeParam(val string) (string) {
    val = strings.Replace(val, "&", "&amp;", -1)
    val = strings.Replace(val, "<", "&lt;", -1)
    val = strings.Replace(val, ">", "&gt;", -1)
    val = strings.Replace(val, "\"", "&quot;", -1)
    val = strings.Replace(val, "'", "&#039;", -1)
    return val
}

// (写)拼凑一个sql语句插入一条记录数据
func (this *DB) Insert(table string, data map[string]string) (bool, error) {
    
    keys := []string{}
    vals := []string{}
    for k, v := range data {
        keys = append(keys, k)
        vals = append(vals, this.AddSlashes(this.StripSlashes(v)))
    }
    keys_sql := "`"+strings.Join(keys, "`, `")+"`"
    vals_sql := "\""+strings.Join(vals, "\", \"")+"\""
    sql := "Insert Into `"+table+"`("+keys_sql+") Values ("+vals_sql+")"
    _, this.res, this.err = this.Query(sql)
    var ok bool = true
    if this.err != nil {
        ok = false
    }
    return ok, this.err
}

// (写)拼凑一个sql语句修改一条记录数据
func (this *DB) Update(table string, data map[string]string, where string) (bool, error) {
    
    sets := []string{}
    for k, v := range data {
        sets = append(sets, "`"+k+"`=\""+this.AddSlashes(this.StripSlashes(v))+"\"")
    }
    sets_sql := strings.Join(sets, ", ")
    sql := "Update `"+table+"` Set "+sets_sql+" "+where
    _, this.res, this.err = this.Query(sql)
    var ok bool = true
    if this.err != nil {
        ok = false
    }
    return ok, this.err
}

// 取得最后一次插入记录的自增ID值
func (this *DB) InsertId() uint64 {
    return this.res.InsertId()
}

// 返回受影响数目
func (this *DB) AffectedRows() uint64 {
    return this.res.AffectedRows()
}
