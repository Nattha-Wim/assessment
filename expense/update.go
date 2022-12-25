package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func UpdateExpenses(c echo.Context) error {
	id := c.Param("id")

	var exp Expense
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := db.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if _, err := stmt.Exec(id, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, exp)
}
