package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

func init() {
	// load env file
	if err := godotenv.Load(); err != nil {
		log.Println("env file not found")
	}
}

func Connection() (db *sql.DB) {
	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")

	db, err := sql.Open(dbDriver, dbUser+":"+"@/"+dbName+"?"+"parseTime=true")

	if err != nil {
		panic(err)
	}

	return db
}
