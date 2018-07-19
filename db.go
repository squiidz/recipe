package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// DB is just a sqlite db wrapper
type DB struct {
	*sql.DB
}

func newDB(dbName string) *DB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	createTableStmt := `
		create table translation (search text, value text);
		create table product (id integer not null primary key, name text, price real, quantity text);
	`

	db.Exec(createTableStmt)
	return &DB{db}
}

func (db *DB) findLike(name string) (*Product, error) {
	stmt, err := db.Prepare("SELECT * FROM product WHERE name LIKE ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var product Product
	err = stmt.QueryRow(name).Scan(&product.ID, &product.Name, &product.Price, &product.Quantity)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
