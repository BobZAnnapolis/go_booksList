package controllers

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "books-list/models"
    "books-list/repository/book"
    "github.com/gorilla/mux"
)

type Controller struct{}

var bookssql []models.Book
var bookRet  models.Book

// error logging
func logFatal(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func (c Controller) GetAllBooksSql(db *sql.DB) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s - was called- return ALL books\n\n", r.URL.Path)
    log.Println("getBooksSql - get all books is called")
    
    var book models.Book
    bookssql = []models.Book{}
    
    bookRepo := bookRepository{}
    bookssql = bookRepo.GetAllBooksSql(db, book, bookssql)
    
    // return the data
    json.NewEncoder(w).Encode(bookssql)
  }
}

func (c Controller) GetBookSql(db *sql.DB) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
     fmt.Fprintf(w, "%s - was called- return 1 book\n\n", r.URL.Path)
     log.Println("getBookSql - get 1 book is called")
 
     var book models.Book
 
     // use built-in func to extract URL vars
     params := mux.Vars(r)
     log.Println(params)

     id, err := strconv.Atoi(params["id"])
     logFatal(err)
 
     bookRepo := bookRepository.BookRepository{}
     bookRet = bookRepo.GetBookSql(db, book, id)    
 
     json.NewEncoder(w).Encode(bookRet)
  }
}
 
// use www.getpostman.com app to send in POST - UPDATE - DELETE requests for ADD, UPDATE, REMOVE handlers
//     POST FORMAT
//          {
//              "id": x,
//              "title": "...",
//              "author": "...",
//              "year": "..."
//          }
//
// after you post, should see a result in POSTMAN - also do a /booksall in local browser to verify new set of books after each
// or select * on elephantsql page
//
func (c Controller) AddBookSql(db *sql.DB) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
     fmt.Fprintf(w, "%s - was called- ADD 1 book\n\n", r.URL.Path)
     log.Println("addBookSql - add 1 book is called")
 
     var book models.Book
     var bookID int
 
     json.NewDecoder(r.Body).Decode(&book)
     log.Println(book)

     bookRepo := bookRepository.BookRepository{}
     bookID = bookRepo.AddBookSql(db, book) 
     
     // return ID to caller
     json.NewEncoder(w).Encode(bookID)
  }
}
 
func (c Controller) UpdateBookSql(db *sql.DB) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
     fmt.Fprintf(w, "%s - was called- UPDATE 1 book\n\n", r.URL.Path)
     log.Println("updateBookSql - update 1 book is called")
 
     var book models.Book

     json.NewDecoder(r.Body).Decode(&book)
 
     bookRepo := bookRepository.BookRepository{}
     rowsUpdated = bookRepo.UpdateBookSql(db, book) 
 
     // return updated set to caller
     json.NewEncoder(w).Encode(rowsUpdated)
  }
}
 
func (c Controller) RemoveBookSql(db *sql.DB) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s - was called- DELETE 1 book\n\n", r.URL.Path)
    log.Println("deleteBookSql - delete 1 book is called")
 
    // extract the id of the book to delete
    params := mux.Vars(r)    

    id, err := strconv.Atoi(params["id"])     
    logFatal(err)
 
    bookRepo := bookRepository.BookRepository{}
    rowsDeleted = bookRepo.RemoveBookSql(db, id) 

    // return updated set to caller
    json.NewEncoder(w).Encode(rowsDeleted)
  }
}
