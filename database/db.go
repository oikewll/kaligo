package database

import (
     "context"
     "database/sql"
     "errors"
     "fmt"
     "regexp"
     //"reflect"
     "strings"
     "sync"
     "time"
     "github.com/owner888/kaligo/config"
     //"github.com/owner888/kaligo/util"
     //"github.com/stretchr/testify/assert"
 )

//type singleton struct {
//}

//var instance *singleton
//var once sync.Once

//func GetInstance() *singleton {
    //once.Do(func() {
        //instance = &singleton{}
    //})
    //return instance
//}

var (
    once sync.Once
    instance *DB
    instances map[string]*DB
)

//var instances = map[string]*DB{}

// QueryType is ...
type QueryType int64

const (
    // SELECT Query select type
    SELECT QueryType = 1
    // INSERT Query insert type
    INSERT QueryType = 2
    // UPDATE Query update type
    UPDATE QueryType = 3
    // DELETE Query delete type
    DELETE QueryType = 4
)

// DB is the struct for MySQL connection handler
type DB struct {
    Error         error         // Global error
    RowsAffected  int64         // For select、update、insert
    LastInsertId  int64         // Only for insert
    InTransaction bool          // 是否正在事务执行中
    query         *Query
    schema        *Schema

    Name          string        // Instance name
    Dialector                   // Dialector database dialector
    StdDB         *sql.DB       // Connection
    StdTx         *sql.Tx       // Connection for Transaction
    initCmds      []string      // SQL commands/queries executed after connect

    timeout       time.Duration // Timeout for connect SetConnMaxLifetime(timeout)
    lastUse       time.Time     // The last use time

    debug         bool          // Debug logging. You may change it at any time.
    logSlowQuery  bool          // 是否记录慢查询
    logSlowTime   int           // 慢查询时长

    queryCount    int           // 执行过多少条SQL

    // Logger
    //Logger logger.Interface

    // NowFunc the function to be used when creating a new timestamp
    NowFunc func() time.Time

    cacheStore   *sync.Map
}

// Open initialize db session based on dialector
// (读+写)连接数据库+选择数据库
func Open(dialector Dialector) (db *DB, err error) {

    db = &DB{InTransaction: false}

    //if db.Logger == nil {
        //db.Logger = logger.Default
    //}

    if db.NowFunc == nil {
        db.NowFunc = func() time.Time { return time.Now().Local() }
    }

    if dialector != nil {
		db.Dialector = dialector
	}

    if db.cacheStore == nil {
		db.cacheStore = &sync.Map{}
	}

    if db.Dialector != nil {
        err = db.Dialector.Initialize(db)
    }

    db.query = &Query{
        DB          : db,
		StdDB       : db.StdDB, // 因为返回的是指针*sql.DB，所以 db.StdDB 和 db.conn.StdDB 是同一个，一个Close()，另一个也会Close()
        tablePrefix : config.Get[string]("database.mysql.table_prefix"),
		Context     : context.Background(),
    }

    // 设置最大初始连接数
    if config.Get[int]("database.mysql.max_open_connections") > 0 {
        db.StdDB.SetMaxOpenConns(config.Get[int]("database.mysql.max_open_connections"))
    }
    // 设置最大空闲连接数
    if config.Get[int]("database.mysql.max_idle_connections") > 0 {
        db.StdDB.SetMaxIdleConns(config.Get[int]("database.mysql.max_idle_connections"))
    }
    // sql.Open 实际上返回了一个数据库抽象，并没有真的连接上
    if err == nil {
        // ping 调用完毕后会马上把连接返回给连接池
        err = db.StdDB.Ping()
    }

    db.setCharset("utf8")

    db.schema = &Schema{
        Name  : db.Name,
        Query : db.query,
    }

    //if err != nil {
        //db.Logger.Error(context.Background(), "failed to initialize database, got error %v", err)
    //}

    return
}

