package bookRepository

import (
    "books-list/models"
    "database/sql"
    "log"
)

type BookRepository struct{}

// error logging
func logFatal(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func (b BookRepository) GetAllBooksSql(db *sql.DB, book models.Book, bookssql []models.Book) []models.Book {
    rows, err := db.Query("select * from books")
    logFatal(err)

    // dealing with a db - close the connection after we get what we need out of there
    defer rows.Close()

    // get each element returned - ASSUME NO ERRS
    for rows.Next() {
        err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
        logFatal(err)
    
        bookssql = append(bookssql, book)
    }

    return bookssql
}


func (b BookRepository) GetBookSql(db *sql.DB, book models,Book, id int) models.Book {
     rows := db.QueryRow("select * from books where id=$1", id) 
 
     err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
     logFatal(err)
 
     return book
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
func (b BookRepository) AddBookSql(db *sql.DB, book models.Book, id int) models.Book {
     // FOR SOME REASON - the  inertion is working but the ID beingused in POSTMAN is BEING IGNORED - might have something to do with
     // the database requiring UNIQUE IDs, don't know - but the insertion succeeds with a new ID and an error gets returned which
     // was causing the program to bomb - not processing the error kept the program up
     //
  //   err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&bookID)
     db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", book.Title, book.Author, book.Year).Scan(&book.ID)
 //    log.Fatal(err)
     
     // return ID to caller
     return book
}

 
func (b BookRepository) UpdateBookSql(db *sql.DB, book models.Book) int64 {
 
     result, _ := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id;", &book.Title, &book.Author, &book.Year, &book.ID)
 
     // NOTE : having same problem here as in above SQL_ADD call - 'something' about the db call is causing a fatal error to be reported 
     //        but the db is being updated with the information being passed in
     //
     rowsUpdated, _ := result.RowsAffected()
//  log.Fatal(err) 
 
     // return updated set to caller
     return rowsUpdated
}
 
func (b BookRepository) RemoveBookSql(db *sql.DB, id int) int64 {
    result, err := db.Exec("delete from books where id = $1;", id)
    logFatal(err)
 
    rowsDeleted, err := result.RowsAffected()
    logFatal(err)
 
    // return updated set to caller
    return rowsDeleted
}
