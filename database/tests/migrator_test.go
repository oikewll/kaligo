package tests

import (
    "testing"

    // "github.com/owner888/kaligo/logs"
    // "github.com/owner888/kaligo/database"
    "github.com/stretchr/testify/assert"
)

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
    assert.NotNil(t, columns)
    // logs.Debug(database.FormatJSON(columns))
}

func TestMigratorListIndexes(t *testing.T) {
    indexes, err := db.Migrator().ListIndexes("user")
    assert.NoError(t, err)
    assert.NotNil(t, indexes)
    // logs.Debug(database.FormatJSON(indexes))
}

func TestMigratorCreateDatabase(t *testing.T) {
    err := db.Migrator().CreateDatabase("demo2")
    assert.NoError(t, err)
}

func TestMigratorDropDatabase(t *testing.T) {
    err := db.Migrator().DropDatabase("demo2")
    assert.NoError(t, err)
}

func TestMigratorCreateTable(t *testing.T) {
    var ok bool    
    var err error

    fields := []map[string]interface{}{
        {
            "name": "id",
            "type": "int",
            "constraint": 11,
            "notnull": true,
            "auto_increment": true,
        },
        {
            "name": "username",
            "type": "varchar",
            "constraint": 50,
        },
        {
            "name": "password",
            "type": "varchar",
            "constraint": 50,
            "default": "mr.",
        },
    }
    err = db.Migrator().CreateTable("migrator_user", fields, []string{"id"})
    assert.NoError(t, err)

    // 添加一个字段
    err = db.Migrator().AddFields("migrator_user", []map[string]interface{}{
        {
            "name": "testfield",
            "type": "varchar",
            "constraint": 50,
        },
    })

    ok, err = db.Migrator().FieldExists("migrator_user", "testfield")
    assert.NoError(t, err)
    assert.True(t, ok)

    err = db.Migrator().ModifyFields("migrator_user", []map[string]interface{}{
        {
            "oldname" : "testfield",
            "name": "testfield111",
            "type": "varchar",
            "constraint": 50,
        },
    })
    assert.NoError(t, err)

    err = db.Migrator().DropFields("migrator_user", "testfield111")
    assert.NoError(t, err)

    // ok = db.Migrator().OptimizeTable("migrator_user")
    // assert.True(t, ok)
    //
    // ok = db.Migrator().CheckTable("migrator_user")
    // assert.True(t, ok)

    // err = db.Migrator().CreateIndex("migrator_user", "username", "name_idx", "UNIQUE")
    // assert.NoError(t, err)
    // err = db.Migrator().RenameIndex("migrator_user", "name_idx", "name_idx222")
    // assert.NoError(t, err)
    // err = db.Migrator().DropIndex("migrator_user", "name_idx")
    // assert.NoError(t, err)
    //
    // err = db.Migrator().AddForeignKey("migrator_user", []map[string]interface{}{
    //     {
    //         "constraint": "fk_uid",
    //         "key": "uid",
    //         "reference": map[string]string {
    //             "table" : "user",   // 要关联的表
    //             "column": "uid",    // 要关联的表的字段
    //         },
    //         "on_update": "CASCADE",
    //         "on_delete": "RESTRICT",
    //     },
    // })
    // assert.NoError(t, err)
    //
    // err = db.Migrator().DropForeignKey("migrator_user", "fk_uid")
    // assert.NoError(t, err)

    err = db.Migrator().RenameTable("migrator_user", "migrator_user_tmp")
    assert.NoError(t, err)

    err = db.Migrator().TruncateTable("migrator_user_tmp")
    assert.NoError(t, err)

    ok, err = db.Migrator().TableExists("migrator_user_tmp")
    assert.NoError(t, err)
    assert.True(t, ok)

    err = db.Migrator().DropTable("migrator_user_tmp")
    assert.NoError(t, err)
}

