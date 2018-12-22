package driver

import (
    "database/sql"
_    "fmt"
    "log"
    "os"

    "github.com/lib/pq"
)

// create our elephantSQL cloud db connection
var db *sql.DB

// error logging
func logFatal(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func ConnectDB() *sql.DB {
    // access and print app-specific environment variables
    log.Println(os.Getenv("APP_ID"))
    log.Println(os.Getenv("APP_SECRET"))
    log.Println(os.Getenv("ELEPHANTSQL_URL"))
 
//  set up db url, connect to it and ping it - report and quit if can't talk to db
    pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
    logFatal(err)
 
    log.Println(pgUrl)
 
//  this has to be an = and not a :=
//  := would be a new local variable declaraion
//  which would override our global def above
//
    db, err = sql.Open("postgres", pgUrl)
    logFatal(err)

    err = db.Ping()
    logFatal(err)

    return db
}

