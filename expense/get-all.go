package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetAllExpenses(c echo.Context) error {

	stmt, err := db.Prepare("SELECT * FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query all expenses"})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query all expenses: " + err.Error()})
	}

	detailExp := []Expense{}
	for rows.Next() {
		var exp Expense
		err = rows.Scan(&exp.Id, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query expenses: " + err.Error()})
		}
		detailExp = append(detailExp, exp)
	}
	return c.JSON(http.StatusOK, detailExp)
}
