package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"rest/internal/authorization"

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

type SignupForm struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"pwd" binding:"required"`
}

type LoginForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"pwd" binding:"required"`
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

//Page navigation functions

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

// Create User

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

// Registration
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func processor(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fname := r.FormValue("name")
	femail := r.FormValue("email")
	fpwd := r.FormValue("pwd")

	d := struct {
		Name  string
		Email string
		Pwd   string
	}{
		Name:  fname,
		Email: femail,
		Pwd:   fpwd + "Питер",
	}

	d.Pwd = base64.StdEncoding.EncodeToString([]byte(d.Pwd))

	db := setupDB()
	var lastInsertID int
	err := db.QueryRow("INSERT INTO users(name, email, password) VALUES($2, $3, $4) returning id;", d.Name, d.Email, d.Pwd).Scan(&lastInsertID)

	// check errors
	checkErr(err)

	tpl.ExecuteTemplate(w, "processor.html", d)
}

// Main function
func main() {
	// Init the mux router
	r := mux.NewRouter()
	authorization.Auth()
	// Registration
	r.HandleFunc("/", index)
	r.HandleFunc("/process", processor)

	// Route handles & endpoints

	// Get all books
	r.Handle("/books", authorization.CheckAuth(GetBooks)).Methods("GET")

	// Get book by id
	r.Handle("/books/{BookID}", authorization.CheckAuth(GetBook)).Methods("GET")

	// Create a book
	r.Handle("/books/", authorization.CheckAuth(CreateBook)).Methods("PUT")

	// Delete a specific book by the ID
	r.HandleFunc("/books/{BookID}", DeleteBook).Methods("DELETE")

	// serve the app
	fmt.Println("Server at 8000")
	fmt.Println("http://localhost:8000")
	http.ListenAndServe(":8000", r)
}
