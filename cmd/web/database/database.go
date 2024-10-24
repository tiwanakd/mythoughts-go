package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func Open(dsn string) (Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return Database{}, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return Database{}, err
	}

	return Database{DB: db}, nil
}

func (d *Database) Close() {
	d.DB.Close()
}
