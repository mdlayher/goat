package goat

import (
	"fmt"
	"log"

	// Bring in the MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// dbConnect connects to MySQL database
func dbConnect() (*sqlx.DB, error) {
	// Generate connection string using configuration
	conn := fmt.Sprintf("%s:%s@/%s", static.Config.DB.Database, static.Config.DB.Username, static.Config.DB.Password)

	// Return connection and associated errors
	return sqlx.Connect("mysql", conn)
}

// Attempt to "ping" the database to verify connectivity
func dbPing() bool {
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	db.Close()
	return true
}
