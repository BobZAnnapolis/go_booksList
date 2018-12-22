package main

// pre-requisites
//    go get -u github.com/gorilla/mux  - http server libs
//    go get github.com/subosito/gotenv - environment variables, for db access
//    go get github.com/lib/pq          - Go postgres driver for elephantsql db
//
//    uses www.getpostman.com to test out the ADD, UPDATE and DELETE api requests
//    uses www.elephantsql.com to use a cloud-based POSTGRESQL db for persistent data tests
//         create table books (id serial, title varchar, author varchar, year varchar);
//         insert into books (title, author, year) values('darya', 'estonia', '2015');
//                insert a bunch of these to create some records to work with
//         select * from books; <- verify you have records
//

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
//    "reflect"
    "strconv"

    "books-list/controllers"
    "books-list/driver"
    "books-list/models"
    "database/sql"
    "github.com/gorilla/mux"
    "github.com/subosito/gotenv"
)

// create an array of book structs; our elephantSQL cloud db connection
var bookssql []models.Book
var bookslocal []models.Book
var db *sql.DB

// error logging
func logFatal(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

// use the environment file to access elephant sql db
func init() {
    gotenv.Load()
}

func main() {

// set up connection to elephantsql cloud postgresql db
    db = driver.ConnectDB()

// HARD-CODED : create some books to access for the app
   bookslocal = append(bookslocal, models.Book{ID: 1, Title: "Golang pointers", Author: "Mr. Golang", Year: "2014"},
       models.Book{ID: 2, Title: "Goroutines", Author: "Ms. Goroutine", Year: "2015"},
       models.Book{ID: 3, Title: "Weekend Devil", Author: "Mistress", Year: "2016"},
       models.Book{ID: 4, Title: "Virus Infect", Author: "???", Year: "2017"},
       models.Book{ID: 5, Title: "Darya Gonchareva", Author: "Riya Albert", Year: "2018"})

// create the endpoints and api handlers
    router := mux.NewRouter()

    controller := controllers.Controller{}

// use these handlers to work with local non-persistent data
    router.HandleFunc("/bookslocal", getBooksLocal).Methods("GET")
    router.HandleFunc("/booklocal/{id}", getBookLocal).Methods("GET")
    router.HandleFunc("/bookslocal", addBookLocal).Methods("POST")
    router.HandleFunc("/bookslocal", updateBookLocal).Methods("PUT")
    router.HandleFunc("/bookslocal/{id}", removeBookLocal).Methods("DELETE")

// use these handlers to demo elephantsql cloud db persistent data
    router.HandleFunc("/bookssql/", controller.GetAllBooksSql(db)).Methods("GET")
    router.HandleFunc("/booksql/{id}", controller.GetBookSql(db)).Methods("GET")
    router.HandleFunc("/bookssql/", controller.AddBookSql(db)).Methods("POST")
    router.HandleFunc("/bookssql/", controller.UpdateBookSql(db)).Methods("PUT")
    router.HandleFunc("/bookssql/{id}", controller.RemoveBookSql(db)).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8080", router))
}

//
// the following handlers take care of local non-persistent data
//

func getBooksLocal(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s - was called- return ALL books\n\n", r.URL.Path)
    log.Println("getBooksLocal - get all books is called")

    json.NewEncoder(w).Encode(bookslocal)

}

func getBookLocal(w http.ResponseWriter, r *http.Request) {
    // use built-in func to extract URL vars
    params := mux.Vars(r)
    log.Println(params)
    
    fmt.Fprintf(w, "%s - was called with params : %s\n", r.URL.Path, params)
    log.Printf("getBookLocal - get book %s\n", params)

    i, _ := strconv.Atoi(params["id"])
//    log.Println(reflect.TypeOf(i))

    for _, booklocal := range bookslocal {
        if booklocal.ID == i {
            json.NewEncoder(w).Encode(&booklocal)
        }
    }
}

// use www.getpostman.com app to send in POST - UPDATE - DELETE requests
//     POST FORMAT
//          {
//              "id": x,
//              "title": "...",
//              "author": "...",
//              "year": "..."
//          }
//
// after you post, should see a result in POSTMAN - also do a /books in local browser to verify new set of books
//
func addBookLocal(w http.ResponseWriter, r *http.Request) {
    var book models.Book

    // decode the request sent in, into our expected format
    // add to existing list
    // return new list to caller
    _ = json.NewDecoder(r.Body).Decode(&book)
    bookslocal = append(bookslocal, book)
    
    // return to caller
    json.NewEncoder(w).Encode(bookslocal)
}

// same input format from POSTMAN as above - except change METHOD=POST AND USE AN EXISTING ID
func updateBookLocal(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s - was called w a PUT METHOD\n", r.URL.Path)
    log.Println("updateBookLocal - update 1 book is called")

    // decode the request body & move it into our expected format - same as addBook EXCEPT
    // we're going to ASSUME on PUTs that the ID will be for an existing record - NO ERROR CHECKING - JUST ASSUME IT'S A VALID ID
    var book models.Book
    json.NewDecoder(r.Body).Decode(&book)

    // find the matching id and replace it
    for i, item := range bookslocal {
        if item.ID == book.ID {
            bookslocal[i] = book
        }
    }

    // return updated set to caller
    json.NewEncoder(w).Encode(bookslocal)
}

// takes an ID from the URL and removes that entry from the current list
// AGAIN NO ERROR PROCESSING DONE - ASSUMES WE HAVE A MATCHING ID
func removeBookLocal(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%s - was called\n", r.URL.Path)
    log.Println("removeBookLocal - delete 1 book is called")

    params := mux.Vars(r)
    
    id, _ := strconv.Atoi(params["id"])

    // find the matching id and delete existing entry
    for i, item := range bookslocal {
        if item.ID == id {
            bookslocal = append(bookslocal[:i], bookslocal[i+1:]...)
        }
    }
    // return updated set to caller
    json.NewEncoder(w).Encode(bookslocal)
}

