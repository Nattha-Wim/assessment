package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h handler) CreateExpense(c echo.Context) error {
	var exp Expense
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := h.db.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id", exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))
	err = row.Scan(&exp.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, exp)
}

// func (h handler) Create(c echo.Context) error {
// 	var exp Expense

// 	err := c.Bind(&exp)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	query := "INSERT INTO expenses (id, title, amount, note, tags) VALUES (?, ?, ?, ?, ?)"
// 	stmt, err := h.db.PrepareContext(ctx, query)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	//log.Println("uUUUU ", exp)
// 	_, err = stmt.ExecContext(ctx, exp.Id, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))
// 	return c.JSON(http.StatusCreated, exp)
// }
