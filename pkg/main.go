package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	apiPath1 = "/apis/v1/books"
)

type library struct {
	dbHost, dbPass, dbName string
}

type Book struct {
	Id, Name, Isbn string
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == " " {
		dbHost = "localhost:3306"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == " " {
		dbPass = "Priya"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == " " {
		apiPath = apiPath1
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == " " {
		dbName = "library"
	}

	l := library{
		dbHost: dbHost,
		dbName: dbName,
		dbPass: dbPass,
	}
	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods(http.MethodGet)
	r.HandleFunc(apiPath, l.postBooks).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)

}

func (l library) postBooks(w http.ResponseWriter, r *http.Request)  {
	book := Book{}
	json.NewEncoder(w).Encode(&book)
	db := l.openConnection()
	insertQuery, err := db.Prepare("insert into book values ( ?,?,?)")
	if err !=nil{
		log.Fatalf("preparing the db query %s/n", err.Error())
	}

	tx,err := db.Begin()
	if err !=nil{
		log.Fatalf("while beginig the transaction %s/n", err.Error())

	}
	
	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Isbn)
	if err !=nil{
		log.Fatalf("executing the insrt command %s/n", err.Error())
	}
    

	err = tx.Commit()
	if err !=nil{
		log.Fatalf("while commiting the transaction %s/n", err.Error())
	}

	l.closeConnection(db)
}

	


func (l library) getBooks(w http.ResponseWriter, r *http.Request) {

	db := l.openConnection()
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("quering the books table %s/n", err.Error())
	}
	books := []Book{}
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("while scanning the row %s/n", err.Error())

		}
		aBook := Book{
			Id:   id,
			Name: name,
			Isbn: isbn,
		}
		books = append(books, aBook)
	}
	json.NewEncoder(w).Encode(books)
	l.closeConnection(db)
}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", l.dbPass, l.dbHost, l.dbName))
	if err != nil {
		log.Fatalf("opening to the connetion database %s/n", err.Error())
	}
	return db

}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("Closing Connection %s/n", err.Error())
	}

}
