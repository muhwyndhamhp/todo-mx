package models

import (
	"html/template"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/muhwyndhamhp/todo-mx/utils/typeext"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Title       string
	Body        pgtype.Text
	EncodedBody template.HTML
	Meta        typeext.JSONB `gorm:"type:jsonb"`
}

func (Todo) TableName() string {
	return "todos"
}

func (u *Todo) BeforeSave(tx *gorm.DB) (err error) {
	a := typeext.JSONB{}

	a["title_label"] = "Title"
	a["title_id"] = "todo-title"
	a["title_name"] = "title"

	u.Meta = a

	return
}
