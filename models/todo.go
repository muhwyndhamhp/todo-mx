package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title string
	Body  pgtype.Text
}

func (Todo) TableName() string {
	return "todos"
}
