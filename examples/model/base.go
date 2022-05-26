package model

import (
    "time"

    "github.com/owner888/kaligo/database"
)

var (
    DB *database.DB
)

type ID int64

type Base struct {
    Id        ID        `db:"id"`          // 自增 ID
    CreatorId ID        `db:"creator_id" ` // 创建用户 ID
    CreatedAt time.Time `db:"created_at"`  // 创建时间
    UpdatorId ID        `db:"updator_id"`  // 更新用户 ID
    UpdatedAt time.Time `db:"updated_at"`  // 更新时间
    DeletorId ID        `db:"deletor_id"`  // 删除用户 ID
    DeletedAt time.Time `db:"deleted_at"`  // 删除时间
}
