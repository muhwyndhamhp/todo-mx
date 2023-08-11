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
	jsonb := typeext.JSONB{}
	u.Meta = jsonb
	return
}

func BuildTodoMeta(formPostPath string, todo *Todo) typeext.JSONB {
	a := typeext.JSONB{}
	a["Title"] = FormMeta{
		Label:    "Title",
		ID:       "todo-title",
		Name:     "title",
		FormPath: formPostPath,
		Value:    todo.Title,
	}
	a["Body"] = FormMeta{
		Label:    "Body",
		ID:       "todo-body",
		Name:     "body",
		FormPath: formPostPath,
		Value:    todo.Body.String,
	}

	return a
}
