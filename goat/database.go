package goat

import (
	"fmt"

	// Bring in the MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DBConnect connects to MySQL database
func DBConnect() (*sqlx.DB, error) {
	// Generate connection string using configuration
	conn := fmt.Sprintf("%s:%s@/%s", Static.Config.DB.Database, Static.Config.DB.Username, Static.Config.DB.Password)

	// Return connection and associated errors
	return sqlx.Connect("mysql", conn)
}
