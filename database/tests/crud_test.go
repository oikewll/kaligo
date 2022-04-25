package tests

import (
    "testing"
)

func TestCRUD(t *testing.T) {
    //db.Transaction(func(tx *database.DB) error {
        //// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
        //sqlStr := "insert into user(name, age, sex) values('test222', '30', '1')"
        ////_, err = db.Exec(sqlStr)
        //db.Query(sqlStr).Execute()
        //if err != nil {
            //t.Logf("%q: %s\n", err, sqlStr)
            //// 返回任何错误都会回滚事务
            //return err
        //}

        //t.Logf("RowsAffected = %d: %d\n", tx.RowsAffected, tx.LastInsertId)
        //// 返回 nil 提交事务
        //return nil
    //})


    // Test Rollback and Rollback
    //db.Begin()
    ////defer db.Rollback()
    //db.Insert("user", []string{"name", "age"}).Values([]string{"test111", "20"}).Execute()
    //db.Rollback()
    //db.Commit()

    // Test SavePoint and RollbackTo
    //db.Begin()
    //db.Insert("user", []string{"name", "age"}).Values([]string{"test111", "20"}).Execute()
    //db.SavePoint("sp1")
    //db.Insert("user", []string{"name", "age"}).Values([]string{"test222", "23"}).Execute()
    //db.RollbackTo("sp1")    // Rollback the user name is test222
    //db.Commit()  // Commit the user name is test111 
}
