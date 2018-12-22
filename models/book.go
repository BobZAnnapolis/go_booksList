package models

// create the model for our book db entries
type Book struct {
    ID      int     `json:id`
    Title   string  `json:title`
    Author  string  `json:author`
    Year    string  `json:year`
}

