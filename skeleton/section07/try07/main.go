package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/tenntenn/sqlite"
)

type Record struct {
	Id          int
	Name        string
	PhoneNumber string
}

var tmpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
	<head><title>電話番号</title></head>
	<body>
		<form action="/add">
			<label for="">新規追加</input>
			<input type="text" name="name">
			<input type="text" name="phoneNumber">
			<input type="submit" value="追加">
		</form>
		<h1>電話番号一覧</h1>
		<ol>{{range .}}
		<li>{{.}}</li>
		{{end}}</ol>
	</body>
</html>`))

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	db, err := sql.Open(sqlite.DriverName, "record.db")
	if err != nil {
		return fmt.Errorf("DBのOpen: %w", err)
	}

	if err = createTable(db); err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		results, err := fetchRecords(db, 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return nil
}

func createTable(db *sql.DB) error {
	sqlStr := `CREATE TABLE IF NOT EXISTS record(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		phoneNumber TEXT NOT NULL
		);`
	_, err := db.Exec(sqlStr)
	if err != nil {
		return fmt.Errorf("テーブル作成時: %w", err)
	}
	return nil
}

func fetchRecords(db *sql.DB, limit int) ([]*Record, error) {
	sqlStr := `SELECT id, name, phoneNumber FROM record LIMIT ?`
	rows, err := db.Query(sqlStr, limit)
	if err != nil {
		return nil, fmt.Errorf("sqlの実行: %w", err)
	}
	defer rows.Close()
	var results []*Record

	for rows.Next() {
		var r Record
		err := rows.Scan(&r.Id, &r.Name, &r.PhoneNumber)
		if err != nil {
			return nil, fmt.Errorf("Scan:%w", err)
		}
		results = append(results, &r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("結果の取得:%w", err)
	}
	return results, nil
}
