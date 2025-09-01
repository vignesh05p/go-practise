package main

//interface defination
type db interface {
	insert(query string, args ...interface{}) error
}

//supabase
type supabase struct {
	db *sql.db
}

func (s supabase) insert(query string, args ...interface{}) error {
	_, err := s.db.Exec(query, args...)
	if err != nil {
		return err
	}
	fmt.Println("Supabase insert successful")
	return nil
}