// AddError add error to db
func (db *DB) AddError(err error) error {
	if db.Error == nil {
		db.Error = err
	} else if err != nil {
		db.Error = fmt.Errorf("%v; %w", db.Error, err)
	}
	return db.Error
}

// DB returns `*sql.DB`
func (db *DB) DB() *sql.DB { return db.StdDB }
// Tx returns `*sql.Tx`
func (db *DB) Tx() *sql.Tx { return db.StdTx }

// Debug start debug mode
func (db *DB) Debug() { db.debug = true }

// Set store value with key into current db instance's context
func (db *DB) Set(key string, value interface{}) *DB {
	db.cacheStore.Store(key, value)
	return db
}

// Get get value with key from current db instance's context
func (db *DB) Get(key string) (interface{}, bool) {
	return db.cacheStore.Load(key)
}

// InstanceSet store value with key into current db instance's context
// db.InstanceSet("kalidb:started_transaction", true)
func (db *DB) InstanceSet(key string, value interface{}) *DB {
    // %p 获取指针地址, ep:[0x140001341b0]
	db.cacheStore.Store(fmt.Sprintf("%p", db) + key, value)
	return db
}

// InstanceGet get value with key from current db instance's context
// if _, ok := db.InstanceGet("kalidb:started_transaction"); ok {
func (db *DB) InstanceGet(key string) (interface{}, bool) {
	return db.cacheStore.Load(fmt.Sprintf("%p", db) + key)
}

// Set the charset
func (db *DB) setCharset(charset string) {
    if db.Dialector.Name() == "mysql" {
        db.StdDB.Query("SET NAMES " + db.Quote(charset));
    } else if db.Dialector.Name() == "sqlite" {
        db.StdDB.Query("PRAGMA encoding = " + db.Quote(charset));
    }
}

// Model specify the model you would like to run db operations
//    // update all users's name to `hello`
//    db.Model(&User{}).Update("name", "hello")
//    // if user's primary key is non-blank, will use it as condition, then will only update the user's name to `hello`
//    db.Model(&user).Update("name", "hello")
//func (db *DB) Model(value interface{}) *Query {
    //db.query = &Query{
        //Model     : value,
        //sqlStr    : "",
        //queryType : 0,
        //DB        : db,
        //StdDB     : db.StdDB,
    //}
    //return db.query
//}

// Query func is use for create a new [*Query]
// Query -> Connection.Query
//     Create a new SELECT query
//     Query("SELECT * FROM users")
//     Create a new DELETE query
//     Query("DELETE FROM users WHERE id = 5")
// @param sqlStr string  SQL statement
// @param queryType int  type:SELECT, UPDATE, etc
// @return *Query
func (db *DB) Query(sqlStr string, args ...QueryType) *Query {
    var queryType QueryType
    if len(args) == 0 {
        queryType = 0
    } else {
        queryType = args[0]
    }

    // 生成一个新的 Query 对象，一个SQL一个 Query 对象
    //q := new(Query)
    db.query = &Query{
        sqlStr    : sqlStr,
        queryType : queryType,
        DB        : db,
        StdDB     : db.StdDB,
    }

    return db.query
}

// Select func is use for create a new [*Select]
// Select -> Where -> Builder -> Query
//     SELECT `id`, `name``
//     Select("id", "username")
//     Select([]string{"id", "username"})
//     SELECT id AS user_id
//     select("id AS user_id")
//     SELECT `id`, `name` FROM `user`
//     Select("id", "name").From("user")
//     SELECT `id`, `name` FROM `user` LEFT JOIN `player` ON `user`.`uid`=`player`.`uid` WHERE `player`.`room_id`="10"
//     Select("id", "name").From("user").Join("player", "LEFT").On("user.uid", "=", "player.uid").Where("player.room_id", "=", "10")
// @param columns []string  columns to select
// @return *Query
func (db *DB) Select(columns ...string) *Query {
    db.query = &Query{
        S: &Select{
            selects : columns,
            distinct: false,
            offset  : 0,
        },
        W: &Where{},
        B: &Builder{},
        sqlStr      : "",
        queryType   : SELECT,
        cryptKey    : config.Get[string]("database.mysql.crypt_key"),
        cryptFields : config.CryptFields,
        DB          : db,
        StdDB       : db.StdDB,
    }
    return db.query
}

