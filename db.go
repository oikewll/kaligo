/**
 * Realize a database operation class
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 */

package kaligo

import (
    "fmt"
    "log"
    "regexp"
    "strings"
    "strconv"
    "time"
    "sort"
    //"container/list"
    //"sync"
    //"reflect"
    "kaligo/conf"
    "kaligo/util"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"  // 空白导入必须在main.go、testing，否则就必须在这里写注释
    //"github.com/owner888/mymysql/autorc"
    //"github.com/owner888/mymysql/mysql"
    ////_ "github.com/ziutek/mymysql/native"    // 普通模式
    //_ "github.com/ziutek/mymysql/thrsafe" // 用了连接池之后连接都是重复利用的，没必要用线程安全模式
)

// DB is the struct for MySQL connection handler
type DB struct {
    logSlowQuery bool
    logSlowTime int64
    //Conn *autorc.Conn
    //res mysql.Result
    Conn *sql.DB    // MySQL connection
    rows sql.Rows
    res sql.Result
}

// New is the function for create an database operation handle
// (读+写)连接数据库+选择数据库
func New() (*DB, error){
    //fmt.Println("InitDB")
    host := conf.GetValue("db", "host")
    port := conf.GetValue("db", "port")
    user := conf.GetValue("db", "user")
    pass := conf.GetValue("db", "pass")
    name := conf.GetValue("db", "name")
    logSlowQuery, _ := strconv.ParseBool(conf.GetValue("db", "log_slow_query")) 
    logSlowTime, _ := strconv.ParseInt(conf.GetValue("db", "log_slow_time"), 0, 64) 

    db := &DB{
        logSlowQuery:logSlowQuery, 
        logSlowTime:logSlowTime,
    }

    connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		user,
		pass,
		host,
		port,
		name,
		"utf8mb4",
	)
	//fmt.Println(connection)
    //conn := autorc.New("tcp", "", address, user, pass, name) 
    //conn.Register("set names utf8")
    // 数据库的抽象，并不是真的数据库连接
    conn, err := sql.Open("mysql", connection)
    //checkErr(err)
    if err != nil {
        return db, nil
    }
    // 推出这个函数时不能关闭链接，否则其他调用的函数就无法执行 Query()、Exel() 方法 了
    //defer db.Close()

    // See "Important settings" section.
    conn.SetConnMaxLifetime(time.Minute * 3)
    conn.SetMaxOpenConns(10)    // 数据库的最大连接数
    conn.SetMaxIdleConns(10)    // 连接池中的保持连接的最大连接数

    // 初始化一个数据库连接，sql.Open 的时候实际上是返回一个链接对象而已，并没有真的和mysql链接上
    err = conn.Ping()
    if err != nil {
        fmt.Println("连接数据库失败", err.Error())
        return db, nil
    }
    conn.Query("set names utf8");

    db.Conn = conn
    return db, nil
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

// slowQueryLog is the function for record the slow query log
// 记录慢查询日志
func (db *DB) slowQueryLog(sql string, queryTime int64) {
    msg := fmt.Sprintf("Time: %d --- %s --- %s \n",
		queryTime,
		time.Now().Format("2006-01-02 15:04:05"),
		sql,
	)
    if ok, err := util.WriteLog("slow_query.log", msg); !ok {
        log.Print(err)
    }
}

// 记录慢查询日志
func (db *DB) errorSQLLog(sql string, err error) {
    msg := fmt.Sprintf("Time: %s --- %s --- %s \n",
		time.Now().Format("2006-01-02 15:04:05"),
		sql,
        err,
	)
    if ok, err := util.WriteLog("error_sql.log", msg); !ok {
        log.Print(err)
    }
}


// Escape ...: Escapes special characters in the txt, so it is safe to place returned string
// to Query method.
func (db *DB) Escape(txt string) string {
	//return db.Conn.Escape(my, txt)
    return txt
}

// Query is the function for query
// 执行一条语句(读 + 写)
//func (db *DB) Query(sql string) ([]mysql.Row, mysql.Result, error){
func (db *DB) Query(sql string) ([]sql.Row, sql.Result, error){
    startTime := time.Now().UnixNano()
    //rows, res, err := db.Conn.Query(sql) 
    rows, res, err := db.Query(sql) 
    if err != nil {
        db.errorSQLLog(sql, err)
    }
    //endTime := time.Now().Unix() - startTime
    //endTime := (time.Now().UnixNano() - 1412524713953787006) / 1000000000
    queryTime := (time.Now().UnixNano() - startTime) / 1000000000

    if queryTime > db.logSlowTime && db.logSlowQuery {
        db.slowQueryLog(sql, queryTime)
    }

    return rows, res, err
}

