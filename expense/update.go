package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h handler) UpdateExpenses(c echo.Context) error {

	rowID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "id should be int " + err.Error()})
	}

	exp := Expense{}
	err = c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	query := "UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id = $1"
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	exp.Id = rowID
	if _, err := stmt.Exec(exp.Id, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags)); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, exp)
}