// Insert func is use for create a new [*Insert]
// Insert -> Builder -> Query
//     INSERT INTO `user` (`name`, `age`) VALUES ("test", "25")
//     Insert("user", []string{"name", "age"}).Values([]string{"test", "25"})
//     Insert("user", []string{"name", "age"}).Values([][]string{{"test", "25"}, {"demo", "30"}})
// @param table   string   table to insert into
// @param columns []string list of column names
// @return *Query
func (db *DB) Insert(table string, args ...[]string) *Query {
    var columns []string
    if len(args) != 0 {
        columns = args[0]
    }

    db.query = &Query{
        I: &Insert{
            table  : table,
            columns: columns,
        },
        //W: &Where{}, // Insert 暂时没有支持 Where 写法
        B: &Builder{},
        sqlStr      : "",
        queryType   : INSERT,
        cryptKey    : config.Get[string]("database.mysql.crypt_key"),
        cryptFields : config.CryptFields,
        DB          : db,
        StdDB       : db.StdDB,
    }
    return db.query
}

// Update func is use for create a new [*Update]
// Update -> Where -> Builder -> Query
//     UPDATE `user` SET `name`="test", `age`="25" WHERE `id`="1"
//     sets := map[string]string{"name":"demo", "age": "25"}
//     Update("user").Set(sets).Where("id", "=", "1")
//     Update("user").Join("player", "LEFT").On("user.uid", "=", "player.uid").Set(sets).Where("player.room_id", "=", "10")
// @param table   string    table to update
// @return *Query
func (db *DB) Update(table string) *Query {
    db.query = &Query{
        U: &Update{
            table : table,
        },
        W: &Where{},
        B: &Builder{},
        sqlStr      : "",
        queryType   : UPDATE,
        cryptKey    : config.Get[string]("database.mysql.crypt_key"),
        cryptFields : config.CryptFields,
        DB          : db,
        StdDB       : db.StdDB,
    }
    return db.query
}

// Delete func is use for create a new [*Delete]
// Delete -> Where -> Builder -> Query
//     DELETE FROM `user` WHERE `id`="1"
//     Delete("user").Where("id", "=", "1")
//     Delete("user").Join("player", "LEFT").On("user.uid", "=", "player.uid").Where("player.id", "=", "test")
// @param table   string    table to delete from
// @return *Query
func (db *DB) Delete(table string) *Query {
    db.query = &Query{
        D: &Delete{
            table : table,
        },
        W: &Where{},
        B: &Builder{},
        sqlStr    : "",
        queryType : DELETE,
        DB        : db,
        StdDB     : db.StdDB,
    }
    return db.query
}

// Schema Database schema operations
// CREATE DATABASE database CHARACTER SET utf-8 DEFAULT utf-8
// Schema.CreateDatabase(/*database*/ database, /*charset*/ 'utf-8', /*ifNotExists*/ true)
func (db *DB) Schema() *Schema {
    return db.schema

    //db.query = &Query{
        //Schema: &Schema{
            //Name : name,
        //},
        //W: &Where{},
        //B: &Builder{},
        //sqlStr    : "",
        //queryType : DELETE,
        //DB        : db,
        //StdDB     : db.StdDB,
    //}
    //return db.query
}

// Expr func is use for create a new [*Expression] which is not escaped. An expression
// is the only way to use SQL functions within query builders.
func (db *DB) Expr(value string) *Expression {
    return &Expression{
        value: value,
    }
}

// TablePrefix Return the table prefix defined in the current configuration.
func (db *DB) TablePrefix(table string) string {
    return db.query.tablePrefix + table
}

