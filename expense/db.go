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

	_, err := h.db.Exec(createTb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
	fmt.Println("create table seccess")
}

func NewApplication(db *sql.DB) *handler {
	return &handler{db}
}