// 提取数据表字段名称
//func (db *DB) getFieldList(str string) ([]string) {
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

// GetOne is the function for get one record
// (读)直接从一个sql语句返回一条记录数据
func (db *DB) GetOne(sql string) (map[string]string, error) {
    // 判断SQL语句是否包含 Limit 1
    reg, _ := regexp.Compile(`(?i:limit)`)
    if (!reg.MatchString(sql)) {
        sql = strings.Trim(sql, " ") 
        reg, _ = regexp.Compile(`(?i:[,;])$`)
        sql = reg.ReplaceAllString(sql, "")
    }
    sql = fmt.Sprintf("%s Limit 1", sql)
    //fmt.Println(sql)

    results, err := db.GetAll(sql)
    //fmt.Println(results[0])
    return results[0], err
}

// GetAll is the function for get all record
// (读)直接从一个sql语句返回多条记录数据
func (db *DB) GetAll(sql string) (map[int]map[string]string, error) {

    // 最后得到的map
    results := make(map[int]map[string]string)

    //rows := db.Conn.QueryRow(sql) // 查询一条
    rows, err := db.Conn.Query(sql) // 查询多条
    if err != nil { 
        fmt.Println("查询数据库失败", err.Error())
        return results, err
    }

    // 读出查询出的列字段名
	cols, _ := rows.Columns()
	// values是每个列的值，这里获取到byte里
	values := make([][]byte, len(cols))
	// rows.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	// 让每一行数据都填充到[][]byte里面
	for i := range values {
		scans[i] = &values[i]
	}

    i := 0
    for rows.Next() { //循环，让游标往下推
        if err := rows.Scan(scans...); err != nil { //rows.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
            fmt.Println(err)
            return results, err
        }

        row := make(map[string]string) //每行数据

        for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
            key := cols[k]
            row[key] = string(v)
        }
        results[i] = row //装入结果集中
        i++
    }

    // 查询出来的数组
    //for k, v := range results {
        //fmt.Println(k, v)
    //}

    return results, err
}

// Insert is the function for insert data
// (写)拼凑一个sql语句插入一条记录数据
func (db *DB) Insert(table string, data map[string]string) (bool, error) {
    
    keys := []string{}
    vals := []string{}
    for k, v := range data {
        keys = append(keys, k)
        vals = append(vals, db.AddSlashes(db.StripSlashes(v)))
    }
    keysSQL := "`"+strings.Join(keys, "`, `")+"`"
    valsSQL := "\""+strings.Join(vals, "\", \"")+"\""
    sql := "Insert Into `"+table+"`("+keysSQL+") Values ("+valsSQL+")"
    //fmt.Println(sql)
    //_, res, err := db.Query(sql)
    res, err := db.Conn.Exec(sql)
    if err != nil {
        return false, err
    }

    db.res = res

    return true, err
}

// InsertBatch is the function for insert data in bulk
// (写)拼凑一个sql语句批量插入多条记录数据
func (db *DB) InsertBatch(table string, data []map[string]string) (bool, error) {

    var keys string
    var vals string
    var valsArr []string
    for _, d := range data {
        keys = ""
        vals = ""
        // slice是无序的，这里是保证他有顺序
        ms := NewSortMap(d)
        sort.Sort(ms)
        for k, v := range ms {
            if k == 0 {
                keys = v.Key
                vals = v.Val
            } else {
                keys = keys+"`,`"+v.Key
                vals = vals+"\",\""+db.AddSlashes(db.StripSlashes(v.Val))
            }
        }
        keys = "`"+keys+"`"
        vals = "(\""+vals+"\")"
        valsArr = append(valsArr, vals)
    }
    sql := "Insert Into `"+table+"`("+keys+") Values "+strings.Join(valsArr, ", ")
    //fmt.Println(sql)

    res, err := db.Conn.Exec(sql)
    if err != nil {
        return false, err
    }

    db.res = res

    return true, err
}

// Update is the function for update data
// (写)拼凑一个sql语句修改一条记录数据
func (db *DB) Update(table string, data map[string]string, where string) (bool, error) {
    
    sets := []string{}
    for k, v := range data {
        sets = append(sets, "`"+k+"`=\""+db.AddSlashes(db.StripSlashes(v))+"\"")
    }
    setsSQL := strings.Join(sets, ", ")
    sql := "Update `"+table+"` Set "+setsSQL+" Where "+where
    //fmt.Println(sql)
    res, err := db.Conn.Exec(sql)
    if err != nil {
        return false, err
    }

    db.res = res

    return true, err
}

