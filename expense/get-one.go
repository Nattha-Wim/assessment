package expense

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h handler) GetExpenseById(c echo.Context) error {
	rowID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "id should be int " + err.Error()})
	}

	row := h.db.QueryRow("SELECT * FROM expenses WHERE id = $1", rowID)

	detailExp := Expense{}
	err = row.Scan(&detailExp.Id, &detailExp.Title, &detailExp.Amount, &detailExp.Note, pq.Array(&detailExp.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't sacan expense: !!! " + err.Error()})
	}

	return c.JSON(http.StatusOK, detailExp)

}
