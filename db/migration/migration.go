package main

import (
	"log"

	"github.com/muhwyndhamhp/todo-mx/db"
	"github.com/muhwyndhamhp/todo-mx/models"
)

func main() {
	db := db.GetDB()

	db.Debug()

	log.Fatal(db.AutoMigrate(&models.Todo{}))
}
