package main

import (
	"log"

	"github.com/muhwyndhamhp/todo-mx/db"
)

func main() {
	log.Fatal(db.SeedDB())
}
