package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/tenntenn/sqlite"
	"go.uber.org/multierr"
)

type Deposit struct {
	ID     int
	Name   string
	Amount int
}

var (
	user1 = Deposit{ID: 1, Name: "user1", Amount: 100}
	user2 = Deposit{ID: 1, Name: "user2", Amount: 20}
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error", err)
	}
}

func run() error {
	db, err := sql.Open(sqlite.DriverName, "deposit.db")
	if err != nil {
		return fmt.Errorf("db open: %w", err)
	}
	if err := createTable(db); err != nil {
		return err
	}
	if err := initTable(db, &user1, &user2); err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT * FROM deposit;`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var d Deposit
		rows.Scan(&d.ID, &d.Name, &d.Amount)
		fmt.Println(d)
	}
	rows.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin:%w", err)
	}

	row1 := tx.QueryRow(`SELECT amount FROM deposit WHERE id=1;`)
	row2 := tx.QueryRow(`SELECT amount FROM deposit WHERE id=2;`)
	var am1, am2 int
	if err := row1.Scan(&am1); err != nil {
		tx.Rollback()
		return err
	}
	if err := row2.Scan(&am2); err != nil {
		tx.Rollback()
		return err
	}
	updateSql := `UPDATE deposit SET amount=? WHERE id=?;`
	if _, err := tx.Exec(updateSql, am1-10, 1); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(updateSql, am2+10, 2); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	rows, err = db.Query(`SELECT * FROM deposit;`)
	if err != nil {
		return err
	}
	fmt.Println(rows)

	return nil
}

func createTable(db *sql.DB) error {
	sqlStr := `CREATE TABLE IF NOT EXISTS deposit (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		amount INTEGER NOT NULL
	);`
	if _, err := db.Exec(sqlStr); err != nil {
		return fmt.Errorf("create db: %w", err)
	}
	return nil
}

func initTable(db *sql.DB, depo ...*Deposit) error {
	sqlStr := `INSERT INTO deposit (name, amount) values (?, ?);`
	var rerr error
	for _, v := range depo {
		_, err := db.Exec(sqlStr, v.Name, v.Amount)
		if err != nil {
			rerr = multierr.Append(rerr, fmt.Errorf("table init: %w", err))
		}
	}
	return rerr
}
