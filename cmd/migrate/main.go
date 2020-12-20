package main

import (
	"os"

	"github.com/mtavano/ghreviews/database"
	_ "github.com/mtavano/ghreviews/database/migrations"
)

func main() {
	// read env vars
	driver := os.Getenv("DATABASE_DRIVER")
	dsn := os.Getenv("DATABASE_URL")

	// database initialization
	db, err := database.NewStore(driver, dsn)
	check(err)

	err = db.Migrate()
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

