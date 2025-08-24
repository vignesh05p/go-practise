package main

//interface defination
type db interface {
	insert(query string, args ...interface{}) error
}

//supabase
type supabase struct {
	db *sql.db
}
