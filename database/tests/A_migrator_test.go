// 数据库工具类测试
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

// 如果存在 demo_user 数据表，则删除
func TestMigratorDropTable(t *testing.T) {
    err := db.Migrator().DropTable("demo_user")
    assert.NoError(t, err)

    err = db.Migrator().DropTable("demo_player")
    assert.NoError(t, err)
}

// 如果 demo_user 表不存在则创建表
func TestMigratorCreateTable(t *testing.T) {
    var err error

//     sqlScript := `CREATE TABLE user (
//   id int unsigned NOT NULL AUTO_INCREMENT,
//   username varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '账号',
//   password varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '密码',
//   realname blob COMMENT '真实姓名',
//   phone blob COMMENT '手机号码',
//   age int DEFAULT '0' COMMENT '年龄',
//   gender tinyint(1) DEFAULT '0' COMMENT '性别',
//   create_at datetime DEFAULT CURRENT_TIMESTAMP,
//   update_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//   delete_at datetime DEFAULT NULL,
//   PRIMARY KEY (id)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci`

    // 生成的 DSL 例子
    // create_table("my_guests", func(t) {
    //     t.Column("id", "uuid", {"primary": true})
    //     t.Column("firstname", "text", {})
    //     t.Column("lastname", "text", {})
    //     t.Column("email", "text", {})
    //     t.Column("reg_date", "timestamp", {})
    // })

    // 通过 Model 生成例子
    // type MyGuest struct {
    //     ID uuid.UUID `json:"id" db:"id"`
    //     CreatedAt time.Time `json:"created_at" db:"created_at"`
    //     UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
    //     Firstname string `json:"firstname" db:"firstname"`
    //     Lastname string `json:"lastname" db:"lastname"`
    //     Email string `json:"email" db:"email"`
    //     RegDate time.Time `json:"reg_date" db:"reg_date"`
    // }

    // 向上 / 向下迁移：
    // soda migrate up
    // soda migrate down

    // 用户表
    fields := []map[string]interface{}{
        {
            "name": "id",
            "type": "int",
            "unsigned": true,
            "notnull": true,
            "auto_increment": true,
            "comment": "ID",
        },
        {
            "name": "username",
            "type": "varchar",
            "constraint": 20,
            "comment": "账号",
        },
        {
            "name": "password",
            "type": "varchar",
            "constraint": 60,
            "comment": "密码",
        },
        {
            "name": "testname",
            "type": "varchar",
            "constraint": 20,
            "comment": "后面会被删除的测试字段",
        },
        {
            "name": "gender",
            "type": "tinyint",
            "constraint": 1,
            "default": "0",
            "comment": "性别",
        },
        {
            "name": "create_at",
            "type": "datetime",
            "default": "CURRENT_TIMESTAMP",
            "comment": "创建时间",
        },
        {
            "name": "update_at",
            "type": "datetime",
            "default": "CURRENT_TIMESTAMP",
            "extra": "ON UPDATE CURRENT_TIMESTAMP",
            "comment": "修改时间",
        },
        {
            "name": "delete_at",
            "type": "datetime",
            "comment": "删除时间",
        },
    }
    err = db.Migrator().CreateTable("demo_user", fields, []string{"id"})
    assert.NoError(t, err)

    // 玩家表
    fields = []map[string]interface{}{
        {
            "name": "id",
            "type": "int",
            "unsigned": true,
            "notnull": true,
            "auto_increment": true,
            "comment": "ID",
        },
        {
            "name": "uid",
            "type": "int",
            "default": "0",
            "comment": "用户ID",
        },
        {
            "name": "room_id",
            "type": "int",
            "default": "0",
            "comment": "房间ID",
        },
        {
            "name": "create_at",
            "type": "datetime",
            "default": "CURRENT_TIMESTAMP",
            "comment": "创建时间",
        },
        {
            "name": "update_at",
            "type": "datetime",
            "default": "CURRENT_TIMESTAMP",
            "extra": "ON UPDATE CURRENT_TIMESTAMP",
            "comment": "修改时间",
        },
        {
            "name": "delete_at",
            "type": "datetime",
            "comment": "删除时间",
        },
    }

    err = db.Migrator().CreateTable("demo_player", fields, []string{"id"})
    assert.NoError(t, err)
}

// 添加字段
func TestMigratorAddFields(t *testing.T) {
    err := db.Migrator().AddFields("demo_user", []map[string]interface{}{
        {
            "name": "realname",
            "type": "blob",
            "comment": "真实姓名",
        },
    })
    assert.NoError(t, err)
}