// Row is the function for query one row
// db.QueryRow() 调用完毕后会将连接传递给sql.Row类型
// 当.Scan()方法调用之后把连接释放回到连接池
func (db *DB) Row(sqlStr string, args ...interface{}) (row *sql.Row) {
    if db.InTransaction {
        row = db.StdTx.QueryRow(sqlStr, args...)
    } else {
        row = db.StdDB.QueryRow(sqlStr, args...)
    }
    return row
}

// Rows is the ...
// db.Query() 调用完毕后会将连接传递给sql.Rows类型
// 当然后者迭代完毕 或者 显性的调用.Close()方法后，连接将会被释放回到连接池
func (db *DB) Rows(sqlStr string, args ...interface{}) (rows *sql.Rows, err error) {
    if db.InTransaction {
        rows, err = db.StdTx.Query(sqlStr, args...)
    } else {
        rows, err = db.StdDB.Query(sqlStr, args...)
    }
    return
}

// Exec is the function for Insert、Update、Delete
// db.Exec() 调用完毕后会马上把连接返回给连接池
// 但是它返回的Result对象还保留这连接的引用，当后面的代码需要处理结果集的时候连接将会被重用
func (db *DB) Exec(sqlStr string, args ...interface{}) (res sql.Result, err error) {
    if db.InTransaction {
        res, err = db.StdTx.Exec(sqlStr, args...)
    } else {
        res, err = db.StdDB.Exec(sqlStr, args...)
    }
    return res, err
}

// Transaction start a transaction as a block, return error will rollback, otherwise to commit.
func (db *DB) Transaction(fc func(tx *DB) error) (err error) {
    panicked := true

    tx := db.Begin()

    defer func() {
        // Make sure to rollback when panic, Block error or Commit error
        if panicked || err != nil {
            tx.Rollback()
        }
    }()

    if err = tx.Error; err == nil {
        err = fc(tx)
    }

    if err == nil {
        err = tx.Commit().Error
    }

    panicked = false
    return
}

// Begin is the function for close db connection
// db.Begin() 调用完毕后将连接传递给sql.Tx类型对象
// 当.Commit()或.Rollback()方法调用后释放连接
func (db *DB) Begin() *DB {
    tx, err := db.StdDB.Begin()
    if err != nil {
        db.AddError(err)
    } else {
        db.StdTx = tx
        db.InTransaction = true
    }
    return db
}

// Commit is the function for close db connection
func (db *DB) Commit() *DB {
    if db.InTransaction == false || db.StdTx == nil {
        db.AddError(ErrInvalidTransaction)
    } else {
        db.InTransaction = false
        err := db.StdTx.Commit()
        if err != nil {
            db.AddError(err)
        }
    }
    return db
}

// Rollback is the function for close db connection
func (db *DB) Rollback() *DB {
    if db.InTransaction == false || db.StdTx == nil {
        db.AddError(ErrInvalidTransaction)
    } else {
        db.InTransaction = false
        err := db.StdTx.Rollback()
        if err != sql.ErrTxDone && err != nil {
            db.AddError(err)
        }
    }
    return db
}

// RollbackTo is the function for close db connection
func (db *DB) RollbackTo(name string) *DB {
    if db.InTransaction == false || db.StdTx == nil {
        db.AddError(ErrInvalidTransaction)
    } else {
        // 还没结束
        //db.InTransaction = false
        _, err := db.StdTx.Exec("ROLLBACK TO SAVEPOINT " + name)
        if err != sql.ErrTxDone && err != nil {
            db.AddError(err)
        }
    }
    return db
}

// SavePoint is the function for close db connection
func (db *DB) SavePoint(name string) *DB {
    if db.InTransaction == false || db.StdTx == nil {
        db.AddError(ErrInvalidTransaction)
    } else {
        // 还没结束
        //db.InTransaction = false
        _, err := db.StdTx.Exec("SAVEPOINT " + name)
        if err != sql.ErrTxDone && err != nil {
            db.AddError(err)
        }
    }
    return db
}

