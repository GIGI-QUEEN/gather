package sqlite

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type DBS struct {
	Conn *sql.DB
}

var DB = &sql.DB{}

const dbName = "social-network.db"

func DataBase() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	DB = db

	dbs := DBS{Conn: db}

	driver, err := sqlite3.WithInstance(dbs.Conn, &sqlite3.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://./pkg/db/migrations", "sqlite3", driver)
	if err != nil {
		return err
	}

	// run the "up" command to apply all migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func dbExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}