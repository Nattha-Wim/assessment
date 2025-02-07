package expense

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func (h handler) InitDB() {
	createTb := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`

	if _, err := h.db.Exec(createTb); err != nil {
		log.Fatal("Can't create table", err)
	}
	fmt.Println("Create table seccess")
}

func NewApplication(db *sql.DB) *handler {
	return &handler{db}
}