// Quote a value for an SQL query.
func (db *DB) Quote(values interface{}) string {
    switch vals := values.(type) {
    case string:
        return db.Escape(vals)
    case []string:
        for k, v := range vals {
            vals[k] = db.Quote(v)
        }
        return "(" + strings.Join(vals, ", ") + ")"
    case *Query:
        // Create a sub-query
        return "(" + vals.Compile() + ")"
    case *Expression:
        // Use a raw expression
        return vals.value
    default:
        return db.Escape(vals.(string))
    }
}

// QuoteTable Quote a database table name and adds the table prefix if needed.
// table = strings.Replace(table, "#DB#", "lrs", 1 )
// @param interface{} value table name or []string{"table", "alias"}
func (db *DB) QuoteTable(values interface{}) string {
    var table string    
    switch vals := values.(type) {
    case *Query:
        // Create a sub-query
        table = "(" + vals.Compile() + ")"
    case string:
        if strings.Index(vals, ".") == -1 {
            // Add the table prefix for tables
            table = db.QuoteIdentifier(db.TablePrefix(vals))
        } else {
            // table.alias 的写法，变成 `table`.`alias`
            parts := regexp.MustCompile(`\.`).Split(vals, 2)
            table = db.QuoteIdentifier(db.QuoteTable(parts[0]) ) + "." + db.QuoteIdentifier(parts[1])
        }
    case []string:
        // Separate the table and alias
        table := vals[0]
        alias := vals[1]
        table = db.QuoteIdentifier(table) + " AS " + db.QuoteIdentifier(alias)
    default:
        table = vals.(string)
    }

    return table
}

// QuoteIdentifier Quote a database identifier, such as a column name. Adds the
// table prefix to the identifier if a table name is present.
// table  ---> `table`
// column ---> `column`
func (db *DB) QuoteIdentifier(values interface{}) string {
    switch vals := values.(type) {
    case string:
        if vals == "*" || strings.Index(vals, "`") != -1 {
            // * 不需要变成 `*`，已经有 `` 包含着的直接返回
            return vals
        } else if strings.Index(vals, ".") != -1 {
            // table.column 的写法，变成 `table`.`column`
            parts := regexp.MustCompile(`\.`).Split(vals, 2)
            return db.QuoteIdentifier(db.QuoteTable(parts[0]) ) + "." + db.QuoteIdentifier(parts[1])
        } else {
            return "`" + vals + "`"
        }
    case []string:
        // Separate the column and alias
        value := vals[0]
        alias := vals[1]
        return db.QuoteIdentifier(value) + " AS " + db.QuoteIdentifier(alias)
    case *Query:
        // Create a sub-query
        return "(" + vals.Compile() + ")"
    case *Expression:
        // Use a raw expression
        return vals.value
    default:
        return vals.(string)
    }
}

// Escape is use for Escapes special characters in the txt, so it is safe to place returned string
func (db *DB) Escape(sql string) string {
    dest := make([]byte, 0, 2*len(sql))
    var escape byte
    for i := 0; i < len(sql); i++ {
        c := sql[i]

        escape = 0

        switch c {
        case 0: /* Must be escaped for 'mysql' */
            escape = '0'
            break
        case '\n': /* Must be escaped for logs */
            escape = 'n'
            break
        case '\r':
            escape = 'r'
            break
        case '\\':
            escape = '\\'
            break
        case '\'':
            escape = '\''
            break
        case '"': /* Better safe than sorry */
            escape = '"'
            break
        case '\032': //十进制26,八进制32,十六进制1a, /* This gives problems on Win32 */
            escape = 'Z'
        }

        if escape != 0 {
            dest = append(dest, '\\', escape)
        } else {
            dest = append(dest, c)
        }
    }

    // SQL standard is to use single-quotes for all values
    return "'" + string(dest) + "'"
}

