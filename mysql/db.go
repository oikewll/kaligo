/**
 * Realize a database operation class
 *
 * @copyright   (C) 2014  seatle
 * @lastmodify  2021-07-06
 *
 */

package mysql

import (
    "fmt"
    "log"
    "regexp"
    "strings"
    "time"
    "sort"
    //"container/list"
    //"sync"
    //"reflect"
    "github.com/owner888/kaligo"
    "kaligo/conf"
    "kaligo/util"
    "database/sql"
    //_ "github.com/go-sql-driver/mysql"  // 空白导入必须在main.go、testing，否则就必须在这里写注释
    //"github.com/owner888/mymysql/autorc"
    //"github.com/owner888/mymysql/mysql"
    ////_ "github.com/ziutek/mymysql/native"    // 普通模式
    //_ "github.com/ziutek/mymysql/thrsafe" // 用了连接池之后连接都是重复利用的，没必要用线程安全模式
)

// DB is the struct for MySQL connection handler
type DB struct {
    Conn *Conn    // MySQL connection
    //rows Rows

    ////Conn *autorc.Conn
    ////res mysql.Result
    //Conn *sql.DB    // MySQL connection
    //rows sql.Rows
    res sql.Result
    initCmds []string   // MySQL commands/queries executed after connect
    logSlowQuery bool
    logSlowTime int64
}


// New is the function for Create new MySQL handler. 
// (读+写)连接数据库+选择数据库
func New() *DB {
    c := NewConn()
    db := DB{
        Conn: c,
        logSlowQuery: kaligo.StrToBool(conf.GetValue("db", "log_slow_query")),
        logSlowTime:  kaligo.StrToInt64(conf.GetValue("db", "log_slow_time")),
    }

    return &db
    //host := conf.GetValue("db", "host")
    //port := conf.GetValue("db", "port")
    //user := conf.GetValue("db", "user")
    //pass := conf.GetValue("db", "pass")
    //name := conf.GetValue("db", "name")
    //logSlowQuery := kaligo.StrToBool( conf.GetValue("db", "log_slow_query"))
    //logSlowTime  := kaligo.StrToInt64(conf.GetValue("db", "log_slow_time"))
    //maxOpenConns := kaligo.StrToInt(  conf.GetValue("db", "max_open_conns"))
    //maxIdleConns := kaligo.StrToInt(  conf.GetValue("db", "max_idle_conns"))

    //db := &DB{
        //logSlowQuery:logSlowQuery, 
        //logSlowTime:logSlowTime,
    //}
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

// 记录错误查询日志
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

// Register registers initialization commands.
// This is workaround, see http://codereview.appspot.com/5706047
func (db *DB) Register(query string) {
	db.initCmds = append(db.initCmds, query)
}

// Query is the function for query
// 执行一条语句(读 + 写)
//func (db *DB) Query(sql string) ([]mysql.Row, mysql.Result, error) {
func (db *DB) Query(sql string) (*sql.Rows, error) {
    startTime := time.Now().UnixNano()
    rows, err := db.Conn.Query(sql) 
    if err != nil {
        db.errorSQLLog(sql, err)
    }
    queryTime := (time.Now().UnixNano() - startTime) / 1000000000

    if queryTime > db.logSlowTime && db.logSlowQuery {
        db.slowQueryLog(sql, queryTime)
    }

    return rows, err
}

//// 提取数据表字段名称
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
func (db *DB) GetOne(sql string) (row map[string]string, err error) {
    // 判断SQL语句是否包含 Limit 1
    reg, _ := regexp.Compile(`(?i:limit)`)
    if (!reg.MatchString(sql)) {
        sql = strings.Trim(sql, " ") 
        reg, _ = regexp.Compile(`(?i:[,;])$`)
        sql = reg.ReplaceAllString(sql, "")
    }
    sql = fmt.Sprintf("%s Limit 1", sql)
    //fmt.Println(sql)

    // 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	//err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
    rows, err := db.GetAll(sql)

    if _, ok := rows[0]; ok {
        row = rows[0]
    }

    //fmt.Println(row)
    return row, err
}

// GetAll is the function for get all record
// (读)直接从一个sql语句返回多条记录数据
func (db *DB) GetAll(sql string) (row map[int]map[string]string, err error) {

    // 最后得到的map
    results := make(map[int]map[string]string)

    //row := db.Conn.QueryRow(sql)  // 查询一条，因为不存在Columns()方法，所以统一用Query吧
    rows, err := db.Conn.Query(sql) // 查询多条
    if err != nil { 
        fmt.Println("查询数据库失败", err.Error())
        return results, err
    }

    // 非常重要：关闭rows释放持有的数据库链接
    defer rows.Close()

    // 读出查询出的列字段名
	cols, _ := rows.Columns()
	// vals是每个列的值，这里获取到byte里
	vals := make([][]byte, len(cols))
	// rows.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	// 让每一行数据都填充到[][]byte里面
	for i := range vals {
		scans[i] = &vals[i]
	}

    i := 0
    // 循环读取结果集中的数据
    for rows.Next() { //循环，让游标往下推
        if err := rows.Scan(scans...); err != nil { //rows.Scan查询出来的不定长值放到scans[i] = &vals[i],也就是每行都放在vals里
            fmt.Println(err)
            return results, err
        }

        row := make(map[string]string) //每行数据

        for k, v := range vals { //每行数据是放在values里面，现在把它挪到row里
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
        ms := util.NewSortMap(d)
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
func (db *DB) Close() (err error) {
    if db.Conn.netConn == nil {
        return nil  // closed before
    }

    // 连接将会被释放回到连接池，而不是真的断开了链接
    err = db.Conn.Close()
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
    err := db.Close()
    return err
}

// Rollback is the function for close db connection
func (db *DB) Rollback() error {
    err := db.Close()
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

