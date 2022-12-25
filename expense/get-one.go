package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetExpenseById(c echo.Context) error {
	id := c.Param("id")

	stmt, err := db.Prepare("SELECT * FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query user statment" + err.Error()})
	}

	res := stmt.QueryRow(id)
	detailExp := Expense{}

	err = res.Scan(&detailExp.Id, &detailExp.Title, &detailExp.Amount, &detailExp.Note, pq.Array(&detailExp.Tags))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't sacan user: " + err.Error()})
	}
	return c.JSON(http.StatusOK, detailExp)

}
