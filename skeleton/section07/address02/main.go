package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/tenntenn/sqlite"
)

type Record struct {
	id    int
	name  string
	phone string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR", err)
	}
}

func run() error {
	db, err := sql.Open(sqlite.DriverName, "address.db")
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	if err := createTable(db); err != nil {
		return err
	}
	for {
		if err := showRecord(db); err != nil {
			return err
		}
		if err := insertRecord(db); err != nil {
			return err
		}
	}
	return nil
}

func createTable(db *sql.DB) error {
	sqlStr := `CREATE TABLE IF NOT EXISTS address (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		phone TEXT NOT NULL
	);`
	if _, err := db.Exec(sqlStr); err != nil {
		return fmt.Errorf("create db: %w", err)
	}
	return nil
}

func showRecord(db *sql.DB) error {
	sqlStr := `SELECT * FROM address;`
	rows, err := db.Query(sqlStr)
	if err != nil {
		return fmt.Errorf("select record: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.id, &r.name, &r.phone); err != nil {
			return fmt.Errorf("scan: %w", err)
		}
		fmt.Printf("[%d] Name: %s Phone: %s \n", r.id, r.name, r.phone)
		fmt.Println("----------")
	}
	return nil
}

func insertRecord(db *sql.DB) error {
	sqlStr := `INSERT INTO address (name, phone) values (?,?);`
	var r Record
	fmt.Println("Name >")
	fmt.Scan(&r.name)
	fmt.Println("Phone >")
	fmt.Scan(&r.phone)
	if _, err := db.Exec(sqlStr, r.name, r.phone); err != nil {
		return fmt.Errorf("insert record: %w", err)
	}
	return nil
}
