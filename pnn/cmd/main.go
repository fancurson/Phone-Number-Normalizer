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
	user     = "postgres"
	password = "Buf_3red"
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
	must(createPhoneNumberTable(db))

	_, err = insertPhoneNumber(db, "1234567890")
	must(err)
	_, err = insertPhoneNumber(db, "123 456 7891")
	must(err)
	_, err = insertPhoneNumber(db, "(123) 456 7892")
	must(err)
	_, err = insertPhoneNumber(db, "(123) 456-7893")
	must(err)
	_, err = insertPhoneNumber(db, "123-456-7894")
	must(err)
	_, err = insertPhoneNumber(db, "123-456-7890")
	must(err)
	_, err = insertPhoneNumber(db, "1234567892")
	must(err)
	_, err = insertPhoneNumber(db, "(123)456-7892")
	must(err)

	phones, err := allPhones(db)
	must(err)
	for _, el := range phones {
		fmt.Printf("id: %d, number: %s", el.id, el.number)
	}
}

type phone struct {
	id     int
	number string
}

func createPhoneNumberTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL,
		value VARCHAR(255)
	)`
	_, err := db.Exec(statement)
	return err
}

func allPhones(db *sql.DB) ([]phone, error) {
	var ret []phone
	rows, err := db.Query("select id, value from phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p phone
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	err := db.QueryRow("SELECT value FROM phone_numbers WHERE id=$1", id).Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func insertPhoneNumber(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil

}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createBD(db, name)
}

func createBD(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
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