// 字段是否存在
func TestMigratorFieldExists(t *testing.T) {
    ok, err := db.Migrator().FieldExists("demo_user", "realname")
    assert.NoError(t, err)
    assert.True(t, ok)
}

// 修改表字段
func TestMigratorModifyFields(t *testing.T) {
    err := db.Migrator().ModifyFields("demo_user", []map[string]interface{}{
        {
            "oldname" : "testname",
            "name": "testname_tmp",
            "type": "varchar",
            "constraint": 50,
        },
    })
    assert.NoError(t, err)
}

// 删除表字段
func TestMigratorDropFields(t *testing.T) {
    err := db.Migrator().DropFields("demo_user", "testname_tmp")
    assert.NoError(t, err)
}

// 优化表
func TestMigratorOptimizeTable(t *testing.T) {
    // ok := db.Migrator().OptimizeTable("demo_user")
    // assert.True(t, ok)
}

// 检查表
func TestMigratorCheckTable(t *testing.T) {
    // ok := db.Migrator().CheckTable("demo_user")
    // assert.True(t, ok)
}

// 添加索引
func TestMigratorCreateIndex(t *testing.T) {
    err := db.Migrator().CreateIndex("demo_user", "username", "username_idx", "UNIQUE")
    assert.NoError(t, err)
}

// 重命名索引
func TestMigratorRenameIndex(t *testing.T) {
    err := db.Migrator().RenameIndex("demo_user", "username_idx", "username_idx_2")
    assert.NoError(t, err)
}

// 删除索引
func TestMigratorDropIndex(t *testing.T) {
    err := db.Migrator().DropIndex("demo_user", "username_idx_2")
    assert.NoError(t, err)
}

// 添加外建
func TestMigratorAddForeignKey(t *testing.T) {
    err := db.Migrator().AddForeignKey("demo_player", []map[string]interface{}{
        {
            "constraint": "fk_uid",
            "key": "uid",
            "reference": map[string]string {
                "table" : "demo_user",    // 要关联的表
                "column": "id",           // 要关联的表的字段
            },
            "on_update": "CASCADE",
            "on_delete": "RESTRICT",
        },
    })
    assert.NoError(t, err)
}

// 删除外键
func TestMigratorDropForeignKey(t *testing.T) {
    // err := db.Migrator().DropForeignKey("demo_user", "fk_uid")
    // assert.NoError(t, err)
}

// 重命名表明
func TestMigratorRenameTable(t *testing.T) {
    err := db.Migrator().RenameTable("demo_user", "demo_user_tmp")
    assert.NoError(t, err)

    err = db.Migrator().RenameTable("demo_user_tmp", "demo_user")
    assert.NoError(t, err)
}

// 清空表数据
func TestMigratorTruncateTable(t *testing.T) {
    err := db.Migrator().TruncateTable("demo_user")
    assert.NoError(t, err)
}

// 表是否存在
func TestMigratorTableExists(t *testing.T) {
    ok, err := db.Migrator().TableExists("demo_user")
    assert.NoError(t, err)
    assert.True(t, ok)
}

// 表列表
func TestMigratorListTables(t *testing.T) {
    tables, err := db.Migrator().ListTables("demo_user")
    assert.NoError(t, err)
    assert.NotNil(t, tables)
    // logs.Debug(database.FormatJSON(tables))
}

// 表字段列表
func TestMigratorListColumns(t *testing.T) {
    columns, err := db.Migrator().ListColumns("demo_user")
    assert.NoError(t, err)
    assert.NotNil(t, columns)
    // logs.Debug(database.FormatJSON(columns))
}

// 表索引列表
func TestMigratorListIndexes(t *testing.T) {
    indexes, err := db.Migrator().ListIndexes("demo_user")
    assert.NoError(t, err)
    assert.NotNil(t, indexes)
    // logs.Debug(database.FormatJSON(indexes))
}

// 删除表
func TestMigratorDropTable2(t *testing.T) {
    // 后面的 select_test.go、insert_test.go、update_test.go、delete_test.go 还需要用到
    // err := db.Migrator().DropTable("demo_user")
    // assert.NoError(t, err)
}

// 添加数据库
func TestMigratorCreateDatabase(t *testing.T) {
    err := db.Migrator().CreateDatabase("demo_schema")
    assert.NoError(t, err)
}

// 删除数据库
func TestMigratorDropDatabase(t *testing.T) {
    err := db.Migrator().DropDatabase("demo_schema")
    assert.NoError(t, err)
}
