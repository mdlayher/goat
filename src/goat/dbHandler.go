package goat

import (
	"fmt"

	// Use the MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DBConnect connects to MySQL database
func DBConnect() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", fmt.Sprintf("%s:%s@/%s", "goat", "goat", "goat"))
}
