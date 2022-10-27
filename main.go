package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// DB set up
func setupDB() *sql.DB {
	db, _ := sql.Open("sqlite3", "./database/DB_Golang")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS books (id INTEGER PRIMARY KEY, title TEXT, firstname TEXT, lastname TEXT)")
	statement.Exec()
	return db
}

type Book struct {
	ID        int    `json:"BookID"`
	Title     string `json:"title"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type JsonResponse struct {
	Type    string `json:"type"`
	Data    []Book `json:"data"`
	Message string `json:"message"`
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Get all books

// response and request handlers
func GetBooks(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	printMessage("Getting books...")

	// Get all books from books table that don't have BookID = "1"
	rows, err := db.Query("SELECT * FROM books")

	// check errors
	checkErr(err)

	var books []Book

	// Foreach book
	for rows.Next() {
		var id int
		var title string
		var firstname string
		var lastname string
		err = rows.Scan(&id, &title, &firstname, &lastname)

		// check errors
		checkErr(err)

		books = append(books, Book{ID: id, Title: title, Firstname: firstname, Lastname: lastname})
	}

	if books != nil {
		var response = JsonResponse{Type: "success", Data: books}
		json.NewEncoder(w).Encode(response)
	} else {
		var response = JsonResponse{Type: "not found", Message: "No such data"}
		json.NewEncoder(w).Encode(response)
	}
}

// Get a book

// response and request handlers
func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params["BookID"]
	var response = JsonResponse{}

	if ID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing ID parameter."}
	} else {
		db := setupDB()
		rows, err := db.Query("SELECT * FROM books WHERE ID = $1", ID)
		var books []Book

		// Foreach movie
		for rows.Next() {
			var id int
			var title string
			var firstname string
			var lastname string
			err = rows.Scan(&id, &title, &firstname, &lastname)
			books = append(books, Book{ID: id, Title: title, Firstname: firstname, Lastname: lastname})
		}
		// check errors
		checkErr(err)

		if books != nil {
			response = JsonResponse{Type: "success", Data: books, Message: "Классная книга, советуем"}
			json.NewEncoder(w).Encode(response)
		} else {
			response = JsonResponse{Type: "not found", Message: "No such data"}
			json.NewEncoder(w).Encode(response)
		}
	}
}

// Create a book

// response and request handlers
func CreateBook(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	var response = JsonResponse{}
	db := setupDB()

	printMessage("Inserting book into DB")

	var lastInsertID int
	err := db.QueryRow("INSERT INTO books(title, firstname, lastname) VALUES($2, $3, $4) returning id;", title, firstname, lastname).Scan(&lastInsertID)
	printID := fmt.Sprintf("Inserting new book with ID: %v", lastInsertID)
	fmt.Println(printID)

	// check errors
	checkErr(err)

	response = JsonResponse{Type: "success", Message: "The movie has been inserted successfully!"}
	json.NewEncoder(w).Encode(response)
}

// Delete a book

// response and request handlers
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID := params["BookID"]

	var response = JsonResponse{}

	if ID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing ID parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting book from DB")

		_, err := db.Exec("DELETE FROM books WHERE ID = $1", ID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

// Main function
func main() {
	// Init the mux router
	router := mux.NewRouter()

	// Route handles & endpoints

	// Get all books
	router.HandleFunc("/books", GetBooks).Methods("GET")

	// Get book by id
	router.HandleFunc("/books/{BookID}", GetBook).Methods("GET")

	// Create a book
	router.HandleFunc("/books/", CreateBook).Methods("PUT")

	// Delete a specific book by the ID
	router.HandleFunc("/books/{BookID}", DeleteBook).Methods("DELETE")

	// serve the app
	fmt.Println("Server at 8080")
	fmt.Println("http://127.0.0.1:8000/books")
	log.Fatal(http.ListenAndServe(":8000", router))
}
