package database

import (
    "time"
)

// Model a basic GoLang struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
// It may be embedded into your model or you may build your own model without it
//    type User struct {
//      mysql.Model
//    }
type Model struct {
    //DB        *DB
    ID        uint      `db:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt time.Time `db:"index"`
    CreatorID uint
    UpdatorID uint
    DeletorID uint
}

// Save is ...
func (model *Model) Save() bool{
    return true
}
