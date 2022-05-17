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
    Id        ID        `db:"id"`
    CreatorId ID        `db:"creator_id"`
    CreatedAt time.Time `db:"created_at"`
    UpdatorId ID        `db:"updator_id"`
    UpdatedAt time.Time `db:"updated_at"`
    DeletorId ID        `db:"deletor_id"`
    DeletedAt time.Time `db:"deleted_at"`
}
