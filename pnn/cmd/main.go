package main

import (
	"bytes"
	"fmt"

	phonedb "github.com/fancurson/Phone-Number-Normalizer/db"
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
	must(phonedb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)

	must(phonedb.Migration("postgres", psqlInfo))
	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)
	defer phonedb.Close(db)

	err = db.Seed()
	must(err)

	phones, err := db.AllPhones()
	must(err)
	fmt.Printf("%+v\n\n", phones)

	for _, el := range phones {
		fmt.Printf("Working on...%+v", el)
		number := normalizing(el.Number)
		if el.Number != number {
			fmt.Println("Updating or removing...")
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				must(db.DeletePhone(el))
			} else {
				el.Number = number
				must(db.UpdatePhone(&el))
			}
		} else {
			fmt.Println("No changes required")
		}
	}

	phones, err = db.AllPhones()
	must(err)
	fmt.Printf("%+v\n\n", phones)

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
