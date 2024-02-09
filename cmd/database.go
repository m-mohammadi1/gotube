package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func openConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Application) getDBConnection() (*sql.DB, error) {
	// get dsn from app and open connection
	conn, err := openConnection(app.DSN)
	if err != nil {
		return nil, err
	}

	log.Printf("connected to db.")

	return conn, nil
}