// ProcessForeignKeys is 
func (db *DB) ProcessForeignKeys(foreignKeys []map[string]interface{}) string {
    var fkList []string    
    for _, definition := range foreignKeys {
        // some sanity checks
        if _, ok := definition["key"]; !ok {
            db.AddError(errors.New("Foreign keys on processForeignKeys() must specify a foreign key name"))
            return ""
        }

        if _, ok := definition["reference"]; !ok {
            db.AddError(errors.New("Foreign keys on processForeignKeys() must specify a foreign key reference"))
            return ""
        }

        reference := definition["reference"].(map[string]string)

        var referenceTable, referenceColumn string    
        if table, ok := reference["table"]; ok {
            referenceTable = table
        } else {
            db.AddError(errors.New("Foreign keys on processForeignKeys() must specify a reference table name"))
            return ""
        }

        if column, ok := reference["column"]; ok {
            referenceColumn = column
        } else {
            db.AddError(errors.New("Foreign keys on processForeignKeys() must specify a reference column name"))
            return ""
        }

        var sqlStr string    
        if table, ok := definition["constraint"]; ok {
            sqlStr += " CONSTRAINT " + db.QuoteTable(table.(string)) 
        }

        sqlStr += " FOREIGN KEY (" + db.QuoteIdentifier(definition["key"].(string)) + ")"
        sqlStr += " REFERENCES " + db.QuoteTable(referenceTable) + " ("
        referenceColumnArr := strings.Split(referenceColumn, ",")
        for k, v := range referenceColumnArr {
            referenceColumnArr[k] = db.QuoteIdentifier(v)
        }
        sqlStr += strings.Join(referenceColumnArr, ", ")
        sqlStr += ")"

        if val, ok := definition["on_update"]; ok {
            sqlStr += " ON UPDATE " + val.(string) 
        }
        if val, ok := definition["on_delete"]; ok {
            sqlStr += " ON DELETE " + val.(string)
        }

        fkList = append(fkList, "\n\t" + strings.TrimLeft(sqlStr, " "))
    }

    return ", " + strings.Join(fkList, ",")
}

// ProcessCharset is
func (db *DB) ProcessCharset(charset string, isDefault bool, args ...string) string {

    var collation string    
    if len(args) > 0 {
        collation = args[0]
    }

    // utf8_unicode_ci
    charsets := strings.Split(charset, "_")
    if collation == "" && len(charsets) > 1 {
        collation = charset     // utf8_unicode_ci
        charset   = charsets[0] // utf8
    }

    charset = " CHARACTER SET " + charset
    if isDefault {
        charset = " DEFAULT " + charset
    }

    if collation != "" {
        if isDefault {
            charset += " DEFAULT"
        }
        charset += " COLLATE " + collation
    }

    return charset
}

// Caching Per connection cache controller setter/getter
//func (c *Connection) Caching() bool { return false }

// slowQueryLog is the function for record the slow query log
// 记录慢查询日志
//func (db *DB) slowQueryLog(sql string, queryTime int64) {
    //msg := fmt.Sprintf("Time: %d --- %s --- %s \n",
		//queryTime,
		//time.Now().Format("2006-01-02 15:04:05"),
		//sql,
	//)
    //if ok, err := util.WriteLog("slow_query.log", msg); !ok {
        //log.Print(err)
    //}
//}

//// 记录错误查询日志
//func (db *DB) errorSQLLog(sql string, err error) {
    //msg := fmt.Sprintf("Time: %s --- %s --- %s \n",
		//time.Now().Format("2006-01-02 15:04:05"),
		//sql,
        //err,
	//)
    //if ok, err := util.WriteLog("error_sql.log", msg); !ok {
        //log.Print(err)
    //}
//}

// Query is the function for query
// 执行一条语句(读 + 写)
//func (db *DB) Query(sql string) ([]mysql.Row, mysql.Result, error) {
//func (db *DB) Query(sql string) (*sql.Rows, error) {
    //startTime := time.Now().UnixNano()
    //rows, err := db.connonn.Query(sql)
    //if err != nil {
        //db.errorSQLLog(sql, err)
    //}
    //queryTime := (time.Now().UnixNano() - startTime) / 1000000000

    //if queryTime > db.logSlowTime && db.logSlowQuery {
        //db.slowQueryLog(sql, queryTime)
    //}

    //return rows, err
