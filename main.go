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
	// We most close db session after programm shutdowning
	// defer db.Close()
}

func createTable(name string) {
	sqlStmt := `
	create table ` + name + ` (id VARCHAR(100) not null primary key, data TEXT);
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
	_, err = stmt.Exec(uuid.New().String(), data)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func readOperationFromDb(tName string) map[string]interface{} {
	var id string
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
	return map[string]interface{}{
		"id":   id,
		"data": name,
	}
}

func main() {
	a := simpleService()
	createConnection("app")
	createTable("userSession")
	storeToDb("userSession", a.(string))
	b := readOperationFromDb("userSession")
	fmt.Println(b["id"])
	fmt.Println(b["data"])
	defer db.Close()
}
