package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
)

func simpleService() interface{} {
	resp, _ := http.Get("https://run.mocky.io/v3/6376ab35-cf8e-4cec-9aaf-f474fcb0f960")
	body, _ := io.ReadAll(resp.Body)
	// fmt.Printf("%s", body)
	return fmt.Sprintf("%s", body)
}

var db *sql.DB

func createConnection(dbName string) {
	// os.Remove("./app.db")

	var err error
	db, err = sql.Open("sqlite3", fmt.Sprintf("./%s.db", dbName))
	if err != nil {
		log.Fatal(err)
	}
	//
	// defer db.Close()
}

func createTable(name string) {
	sqlStmt := `
	create table ` + name + ` (id integer not null primary key, data text);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func storeToDb(tName string, data string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into " + tName + "(id, data) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	// for i := 0; i < 100; i++ {
	// 	_, err = stmt.Exec(i, fmt.Sprintf("text number %03d", i))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	_, err = stmt.Exec(uuid.New().String(), data)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func readOperationFromDb(tName string) interface{} {
	var id int
	var name string
	rows, err := db.Query("select id, data from " + tName)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return struct {
		id   int
		data string
	}{
		id:   id,
		data: name,
	}
}

func main() {
	a := simpleService()
	createConnection("app")
	createTable("userSession")
	storeToDb("userSession", a.(string))
	fmt.Println(readOperationFromDb("userSession"))
	defer db.Close()
}
