package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h handler) UpdateExpenses(c echo.Context) error {
	var err error
	rowID := c.Param("id")

	exp := Expense{}
	err = c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	//stmt, err := h.db.Prepare("UPDATE expenses SET title=?, amount=?, note=?, tags=? WHERE id = ?")
	stmt, err := h.db.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if _, err := stmt.Exec(rowID, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, exp)
}

// func (h handler) Update(c echo.Context) error {
// 	var err error
// 	rowID := c.Param("id")

// 	exp := Expense{}
// 	err = c.Bind(&user)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	//query := "UPDATE expenses SET title = \\$1, amount = \\$2, notes = \\$3, tags = \\$4 WHERE id = \\$5"
// 	query := "UPDATE expenses SET title = ?, amount = ?, note = ?, tags = ? WHERE id = ?"
// 	stmt, err := h.db.PrepareContext(ctx, query)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()
// 	_, err = stmt.ExecContext(ctx, rowID, exp.Title, exp.Amount, exp.Note, user.Tags)

// 	return c.JSON(http.StatusOK, exp)
// }
