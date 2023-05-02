package db

import (
	"database/sql"
	"fmt"
	"github.com/corerid/backend-test/config"
	_ "github.com/lib/pq"
)

type Connection struct {
	RW *sql.DB
}

func New(c config.Config) *Connection {

	return &Connection{
		RW: initDB(c.DBDialect, createDataSource(c)),
	}
}

func initDB(driverName string, dataSourceName string) *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	return db
}

func createDataSource(c config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
	)
}
