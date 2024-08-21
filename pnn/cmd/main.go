package main

import (
	"bytes"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "fancu"
	password = "your-password"
	dbname   = "phone"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user,
		password)
	db, err := sql.Open("postgres", psqlInfo)
	must(err)
	err = resetDB(db, dbname)
	must(err)
	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	must(db.Ping())
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXIST" + name)
	if err != nil {
		return nil
	}
	return createBD(db, name)
}

func createBD(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE" + name)
	if err != nil {
		return err
	}
	return nil
}

func normalizing(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
