package models

import (
	"database/sql"
	"html/template"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/muhwyndhamhp/todo-mx/db"
	"github.com/pressly/goose"
	"gorm.io/gorm"
)

func init() {
	goose.AddMigration(Up001, Down001)
}

type Todo struct {
	gorm.Model
	Title       string
	Body        pgtype.Text
	EncodedBody template.HTML
}

func (Todo) TableName() string {
	return "todos"
}

func Up001(tx *sql.Tx) error {
	return db.GetDB().Migrator().CreateTable(&Todo{})
}

func Down001(tx *sql.Tx) error {
	return db.GetDB().Migrator().DropTable(&Todo{})
}