// UpdateBatch is the function for update data in bulk
// (写)拼凑一个sql语句批量插入多条记录数据
func (db *DB) UpdateBatch(table string, data []map[string]string, index string) (bool, error) {

    sql := "Update `"+table+"` Set "
    ids := []string{}
    rows := map[string][]string {}

    // 下面两段是拆解过程
    //rows := map[string][]string {
        //"channel":[]string {
            //"When `plat_user_name` = 'test111' Then 'kkk5'",
            //"When `plat_user_name` = 'test222' Then '360'",
        //},
        //"plat_name":[]string {
            //"When `plat_user_name` = 'test111' Then 'kkk5_xxx'",
            //"When `plat_user_name` = 'test222' Then '360_xxx'",
        //},
    //}

    //rows["channel"] = []string{}
    //rows["channel"] = append(rows["channel"], "When `plat_user_name` = 'test111' Then 'kkk5'")
    //rows["channel"] = append(rows["channel"], "When `plat_user_name` = 'test222' Then '360'")
    //rows["plat_name"] = []string{}
    //rows["plat_name"] = append(rows["plat_name"], "When `plat_user_name` = 'test111' Then 'kkk5_xxx'")
    //rows["plat_name"] = append(rows["plat_name"], "When `plat_user_name` = 'test222' Then '360_xxx'")

    // 拼凑上面的Map结构出来
    for _, d := range data {
        ids = append(ids, d[index])
        for k, v := range d {
            if k != index {
                str := "When `"+index+"` = '" + d[index]+"' Then '"+v+"'"
                rows[k] = append(rows[k], str)
            }
        }
    }
    // 拼凑批量修改SQL语句
    for k, v := range rows {
        sql += "`"+k+"` = Case "
        for _, vv := range v {
            sql += " "+vv
        }
        sql += " Else `"+k+"` End, "
    }
    // 拼凑Where条件
    join := "'"+strings.Join(ids, "', '")+"'"
    where := " Where `"+index+"` In ("+join+")"
    // 完整的可执行SQL语句
    sql = util.Substr(sql, 0, len(sql)-2) + where

    res, err := db.Conn.Exec(sql)
    if err != nil {
        return false, err
    }

    db.res = res

    return true, err
}

// InsertID is the function for get last insert id
// 取得最后一次插入记录的自增ID值
func (db *DB) InsertID() int64 {
    id, _ := db.res.LastInsertId()
    return id
}

// AffectedRows is the function for return affected rows
// 返回受影响数目
func (db *DB) AffectedRows() int64 {
    rowsAffected, _ := db.res.RowsAffected()
    return rowsAffected
}

// Close is the function for close db connection
func (db *DB) Close() error {
    //err := db.Conn.Raw.Close()
    // 连接将会被释放回到连接池，而不是真的断开了链接
    err := db.Conn.Close()
    return err
}

// Transaction ...
type Transaction struct {
	*DB
}

// Begin is the function for close db connection
//func (db *DB) Begin() error {
    //tx, err := db.Conn.Begin()
    //return err
//}

// Commit is the function for close db connection
func (db *DB) Commit() error {
    err := db.Conn.Close()
    return err
}

// Rollback is the function for close db connection
func (db *DB) Rollback() error {
    err := db.Conn.Close()
    return err
}

// AddSlashes is ...
// 转义：引号、双引号添加反斜杠
func (db *DB) AddSlashes(val string) (string) {
    val = strings.Replace(val, "\"", "\\\"", -1)
    val = strings.Replace(val, "'", "\\'", -1)
    return val
}

// StripSlashes is ...
// 反转义：引号、双引号去除反斜杠
func (db *DB) StripSlashes(val string) (string) {
    val = strings.Replace(val, "\\\"", "\"", -1)
    val = strings.Replace(val, "\\'", "'", -1)
    return val
}

// GetSafeParam is ...
// 防止XSS跨站攻击
func (db *DB) GetSafeParam(val string) (string) {
    val = strings.Replace(val, "&", "&amp;", -1)
    val = strings.Replace(val, "<", "&lt;", -1)
    val = strings.Replace(val, ">", "&gt;", -1)
    val = strings.Replace(val, "\"", "&quot;", -1)
    val = strings.Replace(val, "'", "&#039;", -1)
    return val
}

