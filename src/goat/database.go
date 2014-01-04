package goat

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Connect to MySQL database
func DbConnect() (*sqlx.DB, error) {
	// Generate connection string using configuration
	conn := fmt.Sprintf("%s:%s@/%s", Static.Config.Db.Database, Static.Config.Db.Username, Static.Config.Db.Password)

	// Return connection and associated errors
	return sqlx.Connect("mysql", conn)
}
