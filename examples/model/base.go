package model

import (
    "github.com/owner888/kaligo/database"
)

var (
    DB *database.DB
)

type ID int64

type Base struct {
    Id ID `json:"id"`
    // CreatorId ID
    // CreatedAt time.Time
    // UpdatorId ID
    // UpdatedAt time.Time
    // DeletorId ID
    // DeletedAt time.Time
}
