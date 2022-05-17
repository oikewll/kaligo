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
    Id        ID        `json:"id"`
    CreatorId ID        `json:"creator_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatorId ID        `json:"updator_id"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletorId ID        `json:"deletor_id"`
    DeletedAt time.Time `json:"deleted_at"`
}
