package expense

import "database/sql"

type Expense struct {
	Id     int      `json:id`
	Title  string   `json:title`
	Amount float64  `json:amount`
	Note   string   `json:note`
	Tags   []string `json:tags`
}

type Err struct {
	Message string `json:"message"`
}

type handler struct {
	db *sql.DB
}
