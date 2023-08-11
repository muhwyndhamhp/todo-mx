package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/muhwyndhamhp/todo-mx/utils/errs"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(errs.Wrap(err))
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