//}

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
//func (db *DB) GetOne(sql string) (row map[string]string, err error) {
    //// 判断SQL语句是否包含 Limit 1
    //reg, _ := regexp.Compile(`(?i:limit)`)
    //if !reg.MatchString(sql) {
        //sql = strings.TrimSpace(sql)
        //reg, _ = regexp.Compile(`(?i:[,;])$`)
        //sql = reg.ReplaceAllString(sql, "")
    //}
    //sql = fmt.Sprintf("%s Limit 1", sql)
    ////fmt.Println(sql)

    //// 非常重要：确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	////err := db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)
    //rows, err := db.GetAll(sql)

    //if _, ok := rows[0]; ok {
        //row = rows[0]
    //}

    ////fmt.Println(row)
    //return row, err
//}

// GetAll is the function for get all record
// (读)直接从一个sql语句返回多条记录数据
//func (db *DB) GetAll(sql string) (row map[int]map[string]string, err error) {

    //// 最后得到的map
    //results := make(map[int] map[string]string)

    ////row := db.connonn.QueryRow(sql)  // 查询一条，因为不存在Columns()方法，所以统一用Query吧
    //rows, err := db.conn.Query(sql) // 查询多条
    //if err != nil {
        //fmt.Println("查询数据库失败", err.Error())
        //return results, err
    //}

    //// 非常重要：关闭rows释放持有的数据库链接
    //defer rows.Close()

    //// 读出查询出的列字段名
	//cols, _ := rows.Columns()
	//// vals是每个列的值，这里获取到byte里
	//vals := make([][]byte, len(cols))
	//// rows.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	//scans := make([]interface{}, len(cols))
	//// 让每一行数据都填充到[][]byte里面
	//for i := range vals {
		//scans[i] = &vals[i]
	//}

    //i := 0
    //// 循环读取结果集中的数据
    //for rows.Next() { //循环，让游标往下推
        //if err := rows.Scan(scans...); err != nil { //rows.Scan查询出来的不定长值放到scans[i] = &vals[i],也就是每行都放在vals里
            //fmt.Println(err)
            //return results, err
        //}

        //row := make(map[string]string) //每行数据

        //for k, v := range vals { //每行数据是放在values里面，现在把它挪到row里
            //key := cols[k]
            //row[key] = string(v)
        //}
        //results[i] = row //装入结果集中
        //i++
    //}

    //// 查询出来的数组
    ////for k, v := range results {
        ////fmt.Println(k, v)
    ////}

    //return results, err
//}

// Insert is the function for insert data
// (写)拼凑一个sql语句插入一条记录数据
//func (db *DB) Insert(table string, data map[string]string) (bool, error) {

    //var keys = []string{}
    //var vals = []string{}
    //for k, v := range data {
        //keys = append(keys, k)
        //vals = append(vals, db.AddSlashes(db.StripSlashes(v)))
    //}
    //keysSQL := "`"+strings.Join(keys, "`, `")+"`"
    //valsSQL := "\""+strings.Join(vals, "\", \"")+"\""
    //var sqlStr = "Insert Into `"+table+"`("+keysSQL+") Values ("+valsSQL+")"
    ////fmt.Println(sql)
    ////_, res, err := db.Query(sql)
    //res, err := db.connonn.Exec(sqlStr)
    //if err != nil {
        //return false, err
    //}

    //db.res = res

    //return true, err
//}

