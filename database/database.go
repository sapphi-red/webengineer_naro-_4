package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/srinathgs/mysqlstore"
	"log"
	"os"
)

var (
	db *sqlx.DB
)

func ConnectDB() *sqlx.DB {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	db = _db

	return _db
}

func CreateSessionStore() *mysqlstore.MySQLStore {
	store, err := mysqlstore.NewMySQLStoreFromConnection(
		db.DB,
		"sessions",
		"/",
		60*60*24*14,
		[]byte("secret-token"),
	)
	if err != nil {
		panic(err)
	}
	return store
}
