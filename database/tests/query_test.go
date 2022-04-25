package tests

import (
    "os"
    "testing"

    "github.com/owner888/kaligo/logs"
    "github.com/owner888/kaligo/database"
    mysql "github.com/owner888/kaligo/database/driver/mysql"
    "github.com/stretchr/testify/assert"
)

var db *database.DB

func TestMain(m *testing.M) {
    db, _ = database.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4"))
    code := m.Run()
    os.Exit(code)
}

func TestMigratorCurrentDatabase(t *testing.T) {
    databases := db.Migrator().CurrentDatabase()
    assert.Equal(t, "test", databases)
}

func TestMigratorListDatabases(t *testing.T) {
    databases, err := db.Migrator().ListDatabases("test")
    assert.NoError(t, err)
    assert.Equal(t, []string{"test"}, databases)
}

func TestMigratorListTables(t *testing.T) {
    tables, err := db.Migrator().ListTables("user")
    assert.NoError(t, err)
    assert.Equal(t, tables, []string{"user"})
    // logs.Debug(database.FormatJSON(tables))
}

func TestMigratorListColumns(t *testing.T) {
    columns, err := db.Migrator().ListColumns("user")
    assert.NoError(t, err)
    assert.Equal(t, columns, []database.Column([]database.Column{{Name:"ID", DBName:"id", DataType:"int", Size:0, Precision:0, NotNull:true, DefaultValue:"", Unique:false, PrimaryKey:true, AutoIncrement:true, Comment:"", Readable:true, Creatable:true, Updatable:true, Extra:"auto_increment", CryptKey:""}, {Name:"Username", DBName:"username", DataType:"string", Size:7, Precision:0, NotNull:false, DefaultValue:"", Unique:true, PrimaryKey:false, AutoIncrement:false, Comment:"", Readable:true, Creatable:true, Updatable:true, Extra:"", CryptKey:""}, {Name:"Password", DBName:"password", DataType:"string", Size:12, Precision:0, NotNull:false, DefaultValue:"mr.", Unique:false, PrimaryKey:false, AutoIncrement:false, Comment:"", Readable:true, Creatable:true, Updatable:true, Extra:"", CryptKey:""}, {Name:"Addtime", DBName:"addtime", DataType:"int", Size:0, Precision:0, NotNull:false, DefaultValue:"", Unique:false, PrimaryKey:false, AutoIncrement:false, Comment:"", Readable:true, Creatable:true, Updatable:true, Extra:"", CryptKey:""}}))
    // logs.Debug(database.FormatJSON(tables))
}

func TestUpdate(t *testing.T) {
    q := db.Insert("keywords").Columns([]string{`word`, `creator`}).Values([]string{"电影网站", "1"}).OnDuplicateKeyUpdate(map[string]string{`creator`: "3"})
    logs.Info(q)
}

func TestConfig(t *testing.T) {
    cfg := mysql.NewConfig()
    cfg.Addr    = "localhost:3308"
    cfg.DBName  = "test"
    cfg.User    = "root"
    cfg.Passwd  = "pw"
    dsn := cfg.FormatDSN()
    assert.Equal(t, "root:pw@tcp(localhost:3308)/test", dsn)
}