// InsertBatch is the function for insert data in bulk
// (写)拼凑一个sql语句批量插入多条记录数据
//func (db *DB) InsertBatch(table string, data []map[string]string) (bool, error) {

    //var keys string
    //var vals string
    //var valsArr []string
    //for _, d := range data {
        //keys = ""
        //vals = ""
        //// slice是无序的，这里是保证他有顺序
        //ms := util.NewSortMap(d)
        //sort.Sort(ms)
        //for k, v := range ms {
            //if k == 0 {
                //keys = v.Key
                //vals = v.Val
            //} else {
                //keys = keys+"`,`"+v.Key
                //vals = vals+"\",\""+db.AddSlashes(db.StripSlashes(v.Val))
            //}
        //}
        //keys = "`"+keys+"`"
        //vals = "(\""+vals+"\")"
        //valsArr = append(valsArr, vals)
    //}
    //var sqlStr = "Insert Into `"+table+"`("+keys+") Values "+strings.Join(valsArr, ", ")
    ////fmt.Println(sql)

    //res, err := db.conn.Exec(sqlStr)
    //if err != nil {
        //return false, err
    //}

    //db.res = res

    //return true, err
//}

// Update is the function for update data
// (写)拼凑一个sql语句修改一条记录数据
//func (db *DB) Update(table string, data map[string]string, where string) (bool, error) {

    //var sets []string
    //for k, v := range data {
        //sets = append(sets, "`"+k+"`=\""+db.AddSlashes(db.StripSlashes(v))+"\"")
    //}
    //setsSQL := strings.Join(sets, ", ")
    //var sqlStr = "Update `"+table+"` Set "+setsSQL+" Where "+where
    ////fmt.Println(sql)
    //res, err := db.connonn.Exec(sqlStr)
    //if err != nil {
        //return false, err
    //}

    //db.res = res

    //return true, err
//}

// UpdateBatch is the function for update data in bulk
// (写)拼凑一个sql语句批量插入多条记录数据
//func (db *DB) UpdateBatch(table string, data []map[string]string, index string) (bool, error) {

    //var sqlStr = "Update `"+table+"` Set "
    //ids := []string{}
    //rows := map[string][]string {}

    //// 下面两段是拆解过程
    ////rows := map[string][]string {
        ////"channel":[]string {
            ////"When `plat_user_name` = 'test111' Then 'kkk5'",
            ////"When `plat_user_name` = 'test222' Then '360'",
        ////},
        ////"plat_name":[]string {
            ////"When `plat_user_name` = 'test111' Then 'kkk5_xxx'",
            ////"When `plat_user_name` = 'test222' Then '360_xxx'",
        ////},
    ////}

    ////rows["channel"] = []string{}
    ////rows["channel"] = append(rows["channel"], "When `plat_user_name` = 'test111' Then 'kkk5'")
    ////rows["channel"] = append(rows["channel"], "When `plat_user_name` = 'test222' Then '360'")
    ////rows["plat_name"] = []string{}
    ////rows["plat_name"] = append(rows["plat_name"], "When `plat_user_name` = 'test111' Then 'kkk5_xxx'")
    ////rows["plat_name"] = append(rows["plat_name"], "When `plat_user_name` = 'test222' Then '360_xxx'")

    //// 拼凑上面的Map结构出来
    //for _, d := range data {
        //ids = append(ids, d[index])
        //for k, v := range d {
            //if k != index {
                //str := "When `"+index+"` = '" + d[index]+"' Then '"+v+"'"
                //rows[k] = append(rows[k], str)
            //}
        //}
    //}
    //// 拼凑批量修改SQL语句
    //for k, v := range rows {
        //sqlStr += "`"+k+"` = Case "
        //for _, vv := range v {
            //sqlStr += " "+vv
        //}
        //sqlStr += " Else `"+k+"` End, "
    //}
    //// 拼凑Where条件
    //join := "'"+strings.Join(ids, "', '")+"'"
    //where := " Where `"+index+"` In ("+join+")"
    //// 完整的可执行SQL语句
    //sqlStr = util.Substr(sqlStr, 0, len(sqlStr)-2) + where

    //res, err := db.conn.Exec(sqlStr)
    //if err != nil {
        //return false, err
    //}

    //db.res = res

    //return true, err
//}

//// InsertID is the function for get last insert id
//// 取得最后一次插入记录的自增ID值
//func (db *DB) InsertID() int64 {
    //id, _ := db.res.LastInsertId()
    //return id
//}

