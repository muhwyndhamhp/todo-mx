package models

import (
	"html/template"

	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string
	Body        pgtype.Text
	EncodedBody template.HTML
}

func (Todo) TableName() string {
	return "todos"
}
