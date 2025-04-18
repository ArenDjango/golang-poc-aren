package model

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int64  `bun:",pk,autoincrement"`
	Name          string
	Email         string
}
